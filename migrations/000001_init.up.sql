begin;

create table todo
(
    id          serial primary key,
    name        text not null,
    description text not null
);

commit;
