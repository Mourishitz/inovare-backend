-- Create "catalogs" table
CREATE TABLE "public"."catalogs" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "url" text NOT NULL,
  "package" smallint NOT NULL,
  "approved" boolean NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_catalogs_url" UNIQUE ("url")
);
-- Create index "idx_catalogs_deleted_at" to table: "catalogs"
CREATE INDEX "idx_catalogs_deleted_at" ON "public"."catalogs" ("deleted_at");
-- Create "products" table
CREATE TABLE "public"."products" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NULL,
  "description" text NULL,
  "image_url" text NULL,
  "is_exclusive" boolean NULL DEFAULT false,
  PRIMARY KEY ("id")
);
-- Create index "idx_products_deleted_at" to table: "products"
CREATE INDEX "idx_products_deleted_at" ON "public"."products" ("deleted_at");
-- Create "catalog_products" table
CREATE TABLE "public"."catalog_products" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "price" numeric NULL,
  "is_bought" boolean NULL DEFAULT false,
  "catalog_id" bigint NULL,
  "product_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_catalog_products_catalog" FOREIGN KEY ("catalog_id") REFERENCES "public"."catalogs" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "fk_catalog_products_product" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "idx_catalog_products_deleted_at" to table: "catalog_products"
CREATE INDEX "idx_catalog_products_deleted_at" ON "public"."catalog_products" ("deleted_at");
-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "phone_number" text NOT NULL;
-- Create "comments" table
CREATE TABLE "public"."comments" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "content" text NULL,
  "author_id" bigint NULL,
  "catalog_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_comments_author" FOREIGN KEY ("author_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "fk_comments_catalog" FOREIGN KEY ("catalog_id") REFERENCES "public"."catalogs" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "idx_comments_deleted_at" to table: "comments"
CREATE INDEX "idx_comments_deleted_at" ON "public"."comments" ("deleted_at");
-- Create "preferences" table
CREATE TABLE "public"."preferences" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "style" smallint NULL DEFAULT 1,
  "favorite_colors" smallint[] NULL,
  "preferred_bra" smallint NULL DEFAULT 1,
  "preferred_model" smallint NULL DEFAULT 1,
  "preferred_panties" smallint NULL DEFAULT 1,
  "size" smallint NULL DEFAULT 1,
  "allowed_models" smallint[] NULL,
  "not_allowed_models" text NULL,
  "notes" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_preferences_deleted_at" to table: "preferences"
CREATE INDEX "idx_preferences_deleted_at" ON "public"."preferences" ("deleted_at");
-- Create "showers" table
CREATE TABLE "public"."showers" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "guests" bigint NULL,
  "shower_date" timestamptz NULL,
  "wedding_date" timestamptz NULL,
  "location" text NULL,
  "host_id" bigint NULL,
  "catalog_id" bigint NULL,
  "preferences_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_showers_catalog" FOREIGN KEY ("catalog_id") REFERENCES "public"."catalogs" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "fk_showers_host" FOREIGN KEY ("host_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "fk_showers_preferences" FOREIGN KEY ("preferences_id") REFERENCES "public"."preferences" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "idx_showers_deleted_at" to table: "showers"
CREATE INDEX "idx_showers_deleted_at" ON "public"."showers" ("deleted_at");
