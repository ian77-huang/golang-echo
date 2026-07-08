這是一份幫你整理好、並將名稱改為動態變數的 `README.md`：

````markdown
# Database Migration 使用指南

本專案使用 `golang-migrate` 進行資料庫版本控制。以下為建立新遷移檔案與執行遷移的指令說明。

---

## 1. 建立新的 Migration 檔案

當你需要修改資料庫結構（例如新增資料表或欄位）時，請使用以下指令產生全新的 `up` 和 `down` SQL 檔案。

### 指令

將 `<migration_name>` 替換為你的功能名稱（例如：`create_users_table` 或 `add_phone_to_users`）：

```bash
migrate create -ext sql -dir migrations -seq <migration_name>
```
````

### 範例

如果你要建立訊息資料表，請執行：

```bash
migrate create -ext sql -dir migrations -seq create_messages_table

```

執行後，系統會在 `migrations/` 資料夾下自動生成兩個檔案：

- `00000X_create_messages_table.up.sql`（填入結構變更的 SQL）
- `00000X_create_messages_table.down.sql`（填入還原變更的 SQL）

---

## 2. 執行 Migration (套用變更)

檔案編寫完成後，執行以下 Go 程式來自動檢查並將變更套用到 SQLite 資料庫中：

```bash
go run cmd/migrate/main.go

```

> **附註**：此指令會自動建立 `./databases/` 資料夾（若不存在），並自動略過已執行過的遷移版本，重複執行不會造成資料庫衝突。

```

```
