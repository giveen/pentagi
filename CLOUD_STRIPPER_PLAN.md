# Cloud Stripper: Remove Specialized Cloud LLM Providers

## Objective
Remove Anthropic, Gemini, Bedrock, DeepSeek, GLM, Kimi, and Qwen providers. Keep OpenAI, Ollama, and Custom to support llama.cpp via OpenAI-compatible API.

## Providers to Remove
1. **Anthropic** (Claude models)
2. **Gemini** (Google AI)
3. **Bedrock** (AWS)
4. **DeepSeek** (Chinese AI)
5. **GLM** (Zhipu AI)
6. **Kimi** (Moonshot AI)
7. **Qwen** (Alibaba AI)

## Providers to Keep
1. **OpenAI** - Acts as universal interface (works with llama.cpp via OpenAI-compatible API)
2. **Ollama** - Direct Ollama server support
3. **Custom** - Fallback HTTP endpoint for anything else

---

# Implementation Checklist

## Phase 1: Backend Provider Packages

### Remove Directories (7 total: ~700 LOC)
- [ ] `backend/pkg/providers/anthropic/` - Delete entire directory
  - config.yml, models.yml, anthropic.go
  - EST: 150 LOC
  
- [ ] `backend/pkg/providers/gemini/` - Delete entire directory
  - config.yml, models.yml, gemini.go
  - EST: 160 LOC
  
- [ ] `backend/pkg/providers/bedrock/` - Delete entire directory (complex, has AWS SDK)
  - config.yml, models.yml, bedrock.go
  - EST: 200 LOC
  
- [ ] `backend/pkg/providers/deepseek/` - Delete entire directory
  - config.yml, models.yml, deepseek.go
  - EST: 120 LOC
  
- [ ] `backend/pkg/providers/glm/` - Delete entire directory
  - config.yml, models.yml, glm.go
  - EST: 100 LOC
  
- [ ] `backend/pkg/providers/kimi/` - Delete entire directory
  - config.yml, models.yml, kimi.go
  - EST: 100 LOC
  
- [ ] `backend/pkg/providers/qwen/` - Delete entire directory
  - config.yml, models.yml, qwen.go
  - EST: 100 LOC

### Update Provider Imports
- [ ] `backend/pkg/providers/providers.go`
  - Remove imports (line 24-30):
    ```go
    "pentagi/pkg/providers/anthropic"
    "pentagi/pkg/providers/bedrock"
    "pentagi/pkg/providers/deepseek"
    "pentagi/pkg/providers/gemini"
    "pentagi/pkg/providers/glm"
    "pentagi/pkg/providers/kimi"
    "pentagi/pkg/providers/qwen"
    ```
  - EST: 7 lines

- [ ] `backend/pkg/providers/providers.go` - `NewProviderController()` function
  - Remove 7 conditional provider creation blocks (~60 LOC)
  - Keep only: OpenAI, Anthropic, Gemini, Bedrock checks → delete all
  - Keep: Ollama, Custom checks
  - EST: 60 LOC removed

- [ ] `backend/pkg/providers/providers.go` - `buildProviderFromConfig()` switch
  - Remove 7 cases from switch statement (~30 LOC)
  - EST: 30 LOC removed

- [ ] `backend/pkg/providers/providers.go` - `NewProvider()` switch
  - Remove 7 cases from switch statement (~30 LOC)
  - EST: 30 LOC removed

### Update Provider Types
- [ ] `backend/pkg/providers/provider/provider.go`
  - Remove provider type constants (7 lines):
    ```go
    ProviderAnthropic ProviderType = "anthropic"
    ProviderGemini    ProviderType = "gemini"
    ProviderBedrock   ProviderType = "bedrock"
    ProviderDeepSeek  ProviderType = "deepseek"
    ProviderGLM       ProviderType = "glm"
    ProviderKimi      ProviderType = "kimi"
    ProviderQwen      ProviderType = "qwen"
    ```
  - Remove provider name constants (7 lines)
  - EST: 14 lines

**Phase 1 Total: ~150 LOC removed**

