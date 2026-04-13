CREATE TABLE "hueat_menu_item" (
    "id" VARCHAR(36) PRIMARY KEY NOT NULL,
    "menu_category_id" VARCHAR(36) NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "position" BIGINT NOT NULL,
    "active" BOOLEAN NOT NULL,
    "price" BIGINT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);

ALTER TABLE "hueat_menu_item"
ADD CONSTRAINT "fk_hueat_menu_item_menu_category"
FOREIGN KEY ("menu_category_id")
REFERENCES "hueat_menu_category" ("id")
ON DELETE CASCADE
ON UPDATE CASCADE;

ALTER TABLE "hueat_menu_item" ADD CONSTRAINT "idx_hueat_menu_item_title" UNIQUE ("title");