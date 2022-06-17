CREATE TABLE "ingredients" (
                               "id" bigserial PRIMARY KEY,
                               "name" varchar NOT NULL,
                               "hash" varchar(50) NOT NULL,
                               "user_id" bigint,
                               "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "measures" (
                        "id" bigserial PRIMARY KEY,
                        "name" varchar NOT NULL,
                        "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "units" (
                         "id" bigserial PRIMARY KEY,
                         "user_id" bigint,
                         "ingredient_id" bigint,
                         "amount" float8 NOT NULL,
                         "measure_id" bigint,
                         "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "currencies" (
                              "id" bigserial PRIMARY KEY,
                              "name" varchar NOT NULL,
                              "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "stocks" (
                          "id" bigserial PRIMARY KEY,
                          "unit_id" bigint,
                          "cost" float8 NOT NULL,
                          "currency_id" bigint,
                          "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE "ingredients" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "units" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "units" ADD FOREIGN KEY ("measure_id") REFERENCES "measures" ("id");
ALTER TABLE "units" ADD FOREIGN KEY ("ingredient_id") REFERENCES "ingredients" ("id");

ALTER TABLE "stocks" ADD FOREIGN KEY ("unit_id") REFERENCES "units" ("id");
ALTER TABLE "stocks" ADD FOREIGN KEY ("currency_id") REFERENCES "currencies" ("id");

CREATE INDEX ON "ingredients" ("name");
CREATE INDEX ON "ingredients" ("user_id");

CREATE INDEX ON "units" ("user_id");
CREATE INDEX ON "units" ("ingredient_id");

CREATE INDEX ON "stocks" ("unit_id");