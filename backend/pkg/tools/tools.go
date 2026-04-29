package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"pentagi/pkg/config"
	"pentagi/pkg/database"
	"pentagi/pkg/docker"
	"pentagi/pkg/graphiti"
	"pentagi/pkg/providers/embeddings"
	"pentagi/pkg/schema"

	"github.com/docker/docker/api/types/container"
	"github.com/sirupsen/logrus"
	"github.com/vxcontrol/cloud/anonymizer"
	"github.com/vxcontrol/cloud/anonymizer/patterns"
	"github.com/vxcontrol/langchaingo/llms"
	"github.com/vxcontrol/langchaingo/vectorstores/pgvector"
)

type ExecutorHandler func(ctx context.Context, name string, args json.RawMessage) (string, error)

type SummarizeHandler func(ctx context.Context, result string) (string, error)

type Functions struct {
	Token    *string            `form:"token,omitempty" json:"token,omitempty" validate:"omitempty"`
	Disabled []DisableFunction  `form:"disabled,omitempty" json:"disabled,omitempty" validate:"omitempty,valid"`
	Function []ExternalFunction `form:"functions,omitempty" json:"functions,omitempty" validate:"omitempty,valid"`
}

func (f *Functions) Scan(input any) error {
	switch v := input.(type) {
	case string:
		return json.Unmarshal([]byte(v), f)
	case []byte:
		return json.Unmarshal(v, f)
	case json.RawMessage:
		return json.Unmarshal(v, f)
	}
	return fmt.Errorf("unsupported type of input value to scan")
}

type DisableFunction struct {
	Name    string   `form:"name" json:"name" validate:"required"`
	Context []string `form:"context,omitempty" json:"context,omitempty" validate:"omitempty,dive,oneof=agent adviser coder searcher generator memorist enricher reporter assistant,required"`
}

type ExternalFunction struct {
	Name    string        `form:"name" json:"name" validate:"required"`
	URL     string        `form:"url" json:"url" validate:"required,url" example:"https://example.com/api/v1/function"`
	Timeout *int64        `form:"timeout,omitempty" json:"timeout,omitempty" validate:"omitempty,min=1" example:"60"`
	Context []string      `form:"context,omitempty" json:"context,omitempty" validate:"omitempty,dive,oneof=agent adviser coder searcher generator memorist enricher reporter assistant,required"`
	Schema  schema.Schema `form:"schema" json:"schema" validate:"required" swaggertype:"object"`
}

type FunctionInfo struct {
	Name   string
	Schema string
}

type Tool interface {
	Handle(ctx context.Context, name string, args json.RawMessage) (string, error)
	IsAvailable() bool
}

type ScreenshotProvider interface {
	PutScreenshot(ctx context.Context, name, url string, taskID, subtaskID *int64) (int64, error)
}

type AgentLogProvider interface {
	PutLog(
		ctx context.Context,
		initiator, executor database.MsgchainType,
		task, result string,
		taskID, subtaskID *int64,
	) (int64, error)
}

type MsgLogProvider interface {
	PutMsg(
		ctx context.Context,
		msgType database.MsglogType,
		taskID, subtaskID *int64,
		streamID int64,
		thinking, msg string,
	) (int64, error)
	UpdateMsgResult(
		ctx context.Context,
		msgID, streamID int64,
		result string,
		resultFormat database.MsglogResultFormat,
	) error
}

type SearchLogProvider interface {
	PutLog(
		ctx context.Context,
		initiator database.MsgchainType,
		executor database.MsgchainType,
		engine database.SearchengineType,
		query string,
		result string,
		taskID *int64,
		subtaskID *int64,
	) (int64, error)
}

type TermLogProvider interface {
	PutMsg(
		ctx context.Context,
		msgType database.TermlogType,
		msg string,
		containerID int64,
		taskID, subtaskID *int64,
	) (int64, error)
}

type VectorStoreLogProvider interface {
	PutLog(
		ctx context.Context,
		initiator database.MsgchainType,
		executor database.MsgchainType,
		filter json.RawMessage,
		query string,
		action database.VecstoreActionType,
		result string,
		taskID *int64,
		subtaskID *int64,
	) (int64, error)
}

