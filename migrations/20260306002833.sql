-- Modify "comments" table
ALTER TABLE "public"."comments" DROP CONSTRAINT "fk_comments_catalog", ADD CONSTRAINT "fk_comments_catalog" FOREIGN KEY ("catalog_id") REFERENCES "public"."catalogs" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
