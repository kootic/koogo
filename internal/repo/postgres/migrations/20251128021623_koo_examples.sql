-- Create "koo_users" table
CREATE TABLE "public"."koo_users" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "is_subscribed" boolean NOT NULL DEFAULT false,
  "first_name" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "koo_pets" table
CREATE TABLE "public"."koo_pets" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "owner_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "koo_pets_owner_id_fkey" FOREIGN KEY ("owner_id") REFERENCES "public"."koo_users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
