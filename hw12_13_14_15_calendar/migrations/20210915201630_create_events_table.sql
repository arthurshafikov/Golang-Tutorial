-- +goose Up
-- +goose StatementBegin
create table events (
    id serial primary key,
    title text,
    descr text,
    owner bigint,
    start_at timestamp not null,
    end_at timestamp not null,
    send_notification_at timestamp not null
);
create index owner_idx on events (owner);
create index start_at_idx on events (start_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX start_at_idx;
DROP INDEX owner_idx;
DROP TABLE events;
-- +goose StatementEnd
