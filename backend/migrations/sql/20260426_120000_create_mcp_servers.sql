-- +goose Up
-- +goose StatementBegin
-- Create MCP servers storage (MVP)

CREATE TABLE mcp_servers (
  id            BIGINT        PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  name          TEXT          NOT NULL,
  transport     TEXT          NOT NULL,
  stdio_command TEXT          NULL,
  stdio_args    TEXT          NULL,
  stdio_env     JSONB         NOT NULL DEFAULT '[]'::jsonb,
  sse_url       TEXT          NULL,
  sse_headers   JSONB         NOT NULL DEFAULT '[]'::jsonb,
  tools         JSONB         NOT NULL DEFAULT '[]'::jsonb,
  created_at    TIMESTAMPTZ   DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMPTZ   DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX mcp_servers_name_unique_idx ON mcp_servers(name);
CREATE INDEX mcp_servers_transport_idx ON mcp_servers(transport);
CREATE INDEX mcp_servers_created_at_idx ON mcp_servers(created_at);

-- Add privileges for Admin role (role_id = 1)
INSERT INTO privileges (role_id, name) VALUES
  (1, 'settings.mcp.admin'),
  (1, 'settings.mcp.create'),
  (1, 'settings.mcp.view'),
  (1, 'settings.mcp.edit'),
  (1, 'settings.mcp.delete'),
  (1, 'settings.mcp.subscribe')
  ON CONFLICT DO NOTHING;

-- Add privileges for User role (role_id = 2)
INSERT INTO privileges (role_id, name) VALUES
  (2, 'settings.mcp.create'),
  (2, 'settings.mcp.view'),
  (2, 'settings.mcp.edit'),
  (2, 'settings.mcp.delete'),
  (2, 'settings.mcp.subscribe')
  ON CONFLICT DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM privileges WHERE name IN (
  'settings.mcp.create',
  'settings.mcp.view',
  'settings.mcp.edit',
  'settings.mcp.delete',
  'settings.mcp.admin',
  'settings.mcp.subscribe'
);

DROP INDEX IF EXISTS mcp_servers_name_unique_idx;
DROP TABLE IF EXISTS mcp_servers;
-- +goose StatementEnd
