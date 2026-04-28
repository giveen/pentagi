// Code generated (manually) to mirror sqlc-generated types for MCP servers.
package database

import (
    "database/sql"
    "database/sql/driver"
    "encoding/json"
    "fmt"
)

type McpTransport string

const (
    McpTransportStdio McpTransport = "stdio"
    McpTransportSse   McpTransport = "sse"
)

func (e *McpTransport) Scan(src interface{}) error {
    switch s := src.(type) {
    case []byte:
        *e = McpTransport(s)
    case string:
        *e = McpTransport(s)
    default:
        return fmt.Errorf("unsupported scan type for McpTransport: %T", src)
    }
    return nil
}

type NullMcpTransport struct {
    McpTransport McpTransport `json:"mcp_transport"`
    Valid        bool         `json:"valid"`
}

// Scan implements the Scanner interface.
func (ns *NullMcpTransport) Scan(value interface{}) error {
    if value == nil {
        ns.McpTransport, ns.Valid = "", false
        return nil
    }
    ns.Valid = true
    return ns.McpTransport.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMcpTransport) Value() (driver.Value, error) {
    if !ns.Valid {
        return nil, nil
    }
    return string(ns.McpTransport), nil
}

type McpServer struct {
    ID           int64           `json:"id"`
    Name         string          `json:"name"`
    Transport    McpTransport    `json:"transport"`
    StdioCommand sql.NullString  `json:"stdio_command"`
    StdioArgs    sql.NullString  `json:"stdio_args"`
    StdioEnv     json.RawMessage `json:"stdio_env"`
    SseUrl       sql.NullString  `json:"sse_url"`
    SseHeaders   json.RawMessage `json:"sse_headers"`
    Tools        json.RawMessage `json:"tools"`
    CreatedAt    sql.NullTime    `json:"created_at"`
    UpdatedAt    sql.NullTime    `json:"updated_at"`
}
