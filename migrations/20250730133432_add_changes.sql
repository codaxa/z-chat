-- Modify "messages" table
ALTER TABLE "public"."messages" ALTER COLUMN "created_at" DROP NOT NULL, ALTER COLUMN "created_at" SET DEFAULT now(), ALTER COLUMN "updated_at" DROP NOT NULL, ALTER COLUMN "updated_at" SET DEFAULT now();
-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "created_at" DROP NOT NULL, ALTER COLUMN "created_at" SET DEFAULT now(), ALTER COLUMN "updated_at" DROP NOT NULL, ALTER COLUMN "updated_at" SET DEFAULT now();
