// Package autotuner provides a singleton LLM parameter tuner that maps agent
// roles to sampling presets and enforces a deterministic "sober" failover mode
// when a task is retried after a failure.
//
// Usage:
//
//	ctx = autotuner.WithAttemptCount(ctx, idx+1)
//	opts := autotuner.Get().CallOptions(pconfig.OptionsTypePentester, autotuner.AttemptCountFromCtx(ctx), modelName)
//	// append opts after provider config options so they take precedence
package autotuner

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"pentagi/pkg/providers/pconfig"

	"github.com/vxcontrol/langchaingo/llms"
	"gopkg.in/yaml.v3"
)

// contextKey is an unexported type for context keys in this package.
type contextKey int

const (
	attemptCountKey contextKey = iota
	modelNameKey
)

// WithAttemptCount returns a child context carrying the current attempt count.
// An attempt count of 1 means first try; 2+ activates sober/failover mode.
func WithAttemptCount(ctx context.Context, count int) context.Context {
	return context.WithValue(ctx, attemptCountKey, count)
}

// AttemptCountFromCtx reads the attempt count stored by WithAttemptCount.
// Returns 1 (first attempt) when no value is present.
func AttemptCountFromCtx(ctx context.Context) int {
	if v, ok := ctx.Value(attemptCountKey).(int); ok && v > 0 {
		return v
	}
	return 1
}

// WithModelName returns a child context carrying the model name (with provider prefix).
// The AutoTuner uses this to detect the model family and select the correct preset map.
func WithModelName(ctx context.Context, modelName string) context.Context {
	return context.WithValue(ctx, modelNameKey, modelName)
}

// ModelNameFromCtx reads the model name stored by WithModelName.
// Returns an empty string when no value is present.
func ModelNameFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(modelNameKey).(string); ok {
		return v
	}
	return ""
}

// Profile groups agent roles by their sampling personality.
type Profile int

const (
	ProfileLogic    Profile = iota // Logic & Syntax — deterministic
	ProfileStrategy                // Strategic & Critical — balanced
	ProfileCreative                // Creative Chaos — high entropy
)

func (p Profile) String() string {
	switch p {
	case ProfileLogic:
		return "Logic & Syntax"
	case ProfileStrategy:
		return "Strategic & Critical"
	case ProfileCreative:
		return "Creative Chaos"
	default:
		return "Unknown"
	}
}

// ansiColor returns the ANSI escape sequence for the profile's console colour.
func (p Profile) ansiColor() string {
	switch p {
	case ProfileLogic:
		return "\033[32m" // green
	case ProfileStrategy:
		return "\033[33m" // yellow
	case ProfileCreative:
		return "\033[31m" // red
	default:
		return "\033[0m"
	}
}

const ansiReset = "\033[0m"

// ModelFamily identifies the parameter space of a model family.
// Different model families have significantly different optimal sampling parameters;
// Qwen3 in particular uses presence_penalty=1.5 as its primary diversity mechanism
// (instead of frequency_penalty), top_k=20 for all roles, and repetition_penalty=1.0
// always (its DeltaNet+HybridAttention architecture handles repetition natively).
type ModelFamily int

const (
	FamilyGeneric ModelFamily = iota // Generic / OpenAI-compatible defaults
	FamilyQwen3                      // Qwen3 and Qwen3.5 models
)

// detectFamily returns the ModelFamily for a model name (with or without provider prefix).
// Detection is case-insensitive substring matching, so "Qwen/Qwen3-32B" and
// "vllm/Qwen/Qwen3.5-27B-FP8" both resolve to FamilyQwen3.
func detectFamily(modelName string) ModelFamily {
	lower := strings.ToLower(modelName)
	if strings.Contains(lower, "qwen3") {
		return FamilyQwen3
	}
	return FamilyGeneric
}

