create table user
(
    id       serial
        constraint User_pk
            primary key,
    login    varchar(50)  not null,
    email    varchar(50)  not null,
    password varchar(200) not null,
    name     varchar(50)  not null,
    age      integer      not null,
    deleted  boolean      not null default false
);

alter table user
    owner to postgres;

create unique index user_id_uindex
    on user (id);

create unique index user_login_uindex
    on user (login);