---

## Phase 2: Configuration & Environment

### `backend/pkg/config/config.go` (~70 LOC)
- [ ] Remove Anthropic config block:
  ```go
  AnthropicAPIKey    string
  AnthropicServerURL string
  ```
  - EST: 2 lines

- [ ] Remove Gemini config block:
  ```go
  GeminiAPIKey    string
  GeminiServerURL string
  ```
  - EST: 2 lines

- [ ] Remove Bedrock config block:
  ```go
  BedrockRegion       string
  BedrockDefaultAuth  bool
  BedrockBearerToken  string
  BedrockAccessKey    string
  BedrockSecretKey    string
  BedrockSessionToken string
  BedrockServerURL    string
  ```
  - EST: 7 lines

- [ ] Remove DeepSeek config block:
  ```go
  DeepSeekAPIKey    string
  DeepSeekServerURL string
  DeepSeekProvider  string
  ```
  - EST: 3 lines

- [ ] Remove GLM config block:
  ```go
  GLMAPIKey    string
  GLMServerURL string
  GLMProvider  string
  ```
  - EST: 3 lines

- [ ] Remove Kimi config block:
  ```go
  KimiAPIKey    string
  KimiServerURL string
  KimiProvider  string
  ```
  - EST: 3 lines

- [ ] Remove Qwen config block:
  ```go
  QwenAPIKey    string
  QwenServerURL string
  QwenProvider  string
  ```
  - EST: 3 lines

**Phase 2 Total: ~26 LOC removed**

---

## Phase 3: Installer/Wizard UI

### `backend/cmd/installer/wizard/models/types.go` (~8 LOC)
- [ ] Remove screen IDs:
  ```go
  LLMProviderAnthropicScreen ScreenID = "llm_provider_form§anthropic"
  LLMProviderGeminiScreen    ScreenID = "llm_provider_form§gemini"
  LLMProviderBedrockScreen   ScreenID = "llm_provider_form§bedrock"
  LLMProviderDeepSeekScreen  ScreenID = "llm_provider_form§deepseek"
  LLMProviderGLMScreen       ScreenID = "llm_provider_form§glm"
  LLMProviderKimiScreen      ScreenID = "llm_provider_form§kimi"
  LLMProviderQwenScreen      ScreenID = "llm_provider_form§qwen"
  ```
  - EST: 7 lines

**Phase 3a Total: ~7 LOC removed**

### `backend/cmd/installer/wizard/models/llm_provider_form.go` (~350 LOC)
- [ ] Remove provider form functions:
  ```go
  func (m *LLMProviderFormModel) initAnthropic() { ... }
  func (m *LLMProviderFormModel) initGemini() { ... }
  func (m *LLMProviderFormModel) initBedrock() { ... }
  func (m *LLMProviderFormModel) initDeepSeek() { ... }
  func (m *LLMProviderFormModel) initGLM() { ... }
  func (m *LLMProviderFormModel) initKimi() { ... }
  func (m *LLMProviderFormModel) initQwen() { ... }
  ```
  - EST: ~50 LOC/provider × 7 = 350 LOC

- [ ] Remove cases from `switch` on provider selection
  - EST: ~30 LOC

- [ ] Remove from `getDefaultBaseURL()` switch:
  ```go
  case LLMProviderAnthropic:
      return "https://api.anthropic.com/v1"
  case LLMProviderGemini:
      return "https://generativelanguage.googleapis.com/v1beta"
  case LLMProviderBedrock:
      return ""
  case LLMProviderDeepSeek:
      return "https://api.deepseek.com"
  case LLMProviderGLM:
      return "https://api.z.ai/api/paas/v4"
  case LLMProviderKimi:
      return "https://api.moonshot.ai/v1"
  case LLMProviderQwen:
      return "https://dashscope-us.aliyuncs.com/compatible-mode/v1"
  ```
  - EST: 15 LOC

- [ ] Remove from `getFieldValidators()` switch
  - EST: ~50 LOC

**Phase 3b Total: ~445 LOC removed**

