create table entity
(
    id      serial
        constraint entity_pk
            primary key,
    name    varchar(50) not null,
    age     integer     not null,
    deleted boolean     not null default false
);

alter table entity
    owner to postgres;

create unique index users_id_uindex
    on entity (id);

create unique index entity_name_uindex
    on entity (name);