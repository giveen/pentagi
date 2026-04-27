<todos title="Cloud Stripper: Remove 7 Cloud Providers" rule="Review steps frequently throughout the conversation and DO NOT stop between steps unless they explicitly require it.">
- [ ] 1: Delete anthropic, gemini, bedrock, deepseek, glm, kimi, qwen directories (~700 LOC) 🔴
- [ ] 2: Update imports, remove NewProviderController blocks, refactor switch cases (~120 LOC) 🔴
- [ ] 3: Remove 7 ProviderType and ProviderName constants 🔴
- [ ] 4: Remove config blocks for all 7 cloud providers from config.go (~26 LOC) 🔴
- [ ] 5: Delete screen ID constants from wizard/models/types.go (~7 LOC) 🟡
- [ ] 6: Remove provider init functions, switch cases, field validators (~445 LOC) 🔴
- [ ] 7: Remove provider names, help texts, descriptions from locale.go (~264 LOC) 🟡
- [ ] 8: Delete 7 icon .tsx files and remove from provider-icon.tsx (~42 LOC) 🟡
- [ ] 9: Remove providers from GraphQL ProviderModels, run codegen (~200 LOC) 🔴
- [ ] 10: Create migration to remove 7 provider_type enum values (~20 LOC) 🔴
- [ ] 11: Remove sections from config.md, README.md, .env.example (~380 LOC) 🟢
- [ ] 12: Refactor ctester, ftester, remove provider test cases (~100 LOC) 🟡
- [ ] 13: go build ./cmd/pentagi && go test ./... 🔴
- [ ] 14: npm run graphql:generate && npm run build 🔴
- [ ] 15: Run installer, verify 3 providers (OpenAI, Ollama, Custom), create test flow 🔴
- [ ] 16: Commit all changes with summary of 1,654 LOC removed across 45+ files 🟡
</todos>

<!-- Auto-generated todo section -->
<!-- Add your custom Copilot instructions below -->
