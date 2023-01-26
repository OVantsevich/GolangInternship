create table users
(
    id       serial
        constraint User_pk
            primary key,
    login    varchar(50)                               not null,
    email    varchar(50)                               not null,
    password varchar(200)                              not null,
    name     varchar(50)                               not null,
    age      integer                                   not null,
    token    varchar(200),
    deleted  boolean                                   not null default false,
    created  timestamp(6) default CURRENT_TIMESTAMP(6) not null,
    updated  timestamp(6) default CURRENT_TIMESTAMP(6) not null
);

alter table users
    owner to postgres;

create unique index users_id_uindex
    on users (id);

create unique index users_login_uindex
    on users (login);