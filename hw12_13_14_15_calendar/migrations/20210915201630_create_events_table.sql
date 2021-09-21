-- +goose Up
-- +goose StatementBegin
create table events (
    id serial primary key,
    title text,
    descr text,
    owner bigint,
    start_date date not null,
    start_time time,
    end_date date not null,
    end_time time,
    send_notification_at timestamp not null
);
create index owner_idx on events (owner);
create index start_idx on events using btree (start_date, start_time);

insert into events(title, descr, owner, start_date, start_time, end_date, end_time, send_notification_at)
values('new year', 'watch the irony of fate', 42, '2019-12-31', CURRENT_TIME, '2019-12-31', CURRENT_TIME, NOW());

insert into events(title, descr, owner, start_date, start_time, end_date, end_time, send_notification_at)
values('christmas', 'meet santa', 42, '2020-01-07', CURRENT_TIME, '2020-01-07', CURRENT_TIME, NOW());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX owner_idx;
DROP INDEX start_idx;
DROP TABLE events;
-- +goose StatementEnd
