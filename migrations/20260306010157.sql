-- Modify "products" table
ALTER TABLE "public"."products" ADD COLUMN "catalog_id" bigint NULL, ADD CONSTRAINT "fk_products_catalog" FOREIGN KEY ("catalog_id") REFERENCES "public"."catalogs" ("id") ON UPDATE CASCADE ON DELETE SET NULL;
-- Create index "idx_products_catalog_id" to table: "products"
CREATE INDEX "idx_products_catalog_id" ON "public"."products" ("catalog_id");
-- Create index "idx_products_is_exclusive" to table: "products"
CREATE INDEX "idx_products_is_exclusive" ON "public"."products" ("is_exclusive");
