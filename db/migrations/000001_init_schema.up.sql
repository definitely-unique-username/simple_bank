CREATE TABLE "accounts" (
  "id" BIGSERIAL PRIMARY KEY,
  "owner" varchar NOT NULL DEFAULT (0),
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" BIGSERIAL PRIMARY KEY,
  "account_id" BIGINT NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  FOREIGN KEY ("account_id") REFERENCES "accounts" ("id")
);

CREATE TABLE "transfers" (
  "id" BIGSERIAL PRIMARY KEY,
  "from_account_id" BIGINT NOT NULL,
  "to_account_id" BIGINT NOT NULL,
  "amount" bigint NOT NULL CHECK ("amount" > 0),
  "created_at" timestamp NOT NULL DEFAULT (now()),
  FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id"),
  FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id")
);

CREATE INDEX idx_accounts_owner ON "accounts" ("owner");
CREATE INDEX idx_entries_account_id ON "entries" ("account_id");
CREATE INDEX idx_transfers_from_account_id ON "transfers" ("from_account_id");
CREATE INDEX idx_transfers_to_account_id ON "transfers" ("to_account_id");
CREATE INDEX idx_transfers_from_to_account_id ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';
COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';
