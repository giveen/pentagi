MCP (Model Context Protocol) Server Design
======================================

Goal
----
Add first-class support for MCP servers so PentAGI can configure, persist and invoke local or remote connectors that expose tool-like operations over a small adapter protocol (STDIO or SSE). The initial MVP focuses on safe, local STDIO connectors and a CRUD+test API surface.

Data model (MVP)
----------------
- `mcp_servers`
  - `id` BIGINT PK
  - `name` TEXT UNIQUE
  - `transport` TEXT ("stdio"|"sse")
  - `stdio_command` TEXT NULL
  - `stdio_args` TEXT NULL
  - `stdio_env` JSONB DEFAULT '[]'::jsonb
  - `sse_url` TEXT NULL
  - `sse_headers` JSONB DEFAULT '[]'::jsonb
  - `tools` JSONB DEFAULT '[]'::jsonb
  - `created_at`, `updated_at` TIMESTAMPTZ

Notes:
- We store structured fields (`stdio_env`, `sse_headers`, `tools`) as JSONB for flexibility.
- The frontend expects env as an array of {key,value}; GraphQL will expose typed arrays for convenience and convert to/from JSON at persistence layer.

GraphQL API (MVP)
-----------------
- Queries
  - `mcpServers: [McpServer!]!` — list all servers
  - `mcpServer(mcpServerId: ID!): McpServer` — fetch single server

- Mutations
  - `createMcpServer(input: CreateMcpServerInput!): McpServer`
  - `updateMcpServer(mcpServerId: ID!, input: UpdateMcpServerInput!): McpServer`
  - `deleteMcpServer(mcpServerId: ID!): ResultType!`
  - `testMcpServer(mcpServerId: ID!): ResultType!` — perform a basic connectivity/test invocation

Runtime connector (stdio - MVP)
-------------------------------
- The stdio connector is a local process (binary or node script) launched by PentAGI or referenced by an absolute `command` and `args`.
- Connector input/output uses MCP-style JSON messages on stdin/stdout (simple request/response). For MVP testing we will restrict to a small "ping" or "echo" tool to validate connectivity.
- Environment variables MUST be passed via an explicit `env` map (no shell interpolation).

Security & operational notes
----------------------------
- Do not allow arbitrary commands via the UI in production; restrict to pre-approved commands or run connectors inside isolated containers.
- Secrets (auth headers) must be stored securely (vault/secret store) before enabling public SSE endpoints — MVP will accept plaintext for local-testing only and will require migration to secret storage for production.
- Limit privileges: only admin roles (or roles with `settings.mcp.*`) may create/update/delete MCP servers.

Next steps (implementation plan)
-------------------------------
1. Add DB migration + SQLC models (simple JSONB fields + helper converters).
2. Add backend service + GraphQL resolvers to persist and return typed structs.
3. Implement a safe stdio runtime connector helper to spawn processes with given env map and limited args.
4. Wire the frontend settings page to call the real API and support testing.

Design references
-----------------
- Frontend form: `frontend/src/pages/settings/settings-mcp-server.tsx` — currently mocks the UI; use this shape for GraphQL types.
