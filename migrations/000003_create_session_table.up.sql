-- 建立 Session 資料表 (相容 SQLite 語法)
CREATE TABLE IF NOT EXISTS "session" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "expiresAt" DATETIME NOT NULL,
    "createdAt" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" DATETIME NOT NULL,
    "status" INTEGER NOT NULL DEFAULT 0,          -- status 0: normal, 1: logout, 99: delete
    "countUpdate" INTEGER NOT NULL DEFAULT 0,
    -- 定義外鍵並設定連動刪除
    FOREIGN KEY ("userId") REFERENCES "User" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);

-- 建立一般索引
CREATE INDEX IF NOT EXISTS "session_key_id_status" 
    ON "session" ("id" ASC, "status" ASC);

-- 建立唯一索引
CREATE UNIQUE INDEX IF NOT EXISTS "session_sessionToken_key" 
    ON "session" ("id" ASC);