type MemoryHealth struct {
	Enabled bool
	Reason  string
}

type flowToolsExecutor struct {
	flowID int64
	scp    ScreenshotProvider
	alp    AgentLogProvider
	mlp    MsgLogProvider
	slp    SearchLogProvider
	tlp    TermLogProvider
	vslp   VectorStoreLogProvider

	db             database.Querier
	cfg            *config.Config
	store          *pgvector.Store
	graphitiClient *graphiti.Client
	image          string
	docker         docker.DockerClient
	primaryID      int64
	primaryLID     string
	functions      *Functions
	replacer       anonymizer.Replacer
	memoryHealth   MemoryHealth

	definitions map[string]llms.FunctionDefinition
	handlers    map[string]ExecutorHandler
}

type ContextToolsExecutor interface {
	Tools() []llms.Tool
	Execute(ctx context.Context, streamID int64, id, name, obsName, thinking string, args json.RawMessage) (string, error)
	IsBarrierFunction(name string) bool
	IsFunctionExists(name string) bool
	GetBarrierToolNames() []string
	GetBarrierTools() []FunctionInfo
	GetToolSchema(name string) (*schema.Schema, error)
}

type CustomExecutorConfig struct {
	TaskID      *int64
	SubtaskID   *int64
	Builtin     []string
	Definitions []llms.FunctionDefinition
	Handlers    map[string]ExecutorHandler
	Barriers    []string
	Summarizer  SummarizeHandler
}

type AssistantExecutorConfig struct {
	UseAgents  bool
	Adviser    ExecutorHandler
	Coder      ExecutorHandler
	Installer  ExecutorHandler
	Memorist   ExecutorHandler
	Pentester  ExecutorHandler
	Searcher   ExecutorHandler
	Summarizer SummarizeHandler
}

type PrimaryExecutorConfig struct {
	TaskID     int64
	SubtaskID  int64
	Barrier    ExecutorHandler
	Adviser    ExecutorHandler
	Coder      ExecutorHandler
	Installer  ExecutorHandler
	Memorist   ExecutorHandler
	Pentester  ExecutorHandler
	Searcher   ExecutorHandler
	Summarizer SummarizeHandler
}

type InstallerExecutorConfig struct {
	TaskID            *int64
	SubtaskID         *int64
	Adviser           ExecutorHandler
	Memorist          ExecutorHandler
	Searcher          ExecutorHandler
	MaintenanceResult ExecutorHandler
	Summarizer        SummarizeHandler
}

type CoderExecutorConfig struct {
	TaskID     *int64
	SubtaskID  *int64
	Adviser    ExecutorHandler
	Installer  ExecutorHandler
	Memorist   ExecutorHandler
	Searcher   ExecutorHandler
	CodeResult ExecutorHandler
	Summarizer SummarizeHandler
}

type PentesterExecutorConfig struct {
	TaskID     *int64
	SubtaskID  *int64
	Adviser    ExecutorHandler
	Coder      ExecutorHandler
	Installer  ExecutorHandler
	Memorist   ExecutorHandler
	Searcher   ExecutorHandler
	HackResult ExecutorHandler
	Summarizer SummarizeHandler
}

type SearcherExecutorConfig struct {
	TaskID       *int64
	SubtaskID    *int64
	Memorist     ExecutorHandler
	SearchResult ExecutorHandler
	Summarizer   SummarizeHandler
}

type GeneratorExecutorConfig struct {
	TaskID      int64
	Memorist    ExecutorHandler
	Searcher    ExecutorHandler
	SubtaskList ExecutorHandler
}

type RefinerExecutorConfig struct {
	TaskID       int64
	Memorist     ExecutorHandler
	Searcher     ExecutorHandler
	SubtaskPatch ExecutorHandler
}

type MemoristExecutorConfig struct {
	TaskID       *int64
	SubtaskID    *int64
	SearchResult ExecutorHandler
	Summarizer   SummarizeHandler
}