// Params holds the sampling parameters that the AutoTuner manages.
// RepetitionPenalty is on the llama.cpp / vLLM scale (1.0 = neutral, >1.0 = penalise).
// FrequencyPenalty and PresencePenalty are on the OpenAI scale (0.0 = neutral, >0 = penalise).
// MinP filters tokens whose probability is less than MinP × (probability of the most likely token);
// a value of 0.0 disables min-p filtering (Qwen3 explicitly sets this to 0.0).
// MinLength, MaxLength, and ReasoningEffort are intentionally omitted:
//   - MinLength/MaxLength are task-length concerns, not role-sampling concerns; set them in
//     the provider YAML (max_tokens field) where they can be calibrated per use-case.
//   - ReasoningEffort changes the entire request format and has provider-specific temperature
//     constraints (Anthropic forces temp=1.0; OpenAI ignores temperature); tune it per-agent
//     in the provider YAML rather than applying it globally here.
type Params struct {
	Temperature       float64
	TopP              float64
	TopK              int
	MinP              float64 // 0.0 = disabled; Qwen3 explicitly sets this to 0.0
	FrequencyPenalty  float64
	PresencePenalty   float64
	RepetitionPenalty float64 // 1.0 = neutral; only values > 1.0 apply a penalty
}

// soberParams are applied when attempt_count >= 2 for FamilyGeneric (OpenAI-style) models,
// forcing the model into a highly deterministic mode regardless of profile.
var soberParams = Params{
	Temperature:       0.1,
	TopP:              0.9,
	TopK:              40,
	MinP:              0.0,
	FrequencyPenalty:  0.0,
	PresencePenalty:   0.0,
	RepetitionPenalty: 1.05, // mild — still prevents stuck loops but keeps structured output correct
}

// qwen3SoberParams are applied when attempt_count >= 2 for FamilyQwen3 models.
// Qwen3 "precise coding" parameters are used as the safe/deterministic fallback:
//   - temperature=0.6 is Qwen3's own recommendation for coding/exact tasks (the lowest
//     useful value; going below 0.6 degrades output on a model trained at temp≥0.6)
//   - presence_penalty=0.0 removes all diversity pressure in sober mode
//   - repetition_penalty=1.0 always neutral — Qwen3's architecture handles repetition natively
var qwen3SoberParams = Params{
	Temperature:       0.6,
	TopP:              0.9,
	TopK:              20,
	MinP:              0.0,
	FrequencyPenalty:  0.0,
	PresencePenalty:   0.0,
	RepetitionPenalty: 1.0,
}