### `backend/cmd/installer/wizard/locale/locale.go` (~300 LOC)
- [ ] Remove provider name constants:
  ```go
  LLMProviderAnthropic = "Anthropic"
  LLMProviderGemini = "Google Gemini"
  LLMProviderBedrock = "AWS Bedrock"
  LLMProviderDeepSeek = "DeepSeek"
  LLMProviderGLM = "GLM Zhipu AI"
  LLMProviderKimi = "Kimi Moonshot AI"
  LLMProviderQwen = "Qwen Alibaba Cloud"
  ```
  - EST: 7 lines

- [ ] Remove help text constants (7 × 5-10 lines):
  ```go
  LLMFormAnthropicHelp = "..."
  LLMFormGeminiHelp = "..."
  LLMFormBedrockHelp = "..."
  // ... etc
  ```
  - EST: ~250 LOC

- [ ] Remove description constants:
  ```go
  LLMProviderAnthropicDesc = "..."
  LLMProviderGeminiDesc = "..."
  // ... etc
  ```
  - EST: 7 lines

**Phase 3c Total: ~264 LOC removed**

**Phase 3 Total: ~716 LOC removed**

---

## Phase 4: Frontend

### `frontend/src/components/icons/provider-icon.tsx` (~50 LOC)
- [ ] Remove 7 icon imports:
  ```tsx
  import { AnthropicIcon } from "./anthropic"
  import { GeminiIcon } from "./gemini"
  import { BedrockIcon } from "./bedrock"
  import { DeepSeekIcon } from "./deepseek"
  import { GlmIcon } from "./glm"
  import { KimiIcon } from "./kimi"
  import { QwenIcon } from "./qwen"
  ```
  - EST: 7 lines

- [ ] Remove 7 cases from switch statement
  - EST: ~35 LOC

**Phase 4a Total: ~42 LOC removed**

### Delete Icon Files (7 total: ~20 LOC)
- [ ] `frontend/src/components/icons/anthropic.tsx` - DELETE
- [ ] `frontend/src/components/icons/gemini.tsx` - DELETE
- [ ] `frontend/src/components/icons/bedrock.tsx` - DELETE
- [ ] `frontend/src/components/icons/deepseek.tsx` - DELETE
- [ ] `frontend/src/components/icons/glm.tsx` - DELETE
- [ ] `frontend/src/components/icons/kimi.tsx` - DELETE
- [ ] `frontend/src/components/icons/qwen.tsx` - DELETE

**Phase 4b Total: ~20 LOC removed (from icon files)**

### GraphQL Schema & Types (~200 LOC)
- [ ] `frontend/graphql-schema.graphql`
  - Remove 7 provider entries from `ProviderModels` type
  - EST: ~20 LOC

- [ ] `frontend/src/graphql/types.ts` (auto-generated)
  - Will be regenerated after schema update
  - EST: ~100 LOC reduced

- [ ] Settings UI component
  - Remove 7 provider form sections from `provider-settings.tsx`
  - EST: ~80 LOC

**Phase 4c Total: ~200 LOC removed**

**Phase 4 Total: ~262 LOC removed**

---

## Phase 5: Database & Migrations

### Provider Type Enum (~20 LOC)
- [ ] Create new migration: `backend/migrations/sql/YYYYMMDD_remove_cloud_providers.sql`
  - Alter enum type to remove 7 provider types
  - EST: 15 LOC

- [ ] Update `backend/pkg/database/` models if hardcoded provider type checks exist
  - EST: 5 LOC

**Phase 5 Total: ~20 LOC removed**

---

## Phase 6: Documentation

### `backend/docs/config.md` (~200 LOC)
- [ ] Remove Anthropic section
- [ ] Remove Gemini section
- [ ] Remove Bedrock section
- [ ] Remove DeepSeek section
- [ ] Remove GLM section
- [ ] Remove Kimi section
- [ ] Remove Qwen section
- [ ] Keep OpenAI, Ollama, Custom sections

**Phase 6a Total: ~200 LOC removed**

