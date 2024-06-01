-- +migrate Up

CREATE TABLE IF NOT EXISTS users (
                            id text unique not null default uuid_generate_v4()::text primary key,
                            email citext unique not null,
                            first_name character varying(255),
                            last_name character varying(255),
                            password character varying(60),
                            user_active boolean default false,
                            is_admin boolean default false,
                            created_at timestamp without time zone,
                            updated_at timestamp without time zone
);

-- +migrate Down

DROP TABLE IF EXISTS users;
