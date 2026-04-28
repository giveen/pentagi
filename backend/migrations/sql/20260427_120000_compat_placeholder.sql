-- +goose Up
-- +goose StatementBegin
-- Compatibility placeholder migration.
-- Keeps goose version continuity for databases that already have version 20260427 applied.
SELECT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd
