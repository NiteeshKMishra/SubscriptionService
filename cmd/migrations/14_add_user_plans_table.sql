-- +migrate Up

CREATE TABLE IF NOT EXISTS user_plans (
                            id text unique not null default uuid_generate_v4()::text primary key,
                            user_id text REFERENCES users(id) ON UPDATE RESTRICT ON DELETE CASCADE,
                            plan_id text REFERENCES plans(id) ON UPDATE RESTRICT ON DELETE CASCADE,
                            created_at timestamp without time zone,
                            updated_at timestamp without time zone
);

-- +migrate Down
DROP TABLE IF EXISTS user_plans;
