-- +goose Up
-- SQL in this section is executed when the migration is applied.

create extension pgcrypto;

create table accounts (
    id serial primary key,
    public_id text not null,

    full_name text not null,
    email text not null,

    balance int not null default 0
);

create table tasks (
    id serial primary key,
    public_id text not null default gen_random_uuid(),

    title text not null,
    description text not null,
    cost_done int not null,
    cost_assigne int not null,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table transactions (
    id serial primary key,
    
    account_id int not null,
    task_id int not null,
    event_created_at timestmap not null,

    type text not null, -- debit / credit

    created_at timestamp not null default now()
);

create table billing_cycles (
    id serial primary key,

    date timestamp not null default now()
);

create table payments (
    id serial primary key,

    billing_cycle_id int nor null,
    account_id int not null,
    amount int not null,

    status text "not-paid",

    created_at timestamp not null default now()
);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

drop table accounts;
drop table tasks;
drop table payments;
