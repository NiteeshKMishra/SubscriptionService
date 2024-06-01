-- +migrate Up

create extension if not exists "uuid-ossp";
create extension if not exists "citext";

-- +migrate Down
