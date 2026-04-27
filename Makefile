SHELL := /bin/bash

PROJECT_ROOT := $(CURDIR)
SEARXNG_PATH ?=
SEARXNG_PORT ?= 8445
IMAGE_TAG ?= latest

COMPOSE_FILES := \
	-f docker-compose.yml \
	-f docker-compose-graphiti.yml \
	-f docker-compose.graphiti.override.yml \
	-f docker-compose.graphiti.ports.override.yml \
	-f docker-compose-langfuse.yml \
	-f docker-compose-observability.yml

.PHONY: docker-build docker-up docker-down

docker-build:
	docker build -t local/pentagi:$(IMAGE_TAG) .

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
