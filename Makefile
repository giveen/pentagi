SHELL := /bin/bash

PROJECT_ROOT := $(CURDIR)
SEARXNG_PATH ?=
SEARXNG_PORT ?= 8445
IMAGE_TAG ?= latest
# On-disk BuildKit layer cache — survives docker system prune / builder prune.
# Override with CACHE_DIR=/your/path make docker-build to use a different location.
CACHE_DIR ?= /tmp/pentagi-build-cache

COMPOSE_FILES := \
	-f docker-compose.yml \
	-f docker-compose-graphiti.yml \
	-f docker-compose.graphiti.override.yml \
	-f docker-compose.graphiti.ports.override.yml \
	-f docker-compose-langfuse.yml \
	-f docker-compose-observability.yml

.PHONY: docker-build docker-build-nocache docker-up docker-down kali-build kali-build-nocache

# Build using a persistent local BuildKit cache (fast rebuilds, survives prune).
docker-build:
	docker buildx build \
		--cache-from type=local,src=$(CACHE_DIR) \
		--cache-to   type=local,dest=$(CACHE_DIR),mode=max \
		--load \
		-t local/pentagi:$(IMAGE_TAG) .

# Full rebuild with no cache (equivalent to the old docker build).
docker-build-nocache:
	docker build --no-cache -t local/pentagi:$(IMAGE_TAG) .

# Build the kali-pentest image used as the pentest sandbox container.
kali-build:
	docker buildx build \
		--cache-from type=local,src=$(CACHE_DIR)/kali \
		--cache-to   type=local,dest=$(CACHE_DIR)/kali,mode=max \
		--load \
		-t local/kali-pentest:$(IMAGE_TAG) build/kali-pentest/

kali-build-nocache:
	docker build --no-cache -t local/kali-pentest:$(IMAGE_TAG) build/kali-pentest/

docker-up:
	@if [ -n "$(SEARXNG_PATH)" ] && [ -d "$(SEARXNG_PATH)" ]; then \
		echo "Starting SearXNG from $(SEARXNG_PATH) on host port $(SEARXNG_PORT)"; \
		cd "$(SEARXNG_PATH)" && SEARXNG_PORT="$(SEARXNG_PORT)" docker compose up -d --force-recreate; \
	elif [ -n "$(SEARXNG_PATH)" ]; then \
		echo "Warning: SEARXNG_PATH not found: $(SEARXNG_PATH). Skipping SearXNG startup."; \
	fi
	cd "$(PROJECT_ROOT)" && docker compose $(COMPOSE_FILES) up -d

docker-down:
	cd "$(PROJECT_ROOT)" && docker compose $(COMPOSE_FILES) down --remove-orphans
	@if [ -n "$(SEARXNG_PATH)" ] && [ -d "$(SEARXNG_PATH)" ]; then \
		cd "$(SEARXNG_PATH)" && docker compose down; \
	fi
