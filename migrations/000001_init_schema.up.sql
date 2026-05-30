CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA sh_eff;

CREATE TABLE sh_eff.subscriptions
(
    id         UUID PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    price      SMALLINT NOT NULL CHECK (price >= 0),
    user_id    UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date   DATE CHECK (end_date IS NULL OR end_date > start_date),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_user_id ON sh_eff.subscriptions(user_id);
CREATE INDEX idx_subscriptions_name ON sh_eff.subscriptions(name);
CREATE INDEX idx_subscriptions_dates ON sh_eff.subscriptions(start_date, end_date);