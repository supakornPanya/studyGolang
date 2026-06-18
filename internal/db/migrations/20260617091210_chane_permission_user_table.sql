-- +goose Up
-- rename a column from "can_write" to "is_admin"
ALTER TABLE "public"."users" RENAME COLUMN "can_write" TO "is_admin";
-- rename a column from "can_update" to "is_seller"
ALTER TABLE "public"."users" RENAME COLUMN "can_update" TO "is_seller";
-- rename a column from "can_delete" to "is_customer"
ALTER TABLE "public"."users" RENAME COLUMN "can_delete" TO "is_customer";
-- modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "can_read";
-- create "cart_items" table
CREATE TABLE "public"."cart_items" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "list_product_id" json NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_cart_items_user_id" to table: "cart_items"
CREATE UNIQUE INDEX "idx_cart_items_user_id" ON "public"."cart_items" ("user_id");
-- create "categories" table
CREATE TABLE "public"."categories" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "description" text NOT NULL,
  "parent_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- create "products" table
CREATE TABLE "public"."products" (
  "id" bigserial NOT NULL,
  "sku" text NOT NULL,
  "name" text NOT NULL,
  "description" text NOT NULL,
  "price" numeric NOT NULL,
  "stock" bigint NOT NULL,
  "category" json NULL,
  "tag" json NULL,
  PRIMARY KEY ("id")
);
-- create "tags" table
CREATE TABLE "public"."tags" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  PRIMARY KEY ("id")
);
-- drop "items" table
DROP TABLE "public"."items";

-- +goose Down
-- reverse: drop "items" table
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
-- reverse: create "tags" table
DROP TABLE "public"."tags";
-- reverse: create "products" table
DROP TABLE "public"."products";
-- reverse: create "categories" table
DROP TABLE "public"."categories";
-- reverse: create index "idx_cart_items_user_id" to table: "cart_items"
DROP INDEX "public"."idx_cart_items_user_id";
-- reverse: create "cart_items" table
DROP TABLE "public"."cart_items";
-- reverse: modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "can_read" boolean NULL DEFAULT true;
-- reverse: rename a column from "can_delete" to "is_customer"
ALTER TABLE "public"."users" RENAME COLUMN "is_customer" TO "can_delete";
-- reverse: rename a column from "can_update" to "is_seller"
ALTER TABLE "public"."users" RENAME COLUMN "is_seller" TO "can_update";
-- reverse: rename a column from "can_write" to "is_admin"
ALTER TABLE "public"."users" RENAME COLUMN "is_admin" TO "can_write";