// presetMap maps each ProviderOptionsType to its Profile and default Params.
//
// Hardware note: top_k=40 for Strategic, top_k=100 for Creative are tuned for
// RTX 5090-class VRAM bandwidth.
var presetMap = map[pconfig.ProviderOptionsType]struct {
	profile Profile
	params  Params
}{
	// ── Logic & Syntax ──────────────────────────────────────────────────────
	// RepetitionPenalty=1.0: keywords like `if`, `for`, `{}`, variable names MUST repeat.
	// FrequencyPenalty=0: same reason — penalising token frequency breaks structured output.
	pconfig.OptionsTypeCoder: {ProfileLogic, Params{
		Temperature:       0.1,
		TopP:              0.9,
		TopK:              20,
		FrequencyPenalty:  0.0,
		PresencePenalty:   0.0,
		RepetitionPenalty: 1.0,
	}},
	pconfig.OptionsTypeInstaller: {ProfileLogic, Params{
		Temperature:       0.1,
		TopP:              0.9,
		TopK:              20,
		FrequencyPenalty:  0.0,
		PresencePenalty:   0.0,
		RepetitionPenalty: 1.0,
	}},
	pconfig.OptionsTypeSimpleJSON: {ProfileLogic, Params{
		Temperature:       0.1,
		TopP:              0.9,
		TopK:              20,
		FrequencyPenalty:  0.0,
		PresencePenalty:   0.0,
		RepetitionPenalty: 1.0,
	}},
	pconfig.OptionsTypeEnricher: {ProfileLogic, Params{
		Temperature:       0.2,
		TopP:              0.9,
		TopK:              20,
		FrequencyPenalty:  0.0,
		PresencePenalty:   0.0,
		RepetitionPenalty: 1.0,
	}},

	// ── Strategic & Critical ────────────────────────────────────────────────
	// RepetitionPenalty=1.05: gentle nudge away from repetitive analysis phrasing.
	// FrequencyPenalty=0.1: light push for lexical variety in plans and reports.
	pconfig.OptionsTypePrimaryAgent: {ProfileStrategy, Params{
		Temperature:       0.6,
		TopP:              0.95,
		TopK:              40,
		FrequencyPenalty:  0.1,
		PresencePenalty:   0.1,
		RepetitionPenalty: 1.05,
	}},
	pconfig.OptionsTypeAssistant: {ProfileStrategy, Params{
		Temperature:       0.6,
		TopP:              0.95,
		TopK:              40,
		FrequencyPenalty:  0.1,
		PresencePenalty:   0.1,
		RepetitionPenalty: 1.05,
	}},
	pconfig.OptionsTypeAdviser: {ProfileStrategy, Params{
		Temperature:       0.6,
		TopP:              0.95,
		TopK:              40,
		FrequencyPenalty:  0.1,
		PresencePenalty:   0.1,
		RepetitionPenalty: 1.05,
	}},
	// Reflector gets slightly higher frequency penalty because its job is to critically
	// evaluate and it must NOT just restate the previous agent's words back.
	pconfig.OptionsTypeReflector: {ProfileStrategy, Params{
		Temperature:       0.5,
		TopP:              0.9,
		TopK:              40,
		FrequencyPenalty:  0.2,
		PresencePenalty:   0.1,
		RepetitionPenalty: 1.08,
	}},
	pconfig.OptionsTypeRefiner: {ProfileStrategy, Params{
		Temperature:       0.5,
		TopP:              0.9,
		TopK:              40,
		FrequencyPenalty:  0.1,
		PresencePenalty:   0.1,
		RepetitionPenalty: 1.05,
	}},
	pconfig.OptionsTypeSimple: {ProfileStrategy, Params{
		Temperature:       0.6,
		TopP:              0.95,
		TopK:              40,
		FrequencyPenalty:  0.0,
		PresencePenalty:   0.0,
		RepetitionPenalty: 1.05,
	}},

	// ── Creative Chaos ──────────────────────────────────────────────────────
	// RepetitionPenalty=1.1: strongly discourage looping on the same attack surface.
	// High frequency/presence penalties push the model toward novel attack vectors.
	// MinP=0.05: scales the probability cutoff to the top-token confidence, so the
	// model stays sharp when the path is obvious but can be "brilliantly weird"
	// when uncertain — better than a fixed TopP tail cut for adversarial creativity.
	pconfig.OptionsTypePentester: {ProfileCreative, Params{
		Temperature:       0.75,
		TopP:              0.95,
		TopK:              100,
		MinP:              0.05,
		FrequencyPenalty:  0.3,
		PresencePenalty:   0.3,
		RepetitionPenalty: 1.1,
	}},
	pconfig.OptionsTypeGenerator: {ProfileCreative, Params{
		Temperature:       0.75,
		TopP:              0.95,
		TopK:              100,
		MinP:              0.05,
		FrequencyPenalty:  0.2,
		PresencePenalty:   0.2,
		RepetitionPenalty: 1.1,
	}},
	// Searcher: no MinP — factual retrieval needs stable, predictable token selection.
	pconfig.OptionsTypeSearcher: {ProfileCreative, Params{
		Temperature:       0.5,
		TopP:              0.85,
		TopK:              60,
		MinP:              0.0,
		FrequencyPenalty:  0.1,
		PresencePenalty:   0.1,
		RepetitionPenalty: 1.08,
	}},
}

// defaultPreset is the fallback for FamilyGeneric models when a role has no explicit entry.
var defaultPreset = struct {
	profile Profile
	params  Params
}{
	profile: ProfileStrategy,
	params: Params{
		Temperature:       0.6,
		TopP:              0.95,
		TopK:              40,
		FrequencyPenalty:  0.1,
		PresencePenalty:   0.1,
		RepetitionPenalty: 1.05,
	},
}