type EnricherExecutorConfig struct {
	TaskID         *int64
	SubtaskID      *int64
	EnricherResult ExecutorHandler
	Summarizer     SummarizeHandler
}

type ReporterExecutorConfig struct {
	TaskID       *int64
	SubtaskID    *int64
	ReportResult ExecutorHandler
}

type FlowToolsExecutor interface {
	SetFlowID(flowID int64)
	SetImage(image string)
	SetEmbedder(embedder embeddings.Embedder)
	GetMemoryHealth() MemoryHealth
	SetFunctions(functions *Functions)
	SetScreenshotProvider(sp ScreenshotProvider)
	SetAgentLogProvider(alp AgentLogProvider)
	SetMsgLogProvider(mlp MsgLogProvider)
	SetSearchLogProvider(slp SearchLogProvider)
	SetTermLogProvider(tlp TermLogProvider)
	SetVectorStoreLogProvider(vslp VectorStoreLogProvider)
	SetGraphitiClient(client *graphiti.Client)

	Prepare(ctx context.Context) error
	CancelRunningCommands(ctx context.Context) error
	Release(ctx context.Context) error
	GetCustomExecutor(cfg CustomExecutorConfig) (ContextToolsExecutor, error)
	GetAssistantExecutor(cfg AssistantExecutorConfig) (ContextToolsExecutor, error)
	GetPrimaryExecutor(cfg PrimaryExecutorConfig) (ContextToolsExecutor, error)
	GetInstallerExecutor(cfg InstallerExecutorConfig) (ContextToolsExecutor, error)
	GetCoderExecutor(cfg CoderExecutorConfig) (ContextToolsExecutor, error)
	GetPentesterExecutor(cfg PentesterExecutorConfig) (ContextToolsExecutor, error)
	GetSearcherExecutor(cfg SearcherExecutorConfig) (ContextToolsExecutor, error)
	GetGeneratorExecutor(cfg GeneratorExecutorConfig) (ContextToolsExecutor, error)
	GetRefinerExecutor(cfg RefinerExecutorConfig) (ContextToolsExecutor, error)
	GetMemoristExecutor(cfg MemoristExecutorConfig) (ContextToolsExecutor, error)
	GetEnricherExecutor(cfg EnricherExecutorConfig) (ContextToolsExecutor, error)
	GetReporterExecutor(cfg ReporterExecutorConfig) (ContextToolsExecutor, error)
}

func NewFlowToolsExecutor(
	db database.Querier,
	cfg *config.Config,
	docker docker.DockerClient,
	functions *Functions,
	flowID int64,
) (FlowToolsExecutor, error) {
	allPatterns, err := patterns.LoadPatterns(patterns.PatternListTypeAll)
	if err != nil {
		return nil, fmt.Errorf("failed to load all patterns: %v", err)
	}

	// combine with config secret patterns
	allPatterns.Patterns = append(allPatterns.Patterns, cfg.GetSecretPatterns()...)

	replacer, err := anonymizer.NewReplacer(allPatterns.Regexes(), allPatterns.Names())
	if err != nil {
		return nil, fmt.Errorf("failed to create replacer: %v", err)
	}

	return &flowToolsExecutor{
		db:        db,
		docker:    docker,
		functions: functions,
		replacer:  replacer,
		cfg:       cfg,
		flowID:    flowID,
		memoryHealth: MemoryHealth{
			Enabled: false,
			Reason:  "embedder is not initialized",
		},
		definitions: make(map[string]llms.FunctionDefinition),
		handlers:    make(map[string]ExecutorHandler),
	}, nil
}

func (fte *flowToolsExecutor) SetFlowID(flowID int64) {
	fte.flowID = flowID
}

func (fte *flowToolsExecutor) SetImage(image string) {
	fte.image = image
}

