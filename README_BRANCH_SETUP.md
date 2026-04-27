# PentAGI Branch Setup Guide

This document explains how to set up and run the current checked-out branch of PentAGI on a fresh machine.

It is intentionally branch-specific and reflects how this repository is wired right now, not just the generic upstream setup.

## What This Branch Expects

This branch currently runs with these assumptions baked into `docker-compose.yml`:

- The main application container is started from `local/pentagi:latest`
- The default pentest worker image is `local/kali-pentest:latest`
- The application expects an OpenAI-compatible LLM endpoint on `http://host.docker.internal:8080/v1`
- The application expects an OpenAI-compatible embeddings endpoint on `http://host.docker.internal:8080/v1`
- The application expects a SearXNG instance reachable as `http://searxng-core:8081`
- The compose file expects an external Docker network named `openui_default`

Because of that, a coworker cannot just clone the repo and run `docker compose up` without preparing those dependencies first.

## Prerequisites

Install these first:

- Docker Engine
- Docker Compose plugin
- Git

Recommended host requirements:

- Linux workstation
- Docker socket access for your user
- Enough disk space for Docker images, PostgreSQL data, and local models

Optional but expected by this branch:

- A local OpenAI-compatible LLM server on port `8080`
- A SearXNG container or service reachable from the `openui_default` network as `searxng-core:8081`

## Repository Checkout

Clone the repository and move into it:

```bash
git clone <your-repo-url>
cd pentagi
git checkout main
```

If you are handing this to a coworker on a different branch, replace `main` with the actual branch name.

## Step 1: Create the External Docker Network

This branch's `docker-compose.yml` attaches the application to an external network named `openui_default`.

Create it once on the machine:

```bash
docker network create openui_default
```

If it already exists, Docker will report that and you can continue.

## Step 2: Build the Kali Pentest Worker Image

This branch points `DOCKER_DEFAULT_IMAGE_FOR_PENTEST` at `local/kali-pentest:latest`, so build that image before starting PentAGI:

```bash
docker build -t local/kali-pentest:latest build/kali-pentest
```

Notes:

- The Kali image definition lives in `build/kali-pentest/Dockerfile`
- This branch includes both `openssh-client` and `openssh-server` in that image

## Step 3: Make Sure the LLM and Embedding Endpoint Exists

This branch is configured to call:

```text
http://host.docker.internal:8080/v1
```

for both chat completions and embeddings.

You have two options:

### Option A: Match This Branch Exactly

Run a local OpenAI-compatible model gateway on the host at port `8080`.

That server must expose:

- chat completions
- embeddings

This is the easiest path if you already use `llama-swap`, `vLLM`, or another OpenAI-compatible local stack.

### Option B: Edit `docker-compose.yml`

If your coworker does not have a host-side LLM service on port `8080`, update these environment variables in `docker-compose.yml` before first run:

- `LLM_SERVER_URL`
- `LLM_SERVER_MODEL`
- `EMBEDDING_URL`
- `EMBEDDING_MODEL`
- `EMBEDDING_PROVIDER`

Without a working model endpoint, the UI may load but flows and assistants will not work correctly.

## Step 4: Make Sure SearXNG Is Reachable

This branch sets:

```text
SEARXNG_URL=http://searxng-core:8081
```

That means the PentAGI container expects a service named `searxng-core` on the `openui_default` Docker network.

You have two options:

### Option A: Provide SearXNG on `openui_default`

Run your existing SearXNG stack and attach it to `openui_default`.

### Option B: Change Search Configuration

If your coworker does not have SearXNG, edit `docker-compose.yml` and either:

- point `SEARXNG_URL` to a reachable instance, or
- disable that path and enable a different search provider supported by the repo

PentAGI will still start without SearXNG, but search-related agent behavior may fail.

## Step 5: Build the Main PentAGI Image

This branch's compose file does **not** define a `build:` block for the `pentagi` service. It only references the prebuilt image `local/pentagi:latest`.

Because of that, you must build the image explicitly:

```bash
docker build -t local/pentagi:latest .
```

Important:

- `docker compose up --build pentagi` is not sufficient on this branch
- if code changes, rebuild with `docker build -t local/pentagi:latest .` again

## Step 6: Start the Stack

Start the application and database:

```bash
docker compose up -d
```

Or, if only the application needs to be recreated after a rebuild:

```bash
docker compose up -d pentagi
```

## Step 7: Verify Startup

Check container state:

```bash
docker compose ps
```

You want to see at minimum:

- `pgvector` healthy
- `pentagi` running

Check logs if needed:

```bash
docker compose logs -f pentagi
```

## Step 8: Log In

Open the web UI:

```text
https://localhost:8443
```

Default credentials on this setup are:

```text
admin@pentagi.com
admin
```

If the browser warns about a self-signed certificate, accept it for local development.

## Day-to-Day Workflow

### Rebuild After Code Changes

If backend or frontend code changed:

```bash
docker build -t local/pentagi:latest .
docker compose up -d pentagi
```

### Rebuild the Kali Worker Image

If `build/kali-pentest/Dockerfile` changed:

```bash
docker build -t local/kali-pentest:latest build/kali-pentest
```

### Stop the Stack

```bash
docker compose down
```

### Stop the Stack and Remove Volumes

Use this only if you want to wipe local database state:

```bash
docker compose down -v
```

## First Troubleshooting Checks

### Error: network `openui_default` not found

Create it manually:

```bash
docker network create openui_default
```

### Error: `local/pentagi:latest` not found

Build the app image explicitly:

```bash
docker build -t local/pentagi:latest .
```

### Error: `local/kali-pentest:latest` not found

Build the worker image explicitly:

```bash
docker build -t local/kali-pentest:latest build/kali-pentest
```

### UI starts but flows fail immediately

Check these first:

- the host-side LLM server is running on port `8080`
- embeddings are available from the same endpoint
- the configured model names actually exist on that server

### Search features fail

Check whether `searxng-core` is reachable from the `openui_default` network.

### Need to test VPN-based labs locally

See:

- `HOW_TO_CONNECT_TO_HACKTHEBOX.md`

## Branch-Specific Notes Worth Knowing

- This branch currently uses a locally built image in compose instead of a compose-managed build step
- The compose file contains hardcoded development-oriented values such as local URLs and a fixed cookie salt
- Before using this outside a trusted internal environment, replace development secrets and local placeholder values
- The current branch has custom logic recently added for:
  - correct user-only router registration
  - improved task progress display in the UI
  - fallback mission summary generation when the reporter output is missing

## Quick Start Summary

If your coworker already has a compatible model server on port `8080` and a reachable SearXNG instance, the full setup is:

```bash
git clone <your-repo-url>
cd pentagi
git checkout main
docker network create openui_default
docker build -t local/kali-pentest:latest build/kali-pentest
docker build -t local/pentagi:latest .
docker compose up -d
```

Then open:

```text
https://localhost:8443
```

and log in with:

```text
admin@pentagi.com
admin
```