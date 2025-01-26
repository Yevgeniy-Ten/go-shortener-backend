-- +goose Up
-- +goose StatementBegin

ALTER TABLE urls ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE urls DROP COLUMN is_deleted;

-- +goose StatementEnd