// qwen3PresetMap maps each ProviderOptionsType to Qwen3-specific Profile and Params.
//
// Key differences from the generic presetMap:
//   - top_k=20: Qwen3's official recommendation for ALL roles (wider beams don't improve quality)
//   - presence_penalty=1.5: the primary diversity knob for Qwen3 (replaces frequency_penalty)
//   - repetition_penalty=1.0: always neutral — Qwen3's DeltaNet+HybridAttention architecture
//     handles repetition natively; raising it above 1.0 degrades output quality
//   - frequency_penalty=0.0: not effective for Qwen3; presence_penalty handles diversity
//   - Logic roles: presence_penalty=0.0 — code/scripts must freely repeat keywords and syntax
//   - temp=1.0 for reasoning/creative roles: official Qwen3 thinking-mode recommendation
//   - temp=0.6 for code roles: official Qwen3 "precise coding" recommendation
var qwen3PresetMap = map[pconfig.ProviderOptionsType]struct {
	profile Profile
	params  Params
}{
	// ── Logic & Syntax ──────────────────────────────────────────────────────
	// presence_penalty=0.0: code and shell scripts must repeat tokens (`if`, `for`,
	// variable names, flags, syntax). temp=0.6 is Qwen3's official "precise" setting.
	pconfig.OptionsTypeCoder: {ProfileLogic, Params{
		Temperature:       0.6,
		TopP:              0.95,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   0.0,
		RepetitionPenalty: 1.0,
	}},
	pconfig.OptionsTypeInstaller: {ProfileLogic, Params{
		Temperature:       0.6,
		TopP:              0.95,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   0.0,
		RepetitionPenalty: 1.0,
	}},
	pconfig.OptionsTypeSimpleJSON: {ProfileLogic, Params{
		Temperature:       0.6,
		TopP:              0.95,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   0.0,
		RepetitionPenalty: 1.0,
	}},
	pconfig.OptionsTypeEnricher: {ProfileLogic, Params{
		Temperature:       0.6,
		TopP:              0.95,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   0.0,
		RepetitionPenalty: 1.0,
	}},

	// ── Strategic & Critical ────────────────────────────────────────────────
	// temp=1.0: official Qwen3 recommendation for reasoning/thinking mode agents.
	// presence_penalty=1.5: Qwen3's recommended diversity value for reasoning roles.
	pconfig.OptionsTypePrimaryAgent: {ProfileStrategy, Params{
		Temperature:       1.0,
		TopP:              0.95,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   1.5,
		RepetitionPenalty: 1.0,
	}},
	pconfig.OptionsTypeAssistant: {ProfileStrategy, Params{
		Temperature:       1.0,
		TopP:              0.95,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   1.5,
		RepetitionPenalty: 1.0,
	}},
	pconfig.OptionsTypeAdviser: {ProfileStrategy, Params{
		Temperature:       1.0,
		TopP:              0.95,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   1.5,
		RepetitionPenalty: 1.0,
	}},
	pconfig.OptionsTypeReflector: {ProfileStrategy, Params{
		Temperature:       1.0,
		TopP:              0.95,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   1.5,
		RepetitionPenalty: 1.0,
	}},
	pconfig.OptionsTypeRefiner: {ProfileStrategy, Params{
		Temperature:       1.0,
		TopP:              0.95,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   1.5,
		RepetitionPenalty: 1.0,
	}},
	// Simple/non-thinking roles: temp=0.7, top_p=0.8 per Qwen3 "general tasks" recommendation.
	pconfig.OptionsTypeSimple: {ProfileStrategy, Params{
		Temperature:       0.7,
		TopP:              0.8,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   1.5,
		RepetitionPenalty: 1.0,
	}},

	// ── Creative Chaos ──────────────────────────────────────────────────────
	// Qwen3 does not differentiate further at top_k=20 for creative roles.
	// presence_penalty=1.5 drives exploration toward novel attack vectors and queries.
	// Pentester/Generator: temp=0.75 reduced from 1.0 to curb hallucinated output and
	// raw escape sequences while still allowing creative attack planning.
	// Searcher uses temp=0.5/top_p=0.8 — factual grounding is more important than
	// diversity; presence_penalty reduced to 0.5 to stop it drifting from the target.
	pconfig.OptionsTypePentester: {ProfileCreative, Params{
		Temperature:       0.75,
		TopP:              0.90,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   1.5,
		RepetitionPenalty: 1.0,
	}},
	pconfig.OptionsTypeGenerator: {ProfileCreative, Params{
		Temperature:       0.75,
		TopP:              0.90,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   1.5,
		RepetitionPenalty: 1.0,
	}},
	pconfig.OptionsTypeSearcher: {ProfileCreative, Params{
		Temperature:       0.5,
		TopP:              0.8,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   0.5,
		RepetitionPenalty: 1.0,
	}},
}

