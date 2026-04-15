ALTER TABLE "hueat_menu_category" ADD COLUMN "title_display" VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE "hueat_menu_item" ADD COLUMN "title_display" VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE "hueat_menu_option" ADD COLUMN "title_display" VARCHAR(255) NOT NULL DEFAULT '';

UPDATE hueat_menu_category SET title_display=title;
UPDATE hueat_menu_item SET title_display=title;
UPDATE hueat_menu_option SET title_display=title;