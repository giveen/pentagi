package database

import (
    "context"
    "database/sql"
    "errors"
)

// BeginTx starts a transaction on the underlying *sql.DB if available.
// This is a small helper to enable transactional workflows from callers
// that hold a *Queries instance. If the underlying DB does not support
// BeginTx (for example when Queries was created with a *sql.Tx), an error
// is returned.
func (q *Queries) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
    if db, ok := q.db.(*sql.DB); ok {
        return db.BeginTx(ctx, opts)
    }
    return nil, errors.New("database: underlying DB does not support BeginTx")
}