func (fte *flowToolsExecutor) SetEmbedder(embedder embeddings.Embedder) {
	if embedder == nil || !embedder.IsAvailable() {
		if fte.store != nil {
			fte.store.Close()
			fte.store = nil
		}
		fte.memoryHealth = MemoryHealth{
			Enabled: false,
			Reason:  "embedder is unavailable",
		}
		logrus.WithField("flow_id", fte.flowID).Warn("vector memory disabled: embedder is unavailable")
		return
	}

	if fte.store != nil {
		fte.store.Close()
	}

	store, err := pgvector.New(
		context.Background(),
		pgvector.WithConnectionURL(fte.cfg.DatabaseURL),
		pgvector.WithEmbedder(embedder),
	)
	if err != nil {
		fte.store = nil
		fte.memoryHealth = MemoryHealth{
			Enabled: false,
			Reason:  fmt.Sprintf("pgvector initialization failed: %v", err),
		}
		logrus.WithError(err).WithField("flow_id", fte.flowID).Warn("vector memory disabled: failed to initialize pgvector store")
		return
	}

	fte.store = &store
	fte.memoryHealth = MemoryHealth{Enabled: true}
}

func (fte *flowToolsExecutor) GetMemoryHealth() MemoryHealth {
	return fte.memoryHealth
}

func (fte *flowToolsExecutor) SetFunctions(functions *Functions) {
	fte.functions = functions
}

func (fte *flowToolsExecutor) SetScreenshotProvider(scp ScreenshotProvider) {
	fte.scp = scp
}

func (fte *flowToolsExecutor) SetAgentLogProvider(alp AgentLogProvider) {
	fte.alp = alp
}

func (fte *flowToolsExecutor) SetMsgLogProvider(mlp MsgLogProvider) {
	fte.mlp = mlp
}

func (fte *flowToolsExecutor) SetSearchLogProvider(slp SearchLogProvider) {
	fte.slp = slp
}

func (fte *flowToolsExecutor) SetTermLogProvider(tlp TermLogProvider) {
	fte.tlp = tlp
}

func (fte *flowToolsExecutor) SetVectorStoreLogProvider(vslp VectorStoreLogProvider) {
	fte.vslp = vslp
}

func (fte *flowToolsExecutor) SetGraphitiClient(client *graphiti.Client) {
	fte.graphitiClient = client
}

func (fte *flowToolsExecutor) Prepare(ctx context.Context) error {
	if cnt, err := fte.db.GetFlowPrimaryContainer(ctx, fte.flowID); err == nil {
		currentImage := strings.ToLower(strings.TrimSpace(cnt.Image))
		desiredImage := strings.ToLower(strings.TrimSpace(fte.image))
		imageMatches := currentImage != "" && desiredImage != "" && currentImage == desiredImage

		switch cnt.Status {
		case database.ContainerStatusRunning:
			if imageMatches {
				fte.primaryID = cnt.ID
				fte.primaryLID = cnt.LocalID.String
				return nil
			}

			fte.docker.RemoveContainer(ctx, cnt.LocalID.String, cnt.ID)
		default:
			fte.docker.RemoveContainer(ctx, cnt.LocalID.String, cnt.ID)
		}
	}

	capAdd := []string{"NET_RAW"}
	if fte.cfg.DockerNetAdmin {
		capAdd = append(capAdd, "NET_ADMIN")
	}

	securityOpts := []string{}
	if fte.cfg.DockerSeccompUnconfined {
		securityOpts = append(securityOpts, "seccomp=unconfined")
	}
	if fte.cfg.DockerApparmorUnconfined {
		securityOpts = append(securityOpts, "apparmor=unconfined")
	}

	containerName := PrimaryTerminalName(fte.flowID)
	cnt, err := fte.docker.RunContainer(
		ctx,
		containerName,
		database.ContainerTypePrimary,
		fte.flowID,
		&container.Config{
			Image: fte.image,
			Entrypoint: []string{
				"sh",
				"-lc",
				"if command -v apt-get >/dev/null 2>&1; then " +
					"apt-get update && " +
					"apt-get -y upgrade && " +
					"apt-get install -y --no-install-recommends openssh-client openssh-server && " +
					"rm -rf /var/lib/apt/lists/*; " +
					"fi; " +
					"exec tail -f /dev/null",
			},
		},
		&container.HostConfig{
			CapAdd:      capAdd,
			SecurityOpt: securityOpts,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to launch container '%s': %w", containerName, err)
	}

	fte.primaryID = cnt.ID
	fte.primaryLID = cnt.LocalID.String

	return nil
}

func (fte *flowToolsExecutor) Release(ctx context.Context) error {
	if fte.store != nil {
		fte.store.Close()
	}

	containers, err := fte.db.GetFlowContainers(ctx, fte.flowID)
	if err != nil {
		// DB unavailable — fall back to removing only the primary container.
		// Log the DB error but don't surface it if the removal itself succeeds.
		logrus.WithContext(ctx).WithError(err).Warnf(
			"[Release] failed to get flow containers for flow %d, falling back to primary container removal",
			fte.flowID,
		)
		containerName := PrimaryTerminalName(fte.flowID)
		if removeErr := fte.docker.RemoveContainer(ctx, fte.primaryLID, fte.primaryID); removeErr != nil {
			return fmt.Errorf("failed to purge container '%s': %w", containerName, removeErr)
		}
		return nil
	}

	var errs []error
	for _, cnt := range containers {
		if removeErr := fte.docker.RemoveContainer(ctx, cnt.LocalID.String, cnt.ID); removeErr != nil {
			errs = append(errs, fmt.Errorf("failed to purge container '%s' (id=%d): %w", cnt.Name, cnt.ID, removeErr))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("release flow %d: %w", fte.flowID, errors.Join(errs...))
	}

	var firstErr error
	for _, c := range containers {
		if c.Status == database.ContainerStatusDeleted {
			continue
		}
		localID := c.LocalID.String
		if !c.LocalID.Valid || localID == "" {
			continue
		}
		if removeErr := fte.docker.RemoveContainer(ctx, localID, c.ID); removeErr != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("failed to purge container '%s': %w", c.Name, removeErr)
			}
		}
	}

	return firstErr
}

