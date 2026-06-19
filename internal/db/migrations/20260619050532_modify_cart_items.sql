-- +goose Up
-- drop index "idx_cart_items_user_id" from table: "cart_items"
DROP INDEX "public"."idx_cart_items_user_id";
-- modify "cart_items" table
ALTER TABLE "public"."cart_items" DROP CONSTRAINT "cart_items_pkey", DROP COLUMN "id", ALTER COLUMN "user_id" DROP NOT NULL, ALTER COLUMN "list_product_id" DROP NOT NULL;

-- +goose Down
-- reverse: modify "cart_items" table
ALTER TABLE "public"."cart_items" ALTER COLUMN "list_product_id" SET NOT NULL, ALTER COLUMN "user_id" SET NOT NULL, ADD COLUMN "id" bigserial NOT NULL, ADD PRIMARY KEY ("id");
-- reverse: drop index "idx_cart_items_user_id" from table: "cart_items"
CREATE UNIQUE INDEX "idx_cart_items_user_id" ON "public"."cart_items" ("user_id");
