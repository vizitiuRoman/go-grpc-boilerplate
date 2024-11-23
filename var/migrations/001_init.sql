-- +goose Up
-- +goose StatementBegin

create table todo
(
    id          serial primary key,
    name        text not null,
    description text not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists todo;

-- +goose StatementEnd
