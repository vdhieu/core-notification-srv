-- create "webhook_infos" table
CREATE TABLE "public"."webhook_infos" ("id" bytea NOT NULL, "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, "account_id" text NOT NULL, "callback_url" text NOT NULL, "secret" text NOT NULL, PRIMARY KEY ("id"));
-- create index "idx_webhook_infos_account_id" to table: "webhook_infos"
CREATE INDEX "idx_webhook_infos_account_id" ON "public"."webhook_infos" ("account_id");
