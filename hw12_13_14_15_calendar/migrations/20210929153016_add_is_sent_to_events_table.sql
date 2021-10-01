-- +goose Up
-- +goose StatementBegin
ALTER TABLE events 
ADD is_sent BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
DROP COLUMN is_sent;
-- +goose StatementEnd