func (fte *flowToolsExecutor) CancelRunningCommands(ctx context.Context) error {
	if fte.primaryLID == "" {
		return nil
	}

	return cancelRunningExecCommands(ctx, fte.docker, fte.primaryLID)
}

func (fte *flowToolsExecutor) GetCustomExecutor(cfg CustomExecutorConfig) (ContextToolsExecutor, error) {
	if len(cfg.Definitions) != len(cfg.Handlers) {
		return nil, fmt.Errorf("definitions and handlers must have the same length")
	}

	for _, def := range cfg.Definitions {
		if _, ok := cfg.Handlers[def.Name]; !ok {
			return nil, fmt.Errorf("handler for function %s not found", def.Name)
		}
	}

	for _, builtin := range cfg.Builtin {
		if def, ok := fte.definitions[builtin]; !ok {
			return nil, fmt.Errorf("builtin function %s not found", builtin)
		} else {
			cfg.Definitions = append(cfg.Definitions, def)
			cfg.Handlers[builtin] = fte.handlers[builtin]
		}
	}

	barriers := make(map[string]struct{})
	for _, barrier := range cfg.Barriers {
		if _, ok := fte.handlers[barrier]; !ok {
			return nil, fmt.Errorf("barrier function %s not found", barrier)
		}
		barriers[barrier] = struct{}{}
	}

	return &customExecutor{
		flowID:      fte.flowID,
		taskID:      cfg.TaskID,
		subtaskID:   cfg.SubtaskID,
		mlp:         fte.mlp,
		vslp:        fte.vslp,
		db:          fte.db,
		store:       fte.store,
		definitions: cfg.Definitions,
		handlers:    cfg.Handlers,
		barriers:    barriers,
		summarizer:  cfg.Summarizer,
	}, nil
}

