CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "hash" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT now(),
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00+00'
);


ALTER TABLE accounts 
ALTER COLUMN owner TYPE BIGINT USING owner::bigint;

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");