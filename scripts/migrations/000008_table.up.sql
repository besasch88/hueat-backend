CREATE TABLE "hueat_table" (
    "id" VARCHAR(36) PRIMARY KEY NOT NULL,
    "user_id" VARCHAR(36) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "close" BOOLEAN NOT NULL,
    "payment_method" VARCHAR(255),
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);