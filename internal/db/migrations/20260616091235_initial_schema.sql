-- +goose Up
-- create "items" table
CREATE TABLE "public"."items" (
  "id" bigserial NOT NULL,
  "sku" text NULL,
  "name" text NULL,
  "quantity" bigint NULL,
  "price" numeric NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "username" text NOT NULL,
  "password" text NOT NULL,
  "can_read" boolean NULL DEFAULT true,
  "can_write" boolean NULL DEFAULT false,
  "can_update" boolean NULL DEFAULT false,
  "can_delete" boolean NULL DEFAULT false,
  PRIMARY KEY ("id")
);
-- create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "public"."users" ("username");

-- +goose Down
-- reverse: create index "idx_users_username" to table: "users"
DROP INDEX "public"."idx_users_username";
-- reverse: create "users" table
DROP TABLE "public"."users";
-- reverse: create "items" table
DROP TABLE "public"."items";
