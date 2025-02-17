-- +goose NO TRANSACTION

-- +goose Up
-- +goose StatementBegin
CREATE INDEX inventory_username_idx ON inventory (username);
CREATE INDEX history_username_idx ON history (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX history_username_idx;
DROP INDEX inventory_username_idx;
-- +goose StatementEnd