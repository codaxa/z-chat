-- Create "messages" table
CREATE TABLE "public"."messages" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "sender" character varying(255) NOT NULL,
  "receiver" character varying(255) NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_messages_created_at" to table: "messages"
CREATE INDEX "idx_messages_created_at" ON "public"."messages" ("created_at");
-- Create index "idx_messages_receiver" to table: "messages"
CREATE INDEX "idx_messages_receiver" ON "public"."messages" ("receiver");
-- Create index "idx_messages_sender" to table: "messages"
CREATE INDEX "idx_messages_sender" ON "public"."messages" ("sender");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "username" character varying(50) NOT NULL,
  "email" character varying(255) NOT NULL,
  "password" character varying(64) NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "public"."users" ("username");
