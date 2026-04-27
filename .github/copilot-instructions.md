<todos title="Short-term orchestration retry/backoff mitigation" rule="Review steps frequently throughout the conversation and DO NOT stop between steps unless they explicitly require it.">
- [x] patch-performer-retries: Patch `backend/pkg/providers/performer.go` to shorten agent-chain retries and backoff 🔴
  _Reduce `maxRetriesToCallAgentChain` from 3→2 and `delayBetweenRetries` from 5s→2s to fail faster on bad LLM responses_
- [x] gofmt-commit: Run `gofmt` and commit the change 🟡
- [x] build-docker-image: Build Docker image `local/pentagi:latest` 🟡
- [x] restart-pentagi: Restart `pentagi` service via Docker Compose 🟡
- [-] tail-logs: Tail `pentagi` logs to verify reduced retries and failures 🟡
</todos>

<!-- Auto-generated todo section -->
<!-- Add your custom Copilot instructions below -->