// qwen3DefaultPreset is the fallback for FamilyQwen3 models with unknown roles.
// Uses the "general tasks" recommendation from the official Qwen3 inference guide.
var qwen3DefaultPreset = struct {
	profile Profile
	params  Params
}{
	profile: ProfileStrategy,
	params: Params{
		Temperature:       0.7,
		TopP:              0.8,
		TopK:              20,
		MinP:              0.0,
		FrequencyPenalty:  0.0,
		PresencePenalty:   1.5,
		RepetitionPenalty: 1.0,
	},
}

// overridePatch holds a partial parameter patch loaded from a YAML config file.
// Only non-nil fields are applied on top of the built-in preset; nil fields keep
// the preset value. This lets a config file say just "temperature: 0.4" without
// having to repeat every other parameter.
type overridePatch struct {
	Temperature       *float64
	TopP              *float64
	TopK              *int
	MinP              *float64
	FrequencyPenalty  *float64
	PresencePenalty   *float64
	RepetitionPenalty *float64
}

// fileConfigRole is the per-role section of the YAML config file.
type fileConfigRole struct {
	Temperature       *float64 `yaml:"temperature"`
	TopP              *float64 `yaml:"top_p"`
	TopK              *int     `yaml:"top_k"`
	MinP              *float64 `yaml:"min_p"`
	FrequencyPenalty  *float64 `yaml:"frequency_penalty"`
	PresencePenalty   *float64 `yaml:"presence_penalty"`
	RepetitionPenalty *float64 `yaml:"repetition_penalty"`
}

// fileConfig is the top-level structure of the autotuner YAML config file.
//
// Example file (autotuner.yml):
//
//	roles:
//	  searcher:
//	    temperature: 0.4
//	    presence_penalty: 0.3
//	  generator:
//	    temperature: 0.7
type fileConfig struct {
	Roles map[string]fileConfigRole `yaml:"roles"`
}

// applyPatch overlays the non-nil fields from patch onto p and returns the result.
func applyPatch(p Params, patch overridePatch) Params {
	if patch.Temperature != nil {
		p.Temperature = *patch.Temperature
	}
	if patch.TopP != nil {
		p.TopP = *patch.TopP
	}
	if patch.TopK != nil {
		p.TopK = *patch.TopK
	}
	if patch.MinP != nil {
		p.MinP = *patch.MinP
	}
	if patch.FrequencyPenalty != nil {
		p.FrequencyPenalty = *patch.FrequencyPenalty
	}
	if patch.PresencePenalty != nil {
		p.PresencePenalty = *patch.PresencePenalty
	}
	if patch.RepetitionPenalty != nil {
		p.RepetitionPenalty = *patch.RepetitionPenalty
	}

	return p
}

// AutoTuner is the singleton LLM parameter tuner.
type AutoTuner struct {
	mu        sync.RWMutex
	overrides map[pconfig.ProviderOptionsType]Params        // full-replace runtime overrides (used in tests)
	patches   map[pconfig.ProviderOptionsType]overridePatch // partial overrides from config file
}

var (
	instance *AutoTuner
	once     sync.Once
)

// Get returns the process-wide AutoTuner singleton.
func Get() *AutoTuner {
	once.Do(func() {
		instance = &AutoTuner{
			overrides: make(map[pconfig.ProviderOptionsType]Params),
			patches:   make(map[pconfig.ProviderOptionsType]overridePatch),
		}
	})

	return instance
}

