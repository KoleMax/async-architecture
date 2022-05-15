-- +goose Up
-- SQL in this section is executed when the migration is applied.

create extension pgcrypto;

create table accounts (
    id serial primary key,

    email text not null unique,
    password text not null,

    public_id text not null default gen_random_uuid(),

    full_name text not null,
    position text not null
);

insert into accounts (email, password, full_name, position) values ('p1@pmail.po', 'pass', 'popug1', 'admin');
insert into accounts (email, password, full_name, position) values ('p2@pmail.po', 'pass', 'popug2', 'manager');
insert into accounts (email, password, full_name, position) values ('p3@pmail.po', 'pass', 'popug3', 'accountant');
insert into accounts (email, password, full_name, position) values ('p4@pmail.po', 'pass', 'popug4', 'worker');
insert into accounts (email, password, full_name, position) values ('p5@pmail.po', 'pass', 'popug5', 'worker');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

drop table accounts;
