CREATE TABLE "testers" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX ON "testers" ("email", "user_id")