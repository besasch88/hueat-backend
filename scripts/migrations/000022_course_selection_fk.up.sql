ALTER TABLE "hueat_course_selection" DROP CONSTRAINT "fk_hueat_course_selection_menu_item";

ALTER TABLE "hueat_course_selection"
ADD CONSTRAINT "fk_hueat_course_selection_menu_item"
FOREIGN KEY ("menu_item_id")
REFERENCES "hueat_menu_item" ("id")
ON DELETE CASCADE
ON UPDATE CASCADE;