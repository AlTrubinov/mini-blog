-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN password TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN IF EXISTS password;
-- +goose StatementEnd
