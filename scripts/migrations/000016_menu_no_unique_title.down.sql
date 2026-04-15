ALTER TABLE "hueat_menu_category" ADD CONSTRAINT "idx_hueat_menu_category_title" UNIQUE ("title");
ALTER TABLE "hueat_menu_item" ADD CONSTRAINT "idx_hueat_menu_item_title" UNIQUE ("title");
ALTER TABLE "hueat_menu_option" ADD CONSTRAINT "idx_hueat_menu_option_title" UNIQUE ("title");