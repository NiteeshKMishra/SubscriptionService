-- +migrate Up

CREATE TABLE IF NOT EXISTS plans (
                            id text unique not null default uuid_generate_v4()::text primary key,
                            plan_name citext unique not null,
                            plan_amount integer,
                            created_at timestamp without time zone,
                            updated_at timestamp without time zone
);

-- +migrate Down

DROP TABLE IF EXISTS plans;
