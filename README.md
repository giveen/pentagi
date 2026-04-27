# PentAGI

<div align="center" style="font-size: 1.5em; margin: 20px 0;">
    <strong>P</strong>enetration testing <strong>A</strong>rtificial <strong>G</strong>eneral <strong>I</strong>ntelligence
</div>
<br>
<div align="center">

> **Join the Community!** Connect with security researchers, AI enthusiasts, and fellow ethical hackers. Get support, share insights, and stay updated with the latest PentAGI developments.

[![Discord](https://img.shields.io/badge/Discord-7289DA?logo=discord&logoColor=white)](https://discord.gg/2xrMh7qX6m)⠀[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?logo=telegram&logoColor=white)](https://t.me/+Ka9i6CNwe71hMWQy)

<a href="https://trendshift.io/repositories/15161" target="_blank"><img src="https://trendshift.io/api/badge/repositories/15161" alt="vxcontrol%2Fpentagi | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>

</div>

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
  - [Agent Supervision](#advanced-agent-supervision)
- [Quick Start](#quick-start)
- [API Access](#api-access)
  - [LLM Provider Configuration](#custom-llm-provider-configuration)
    - [Ollama](#ollama-provider-configuration)
    - [OpenAI](#openai-provider-configuration)
- [Advanced Setup](#advanced-setup)
  - [Langfuse Integration](#langfuse-integration)
  - [Monitoring and Observability](#monitoring-and-observability)
  - [Knowledge Graph (Graphiti)](#knowledge-graph-integration-graphiti)
  - [OAuth Integration](#github-and-google-oauth-integration)
  - [Docker Image Configuration](#docker-image-configuration)
- [Development](#development)
- [Testing LLM Agents](#testing-llm-agents)
- [Embedding Configuration and Testing](#embedding-configuration-and-testing)
- [Function Testing with ftester](#function-testing-with-ftester)
- [Building](#building)
- [Credits](#credits)
-- [License](#license)

- [Ports](#ports)

## Ports

- **Pentagi (web UI/API):** host 8443 → container 8443 (HTTPS)
- **Graphiti (knowledge graph API):** host 8444 → container 8000 (HTTP)
- **pgvector (Postgres):** host 8445 → container 5432 (Postgres)
- **Neo4j (HTTP):** host 8446 → container 7474 (HTTP)
- **Neo4j (Bolt):** host 8447 → container 7687 (Bolt)
- **Langfuse (web):** host 8448 → container 3000 (web) — if enabled
- **Grafana (web):** host 8449 → container web UI (observability)
- **OTEL / Collector endpoints:** host 8450/8451 → container OTEL ports

> Note: Pentagi uses Docker-internal networking by default. The service reaches Graphiti via `http://graphiti:8000` inside the Docker network; host-exposed ports are for external access only.

## Overview

PentAGI is an innovative tool for automated security testing that leverages cutting-edge artificial intelligence technologies. The project is designed for information security professionals, researchers, and enthusiasts who need a powerful and flexible solution for conducting penetration tests.

You can watch the video **PentAGI overview**:
[![PentAGI Overview Video](https://github.com/user-attachments/assets/0828dc3e-15f1-4a1d-858e-9696a146e478)](https://youtu.be/R70x5Ddzs1o)

## Features

- Secure & Isolated. All operations are performed in a sandboxed Docker environment with complete isolation.
- Fully Autonomous. AI-powered agent that automatically determines and executes penetration testing steps with optional execution monitoring and intelligent task planning for enhanced reliability.
- Professional Pentesting Tools. Built-in suite of 20+ professional security tools including nmap, metasploit, sqlmap, and more.
- Smart Memory System. Long-term storage of research results and successful approaches for future use.
- Knowledge Graph Integration. Graphiti-powered knowledge graph using Neo4j for semantic relationship tracking and advanced context understanding.
- Web Intelligence. Built-in browser via [scraper](https://hub.docker.com/r/vxcontrol/scraper) for gathering latest information from web sources.
- External Search Systems. Integration with advanced search APIs including [Tavily](https://tavily.com), [Traversaal](https://traversaal.ai), [Perplexity](https://www.perplexity.ai), [DuckDuckGo](https://duckduckgo.com/), [Google Custom Search](https://programmablesearchengine.google.com/), [Sploitus Search](https://sploitus.com) and [Searxng](https://searxng.org) for comprehensive information gathering.
- Team of Specialists. Delegation system with specialized AI agents for research, development, and infrastructure tasks, enhanced with optional execution monitoring and intelligent task planning for optimal performance with smaller models.
- Comprehensive Monitoring. Detailed logging and integration with Grafana/Prometheus for real-time system observation.
- Detailed Reporting. Generation of thorough vulnerability reports with exploitation guides.
- Smart Container Management. Automatic Docker image selection based on specific task requirements.
- Modern Interface. Clean and intuitive web UI for system management and monitoring.
- Comprehensive APIs. Full-featured REST and GraphQL APIs with Bearer token authentication for automation and integration.
- Persistent Storage. All commands and outputs are stored in PostgreSQL with [pgvector](https://hub.docker.com/r/vxcontrol/pgvector) extension.
- Scalable Architecture. Microservices-based design supporting horizontal scaling.
- Self-Hosted Solution. Complete control over your deployment and data.
- Flexible Authentication. Support for 3 LLM providers ([OpenAI](https://platform.openai.com/), [Ollama](https://ollama.com/), Custom/OpenAI-compatible endpoints). Easily integrate with OpenRouter, vLLM, LocalAI, and other OpenAI-compatible services. For production local deployments, see examples in [examples/configs/](examples/configs/).
- API Token Authentication. Secure Bearer token system for programmatic access to REST and GraphQL APIs.
- Quick Deployment. Easy setup through [Docker Compose](https://docs.docker.com/compose/) with comprehensive environment configuration.

## Architecture

### System Context

```mermaid
flowchart TB
    classDef person fill:#08427B,stroke:#073B6F,color:#fff
    classDef system fill:#1168BD,stroke:#0B4884,color:#fff
    classDef external fill:#666666,stroke:#0B4884,color:#fff

    pentester["👤 Security Engineer
    (User of the system)"]

    pentagi["✨ PentAGI
    (Autonomous penetration testing system)"]

    target["🎯 target-system
    (System under test)"]
    llm["🧠 llm-provider
    (OpenAI/Anthropic/Ollama/Bedrock/Gemini/Custom)"]
    search["🔍 search-systems
    (Google/DuckDuckGo/Tavily/Traversaal/Perplexity/Sploitus/Searxng)"]
    langfuse["📊 langfuse-ui
    (LLM Observability Dashboard)"]
    grafana["📈 grafana
    (System Monitoring Dashboard)"]

    pentester --> |Uses HTTPS| pentagi
    pentester --> |Monitors AI HTTPS| langfuse
    pentester --> |Monitors System HTTPS| grafana
    pentagi --> |Tests Various protocols| target
    pentagi --> |Queries HTTPS| llm
    pentagi --> |Searches HTTPS| search
    pentagi --> |Reports HTTPS| langfuse
    pentagi --> |Reports HTTPS| grafana

    class pentester person
    class pentagi system
    class target,llm,search,langfuse,grafana external

    linkStyle default stroke:#ffffff,color:#ffffff
```

<details>
<summary><b>Container Architecture</b> (click to expand)</summary>

```mermaid
graph TB
    subgraph Core Services
        UI[Frontend UI<br/>React + TypeScript]
        API[Backend API<br/>Go + GraphQL]
        DB[(Vector Store<br/>PostgreSQL + pgvector)]
        MQ[Task Queue<br/>Async Processing]
        Agent[AI Agents<br/>Multi-Agent System]
    end

    subgraph Knowledge Graph
        Graphiti[Graphiti<br/>Knowledge Graph API]
        Neo4j[(Neo4j<br/>Graph Database)]
    end

    subgraph Monitoring
        Grafana[Grafana<br/>Dashboards]
        VictoriaMetrics[VictoriaMetrics<br/>Time-series DB]
        Jaeger[Jaeger<br/>Distributed Tracing]
        Loki[Loki<br/>Log Aggregation]
        OTEL[OpenTelemetry<br/>Data Collection]
    end

    subgraph Analytics
        Langfuse[Langfuse<br/>LLM Analytics]
        ClickHouse[ClickHouse<br/>Analytics DB]
        Redis[Redis<br/>Cache + Rate Limiter]
        MinIO[MinIO<br/>S3 Storage]
    end

    subgraph Security Tools
        Scraper[Web Scraper<br/>Isolated Browser]
        PenTest[Security Tools<br/>20+ Pro Tools<br/>Sandboxed Execution]
    end

    UI --> |HTTP/WS| API
    API --> |SQL| DB
    API --> |Events| MQ
    MQ --> |Tasks| Agent
    Agent --> |Commands| PenTest
    Agent --> |Queries| DB
    Agent --> |Knowledge| Graphiti
    Graphiti --> |Graph| Neo4j

    API --> |Telemetry| OTEL
    OTEL --> |Metrics| VictoriaMetrics
    OTEL --> |Traces| Jaeger
    OTEL --> |Logs| Loki

    Grafana --> |Query| VictoriaMetrics
    Grafana --> |Query| Jaeger
    Grafana --> |Query| Loki

    API --> |Analytics| Langfuse
    Langfuse --> |Store| ClickHouse
    Langfuse --> |Cache| Redis
    Langfuse --> |Files| MinIO

    classDef core fill:#f9f,stroke:#333,stroke-width:2px,color:#000
    classDef knowledge fill:#ffa,stroke:#333,stroke-width:2px,color:#000
    classDef monitoring fill:#bbf,stroke:#333,stroke-width:2px,color:#000
    classDef analytics fill:#bfb,stroke:#333,stroke-width:2px,color:#000
    classDef tools fill:#fbb,stroke:#333,stroke-width:2px,color:#000

    class UI,API,DB,MQ,Agent core
    class Graphiti,Neo4j knowledge
    class Grafana,VictoriaMetrics,Jaeger,Loki,OTEL monitoring
    class Langfuse,ClickHouse,Redis,MinIO analytics
    class Scraper,PenTest tools
```

</details>

<details>
<summary><b>Entity Relationship</b> (click to expand)</summary>

```mermaid
erDiagram
    Flow ||--o{ Task : contains
    Task ||--o{ SubTask : contains
    SubTask ||--o{ Action : contains
    Action ||--o{ Artifact : produces
    Action ||--o{ Memory : stores

    Flow {
        string id PK
        string name "Flow name"
        string description "Flow description"
        string status "active/completed/failed"
        json parameters "Flow parameters"
        timestamp created_at
        timestamp updated_at
    }

    Task {
        string id PK
        string flow_id FK
        string name "Task name"
        string description "Task description"
        string status "pending/running/done/failed"
        json result "Task results"
        timestamp created_at
        timestamp updated_at
    }

    SubTask {
        string id PK
        string task_id FK
        string name "Subtask name"
        string description "Subtask description"
        string status "queued/running/completed/failed"
        string agent_type "researcher/developer/executor"
        json context "Agent context"
        timestamp created_at
        timestamp updated_at
    }

    Action {
        string id PK
        string subtask_id FK
        string type "command/search/analyze/etc"
        string status "success/failure"
        json parameters "Action parameters"
        json result "Action results"
        timestamp created_at
    }

    Artifact {
        string id PK
        string action_id FK
        string type "file/report/log"
        string path "Storage path"
        json metadata "Additional info"
        timestamp created_at
    }

    Memory {
        string id PK
        string action_id FK
        string type "observation/conclusion"
        vector embedding "Vector representation"
        text content "Memory content"
        timestamp created_at
    }
```

</details>

<details>
<summary><b>Agent Interaction</b> (click to expand)</summary>

```mermaid
sequenceDiagram
    participant O as Orchestrator
    participant R as Researcher
    participant D as Developer
    participant E as Executor
    participant VS as Vector Store
    participant KB as Knowledge Base

    Note over O,KB: Flow Initialization
    O->>VS: Query similar tasks
    VS-->>O: Return experiences
    O->>KB: Load relevant knowledge
    KB-->>O: Return context

    Note over O,R: Research Phase
    O->>R: Analyze target
    R->>VS: Search similar cases
    VS-->>R: Return patterns
    R->>KB: Query vulnerabilities
    KB-->>R: Return known issues
    R->>VS: Store findings
    R-->>O: Research results

    Note over O,D: Planning Phase
    O->>D: Plan attack
    D->>VS: Query exploits
    VS-->>D: Return techniques
    D->>KB: Load tools info
    KB-->>D: Return capabilities
    D-->>O: Attack plan

    Note over O,E: Execution Phase
    O->>E: Execute plan
    E->>KB: Load tool guides
    KB-->>E: Return procedures
    E->>VS: Store results
    E-->>O: Execution status
```

</details>

<details>
<summary><b>Memory System</b> (click to expand)</summary>

```mermaid
graph TB
    subgraph "Long-term Memory"
        VS[(Vector Store<br/>Embeddings DB)]
        KB[Knowledge Base<br/>Domain Expertise]
        Tools[Tools Knowledge<br/>Usage Patterns]
    end

    subgraph "Working Memory"
        Context[Current Context<br/>Task State]
        Goals[Active Goals<br/>Objectives]
        State[System State<br/>Resources]
    end

    subgraph "Episodic Memory"
        Actions[Past Actions<br/>Commands History]
        Results[Action Results<br/>Outcomes]
        Patterns[Success Patterns<br/>Best Practices]
    end

    Context --> |Query| VS
    VS --> |Retrieve| Context

    Goals --> |Consult| KB
    KB --> |Guide| Goals

    State --> |Record| Actions
    Actions --> |Learn| Patterns
    Patterns --> |Store| VS

    Tools --> |Inform| State
    Results --> |Update| Tools

    VS --> |Enhance| KB
    KB --> |Index| VS

    classDef ltm fill:#f9f,stroke:#333,stroke-width:2px,color:#000
    classDef wm fill:#bbf,stroke:#333,stroke-width:2px,color:#000
    classDef em fill:#bfb,stroke:#333,stroke-width:2px,color:#000

    class VS,KB,Tools ltm
    class Context,Goals,State wm
    class Actions,Results,Patterns em
```

</details>

<details>
<summary><b>Chain Summarization</b> (click to expand)</summary>

The chain summarization system manages conversation context growth by selectively summarizing older messages. This is critical for preventing token limits from being exceeded while maintaining conversation coherence.

```mermaid
flowchart TD
    A[Input Chain] --> B{Needs Summarization?}
    B -->|No| C[Return Original Chain]
    B -->|Yes| D[Convert to ChainAST]
    D --> E[Apply Section Summarization]
    E --> F[Process Oversized Pairs]
    F --> G[Manage Last Section Size]
    G --> H[Apply QA Summarization]
    H --> I[Rebuild Chain with Summaries]
    I --> J{Is New Chain Smaller?}
    J -->|Yes| K[Return Optimized Chain]
    J -->|No| C

    classDef process fill:#bbf,stroke:#333,stroke-width:2px,color:#000
    classDef decision fill:#bfb,stroke:#333,stroke-width:2px,color:#000
    classDef output fill:#fbb,stroke:#333,stroke-width:2px,color:#000

    class A,D,E,F,G,H,I process
    class B,J decision
    class C,K output
```

The algorithm operates on a structured representation of conversation chains (ChainAST) that preserves message types including tool calls and their responses. All summarization operations maintain critical conversation flow while reducing context size.

### Global Summarizer Configuration Options

| Parameter             | Environment Variable             | Default | Description                                                |
| --------------------- | -------------------------------- | ------- | ---------------------------------------------------------- |
| Preserve Last         | `SUMMARIZER_PRESERVE_LAST`       | `true`  | Whether to keep all messages in the last section intact    |
| Use QA Pairs          | `SUMMARIZER_USE_QA`              | `true`  | Whether to use QA pair summarization strategy              |
| Summarize Human in QA | `SUMMARIZER_SUM_MSG_HUMAN_IN_QA` | `false` | Whether to summarize human messages in QA pairs            |
| Last Section Size     | `SUMMARIZER_LAST_SEC_BYTES`      | `51200` | Maximum byte size for last section (50KB)                  |
| Max Body Pair Size    | `SUMMARIZER_MAX_BP_BYTES`        | `16384` | Maximum byte size for a single body pair (16KB)            |
| Max QA Sections       | `SUMMARIZER_MAX_QA_SECTIONS`     | `10`    | Maximum QA pair sections to preserve                       |
| Max QA Size           | `SUMMARIZER_MAX_QA_BYTES`        | `65536` | Maximum byte size for QA pair sections (64KB)              |
| Keep QA Sections      | `SUMMARIZER_KEEP_QA_SECTIONS`    | `1`     | Number of recent QA sections to keep without summarization |

### Assistant Summarizer Configuration Options

Assistant instances can use customized summarization settings to fine-tune context management behavior:

| Parameter          | Environment Variable                    | Default | Description                                                          |
| ------------------ | --------------------------------------- | ------- | -------------------------------------------------------------------- |
| Preserve Last      | `ASSISTANT_SUMMARIZER_PRESERVE_LAST`    | `true`  | Whether to preserve all messages in the assistant's last section     |
| Last Section Size  | `ASSISTANT_SUMMARIZER_LAST_SEC_BYTES`   | `76800` | Maximum byte size for assistant's last section (75KB)                |
| Max Body Pair Size | `ASSISTANT_SUMMARIZER_MAX_BP_BYTES`     | `16384` | Maximum byte size for a single body pair in assistant context (16KB) |
| Max QA Sections    | `ASSISTANT_SUMMARIZER_MAX_QA_SECTIONS`  | `7`     | Maximum QA sections to preserve in assistant context                 |
| Max QA Size        | `ASSISTANT_SUMMARIZER_MAX_QA_BYTES`     | `76800` | Maximum byte size for assistant's QA sections (75KB)                 |
| Keep QA Sections   | `ASSISTANT_SUMMARIZER_KEEP_QA_SECTIONS` | `3`     | Number of recent QA sections to preserve without summarization       |

The assistant summarizer configuration provides more memory for context retention compared to the global settings, preserving more recent conversation history while still ensuring efficient token usage.

### Summarizer Environment Configuration

```bash
# Default values for global summarizer logic
SUMMARIZER_PRESERVE_LAST=true
SUMMARIZER_USE_QA=true
SUMMARIZER_SUM_MSG_HUMAN_IN_QA=false
SUMMARIZER_LAST_SEC_BYTES=51200
SUMMARIZER_MAX_BP_BYTES=16384
SUMMARIZER_MAX_QA_SECTIONS=10
SUMMARIZER_MAX_QA_BYTES=65536
SUMMARIZER_KEEP_QA_SECTIONS=1

# Default values for assistant summarizer logic
ASSISTANT_SUMMARIZER_PRESERVE_LAST=true
ASSISTANT_SUMMARIZER_LAST_SEC_BYTES=76800
ASSISTANT_SUMMARIZER_MAX_BP_BYTES=16384
ASSISTANT_SUMMARIZER_MAX_QA_SECTIONS=7
ASSISTANT_SUMMARIZER_MAX_QA_BYTES=76800
ASSISTANT_SUMMARIZER_KEEP_QA_SECTIONS=3
```

</details>

<a id="advanced-agent-supervision"></a>
<details>
<summary><b>Advanced Agent Supervision</b> (click to expand)</summary>

PentAGI includes sophisticated multi-layered agent supervision mechanisms to ensure efficient task execution, prevent infinite loops, and provide intelligent recovery from stuck states:

### Execution Monitoring (Beta)
- **Automatic Mentor Intervention**: Adviser agent (mentor) is automatically invoked when execution patterns indicate potential issues
- **Pattern Detection**: Monitors identical tool calls (threshold: 5, configurable) and total tool calls (threshold: 10, configurable)
- **Progress Analysis**: Evaluates whether agent advances toward subtask objective, detects loops and inefficiencies
- **Alternative Strategies**: Recommends different approaches when current strategy fails
- **Information Retrieval Guidance**: Suggests searching for established solutions instead of reinventing
- **Enhanced Response Format**: Tool responses include both `<original_result>` and `<mentor_analysis>` sections
- **Configurable**: Enable via `EXECUTION_MONITOR_ENABLED` (default: false), customize thresholds with `EXECUTION_MONITOR_SAME_TOOL_LIMIT` and `EXECUTION_MONITOR_TOTAL_TOOL_LIMIT`

**Best for**: Smaller models (< 32B parameters), complex attack scenarios requiring continuous guidance, preventing agents from getting stuck on single approach

**Performance Impact**: 2-3x increase in execution time and token usage, but delivers **2x improvement in result quality** based on testing with Qwen3.5-27B-FP8

### Intelligent Task Planning (Beta)
- **Automated Decomposition**: Planner (adviser in planning mode) generates 3-7 specific, actionable steps before specialist agents begin work
- **Context-Aware Plans**: Analyzes full execution context via enricher agent to create informed plans
- **Structured Assignment**: Original request wrapped in `<task_assignment>` structure with execution plan and instructions
- **Scope Management**: Prevents scope creep by keeping agents focused on current subtask only
- **Enriched Instructions**: Plans highlight critical actions, potential pitfalls, and verification points
- **Configurable**: Enable via `AGENT_PLANNING_STEP_ENABLED` (default: false)

**Best for**: Models < 32B parameters, complex penetration testing workflows, improving success rates on sophisticated tasks

**Enhanced Adviser Configuration**: Works exceptionally well when adviser agent uses stronger model or enhanced settings. Example: using same base model with maximum reasoning mode for adviser (see [`vllm-qwen3.5-27b-fp8.provider.yml`](examples/configs/vllm-qwen3.5-27b-fp8.provider.yml)) enables comprehensive task analysis and strategic planning from identical model architecture.

**Performance Impact**: Adds planning overhead but significantly improves completion rates and reduces redundant work

### Tool Call Limits (Always Active)
- **Hard Limits**: Prevent runaway executions regardless of supervision mode status
- **Differentiated by Agent Type**:
  - General agents (Assistant, Primary Agent, Pentester, Coder, Installer): `MAX_GENERAL_AGENT_TOOL_CALLS` (default: 100)
  - Limited agents (Searcher, Enricher, Memorist, Generator, Reporter, Adviser, Reflector, Planner): `MAX_LIMITED_AGENT_TOOL_CALLS` (default: 20)
- **Graceful Termination**: Reflector guides agents to proper completion when approaching limits
- **Resource Protection**: Ensures system stability and prevents resource exhaustion

### Reflector Integration (Always Active)
- **Automatic Correction**: Invoked when LLM fails to generate tool calls after 3 attempts
- **Strategic Guidance**: Analyzes failures and guides agents toward proper tool usage or barrier tools (`done`, `ask`)
- **Recovery Mechanism**: Provides contextual guidance based on specific failure patterns
- **Limit Enforcement**: Coordinates graceful termination when tool call limits are reached

### Recommendations for Open Source Models

**Must-Have for Models < 32B Parameters**:
Testing with Qwen3.5-27B-FP8 demonstrates that enabling both Execution Monitoring and Task Planning is **essential** for smaller open source models:
- **Quality Improvement**: 2x better results compared to baseline execution without supervision
- **Loop Prevention**: Significantly reduces infinite loops and redundant work
- **Attack Diversity**: Encourages exploration of multiple attack vectors instead of fixating on single approach
- **Air-Gapped Deployments**: Enables production-grade autonomous pentesting in closed network environments with local LLM inference

**Trade-offs**:
- Token consumption: 2-3x increase due to mentor/planner invocations
- Execution time: 2-3x longer due to analysis and planning steps
- Result quality: 2x improvement in completeness, accuracy, and attack coverage
- Model requirements: Works best when adviser uses enhanced configuration (higher reasoning parameters, stronger model variant, or different model)

**Configuration Strategy**:
For optimal performance with smaller models, configure adviser agent with enhanced settings:
- Use same model with maximum reasoning mode (example: [`vllm-qwen3.5-27b-fp8.provider.yml`](examples/configs/vllm-qwen3.5-27b-fp8.provider.yml))
- Or use stronger model for adviser while keeping base model for other agents
- Adjust monitoring thresholds based on task complexity and model capabilities



</details>

The architecture of PentAGI is designed to be modular, scalable, and secure. Here are the key components:

1. **Core Services**
   - Frontend UI: React-based web interface with TypeScript for type safety
   - Backend API: Go-based REST and GraphQL APIs with Bearer token authentication for programmatic access
   - Vector Store: PostgreSQL with pgvector for semantic search and memory storage
   - Task Queue: Async task processing system for reliable operation
   - AI Agent: Multi-agent system with specialized roles for efficient testing

2. **Knowledge Graph**
   - Graphiti: Knowledge graph API for semantic relationship tracking and contextual understanding
   - Neo4j: Graph database for storing and querying relationships between entities, actions, and outcomes
   - Automatic capturing of agent responses and tool executions for building comprehensive knowledge base

3. **Monitoring Stack**
   - OpenTelemetry: Unified observability data collection and correlation
   - Grafana: Real-time visualization and alerting dashboards
   - VictoriaMetrics: High-performance time-series metrics storage
   - Jaeger: End-to-end distributed tracing for debugging
   - Loki: Scalable log aggregation and analysis

4. **Analytics Platform**
   - Langfuse: Advanced LLM observability and performance analytics
   - ClickHouse: Column-oriented analytics data warehouse
   - Redis: High-speed caching and rate limiting
   - MinIO: S3-compatible object storage for artifacts

5. **Security Tools**
   - Web Scraper: Isolated browser environment for safe web interaction
   - Pentesting Tools: Comprehensive suite of 20+ professional security tools
   - Sandboxed Execution: All operations run in isolated containers

6. **Memory Systems**
   - Long-term Memory: Persistent storage of knowledge and experiences
   - Working Memory: Active context and goals for current operations
   - Episodic Memory: Historical actions and success patterns
   - Knowledge Base: Structured domain expertise and tool capabilities
   - Context Management: Intelligently manages growing LLM context windows using chain summarization

The system uses Docker containers for isolation and easy deployment, with separate networks for core services, monitoring, and analytics to ensure proper security boundaries. Each component is designed to scale horizontally and can be configured for high availability in production environments.

## Quick Start

### System Requirements

- Docker and Docker Compose (or Podman - see [Podman configuration](#running-pentagi-with-podman))
- Minimum 2 vCPU
- Minimum 4GB RAM
- 20GB free disk space
- Internet access for downloading images and updates

### Using Installer (Recommended)

PentAGI provides an interactive installer with a terminal-based UI for streamlined configuration and deployment. The installer guides you through system checks, LLM provider setup, search engine configuration, and security hardening.

**Supported Platforms:**
- **Linux**: amd64 [download](https://pentagi.com/downloads/linux/amd64/installer-latest.zip) | arm64 [download](https://pentagi.com/downloads/linux/arm64/installer-latest.zip)
- **Windows**: amd64 [download](https://pentagi.com/downloads/windows/amd64/installer-latest.zip)
- **macOS**: amd64 (Intel) [download](https://pentagi.com/downloads/darwin/amd64/installer-latest.zip) | arm64 (M-series) [download](https://pentagi.com/downloads/darwin/arm64/installer-latest.zip)

**Quick Installation (Linux amd64):**

```bash
# Create installation directory
mkdir -p pentagi && cd pentagi

# Download installer
wget -O installer.zip https://pentagi.com/downloads/linux/amd64/installer-latest.zip

# Extract
unzip installer.zip

# Run interactive installer
./installer
```

**Prerequisites & Permissions:**

The installer requires appropriate privileges to interact with the Docker API for proper operation. By default, it uses the Docker socket (`/var/run/docker.sock`) which requires either:

- **Option 1 (Recommended for production):** Run the installer as root:
  ```bash
  sudo ./installer
  ```

- **Option 2 (Development environments):** Grant your user access to the Docker socket by adding them to the `docker` group:
  ```bash
  # Add your user to the docker group
  sudo usermod -aG docker $USER
  
  # Log out and log back in, or activate the group immediately
  newgrp docker
  
  # Verify Docker access (should run without sudo)
  docker ps
  ```

  ⚠️ **Security Note:** Adding a user to the `docker` group grants root-equivalent privileges. Only do this for trusted users in controlled environments. For production deployments, consider using rootless Docker mode or running the installer with sudo.

The installer will:
1. **System Checks**: Verify Docker, network connectivity, and system requirements
2. **Environment Setup**: Create and configure `.env` file with optimal defaults
3. **Provider Configuration**: Set up LLM providers (OpenAI, Anthropic, Gemini, Bedrock, Ollama, Custom)
4. **Search Engines**: Configure DuckDuckGo, Google, Tavily, Traversaal, Perplexity, Sploitus, Searxng
5. **Security Hardening**: Generate secure credentials and configure SSL certificates
6. **Deployment**: Start PentAGI with docker-compose

**For Production & Enhanced Security:**

For production deployments or security-sensitive environments, we **strongly recommend** using a distributed two-node architecture where worker operations are isolated on a separate server. This prevents untrusted code execution and network access issues on your main system.

**See detailed guide**: [Worker Node Setup](examples/guides/worker_node.md)

The two-node setup provides:
- **Isolated Execution**: Worker containers run on dedicated hardware
- **Network Isolation**: Separate network boundaries for penetration testing
- **Security Boundaries**: Docker-in-Docker with TLS authentication
- **OOB Attack Support**: Dedicated port ranges for out-of-band techniques

### Manual Installation

1. Create a working directory or clone the repository:

```bash
mkdir pentagi && cd pentagi
```

2. Copy `.env.example` to `.env` or download it:

```bash
curl -o .env https://raw.githubusercontent.com/vxcontrol/pentagi/master/.env.example
```

3. Touch examples files (`example.custom.provider.yml`, `example.ollama.provider.yml`) or download it:

```bash
curl -o example.custom.provider.yml https://raw.githubusercontent.com/vxcontrol/pentagi/master/examples/configs/custom-openai.provider.yml
curl -o example.ollama.provider.yml https://raw.githubusercontent.com/vxcontrol/pentagi/master/examples/configs/ollama-llama318b.provider.yml
```

4. Fill in the required API keys in `.env` file.

```bash
# Required: At least one of these LLM providers
OPEN_AI_KEY=your_openai_key
ANTHROPIC_API_KEY=your_anthropic_key
GEMINI_API_KEY=your_gemini_key

# Optional: AWS Bedrock provider (enterprise-grade models)
BEDROCK_REGION=us-east-1
# Choose one authentication method:
BEDROCK_DEFAULT_AUTH=true                        # Option 1: Use AWS SDK default credential chain (recommended for EC2/ECS)
# BEDROCK_BEARER_TOKEN=your_bearer_token         # Option 2: Bearer token authentication
# BEDROCK_ACCESS_KEY_ID=your_aws_access_key      # Option 3: Static credentials
# BEDROCK_SECRET_ACCESS_KEY=your_aws_secret_key

# Optional: Ollama provider (local or cloud)
# OLLAMA_SERVER_URL=http://ollama-server:11434   # Local server
# OLLAMA_SERVER_URL=https://ollama.com           # Cloud service
# OLLAMA_SERVER_API_KEY=your_ollama_cloud_key    # Required for cloud, empty for local

# Optional: Chinese AI providers
# DEEPSEEK_API_KEY=your_deepseek_key             # DeepSeek (strong reasoning)
# GLM_API_KEY=your_glm_key                       # GLM (Zhipu AI)
# KIMI_API_KEY=your_kimi_key                     # Kimi (Moonshot AI, ultra-long context)
# QWEN_API_KEY=your_qwen_key                     # Qwen (Alibaba Cloud, multimodal)

# Optional: Local LLM provider (zero-cost inference)
OLLAMA_SERVER_URL=http://localhost:11434
OLLAMA_SERVER_MODEL=your_model_name

# Optional: Additional search capabilities
DUCKDUCKGO_ENABLED=true
DUCKDUCKGO_REGION=us-en
DUCKDUCKGO_SAFESEARCH=
DUCKDUCKGO_TIME_RANGE=
SPLOITUS_ENABLED=true
GOOGLE_API_KEY=your_google_key
GOOGLE_CX_KEY=your_google_cx
TAVILY_API_KEY=your_tavily_key
TRAVERSAAL_API_KEY=your_traversaal_key
PERPLEXITY_API_KEY=your_perplexity_key
PERPLEXITY_MODEL=sonar-pro
PERPLEXITY_CONTEXT_SIZE=medium

# Searxng meta search engine (aggregates results from multiple sources)
SEARXNG_URL=http://your-searxng-instance:8080
SEARXNG_CATEGORIES=general
SEARXNG_LANGUAGE=
SEARXNG_SAFESEARCH=0
SEARXNG_TIME_RANGE=
SEARXNG_TIMEOUT=

## Graphiti knowledge graph settings
GRAPHITI_ENABLED=true
GRAPHITI_TIMEOUT=30
GRAPHITI_URL=http://graphiti:8000
GRAPHITI_MODEL_NAME=gpt-5-mini

# Neo4j settings (used by Graphiti stack)
NEO4J_USER=neo4j
NEO4J_DATABASE=neo4j
NEO4J_PASSWORD=devpassword
NEO4J_URI=bolt://neo4j:7687

# Assistant configuration
ASSISTANT_USE_AGENTS=false         # Default value for agent usage when creating new assistants
```

5. Change all security related environment variables in `.env` file to improve security.

<details>
    <summary>Security related environment variables</summary>

### Main Security Settings
- `COOKIE_SIGNING_SALT` - Salt for cookie signing, change to random value
- `PUBLIC_URL` - Public URL of your server (eg. `https://pentagi.example.com`)
- `SERVER_SSL_CRT` and `SERVER_SSL_KEY` - Custom paths to your existing SSL certificate and key for HTTPS (these paths should be used in the docker-compose.yml file to mount as volumes)

### Scraper Access
- `SCRAPER_PUBLIC_URL` - Public URL for scraper if you want to use different scraper server for public URLs
- `SCRAPER_PRIVATE_URL` - Private URL for scraper (local scraper server in docker-compose.yml file to access it to local URLs)

### Access Credentials
- `PENTAGI_POSTGRES_USER` and `PENTAGI_POSTGRES_PASSWORD` - PostgreSQL credentials
- `NEO4J_USER` and `NEO4J_PASSWORD` - Neo4j credentials (for Graphiti knowledge graph)

</details>

6. Remove all inline comments from `.env` file if you want to use it in VSCode or other IDEs as a envFile option:

```bash
perl -i -pe 's/\s+#.*$//' .env
```

7. Run the PentAGI stack:

```bash
curl -O https://raw.githubusercontent.com/vxcontrol/pentagi/master/docker-compose.yml
docker compose up -d
```

Visit [localhost:8443](https://localhost:8443) to access PentAGI Web UI (default is `admin@pentagi.com` / `admin`)

> [!NOTE]
> If you caught an error about `pentagi-network` or `observability-network` or `langfuse-network` you need to run `docker-compose.yml` firstly to create these networks and after that run `docker-compose-langfuse.yml`, `docker-compose-graphiti.yml`, and `docker-compose-observability.yml` to use Langfuse, Graphiti, and Observability services.
>
> You have to set at least one Language Model provider (OpenAI, Anthropic, Gemini, AWS Bedrock, or Ollama) to use PentAGI. AWS Bedrock provides enterprise-grade access to multiple foundation models from leading AI companies, while Ollama provides zero-cost local inference if you have sufficient computational resources. Additional API keys for search engines are optional but recommended for better results.
>
> **For fully local deployment with advanced models**: See our comprehensive guide on [Running PentAGI with vLLM and Qwen3.5-27B-FP8](examples/guides/vllm-qwen35-27b-fp8.md) for a production-grade local LLM setup. This configuration achieves ~13,000 TPS for prompt processing and ~650 TPS for completion on 4× RTX 5090 GPUs, supporting 12+ concurrent flows with complete independence from cloud providers.
>
> `LLM_SERVER_*` environment variables are experimental feature and will be changed in the future. Right now you can use them to specify custom LLM server URL and one model for all agent types.
>
> `PROXY_URL` is a global proxy URL for all LLM providers and external search systems. You can use it for isolation from external networks.
>
> The `docker-compose.yml` file runs the PentAGI service as root user because it needs access to docker.sock for container management. If you're using TCP/IP network connection to Docker instead of socket file, you can remove root privileges and use the default `pentagi` user for better security.

### Accessing PentAGI from External Networks

By default, PentAGI binds to `127.0.0.1` (localhost only) for security. To access PentAGI from other machines on your network, you need to configure external access.

#### Configuration Steps

1. **Update `.env` file** with your server's IP address:

```bash
# Network binding - allow external connections
PENTAGI_LISTEN_IP=0.0.0.0
PENTAGI_LISTEN_PORT=8443

# Public URL - use your actual server IP or hostname
# Replace 192.168.1.100 with your server's IP address
PUBLIC_URL=https://192.168.1.100:8443

# CORS origins - list all URLs that will access PentAGI
# Include localhost for local access AND your server IP for external access
CORS_ORIGINS=https://localhost:8443,https://192.168.1.100:8443
```

> [!IMPORTANT]
> - Replace `192.168.1.100` with your actual server's IP address
> - Do NOT use `0.0.0.0` in `PUBLIC_URL` or `CORS_ORIGINS` - use the actual IP address
> - Include both localhost and your server IP in `CORS_ORIGINS` for flexibility

2. **Recreate containers** to apply the changes:

```bash
docker compose down
docker compose up -d --force-recreate
```

3. **Verify port binding:**

```bash
docker ps | grep pentagi
```

You should see `0.0.0.0:8443->8443/tcp` or `:::8443->8443/tcp`.

If you see `127.0.0.1:8443->8443/tcp`, the environment variable wasn't picked up. In this case, directly edit `docker-compose.yml` line 31:

```yaml
ports:
  - "0.0.0.0:8443:8443"
```

Then recreate containers again.

4. **Configure firewall** to allow incoming connections on port 8443:

```bash
# Ubuntu/Debian with UFW
sudo ufw allow 8443/tcp
sudo ufw reload

# CentOS/RHEL with firewalld
sudo firewall-cmd --permanent --add-port=8443/tcp
sudo firewall-cmd --reload
```

5. **Access PentAGI:**

- **Local access:** `https://localhost:8443`
- **Network access:** `https://your-server-ip:8443`

> [!NOTE]
> You'll need to accept the self-signed SSL certificate warning in your browser when accessing via IP address.

---

### Running PentAGI with Podman

PentAGI fully supports Podman as a Docker alternative. However, when using **Podman in rootless mode**, the scraper service requires special configuration because rootless containers cannot bind privileged ports (ports below 1024).

#### Podman Rootless Configuration

The default scraper configuration uses port 443 (HTTPS), which is a privileged port. For Podman rootless, reconfigure the scraper to use a non-privileged port:

**1. Edit `docker-compose.yml`** - modify the `scraper` service (around line 199):

```yaml
scraper:
  image: vxcontrol/scraper:latest
  restart: unless-stopped
  container_name: scraper
  hostname: scraper
  expose:
    - 3000/tcp  # Changed from 443 to 3000
  ports:
    - "${SCRAPER_LISTEN_IP:-127.0.0.1}:${SCRAPER_LISTEN_PORT:-9443}:3000"  # Map to port 3000
  environment:
    - MAX_CONCURRENT_SESSIONS=${LOCAL_SCRAPER_MAX_CONCURRENT_SESSIONS:-10}
    - USERNAME=${LOCAL_SCRAPER_USERNAME:-someuser}
    - PASSWORD=${LOCAL_SCRAPER_PASSWORD:-somepass}
  logging:
    options:
      max-size: 50m
      max-file: "7"
  volumes:
    - scraper-ssl:/usr/src/app/ssl
  networks:
    - pentagi-network
  shm_size: 2g
```

**2. Update `.env` file** - change the scraper URL to use HTTP and port 3000:

```bash
# Scraper configuration for Podman rootless
SCRAPER_PRIVATE_URL=http://someuser:somepass@scraper:3000/
LOCAL_SCRAPER_USERNAME=someuser
LOCAL_SCRAPER_PASSWORD=somepass
```

> [!IMPORTANT]
> Key changes for Podman:
> - Use **HTTP** instead of HTTPS for `SCRAPER_PRIVATE_URL`
> - Use port **3000** instead of 443
> - Change internal `expose` to `3000/tcp`
> - Update port mapping to target `3000` instead of `443`

**3. Recreate containers:**

```bash
podman-compose down
podman-compose up -d --force-recreate
```

**4. Test scraper connectivity:**

```bash
# Test from within the pentagi container
podman exec -it pentagi wget -O- "http://someuser:somepass@scraper:3000/html?url=http://example.com"
```

If you see HTML output, the scraper is working correctly.

#### Podman Rootful Mode

If you're running Podman in rootful mode (with sudo), you can use the default configuration without modifications. The scraper will work on port 443 as intended.

#### Docker Compatibility

All Podman configurations remain fully compatible with Docker. The non-privileged port approach works identically on both container runtimes.

### Assistant Configuration

PentAGI allows you to configure default behavior for assistants:

| Variable               | Default | Description                                                             |
| ---------------------- | ------- | ----------------------------------------------------------------------- |
| `ASSISTANT_USE_AGENTS` | `false` | Controls the default value for agent usage when creating new assistants |

The `ASSISTANT_USE_AGENTS` setting affects the initial state of the "Use Agents" toggle when creating a new assistant in the UI:
- `false` (default): New assistants are created with agent delegation disabled by default
- `true`: New assistants are created with agent delegation enabled by default

Note that users can always override this setting by toggling the "Use Agents" button in the UI when creating or editing an assistant. This environment variable only controls the initial default state.

## API Access

PentAGI provides comprehensive programmatic access through both REST and GraphQL APIs, allowing you to integrate penetration testing workflows into your automation pipelines, CI/CD processes, and custom applications.

### Generating API Tokens

API tokens are managed through the PentAGI web interface:

1. Navigate to **Settings** → **API Tokens** in the web UI
2. Click **Create Token** to generate a new API token
3. Configure token properties:
   - **Name** (optional): A descriptive name for the token
   - **Expiration Date**: When the token will expire (minimum 1 minute, maximum 3 years)
4. Click **Create** and **copy the token immediately** - it will only be shown once for security reasons
5. Use the token as a Bearer token in your API requests

Each token is associated with your user account and inherits your role's permissions.

### Using API Tokens

Include the API token in the `Authorization` header of your HTTP requests:

```bash
# GraphQL API example
curl -X POST https://your-pentagi-instance:8443/api/v1/graphql \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"query": "{ flows { id title status } }"}'

# REST API example
curl https://your-pentagi-instance:8443/api/v1/flows \
  -H "Authorization: Bearer YOUR_API_TOKEN"
```

### API Exploration and Testing

PentAGI provides interactive documentation for exploring and testing API endpoints:

#### GraphQL Playground

Access the GraphQL Playground at `https://your-pentagi-instance:8443/api/v1/graphql/playground`

1. Click the **HTTP Headers** tab at the bottom
2. Add your authorization header:
   ```json
   {
     "Authorization": "Bearer YOUR_API_TOKEN"
   }
   ```
3. Explore the schema, run queries, and test mutations interactively

#### Swagger UI

Access the REST API documentation at `https://your-pentagi-instance:8443/api/v1/swagger/index.html`

1. Click the **Authorize** button
2. Enter your token in the format: `Bearer YOUR_API_TOKEN`
3. Click **Authorize** to apply
4. Test endpoints directly from the Swagger UI

### Generating API Clients

You can generate type-safe API clients for your preferred programming language using the schema files included with PentAGI:

#### GraphQL Clients

The GraphQL schema is available at:
- **Web UI**: Navigate to Settings to download `schema.graphqls`
- **Direct file**: `backend/pkg/graph/schema.graphqls` in the repository

Generate clients using tools like:
- **GraphQL Code Generator** (JavaScript/TypeScript): [https://the-guild.dev/graphql/codegen](https://the-guild.dev/graphql/codegen)
- **genqlient** (Go): [https://github.com/Khan/genqlient](https://github.com/Khan/genqlient)
- **Apollo iOS** (Swift): [https://www.apollographql.com/docs/ios](https://www.apollographql.com/docs/ios)

#### REST API Clients

The OpenAPI specification is available at:
- **Swagger JSON**: `https://your-pentagi-instance:8443/api/v1/swagger/doc.json`
- **Swagger YAML**: Available in `backend/pkg/server/docs/swagger.yaml`

Generate clients using:
- **OpenAPI Generator**: [https://openapi-generator.tech](https://openapi-generator.tech)
  ```bash
  openapi-generator-cli generate \
    -i https://your-pentagi-instance:8443/api/v1/swagger/doc.json \
    -g python \
    -o ./pentagi-client
  ```

- **Swagger Codegen**: [https://github.com/swagger-api/swagger-codegen](https://github.com/swagger-api/swagger-codegen)
  ```bash
  swagger-codegen generate \
    -i https://your-pentagi-instance:8443/api/v1/swagger/doc.json \
    -l typescript-axios \
    -o ./pentagi-client
  ```

- **swagger-typescript-api** (TypeScript): [https://github.com/acacode/swagger-typescript-api](https://github.com/acacode/swagger-typescript-api)
  ```bash
  npx swagger-typescript-api \
    -p https://your-pentagi-instance:8443/api/v1/swagger/doc.json \
    -o ./src/api \
    -n pentagi-api.ts
  ```

### API Usage Examples

<details>
<summary><b>Creating a New Flow (GraphQL)</b></summary>

```graphql
mutation CreateFlow {
  createFlow(
    modelProvider: "openai"
    input: "Test the security of https://example.com"
  ) {
    id
    title
    status
    createdAt
  }
}
```

</details>

<details>
<summary><b>Listing Flows (REST API)</b></summary>

```bash
curl https://your-pentagi-instance:8443/api/v1/flows \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  | jq '.flows[] | {id, title, status}'
```

</details>

<details>
<summary><b>Python Client Example</b></summary>

```python
import requests

class PentAGIClient:
    def __init__(self, base_url, api_token):
        self.base_url = base_url
        self.headers = {
            "Authorization": f"Bearer {api_token}",
            "Content-Type": "application/json"
        }
    
    def create_flow(self, provider, target):
        query = """
        mutation CreateFlow($provider: String!, $input: String!) {
          createFlow(modelProvider: $provider, input: $input) {
            id
            title
            status
          }
        }
        """
        response = requests.post(
            f"{self.base_url}/api/v1/graphql",
            json={
                "query": query,
                "variables": {
                    "provider": provider,
                    "input": target
                }
            },
            headers=self.headers
        )
        return response.json()
    
    def get_flows(self):
        response = requests.get(
            f"{self.base_url}/api/v1/flows",
            headers=self.headers
        )
        return response.json()

# Usage
client = PentAGIClient(
    "https://your-pentagi-instance:8443",
    "your_api_token_here"
)

# Create a new flow
flow = client.create_flow("openai", "Scan https://example.com for vulnerabilities")
print(f"Created flow: {flow}")

# List all flows
flows = client.get_flows()
print(f"Total flows: {len(flows['flows'])}")
```

</details>

<details>
<summary><b>TypeScript Client Example</b></summary>

```typescript
import axios, { AxiosInstance } from 'axios';

interface Flow {
  id: string;
  title: string;
  status: string;
  createdAt: string;
}

class PentAGIClient {
  private client: AxiosInstance;

  constructor(baseURL: string, apiToken: string) {
    this.client = axios.create({
      baseURL: `${baseURL}/api/v1`,
      headers: {
        'Authorization': `Bearer ${apiToken}`,
        'Content-Type': 'application/json',
      },
    });
  }

  async createFlow(provider: string, input: string): Promise<Flow> {
    const query = `
      mutation CreateFlow($provider: String!, $input: String!) {
        createFlow(modelProvider: $provider, input: $input) {
          id
          title
          status
          createdAt
        }
      }
    `;

    const response = await this.client.post('/graphql', {
      query,
      variables: { provider, input },
    });

    return response.data.data.createFlow;
  }

  async getFlows(): Promise<Flow[]> {
    const response = await this.client.get('/flows');
    return response.data.flows;
  }

  async getFlow(flowId: string): Promise<Flow> {
    const response = await this.client.get(`/flows/${flowId}`);
    return response.data;
  }
}

// Usage
const client = new PentAGIClient(
  'https://your-pentagi-instance:8443',
  'your_api_token_here'
);

// Create a new flow
const flow = await client.createFlow(
  'openai',
  'Perform penetration test on https://example.com'
);
console.log('Created flow:', flow);

// List all flows
const flows = await client.getFlows();
console.log(`Total flows: ${flows.length}`);
```

</details>

### Security Best Practices

When working with API tokens:

- **Never commit tokens to version control** - use environment variables or secrets management
- **Rotate tokens regularly** - set appropriate expiration dates and create new tokens periodically
- **Use separate tokens for different applications** - makes it easier to revoke access if needed
- **Monitor token usage** - review API token activity in the Settings page
- **Revoke unused tokens** - disable or delete tokens that are no longer needed
- **Use HTTPS only** - never send API tokens over unencrypted connections

### Token Management

- **View tokens**: See all your active tokens in Settings → API Tokens
- **Edit tokens**: Update token names or revoke tokens
- **Delete tokens**: Permanently remove tokens (this action cannot be undone)
- **Token ID**: Each token has a unique ID that can be copied for reference

The token list shows:
- Token name (if provided)
- Token ID (unique identifier)
- Status (active/revoked/expired)
- Creation date
- Expiration date

### Custom LLM Provider Configuration

When using custom LLM providers with the `LLM_SERVER_*` variables, you can fine-tune the reasoning format used in requests.

> [!TIP]
> For production-grade local deployments, consider using **vLLM** with **Qwen3.5-27B-FP8** for optimal performance. See our [comprehensive deployment guide](examples/guides/vllm-qwen35-27b-fp8.md) which includes hardware requirements, configuration templates ([thinking mode](examples/configs/vllm-qwen3.5-27b-fp8.provider.yml) and [non-thinking mode](examples/configs/vllm-qwen3.5-27b-fp8-no-think.provider.yml)), and performance benchmarks showing 13K TPS prompt processing on 4× RTX 5090 GPUs.

| Variable                        | Default | Description                                                                             |
| ------------------------------- | ------- | --------------------------------------------------------------------------------------- |
| `LLM_SERVER_URL`                |         | Base URL for the custom LLM API endpoint                                                |
| `LLM_SERVER_KEY`                |         | API key for the custom LLM provider                                                     |
| `LLM_SERVER_MODEL`              |         | Default model to use (can be overridden in provider config)                             |
| `LLM_SERVER_CONFIG_PATH`        |         | Path to the YAML configuration file for agent-specific models                           |
| `LLM_SERVER_PROVIDER`           |         | Provider name prefix for model names (e.g., `openrouter`, `deepseek` for LiteLLM proxy) |
| `LLM_SERVER_LEGACY_REASONING`   | `false` | Controls reasoning format in API requests                                               |
| `LLM_SERVER_PRESERVE_REASONING` | `false` | Preserve reasoning content in multi-turn conversations (required by some providers)     |

The `LLM_SERVER_PROVIDER` setting is particularly useful when using **LiteLLM proxy**, which adds a provider prefix to model names. For example, when connecting to Moonshot API through LiteLLM, models like `kimi-2.5` become `moonshot/kimi-2.5`. By setting `LLM_SERVER_PROVIDER=moonshot`, you can use the same provider configuration file for both direct API access and LiteLLM proxy access without modifications.

The `LLM_SERVER_LEGACY_REASONING` setting affects how reasoning parameters are sent to the LLM:
- `false` (default): Uses modern format where reasoning is sent as a structured object with `max_tokens` parameter
- `true`: Uses legacy format with string-based `reasoning_effort` parameter

This setting is important when working with different LLM providers as they may expect different reasoning formats in their API requests. If you encounter reasoning-related errors with custom providers, try changing this setting.

The `LLM_SERVER_PRESERVE_REASONING` setting controls whether reasoning content is preserved in multi-turn conversations:
- `false` (default): Reasoning content is not preserved in conversation history
- `true`: Reasoning content is preserved and sent in subsequent API calls

This setting is required by some LLM providers (e.g., Moonshot) that return errors like "thinking is enabled but reasoning_content is missing in assistant tool call message" when reasoning content is not included in multi-turn conversations. Enable this setting if your provider requires reasoning content to be preserved.

### Ollama Provider Configuration

PentAGI supports Ollama for both local LLM inference (zero-cost, enhanced privacy) and Ollama Cloud (managed service with free tier).

#### Configuration Variables

| Variable                            | Default     | Description                               |
| ----------------------------------- | ----------- | ----------------------------------------- |
| `OLLAMA_SERVER_URL`                 |             | URL of your Ollama server or Ollama Cloud |
| `OLLAMA_SERVER_API_KEY`             |             | API key for Ollama Cloud authentication   |
| `OLLAMA_SERVER_MODEL`               |             | Default model for inference               |
| `OLLAMA_SERVER_CONFIG_PATH`         |             | Path to custom agent configuration file   |
| `OLLAMA_SERVER_PULL_MODELS_TIMEOUT` | `600`       | Timeout for model downloads (seconds)     |
| `OLLAMA_SERVER_PULL_MODELS_ENABLED` | `false`     | Auto-download models on startup           |
| `OLLAMA_SERVER_LOAD_MODELS_ENABLED` | `false`     | Query server for available models         |

#### Ollama Cloud Configuration

Ollama Cloud provides managed inference with a generous free tier and scalable paid plans.

**Free Tier Setup (Single Model)**

```bash
# Free tier allows one model at a time
OLLAMA_SERVER_URL=https://ollama.com
OLLAMA_SERVER_API_KEY=your_ollama_cloud_api_key
OLLAMA_SERVER_MODEL=gpt-oss:120b  # Example: OpenAI OSS 120B model
```

**Paid Tier Setup (Multi-Model with Pre-built Configuration)**

For paid tiers supporting multiple concurrent models, use the pre-built Ollama Cloud configuration:

```bash
# Using pre-built Ollama Cloud configuration (included in Docker image)
OLLAMA_SERVER_URL=https://ollama.com
OLLAMA_SERVER_API_KEY=your_ollama_cloud_api_key
OLLAMA_SERVER_CONFIG_PATH=/opt/pentagi/conf/ollama-cloud.provider.yml
```

The pre-built `ollama-cloud.provider.yml` configuration includes optimized model assignments for all agent types:
- **Simple/Assistant**: `nemotron-3-super:cloud` - Fast general-purpose model
- **Primary Agent**: `qwen3-coder-next:cloud` - Advanced reasoning with high effort mode
- **Coder/Pentester**: `qwen3-coder-next:cloud` - Specialized coding models
- **Searcher**: `qwen3.5:397b-cloud` - Large context for information gathering
- **Refiner/Refactor**: `glm-5:cloud` - High-quality text refinement
- **Adviser/Enricher**: `minimax-m2.7:cloud` - Efficient advisory tasks
- **Installer**: `devstral-2:123b-cloud` - Installation and setup tasks

**Custom Configuration (Advanced)**

To create your own agent configuration, mount a custom file from your host filesystem:

```bash
# Using custom provider configuration
OLLAMA_SERVER_URL=https://ollama.com
OLLAMA_SERVER_API_KEY=your_ollama_cloud_api_key
OLLAMA_SERVER_CONFIG_PATH=/opt/pentagi/conf/ollama.provider.yml

# Mount custom configuration from host filesystem (in .env or docker-compose override)
PENTAGI_OLLAMA_SERVER_CONFIG_PATH=/path/on/host/my-ollama-config.yml
```

The `PENTAGI_OLLAMA_SERVER_CONFIG_PATH` environment variable maps your host configuration file to `/opt/pentagi/conf/ollama.provider.yml` inside the container.

**Example custom configuration** (`my-ollama-config.yml`):

```yaml
primary_agent:
  model: "qwen3-coder-next:cloud"
  temperature: 1.0
  top_p: 0.9
  max_tokens: 32768
  reasoning:
    effort: high

coder:
  model: "qwen3-coder:32b"
  temperature: 1.0
  max_tokens: 20480
```

#### Local Ollama Configuration

For self-hosted Ollama instances:

```bash
# Basic local Ollama setup
OLLAMA_SERVER_URL=http://localhost:11434
OLLAMA_SERVER_MODEL=llama3.1:8b-instruct-q8_0

# Production setup with auto-pull and model discovery
OLLAMA_SERVER_URL=http://ollama-server:11434
OLLAMA_SERVER_PULL_MODELS_ENABLED=true
OLLAMA_SERVER_PULL_MODELS_TIMEOUT=900
OLLAMA_SERVER_LOAD_MODELS_ENABLED=true

# Using pre-built configurations from Docker image
OLLAMA_SERVER_CONFIG_PATH=/opt/pentagi/conf/ollama-llama318b.provider.yml
# or
OLLAMA_SERVER_CONFIG_PATH=/opt/pentagi/conf/ollama-qwen332b-fp16-tc.provider.yml
# or
OLLAMA_SERVER_CONFIG_PATH=/opt/pentagi/conf/ollama-qwq32b-fp16-tc.provider.yml
```

**Performance Considerations:**

- **Model Discovery** (`OLLAMA_SERVER_LOAD_MODELS_ENABLED=true`): Adds 1-2s startup latency querying Ollama API
- **Auto-pull** (`OLLAMA_SERVER_PULL_MODELS_ENABLED=true`): First startup may take several minutes downloading models
- **Pull timeout** (`OLLAMA_SERVER_PULL_MODELS_TIMEOUT=900`): 15 minutes in seconds
- **Static Config**: Disable both flags and specify models in config file for fastest startup

#### Creating Custom Ollama Models with Extended Context

PentAGI requires models with larger context windows than the default Ollama configurations. You need to create custom models with increased `num_ctx` parameter through Modelfiles. While typical agent workflows consume around 64K tokens, PentAGI uses 110K context size for safety margin and handling complex penetration testing scenarios.

**Important**: The `num_ctx` parameter can only be set during model creation via Modelfile - it cannot be changed after model creation or overridden at runtime.

##### Example: Qwen3 32B FP16 with Extended Context

Create a Modelfile named `Modelfile_qwen3_32b_fp16_tc`:

```dockerfile
FROM qwen3:32b-fp16
PARAMETER num_ctx 110000
PARAMETER temperature 0.3
PARAMETER top_p 0.8
PARAMETER min_p 0.0
PARAMETER top_k 20
PARAMETER repeat_penalty 1.1
```

Build the custom model:

```bash
ollama create qwen3:32b-fp16-tc -f Modelfile_qwen3_32b_fp16_tc
```

##### Example: QwQ 32B FP16 with Extended Context

Create a Modelfile named `Modelfile_qwq_32b_fp16_tc`:

```dockerfile
FROM qwq:32b-fp16
PARAMETER num_ctx 110000
PARAMETER temperature 0.2
PARAMETER top_p 0.7
PARAMETER min_p 0.0
PARAMETER top_k 40
PARAMETER repeat_penalty 1.2
```

Build the custom model:

```bash
ollama create qwq:32b-fp16-tc -f Modelfile_qwq_32b_fp16_tc
```

> **Note**: The QwQ 32B FP16 model requires approximately **71.3 GB VRAM** for inference. Ensure your system has sufficient GPU memory before attempting to use this model.

These custom models are referenced in the pre-built provider configuration files (`ollama-qwen332b-fp16-tc.provider.yml` and `ollama-qwq32b-fp16-tc.provider.yml`) that are included in the Docker image at `/opt/pentagi/conf/`.

### OpenAI Provider Configuration

PentAGI integrates with OpenAI's comprehensive model lineup, featuring advanced reasoning capabilities with extended chain-of-thought, agentic models with enhanced tool integration, and specialized code models for security engineering.

#### Configuration Variables

| Variable             | Default                     | Description                 |
| -------------------- | --------------------------- | --------------------------- |
| `OPEN_AI_KEY`        |                             | API key for OpenAI services |
| `OPEN_AI_SERVER_URL` | `https://api.openai.com/v1` | OpenAI API endpoint         |

#### Configuration Examples

```bash
# Basic OpenAI setup
OPEN_AI_KEY=your_openai_api_key
OPEN_AI_SERVER_URL=https://api.openai.com/v1

# Using with proxy for enhanced security
OPEN_AI_KEY=your_openai_api_key
PROXY_URL=http://your-proxy:8080
```

#### Supported Models

PentAGI supports 31 OpenAI models with tool calling, streaming, reasoning modes, and prompt caching. Models marked with `*` are used in default configuration.

**GPT-5.2 Series - Latest Flagship Agentic (December 2025)**

| Model ID              | Thinking | Price (Input/Output/Cache) | Use Case                                        |
| --------------------- | -------- | -------------------------- | ----------------------------------------------- |
| `gpt-5.2`*            | ✅        | $1.75/$14.00/$0.18         | Latest flagship with enhanced reasoning and tool integration, autonomous security research |
| `gpt-5.2-pro`         | ✅        | $21.00/$168.00/$0.00       | Premium version with superior agentic coding, mission-critical security research, zero-day discovery |
| `gpt-5.2-codex`       | ✅        | $1.75/$14.00/$0.18         | Most advanced code-specialized, context compaction, strong cybersecurity capabilities |

**GPT-5/5.1 Series - Advanced Agentic Models**

| Model ID              | Thinking | Price (Input/Output/Cache) | Use Case                                        |
| --------------------- | -------- | -------------------------- | ----------------------------------------------- |
| `gpt-5`               | ✅        | $1.25/$10.00/$0.13         | Premier agentic with advanced reasoning, autonomous security research, exploit chain development |
| `gpt-5.1`             | ✅        | $1.25/$10.00/$0.13         | Enhanced agentic with adaptive reasoning, balanced penetration testing with strong tool coordination |
| `gpt-5-pro`           | ✅        | $15.00/$120.00/$0.00       | Premium version with major reasoning improvements, reduced hallucinations, critical security operations |
| `gpt-5-mini`          | ✅        | $0.25/$2.00/$0.03          | Efficient balancing speed and intelligence, automated vulnerability analysis, exploit generation |
| `gpt-5-nano`          | ✅        | $0.05/$0.40/$0.01          | Fastest for high-throughput scanning, reconnaissance, bulk vulnerability detection |

**GPT-5/5.1 Codex Series - Code-Specialized**

| Model ID              | Thinking | Price (Input/Output/Cache) | Use Case                                        |
| --------------------- | -------- | -------------------------- | ----------------------------------------------- |
| `gpt-5.1-codex-max`   | ✅        | $1.25/$10.00/$0.13         | Enhanced reasoning for sophisticated coding, proven CVE findings, systematic exploit development |
| `gpt-5.1-codex`       | ✅        | $1.25/$10.00/$0.13         | Standard code-optimized with strong reasoning, exploit generation, vulnerability analysis |
| `gpt-5-codex`         | ✅        | $1.25/$10.00/$0.13         | Foundational code-specialized, vulnerability scanning, basic exploit generation |
| `gpt-5.1-codex-mini`  | ✅        | $0.25/$2.00/$0.03          | Compact high-performance, 4x higher capacity, rapid vulnerability detection |
| `codex-mini-latest`   | ✅        | $1.50/$6.00/$0.38          | Latest compact code model, automated code review, basic vulnerability analysis |

**GPT-4.1 Series - Enhanced Intelligence**

| Model ID              | Thinking | Price (Input/Output/Cache) | Use Case                                        |
| --------------------- | -------- | -------------------------- | ----------------------------------------------- |
| `gpt-4.1`             | ❌        | $2.00/$8.00/$0.50          | Enhanced flagship with superior function calling, complex threat analysis, sophisticated exploit development |
| `gpt-4.1-mini`*       | ❌        | $0.40/$1.60/$0.10          | Balanced performance with improved efficiency, routine security assessments, automated code analysis |
| `gpt-4.1-nano`        | ❌        | $0.10/$0.40/$0.03          | Ultra-fast lightweight, bulk security scanning, rapid reconnaissance, continuous monitoring |

**GPT-4o Series - Multimodal Flagship**

| Model ID              | Thinking | Price (Input/Output/Cache) | Use Case                                        |
| --------------------- | -------- | -------------------------- | ----------------------------------------------- |
| `gpt-4o`              | ❌        | $2.50/$10.00/$1.25         | Multimodal flagship with vision, image analysis, web UI assessment, multi-tool orchestration |
| `gpt-4o-mini`         | ❌        | $0.15/$0.60/$0.08          | Compact multimodal with strong function calling, high-frequency scanning, cost-effective bulk operations |

**o-Series - Advanced Reasoning Models**

| Model ID              | Thinking | Price (Input/Output/Cache) | Use Case                                        |
| --------------------- | -------- | -------------------------- | ----------------------------------------------- |
| `o4-mini`*            | ✅        | $1.10/$4.40/$0.28          | Next-gen reasoning with enhanced speed, methodical security assessments, systematic exploit development |
| `o3`*                 | ✅        | $2.00/$8.00/$0.50          | Advanced reasoning powerhouse, multi-stage attack chains, deep vulnerability analysis |
| `o3-mini`             | ✅        | $1.10/$4.40/$0.55          | Compact reasoning with extended thinking, step-by-step attack planning, logical vulnerability chaining |
| `o1`                  | ✅        | $15.00/$60.00/$7.50        | Premier reasoning with maximum depth, advanced penetration testing, novel exploit research |
| `o3-pro`              | ✅        | $20.00/$80.00/$0.00        | Most advanced reasoning, 80% cheaper than o1-pro, zero-day research, critical security investigations |
| `o1-pro`              | ✅        | $150.00/$600.00/$0.00      | Previous-gen premium reasoning, exhaustive security analysis, mission-critical challenges |

**Prices**: Per 1M tokens. Reasoning models include thinking tokens in output pricing.

> [!WARNING]
> **GPT-5* Models - Trusted Access Required**
>
> All GPT-5 series models (`gpt-5`, `gpt-5.1`, `gpt-5.2`, `gpt-5-pro`, `gpt-5.2-pro`, and all Codex variants) work **unstably with PentAGI** and may trigger OpenAI's cybersecurity safety mechanisms without verified access.
>
> **To use GPT-5* models reliably:**
> 1. **Individual users**: Verify your identity at [chatgpt.com/cyber](https://chatgpt.com/cyber)
> 2. **Enterprise teams**: Request trusted access through your OpenAI representative
> 3. **Security researchers**: Apply for the [Cybersecurity Grant Program](https://openai.com/form/cybersecurity-grant-program/) (includes $10M in API credits)
>
> **Recommended alternatives without verification:**
> - Use `o-series` models (o3, o4-mini, o1) for reasoning tasks
> - Use `gpt-4.1` series for general intelligence and function calling
> - All o-series and gpt-4.x models work reliably without special access

**Reasoning Effort Levels**:
- **High**: Maximum reasoning depth (refiner - o3 with high effort)
- **Medium**: Balanced reasoning (primary_agent, assistant, reflector - o4-mini/o3 with medium effort)
- **Low**: Efficient targeted reasoning (coder, installer, pentester - o3/o4-mini with low effort; adviser - gpt-5.2 with low effort)

**Key Features**:
- **Extended Reasoning**: o-series models with chain-of-thought for complex security analysis
- **Agentic Intelligence**: GPT-5/5.1/5.2 series with enhanced tool integration and autonomous capabilities
- **Prompt Caching**: Cost reduction on repeated context (10-50% of input price)
- **Code Specialization**: Dedicated Codex models for vulnerability discovery and exploit development
- **Multimodal Support**: GPT-4o series for vision-based security assessments
- **Tool Calling**: Robust function calling across all models for pentesting tool orchestration
- **Streaming**: Real-time response streaming for interactive workflows
- **Proven Track Record**: Industry-leading models with CVE discoveries and real-world security applications

