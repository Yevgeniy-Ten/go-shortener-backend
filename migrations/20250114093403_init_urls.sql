-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls (
    short_url TEXT PRIMARY KEY,
    url TEXT NOT NULL UNIQUE
);
CREATE INDEX IF NOT EXISTS urls_url_idx ON urls (short_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd
