ALTER TABLE "hueat_course_selection" DROP CONSTRAINT "fk_hueat_course_selection_menu_option";
ALTER TABLE "hueat_course_selection" DROP CONSTRAINT "fk_hueat_course_selection_menu_item";
ALTER TABLE "hueat_course_selection" DROP CONSTRAINT "fk_hueat_course_selection_course";
DROP TABLE IF EXISTS "hueat_course_selection";

ALTER TABLE "hueat_course" DROP CONSTRAINT "fk_hueat_course_order";
DROP TABLE IF EXISTS "hueat_course";


ALTER TABLE "hueat_order" DROP CONSTRAINT "fk_hueat_order_table";
DROP TABLE IF EXISTS "hueat_order";
