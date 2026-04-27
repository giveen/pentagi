-- name: GetMcpServers :many
SELECT
  id, name, transport, stdio_command, stdio_args, stdio_env, sse_url, sse_headers, tools, created_at, updated_at
FROM mcp_servers
ORDER BY created_at DESC;

-- name: GetMcpServer :one
SELECT
  id, name, transport, stdio_command, stdio_args, stdio_env, sse_url, sse_headers, tools, created_at, updated_at
FROM mcp_servers
WHERE id = $1;

-- name: GetMcpServerByName :one
SELECT
  id, name, transport, stdio_command, stdio_args, stdio_env, sse_url, sse_headers, tools, created_at, updated_at
FROM mcp_servers
WHERE name = $1;

-- name: CreateMcpServer :one
INSERT INTO mcp_servers (
  name, transport, stdio_command, stdio_args, stdio_env, sse_url, sse_headers, tools
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, name, transport, stdio_command, stdio_args, stdio_env, sse_url, sse_headers, tools, created_at, updated_at;

-- name: UpdateMcpServer :one
UPDATE mcp_servers
SET name = $1,
    transport = $2,
    stdio_command = $3,
    stdio_args = $4,
    stdio_env = $5,
    sse_url = $6,
    sse_headers = $7,
    tools = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $9
RETURNING id, name, transport, stdio_command, stdio_args, stdio_env, sse_url, sse_headers, tools, created_at, updated_at;

-- name: DeleteMcpServer :exec
DELETE FROM mcp_servers
WHERE id = $1;
