create table users (
    id            UUID primary key,
    username      varchar(255) not null unique,
    email         varchar(255) not null unique,
    password_hash varchar(60) not null,
    created_at    timestamp with time zone not null,
    updated_at    timestamp with time zone not null,
    deleted_at    timestamp with time zone default null
)