func (fte *flowToolsExecutor) GetAssistantExecutor(cfg AssistantExecutorConfig) (ContextToolsExecutor, error) {
	if cfg.Adviser == nil {
		return nil, fmt.Errorf("adviser handler is required")
	}

	if cfg.Coder == nil {
		return nil, fmt.Errorf("coder handler is required")
	}

	if cfg.Installer == nil {
		return nil, fmt.Errorf("installer handler is required")
	}

	if cfg.Memorist == nil {
		return nil, fmt.Errorf("memorist handler is required")
	}

	if cfg.Pentester == nil {
		return nil, fmt.Errorf("pentester handler is required")
	}

	if cfg.Searcher == nil {
		return nil, fmt.Errorf("searcher handler is required")
	}

	container, err := fte.db.GetFlowPrimaryContainer(context.Background(), fte.flowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get container %d: %w", fte.flowID, err)
	}

	term := NewTerminalTool(
		fte.flowID, nil, nil,
		container.ID,
		container.LocalID.String,
		fte.docker,
		fte.tlp,
	)

	ts := fte.newToolSet(nil, nil, cfg.Summarizer).
		addTerminalFrom(term).
		addBrowser().
		addHttpClient()

	if cfg.UseAgents {
		ts.add(AdviceToolName, cfg.Adviser).
			add(CoderToolName, cfg.Coder).
			add(MaintenanceToolName, cfg.Installer).
			add(MemoristToolName, cfg.Memorist).
			add(PentesterToolName, cfg.Pentester).
			add(SearchToolName, cfg.Searcher)
	} else {
		ts.addMemory().
			addGuideSearch().
			addAnswerSearch().
			addCodeSearch().
			addWebSearch().
			addSploitus()
	}

	return ts.build(), nil
}

func (fte *flowToolsExecutor) GetPrimaryExecutor(cfg PrimaryExecutorConfig) (ContextToolsExecutor, error) {
	if cfg.Barrier == nil {
		return nil, fmt.Errorf("barrier (done) handler is required")
	}

	if cfg.Adviser == nil {
		return nil, fmt.Errorf("adviser handler is required")
	}

	if cfg.Coder == nil {
		return nil, fmt.Errorf("coder handler is required")
	}

	if cfg.Installer == nil {
		return nil, fmt.Errorf("installer handler is required")
	}

	if cfg.Memorist == nil {
		return nil, fmt.Errorf("memorist handler is required")
	}

	if cfg.Pentester == nil {
		return nil, fmt.Errorf("pentester handler is required")
	}

	if cfg.Searcher == nil {
		return nil, fmt.Errorf("searcher handler is required")
	}

	ts := fte.newToolSet(&cfg.TaskID, &cfg.SubtaskID, cfg.Summarizer).
		addBarrier(FinalyToolName, cfg.Barrier).
		add(AdviceToolName, cfg.Adviser).
		add(CoderToolName, cfg.Coder).
		add(MaintenanceToolName, cfg.Installer).
		add(MemoristToolName, cfg.Memorist).
		add(PentesterToolName, cfg.Pentester).
		add(SearchToolName, cfg.Searcher)

	if fte.cfg.AskUser {
		ts.addBarrier(AskUserToolName, cfg.Barrier)
	}

	return ts.build(), nil
}

func (fte *flowToolsExecutor) GetInstallerExecutor(cfg InstallerExecutorConfig) (ContextToolsExecutor, error) {
	if cfg.MaintenanceResult == nil {
		return nil, fmt.Errorf("maintenance result handler is required")
	}

	if cfg.Adviser == nil {
		return nil, fmt.Errorf("adviser handler is required")
	}

	if cfg.Memorist == nil {
		return nil, fmt.Errorf("memorist handler is required")
	}

	if cfg.Searcher == nil {
		return nil, fmt.Errorf("searcher handler is required")
	}

	container, err := fte.db.GetFlowPrimaryContainer(context.Background(), fte.flowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get container %d: %w", fte.flowID, err)
	}

	term := NewTerminalTool(
		fte.flowID,
		cfg.TaskID,
		cfg.SubtaskID,
		container.ID,
		container.LocalID.String,
		fte.docker,
		fte.tlp,
	)

	return fte.newToolSet(cfg.TaskID, cfg.SubtaskID, cfg.Summarizer).
		addBarrier(MaintenanceResultToolName, cfg.MaintenanceResult).
		add(AdviceToolName, cfg.Adviser).
		add(MemoristToolName, cfg.Memorist).
		add(SearchToolName, cfg.Searcher).
		addTerminalFrom(term).
		addBrowser().
		addGuide().
		build(), nil
}

