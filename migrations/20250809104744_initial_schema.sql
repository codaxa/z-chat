-- Create "rooms" table
CREATE TABLE "public"."rooms" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(50) NOT NULL,
  "created_by" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_rooms_name" to table: "rooms"
CREATE UNIQUE INDEX "idx_rooms_name" ON "public"."rooms" ("name");
-- Create "messages" table
CREATE TABLE "public"."messages" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "sender" character varying(255) NOT NULL,
  "receiver" character varying(255) NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NULL DEFAULT now(),
  "updated_at" timestamptz NULL DEFAULT now(),
  "room_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_rooms_messages" FOREIGN KEY ("room_id") REFERENCES "public"."rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
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
  "created_at" timestamptz NULL DEFAULT now(),
  "updated_at" timestamptz NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "public"."users" ("username");
-- Create "room_admins" table
CREATE TABLE "public"."room_admins" (
  "room_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL DEFAULT gen_random_uuid(),
  PRIMARY KEY ("room_id", "user_id"),
  CONSTRAINT "fk_room_admins_room" FOREIGN KEY ("room_id") REFERENCES "public"."rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_room_admins_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "room_members" table
CREATE TABLE "public"."room_members" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "room_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "joined_at" timestamptz NULL,
  "role" character varying(20) NULL DEFAULT 'member',
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_room_members_room" FOREIGN KEY ("room_id") REFERENCES "public"."rooms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_room_members_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
