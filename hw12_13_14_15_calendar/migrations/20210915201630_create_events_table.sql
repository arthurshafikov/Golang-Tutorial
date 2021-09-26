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

insert into events(title, descr, owner, start_at, end_at, send_notification_at)
values('new year', 'watch the irony of fate', 42, '2019-12-31:10:10:00', '2019-12-31:12:01:02', NOW());

insert into events(title, descr, owner, start_at, end_at, send_notification_at)
values('christmas', 'meet santa', 42, '2020-01-07:10:10:00', '2020-01-07:12:01:02', NOW());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX owner_idx;
DROP INDEX start_at_idx;
DROP TABLE events;
-- +goose StatementEnd