func (fte *flowToolsExecutor) GetCoderExecutor(cfg CoderExecutorConfig) (ContextToolsExecutor, error) {
	if cfg.CodeResult == nil {
		return nil, fmt.Errorf("code result handler is required")
	}

	if cfg.Adviser == nil {
		return nil, fmt.Errorf("adviser handler is required")
	}

	if cfg.Installer == nil {
		return nil, fmt.Errorf("installer handler is required")
	}

	if cfg.Memorist == nil {
		return nil, fmt.Errorf("memorist handler is required")
	}

	if cfg.Searcher == nil {
		return nil, fmt.Errorf("searcher handler is required")
	}

	return fte.newToolSet(cfg.TaskID, cfg.SubtaskID, cfg.Summarizer).
		addBarrier(CodeResultToolName, cfg.CodeResult).
		add(AdviceToolName, cfg.Adviser).
		add(MaintenanceToolName, cfg.Installer).
		add(MemoristToolName, cfg.Memorist).
		add(SearchToolName, cfg.Searcher).
		addBrowser().
		addCode().
		addGraphiti().
		build(), nil
}

func (fte *flowToolsExecutor) GetPentesterExecutor(cfg PentesterExecutorConfig) (ContextToolsExecutor, error) {
	if cfg.HackResult == nil {
		return nil, fmt.Errorf("hack result handler is required")
	}

	if cfg.Adviser == nil {
		return nil, fmt.Errorf("adviser handler is required")
	}

	if cfg.Coder == nil {
		return nil, fmt.Errorf("coder handler is required")
	}

	if cfg.Installer == nil {
		return nil, fmt.Errorf("installer handler is required")
	}

	if cfg.Memorist == nil {
		return nil, fmt.Errorf("memorist handler is required")
	}

	if cfg.Searcher == nil {
		return nil, fmt.Errorf("searcher handler is required")
	}

	container, err := fte.db.GetFlowPrimaryContainer(context.Background(), fte.flowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get container %d: %w", fte.flowID, err)
	}

	term := NewTerminalTool(
		fte.flowID,
		cfg.TaskID,
		cfg.SubtaskID,
		container.ID,
		container.LocalID.String,
		fte.docker,
		fte.tlp,
	)

	return fte.newToolSet(cfg.TaskID, cfg.SubtaskID, cfg.Summarizer).
		addBarrier(HackResultToolName, cfg.HackResult).
		add(AdviceToolName, cfg.Adviser).
		add(CoderToolName, cfg.Coder).
		add(MaintenanceToolName, cfg.Installer).
		add(MemoristToolName, cfg.Memorist).
		add(SearchToolName, cfg.Searcher).
		addTerminalFrom(term).
		addBrowser().
		addGuide().
		addGraphiti().
		addSploitus().
		addHttpClient().
		build(), nil
}

func (fte *flowToolsExecutor) GetSearcherExecutor(cfg SearcherExecutorConfig) (ContextToolsExecutor, error) {
	if cfg.SearchResult == nil {
		return nil, fmt.Errorf("search result handler is required")
	}

	if cfg.Memorist == nil {
		return nil, fmt.Errorf("memorist handler is required")
	}

	return fte.newToolSet(cfg.TaskID, cfg.SubtaskID, cfg.Summarizer).
		addBarrier(SearchResultToolName, cfg.SearchResult).
		add(MemoristToolName, cfg.Memorist).
		addBrowser().
		addWebSearch().
		addSploitus().
		addAnswers().
		build(), nil
}

func (fte *flowToolsExecutor) GetGeneratorExecutor(cfg GeneratorExecutorConfig) (ContextToolsExecutor, error) {
	if cfg.SubtaskList == nil {
		return nil, fmt.Errorf("subtask list handler is required")
	}

	if cfg.Memorist == nil {
		return nil, fmt.Errorf("memorist handler is required")
	}

	container, err := fte.db.GetFlowPrimaryContainer(context.Background(), fte.flowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get container %d: %w", fte.flowID, err)
	}

	term := NewTerminalTool(
		fte.flowID,
		&cfg.TaskID,
		nil,
		container.ID,
		container.LocalID.String,
		fte.docker,
		fte.tlp,
	)

	return fte.newToolSet(&cfg.TaskID, nil, nil).
		add(MemoristToolName, cfg.Memorist).
		add(SearchToolName, cfg.Searcher).
		addBarrier(SubtaskListToolName, cfg.SubtaskList).
		addTerminalFrom(term).
		addBrowser().
		build(), nil
}

