-- Modify "preferences" table
ALTER TABLE "public"."preferences" ADD COLUMN "measurements_bust" numeric NULL, ADD COLUMN "measurements_under_bust" numeric NULL, ADD COLUMN "measurements_waist" numeric NULL, ADD COLUMN "measurements_hip" numeric NULL;
