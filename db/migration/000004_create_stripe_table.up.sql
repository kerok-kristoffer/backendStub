CREATE TABLE "stripe" (
  "id" uuid PRIMARY KEY,
  "user_id" bigserial NOT NULL,
  "stripe_customer_id" varchar,
  "stripe_plan_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "stripe_plans" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "user_access_id" int4 NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "stripe" ADD FOREIGN KEY ("user_id") REFERENCES "users" (id);
ALTER TABLE "stripe" ADD FOREIGN KEY ("stripe_plan_id") REFERENCES "stripe_plans" (id);

CREATE INDEX ON "stripe" ("id", "user_id", "stripe_plan_id");
CREATE INDEX ON "stripe_plans" ("id", "user_access_id");
