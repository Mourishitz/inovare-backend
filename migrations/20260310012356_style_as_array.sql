-- Modify "preferences" table
ALTER TABLE "public"."preferences"
  ALTER COLUMN "style" TYPE smallint[] USING ARRAY["style"]::smallint[],
  ALTER COLUMN "style" DROP DEFAULT;
