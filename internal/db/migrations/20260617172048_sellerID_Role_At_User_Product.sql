-- +goose Up
-- modify "products" table
ALTER TABLE "public"."products" ADD COLUMN "created_by" bigint NOT NULL;
-- modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "is_admin", DROP COLUMN "is_seller", DROP COLUMN "is_customer", ADD COLUMN "role" text NOT NULL;

-- +goose Down
-- reverse: modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "role", ADD COLUMN "is_customer" boolean NULL DEFAULT false, ADD COLUMN "is_seller" boolean NULL DEFAULT false, ADD COLUMN "is_admin" boolean NULL DEFAULT false;
-- reverse: modify "products" table
ALTER TABLE "public"."products" DROP COLUMN "created_by";
