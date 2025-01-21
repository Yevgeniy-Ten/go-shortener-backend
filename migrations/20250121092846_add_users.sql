-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users
(
    ID
    SERIAL
    PRIMARY
    KEY
);
ALTER TABLE urls
    ADD COLUMN user_id INTEGER NOT NULL REFERENCES users (ID);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE urls DROP COLUMN user_id;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
