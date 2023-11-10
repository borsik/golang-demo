-- +migrate Up
create table if not exists users
(
    id         uuid                     default gen_random_uuid() not null primary key,
    first_name text,
    last_name  text,
    nickname   text
        constraint idx_users_nickname unique,
    password   text,
    email      text,
    country    text,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

-- +migrate Down
drop table users;