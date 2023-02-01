create table if not exists roles
(
    id      serial
        constraint Role_pk
            primary key,
    name    varchar(50)                               not null,
    deleted boolean                                   not null default false,
    created timestamp(6) default CURRENT_TIMESTAMP(6) not null,
    updated timestamp(6) default CURRENT_TIMESTAMP(6) not null
);

create unique index if not exists roles_name_index
    on roles (name);

create table if not exists l_role_user
(
    id      serial
        constraint l_role_user_pk
            primary key,
    user_id int not null
        constraint l_role_user_user_id_fk
            references users
            on update cascade on delete cascade,
    role_id int not null
        constraint l_role_user_role_id_fk
            references roles
            on update cascade on delete cascade
);

alter table l_role_user
    owner to postgres;

create unique index if not exists l_role_user_id_uindex
    on l_role_user (id);

create index if not exists l_role_user_user_id_role_id_index
    on l_role_user (user_id, role_id);

