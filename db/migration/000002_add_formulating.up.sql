CREATE TABLE "ingredient_functions" (
                                        "id" bigserial PRIMARY KEY,
                                        "name" varchar(25) NOT NULL,
                                        "user_id" bigint NOT NULL,
                                        "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "ingredients" (
                           "id" bigserial PRIMARY KEY,
                           "name" varchar(25) NOT NULL,
                           "inci" varchar(50) NOT NULL,
                           "hash" varchar(50) NOT NULL,
                           "user_id" bigint NOT NULL,
                           "function_id" bigint,
                           "cost" int4 DEFAULT 0,
                           "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE "ingredient_tags" (
                               "id" bigserial PRIMARY KEY,
                               "name" varchar(25) NOT NULL,
                               "user_id" bigint NOT NULL,
                               "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "ingredient_tag_maps" (
                                       "id" bigserial PRIMARY KEY,
                                       "ingredient_tag_id" bigint NOT NULL,
                                       "ingredient_id" bigint NOT NULL
);

CREATE TABLE "formulas" (
                        "id" bigserial PRIMARY KEY,
                        "name" varchar(25) NOT NULL,
                        "default_amount" float4 NOT NULL,
                        "default_amount_oz" float4 NOT NULL,
                        "description" varchar(500) NOT NULL,
                        "user_id" bigint NOT NULL,
                        "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "formula_tags" (
                        "id" bigserial PRIMARY KEY,
                        "user_id" bigint NOT NULL,
                        "name" varchar(25) NOT NULL

);

CREATE TABLE "formula_tag_maps" (
                           "id" bigserial PRIMARY KEY,
                           "formula_id" bigint NOT NULL,
                           "formula_tag_id" bigint NOT NULL
);

CREATE TABLE "phases" (
                          "id" bigserial PRIMARY KEY,
                          "name" varchar(25) NOT NULL,
                          "description" varchar(500) NOT NULL,
                          "formula_id" bigint NOT NULL,
                          "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "formula_ingredients" (
                          "id" bigserial PRIMARY KEY,
                          "ingredient_id" bigint NOT NULL,
                          "percentage" int NOT NULL,
                          "description" varchar(250),
                          "phase_id" bigint NOT NULL,
                          "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "currencies" (
                              "id" bigserial PRIMARY KEY,
                              "name" varchar(25) NOT NULL,
                              "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                              "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "inventory_items" (
                                "id" bigserial PRIMARY KEY,
                                "user_id" bigint NOT NULL,
                                "ingredient_id" bigint NOT NULL,
                                "amount_in_grams" int NOT NULL,
                                "cost_per_gram" float4 NOT NULL,
                                "currency_id" bigint NOT NULL,
                                "expiry_date" timestamptz NOT NULL,
                                "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE "batches" (
                           "id" bigserial PRIMARY KEY,
                           "user_id" bigint NOT NULL,
                           "name" varchar(25) NOT NULL,
                           "description" varchar(500),
                           "production_date" timestamptz default CURRENT_TIMESTAMP,
                           "expiry_date" timestamptz NOT NULL,
                           "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "batch_items" (
                            "id" bigserial PRIMARY KEY,
                            "amount" int NOT NULL,
                            "inventory_item_id" bigint NOT NULL,
                            "formula_ingredient_id" bigint NOT NULL,
                            "batch_id" bigint NOT NULL,
                            "user_id" bigint NOT NULL,
                            "description" varchar(250) NOT NULL,
                            "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "ingredients" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "ingredients" ADD FOREIGN KEY ("function_id") REFERENCES "ingredient_functions" ("id");
ALTER TABLE "ingredient_tags" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "ingredient_tag_maps" ADD FOREIGN KEY ("ingredient_id") REFERENCES "ingredients" ("id");
ALTER TABLE "ingredient_tag_maps" ADD FOREIGN KEY ("ingredient_tag_id") REFERENCES "ingredient_tags" ("id");

ALTER TABLE "formulas" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "formula_tags" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "formula_tag_maps" ADD FOREIGN KEY ("formula_id") REFERENCES "formulas" ("id");
ALTER TABLE "formula_tag_maps" ADD FOREIGN KEY ("formula_tag_id") REFERENCES "formula_tags" ("id");
ALTER TABLE "phases" ADD FOREIGN KEY ("formula_id") REFERENCES "formulas"("id");
ALTER TABLE "formula_ingredients" ADD FOREIGN KEY ("phase_id") REFERENCES "phases" ("id");
ALTER TABLE "formula_ingredients" ADD FOREIGN KEY ("ingredient_id") REFERENCES "ingredients" ("id");

ALTER TABLE "inventory_items" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "inventory_items" ADD FOREIGN KEY ("ingredient_id") REFERENCES "ingredients" ("id");
ALTER TABLE "inventory_items" ADD FOREIGN KEY ("currency_id") REFERENCES "currencies" ("id");
ALTER TABLE "batch_items" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "batch_items" ADD FOREIGN KEY ("batch_id") REFERENCES "batches" ("id");
ALTER TABLE "batch_items" ADD FOREIGN KEY ("inventory_item_id") REFERENCES "inventory_items" ("id");
ALTER TABLE "batch_items" ADD FOREIGN KEY ("formula_ingredient_id") REFERENCES "formula_ingredients" ("id");

CREATE INDEX ON "ingredients" ("name", "user_id");
CREATE INDEX ON "ingredient_tags" ("id", "user_id");
CREATE INDEX ON "ingredient_tag_maps" ("ingredient_tag_id", "ingredient_id");

CREATE INDEX ON "formulas" ("id", "user_id");
CREATE INDEX ON "formula_tags" ("id", "user_id");
CREATE INDEX ON "formula_tag_maps" ("formula_tag_id", "formula_id");

CREATE INDEX ON "phases" ("id", "formula_id");
CREATE INDEX ON "formula_ingredients" ("id", "phase_id", "ingredient_id");

CREATE INDEX ON "inventory_items" ("id", "user_id", "ingredient_id");
CREATE INDEX ON "batches" ("id", "user_id");
CREATE INDEX ON "batch_items" ("id", "user_id", "inventory_item_id", "formula_ingredient_id", "batch_id");

