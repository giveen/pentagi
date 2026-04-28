SHELL := /bin/bash

IMAGE_TAG ?= latest
PENTAGI_READY_URL ?= https://localhost:8443/api/v1/info
PENTAGI_READY_TIMEOUT ?= 180
PENTAGI_READY_INTERVAL ?= 2
PENTAGI_PERSIST_DIR ?= ./data

.PHONY: docker-build docker-up docker-up-persist docker-down

define wait_for_pentagi_ready
	@echo "Waiting for PentAGI API readiness at $(PENTAGI_READY_URL) ..."
	@elapsed=0; \
	while true; do \
		if curl -ksS --max-time 5 "$(PENTAGI_READY_URL)" > /dev/null; then \
			echo "PentAGI API is ready."; \
			break; \
		fi; \
		if [[ $$elapsed -ge $(PENTAGI_READY_TIMEOUT) ]]; then \
			echo "Timed out waiting for PentAGI API readiness after $(PENTAGI_READY_TIMEOUT)s."; \
			exit 1; \
		fi; \
		sleep $(PENTAGI_READY_INTERVAL); \
		elapsed=$$((elapsed + $(PENTAGI_READY_INTERVAL))); \
	done
endef

docker-build:
	docker build -t local/pentagi:$(IMAGE_TAG) .

docker-up:
	docker compose up -d
	$(call wait_for_pentagi_ready)

docker-up-persist:
	@mkdir -p "$(PENTAGI_PERSIST_DIR)/pentagi" "$(PENTAGI_PERSIST_DIR)/ssl" "$(PENTAGI_PERSIST_DIR)/ollama" "$(PENTAGI_PERSIST_DIR)/postgres" "$(PENTAGI_PERSIST_DIR)/scraper-ssl"
	PENTAGI_DATA_DIR=$(PENTAGI_PERSIST_DIR)/pentagi \
	PENTAGI_SSL_DIR=$(PENTAGI_PERSIST_DIR)/ssl \
	PENTAGI_OLLAMA_DIR=$(PENTAGI_PERSIST_DIR)/ollama \
	PENTAGI_POSTGRES_DATA_DIR=$(PENTAGI_PERSIST_DIR)/postgres \
	SCRAPER_SSL_DIR=$(PENTAGI_PERSIST_DIR)/scraper-ssl \
	docker compose up -d
	$(call wait_for_pentagi_ready)

docker-down:
	docker compose down --remove-orphans