// GetParams returns the LLM sampling Params for the given agent role, attempt
// count, and model name.
//
// Rules:
//   - Role lookup is case-insensitive (normalised to lower snake_case).
//   - Unknown roles fall back to the family-appropriate default preset.
//   - attempt_count >= 2 activates family-specific "sober mode" (deterministic fallback).
//   - Model family is detected from modelName (e.g. "Qwen/Qwen3-32B" → FamilyQwen3).
func (t *AutoTuner) GetParams(role pconfig.ProviderOptionsType, attemptCount int, modelName string) Params {
	r := t.resolve(role, attemptCount, modelName)
	t.log(role, r.profile, r.params, attemptCount)

	return r.params
}

// resolvedResult bundles the fully computed params with their display metadata.
type resolvedResult struct {
	params  Params
	profile Profile
	family  ModelFamily
}

// resolve is the internal core that computes params+metadata without side effects.
func (t *AutoTuner) resolve(role pconfig.ProviderOptionsType, attemptCount int, modelName string) resolvedResult {
	normalised := pconfig.ProviderOptionsType(strings.ToLower(string(role)))
	family := detectFamily(modelName)

	t.mu.RLock()
	override, hasOverride := t.overrides[normalised]
	t.mu.RUnlock()

	var (
		entry = defaultPreset
		sober = soberParams
	)

	switch family {
	case FamilyQwen3:
		entry.profile = qwen3DefaultPreset.profile
		entry.params = qwen3DefaultPreset.params
		sober = qwen3SoberParams
		if preset, ok := qwen3PresetMap[normalised]; ok {
			entry.profile = preset.profile
			entry.params = preset.params
		}
	default:
		if preset, ok := presetMap[normalised]; ok {
			entry.profile = preset.profile
			entry.params = preset.params
		}
	}

	profile := entry.profile
	params := entry.params

	// Apply partial config-file patch on top of the family preset.
	t.mu.RLock()
	patch, hasPatch := t.patches[normalised]
	t.mu.RUnlock()
	if hasPatch {
		params = applyPatch(params, patch)
	}

	// Apply any full runtime override set via Reset() (mainly used in tests).
	if hasOverride {
		params = override
	}

	// Failover "sober" mode: clamp to family-specific deterministic values on retry.
	if attemptCount >= 2 {
		params.Temperature = sober.Temperature
		params.TopP = sober.TopP
		params.TopK = sober.TopK
		params.MinP = sober.MinP
		params.RepetitionPenalty = sober.RepetitionPenalty
		// Leave FrequencyPenalty/PresencePenalty at their profile values — zeroing
		// them on sober mode would remove all diversity nudges and risks degenerate
		// looping in outputs that already contain many repeated tokens.
	}

	return resolvedResult{params: params, profile: profile, family: family}
}

// GetParamsWithMeta is like GetParams but also returns the human-readable profile
// name (e.g. "Logic & Syntax") and model family name (e.g. "Qwen3" or "Generic").
// Intended for display in UI tooltips and settings previews.
func (t *AutoTuner) GetParamsWithMeta(role pconfig.ProviderOptionsType, attemptCount int, modelName string) (Params, string, string) {
	r := t.resolve(role, attemptCount, modelName)
	t.log(role, r.profile, r.params, attemptCount)

	familyStr := "Generic"
	if r.family == FamilyQwen3 {
		familyStr = "Qwen3"
	}

	return r.params, r.profile.String(), familyStr
}

// CallOptions returns the Params as a slice of llms.CallOption ready to be
// appended after provider-config options so they take precedence.
// MinP is only emitted when > 0 to avoid sending the field to providers that
// do not support it (e.g. the real OpenAI API rejects unknown sampling fields).
func (t *AutoTuner) CallOptions(role pconfig.ProviderOptionsType, attemptCount int, modelName string) []llms.CallOption {
	p := t.GetParams(role, attemptCount, modelName)

	opts := make([]llms.CallOption, 0, 7)
	opts = append(opts,
		llms.WithTemperature(p.Temperature),
		llms.WithTopP(p.TopP),
		llms.WithTopK(p.TopK),
		llms.WithFrequencyPenalty(p.FrequencyPenalty),
		llms.WithPresencePenalty(p.PresencePenalty),
		llms.WithRepetitionPenalty(p.RepetitionPenalty),
	)
	if p.MinP > 0 {
		opts = append(opts, llms.WithMinP(p.MinP))
	}

	return opts
}

