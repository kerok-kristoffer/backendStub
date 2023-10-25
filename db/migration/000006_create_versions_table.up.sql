CREATE TABLE "versions" (
                           "id" bigserial PRIMARY KEY,
                           "number" float4,
                           "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX ON "versions" ("number")