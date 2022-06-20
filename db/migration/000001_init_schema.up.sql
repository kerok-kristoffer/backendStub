CREATE TABLE "users" (
                         "id" bigserial PRIMARY KEY,
                         "user_name" varchar UNIQUE NOT NULL,
                         "email" varchar UNIQUE NOT NULL,
                         "full_name" varchar NOT NULL,
                         "hash" varchar(60) NOT NULL,
                         "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE "transfers" (
                        "id" bigserial PRIMARY KEY,
                        "from_user_id" bigint,
                        "to_user_id" bigint,
                        "amount" bigint,
                        "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "entries" (
                        "id" bigserial PRIMARY KEY,
                        "user_id" bigint,
                        "amount" bigint,
                        "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE "transfers" ADD FOREIGN KEY ("from_user_id") REFERENCES "users" ("id");
ALTER TABLE "transfers" ADD FOREIGN KEY ("to_user_id") REFERENCES "users" ("id");

ALTER TABLE "entries" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

CREATE INDEX ON "users" ("full_name");
CREATE INDEX ON "users" ("id");
CREATE INDEX ON "users" ("email");
