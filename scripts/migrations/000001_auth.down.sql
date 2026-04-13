ALTER TABLE "hueat_auth_session" DROP CONSTRAINT "idx_hueat_auth_session_refresh_token";

DROP TABLE IF EXISTS "hueat_auth_session";