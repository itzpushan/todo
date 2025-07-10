-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "name" character varying NOT NULL,
  "email" character varying NOT NULL,
  "password" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "public"."users" ("email");
-- Create "todos" table
CREATE TABLE "public"."todos" (
  "id" uuid NOT NULL,
  "title" character varying NOT NULL,
  "description" character varying NULL,
  "author_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "todos_users_todos" FOREIGN KEY ("author_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