### `README.md` (~150 LOC)
- [ ] Remove provider setup examples (7 × 15-25 LOC)
- [ ] Simplify provider section
- [ ] Simplify cost/performance comparison

**Phase 6b Total: ~150 LOC removed**

### `.env.example` (~30 LOC)
- [ ] Remove 7 provider config blocks

**Phase 6c Total: ~30 LOC removed**

**Phase 6 Total: ~380 LOC removed**

---

## Phase 7: Tests & Utilities

### `backend/cmd/ctester/main.go` (~50 LOC)
- [ ] Remove 7 provider cases from `createProvider()` function

**Phase 7a Total: ~50 LOC removed**

### `backend/cmd/ftester/main.go` (~30 LOC)
- [ ] Remove 7 provider initialization blocks

**Phase 7b Total: ~30 LOC removed**

### Unit Tests
- [ ] Remove test files for 7 providers (if any exist in `*_test.go` files)
- [ ] Update integration tests to exclude cloud providers

**Phase 7c Total: ~20 LOC removed**

**Phase 7 Total: ~100 LOC removed**

---

## Summary: Code Reduction Estimates

| Phase | Component | LOC Removed | Notes |
|-------|-----------|------------|-------|
| 1 | Backend Provider Packages | ~150 | Delete 7 dirs (~700 LOC) + refactor imports |
| 2 | Config & Environment | ~26 | Remove env var blocks |
| 3 | Installer/Wizard | ~716 | Remove forms + help text + locales |
| 4 | Frontend | ~262 | Remove icons + GraphQL schema entries |
| 5 | Database | ~20 | Update enum migration |
| 6 | Documentation | ~380 | Remove provider guides + examples |
| 7 | Tests & Utils | ~100 | Update test utilities |
| **TOTAL** | **ALL** | **~1,654 LOC** | **~45% reduction from ~3,700 cloud LOC** |

### Files to Delete (7 directories + 10 files)
- `backend/pkg/providers/anthropic/` (entire)
- `backend/pkg/providers/gemini/` (entire)
- `backend/pkg/providers/bedrock/` (entire)
- `backend/pkg/providers/deepseek/` (entire)
- `backend/pkg/providers/glm/` (entire)
- `backend/pkg/providers/kimi/` (entire)
- `backend/pkg/providers/qwen/` (entire)
- `frontend/src/components/icons/anthropic.tsx`
- `frontend/src/components/icons/gemini.tsx`
- `frontend/src/components/icons/bedrock.tsx`
- `frontend/src/components/icons/deepseek.tsx`
- `frontend/src/components/icons/glm.tsx`
- `frontend/src/components/icons/kimi.tsx`
- `frontend/src/components/icons/qwen.tsx`

### Files to Modify (30+ files)
- Backend: 8 files
- Frontend: 5 files
- Docs: 3 files
- Tests: 4 files
- Config/Migrations: 2 files

---

## Implementation Order (Recommended)

1. **Phase 1** (Backend packages) - Safest, clearest diffs
2. **Phase 2** (Config) - Low risk, straightforward
3. **Phase 5** (Database) - Do before Phase 3 to avoid runtime enum errors
4. **Phase 3** (Installer) - UI changes, moderate complexity
5. **Phase 4** (Frontend) - Requires GraphQL codegen
6. **Phase 7** (Tests) - Last, after core is working
7. **Phase 6** (Docs) - Final polish

## Testing Strategy

- [ ] Build backend: `go build ./cmd/pentagi`
- [ ] Run tests: `go test ./...`
- [ ] Run installer: Test LLM provider screen shows only 3 options
- [ ] Frontend: `npm run graphql:generate && npm run build`
- [ ] Manual: Create new flow and verify only OpenAI/Ollama/Custom appear

## Notes

- **No breaking changes to API** - Provider abstraction stays intact
- **Backward compat**: Existing flows using removed providers will error on load (handled gracefully)
- **Rollback**: Each phase is independently commitable
- **Dependencies**: Bedrock removal eliminates AWS SDK (~5MB from binary)