// LoadConfig reads the YAML file at path and atomically replaces all config-file patches.
// It can be called multiple times (e.g., on file change detection by WatchConfig).
// An empty roles map in the file clears all patches.
func (t *AutoTuner) LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("autotuner: read config %q: %w", path, err)
	}

	var cfg fileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("autotuner: parse config %q: %w", path, err)
	}

	newPatches := make(map[pconfig.ProviderOptionsType]overridePatch, len(cfg.Roles))
	for role, r := range cfg.Roles {
		norm := pconfig.ProviderOptionsType(strings.ToLower(role))
		newPatches[norm] = overridePatch{
			Temperature:       r.Temperature,
			TopP:              r.TopP,
			TopK:              r.TopK,
			MinP:              r.MinP,
			FrequencyPenalty:  r.FrequencyPenalty,
			PresencePenalty:   r.PresencePenalty,
			RepetitionPenalty: r.RepetitionPenalty,
		}
	}

	t.mu.Lock()
	t.patches = newPatches
	t.mu.Unlock()

	fmt.Printf("[Auto-Tuner] Config loaded from %q (%d role overrides)\n", path, len(newPatches))

	return nil
}

// WatchConfig loads the config at path immediately, then polls for file changes every 5 seconds.
// When the mtime changes, the config is reloaded live — no restart required.
// If the file is deleted, all patches are cleared. The goroutine exits when ctx is cancelled.
// Intended to be launched as a goroutine:
//
//	go autotuner.Get().WatchConfig(ctx, os.Getenv("AUTOTUNER_CONFIG"))
func (t *AutoTuner) WatchConfig(ctx context.Context, path string) {
	if err := t.LoadConfig(path); err != nil {
		fmt.Printf("[Auto-Tuner] WARNING: initial config load failed: %v\n", err)
	}

	var lastMod time.Time
	if info, err := os.Stat(path); err == nil {
		lastMod = info.ModTime()
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			info, err := os.Stat(path)
			if err != nil {
				if os.IsNotExist(err) {
					t.mu.Lock()
					if len(t.patches) > 0 {
						t.patches = make(map[pconfig.ProviderOptionsType]overridePatch)
						fmt.Printf("[Auto-Tuner] Config file %q removed — patches cleared\n", path)
					}
					t.mu.Unlock()
					lastMod = time.Time{}
				}
				continue
			}
			if info.ModTime().After(lastMod) {
				lastMod = info.ModTime()
				if err := t.LoadConfig(path); err != nil {
					fmt.Printf("[Auto-Tuner] WARNING: config reload failed: %v\n", err)
				}
			}
		}
	}
}

// Reset restores a role's parameters to the built-in preset values, discarding
// any previous runtime override.  Calling Reset on an unknown role is a no-op.
func (t *AutoTuner) Reset(role pconfig.ProviderOptionsType) {
	normalised := pconfig.ProviderOptionsType(strings.ToLower(string(role)))

	t.mu.Lock()
	delete(t.overrides, normalised)
	t.mu.Unlock()
}

// log writes a coloured one-liner to stdout whenever parameters are resolved.
func (t *AutoTuner) log(
	role pconfig.ProviderOptionsType,
	profile Profile,
	params Params,
	attemptCount int,
) {
	color := profile.ansiColor()
	suffix := ""
	if attemptCount >= 2 {
		suffix = " [SOBER MODE]"
	}

	fmt.Printf(
		"%s[Auto-Tuner] Role: %-14s | Profile: %-22s | Temp: %.2f | TopP: %.2f | TopK: %3d | MinP: %.2f | RepPen: %.2f | FreqPen: %.2f | PresPen: %.2f | Attempt: %d%s%s\n",
		color,
		role,
		profile.String(),
		params.Temperature,
		params.TopP,
		params.TopK,
		params.MinP,
		params.RepetitionPenalty,
		params.FrequencyPenalty,
		params.PresencePenalty,
		attemptCount,
		suffix,
		ansiReset,
	)
}
