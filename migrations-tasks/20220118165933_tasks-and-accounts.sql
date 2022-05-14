-- +goose Up
-- SQL in this section is executed when the migration is applied.

create extension pgcrypto;

create table tasks (
    id serial primary key,
    public_id text not null default gen_random_uuid(),

    description text not null,
    status text not null default 'active',
    version int not null default 1,

    assignee_id int not null,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table accounts (
    id serial primary key,
    public_id text not null,

    full_name text not null,
    position text not null
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

drop table tasks;
drop table accounts;
