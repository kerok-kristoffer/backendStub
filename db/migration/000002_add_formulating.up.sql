CREATE TABLE "ingredients" (
                           "id" bigserial PRIMARY KEY,
                           "name" varchar(25) NOT NULL,
                           "inci" varchar(50) NOT NULL,
                           "hash" varchar(50) NOT NULL,
                           "user_id" bigint,
                           "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE "ingredient_tags" (
                               "id" bigserial PRIMARY KEY,
                               "user_id" bigint,
                               "ingredient_id" bigint,
                               "label" varchar(25),
                               "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE "recipes" (
                        "id" bigserial PRIMARY KEY,
                        "name" varchar(25),
                        "default_amount" int,
                        "description" varchar(250),
                        "user_id" bigint,
                        "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "phases" (
                          "id" bigserial PRIMARY KEY,
                          "name" varchar(25),
                          "description" varchar(250),
                          "recipe_id" bigint,
                          "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "recipe_ingredients" (
                          "id" bigserial PRIMARY KEY,
                          "ingredient_id" bigint,
                          "percentage" int,
                          "description" varchar(250),
                          "phase_id" bigint,
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
                                "user_id" bigint,
                                "ingredient_id" bigint,
                                "amount_in_grams" int,
                                "cost_per_gram" float8,
                                "currency_id" bigint,
                                "expiry_date" timestamptz,
                                "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE "batches" (
                           "id" bigserial PRIMARY KEY,
                           "user_id" bigint,
                           "name" varchar(25),
                           "description" varchar(250),
                           "production_date" timestamptz default CURRENT_TIMESTAMP,
                           "expiry_date" timestamptz,
                           "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "batch_items" (
                            "id" bigserial PRIMARY KEY,
                            "amount" int,
                            "inventory_item_id" bigint,
                            "recipe_ingredient_id" bigint,
                            "batch_id" bigint,
                            "user_id" bigint,
                            "description" varchar(250),
                            "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "ingredients" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "ingredient_tags" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "ingredient_tags" ADD FOREIGN KEY ("ingredient_id") REFERENCES "ingredients" ("id");

ALTER TABLE "recipes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "phases" ADD FOREIGN KEY ("recipe_id") REFERENCES "recipes"("id");
ALTER TABLE "recipe_ingredients" ADD FOREIGN KEY ("phase_id") REFERENCES "phases" ("id");
ALTER TABLE "recipe_ingredients" ADD FOREIGN KEY ("ingredient_id") REFERENCES "ingredients" ("id");

ALTER TABLE "inventory_items" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "inventory_items" ADD FOREIGN KEY ("ingredient_id") REFERENCES "ingredients" ("id");
ALTER TABLE "inventory_items" ADD FOREIGN KEY ("currency_id") REFERENCES "currencies" ("id");
ALTER TABLE "batch_items" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "batch_items" ADD FOREIGN KEY ("batch_id") REFERENCES "batches" ("id");
ALTER TABLE "batch_items" ADD FOREIGN KEY ("inventory_item_id") REFERENCES "inventory_items" ("id");
ALTER TABLE "batch_items" ADD FOREIGN KEY ("recipe_ingredient_id") REFERENCES "recipe_ingredients" ("id");

CREATE INDEX ON "ingredients" ("name");
CREATE INDEX ON "ingredients" ("user_id");

CREATE INDEX ON "ingredient_tags" ("id");
CREATE INDEX ON "ingredient_tags" ("user_id");
CREATE INDEX ON "ingredient_tags" ("ingredient_id");

CREATE INDEX ON "recipes" ("id");
CREATE INDEX ON "recipes" ("user_id");

CREATE INDEX ON "phases" ("id");
CREATE INDEX ON "phases" ("recipe_id");

CREATE INDEX ON "recipe_ingredients" ("id");
CREATE INDEX ON "recipe_ingredients" ("phase_id");
CREATE INDEX ON "recipe_ingredients" ("ingredient_id");

CREATE INDEX ON "inventory_items" ("id");
CREATE INDEX ON "inventory_items" ("user_id");
CREATE INDEX ON "inventory_items" ("ingredient_id");

CREATE INDEX ON "batches" ("id");
CREATE INDEX ON "batches" ("user_id");

CREATE INDEX ON "batch_items" ("id");
CREATE INDEX ON "batch_items" ("user_id");
CREATE INDEX ON "batch_items" ("inventory_item_id");
CREATE INDEX ON "batch_items" ("recipe_ingredient_id");
CREATE INDEX ON "batch_items" ("batch_id");

