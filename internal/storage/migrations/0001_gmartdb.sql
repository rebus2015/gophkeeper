-- +goose Up
-- +goose StatementBegin

create table if not exists  users
(
    id    uuid    not null
        constraint users_pk
            primary key,
    hash  bytea   not null,
    login varchar not null
        constraint users_un
            unique
);

create unique index users_id_idx
    on users (id);


create or replace function user_add(_login character varying, _hash bytea) returns character varying
    language sql
as
$$
   INSERT INTO users (id, hash, login)
values (gen_random_uuid (),_hash,_login)
ON CONFLICT on constraint users_un
do nothing
returning cast(id as varchar);
$$;

create or replace function public.user_check(_login character varying, OUT id character varying, OUT hash bytea) returns record
    language sql
as
$$
   select
       cast(u.id as varchar) as id,
       u.hash AS hash
   from users u
   where login =_login
$$;

-- +goose StatementEnd