func (fte *flowToolsExecutor) GetRefinerExecutor(cfg RefinerExecutorConfig) (ContextToolsExecutor, error) {
	if cfg.SubtaskPatch == nil {
		return nil, fmt.Errorf("subtask patch handler is required")
	}

	if cfg.Memorist == nil {
		return nil, fmt.Errorf("memorist handler is required")
	}

	container, err := fte.db.GetFlowPrimaryContainer(context.Background(), fte.flowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get container %d: %w", fte.flowID, err)
	}

	term := NewTerminalTool(
		fte.flowID,
		&cfg.TaskID,
		nil,
		container.ID,
		container.LocalID.String,
		fte.docker,
		fte.tlp,
	)

	return fte.newToolSet(&cfg.TaskID, nil, nil).
		add(MemoristToolName, cfg.Memorist).
		add(SearchToolName, cfg.Searcher).
		addBarrier(SubtaskPatchToolName, cfg.SubtaskPatch).
		addTerminalFrom(term).
		addBrowser().
		build(), nil
}

func (fte *flowToolsExecutor) GetMemoristExecutor(cfg MemoristExecutorConfig) (ContextToolsExecutor, error) {
	if cfg.SearchResult == nil {
		return nil, fmt.Errorf("search result handler is required")
	}

	container, err := fte.db.GetFlowPrimaryContainer(context.Background(), fte.flowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get container %d: %w", fte.flowID, err)
	}

	term := NewTerminalTool(
		fte.flowID,
		cfg.TaskID,
		cfg.SubtaskID,
		container.ID,
		container.LocalID.String,
		fte.docker,
		fte.tlp,
	)

	return fte.newToolSet(cfg.TaskID, cfg.SubtaskID, cfg.Summarizer).
		addBarrier(MemoristResultToolName, cfg.SearchResult).
		addTerminalFrom(term).
		addMemory().
		addGraphiti().
		build(), nil
}

func (fte *flowToolsExecutor) GetEnricherExecutor(cfg EnricherExecutorConfig) (ContextToolsExecutor, error) {
	if cfg.EnricherResult == nil {
		return nil, fmt.Errorf("enricher result handler is required")
	}

	container, err := fte.db.GetFlowPrimaryContainer(context.Background(), fte.flowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get container %d: %w", fte.flowID, err)
	}

	term := NewTerminalTool(
		fte.flowID,
		cfg.TaskID,
		cfg.SubtaskID,
		container.ID,
		container.LocalID.String,
		fte.docker,
		fte.tlp,
	)

	return fte.newToolSet(cfg.TaskID, cfg.SubtaskID, cfg.Summarizer).
		addBarrier(EnricherResultToolName, cfg.EnricherResult).
		addTerminalFrom(term).
		addMemory().
		addGraphiti().
		addBrowser().
		build(), nil
}

func (fte *flowToolsExecutor) GetReporterExecutor(cfg ReporterExecutorConfig) (ContextToolsExecutor, error) {
	if cfg.ReportResult == nil {
		return nil, fmt.Errorf("report result handler is required")
	}

	return fte.newToolSet(cfg.TaskID, cfg.SubtaskID, nil).
		addBarrier(ReportResultToolName, cfg.ReportResult).
		build(), nil
}

func enrichLogrusFields(flowID int64, taskID, subtaskID *int64, fields logrus.Fields) logrus.Fields {
	if fields == nil {
		fields = logrus.Fields{}
	}

	fields["flow_id"] = flowID
	if taskID != nil {
		fields["task_id"] = *taskID
	}
	if subtaskID != nil {
		fields["subtask_id"] = *subtaskID
	}

	return fields
}
