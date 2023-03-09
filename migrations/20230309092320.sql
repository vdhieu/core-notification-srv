-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS webhook_infos (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "callback_url" text,
    "account_id" text,
    "secret" text,

    "created_at" timestamptz NOT NULL DEFAULT now()
);
-- +migrate Down
-- DROP TABLE IF EXISTS "webhook_infos";