-- 刪除索引
DROP INDEX IF EXISTS "session_sessionToken_key";
DROP INDEX IF EXISTS "session_key_id_status";

-- 刪除資料表
DROP TABLE IF EXISTS "session";