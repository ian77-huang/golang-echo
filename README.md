# Go Echo Web 專案開發指南

本專案是一個基於 **Go (Echo v5)** 框架開發的 Web 應用程式，整合了 **SQLite** 資料庫、**golang-migrate** 遷移控制、自訂 HTML 模板渲染引擎（支援多層級 Layout）、多語系（i18n）支援，以及表單驗證機制。

---

## 🛠 核心技術棧 (Tech Stack)

*   **後端框架**：[Echo v5](https://github.com/labstack/echo) (高效能、極簡的 Go Web 框架)
*   **資料庫**：SQLite (搭配 Go 標準庫 `database/sql` 及 `github.com/mattn/go-sqlite3`)
*   **資料庫遷移**：[golang-migrate](https://github.com/golang-migrate/migrate) (版本控制資料庫結構)
*   **視圖引擎**：自訂 HTML 模板渲染（支援多層級 Layout，如 `base -> frontend -> users -> page`）
*   **多語系支援 (i18n)**：基於 TOML 翻譯檔的多語系中介軟體 (支援繁體中文 `zh-TW` 與英文 `en`)
*   **資料驗證**：[go-playground/validator/v10](https://github.com/go-playground/validator) (強大的結構體/表單驗證器)
*   **熱重載**：[Air](https://github.com/air-verse/air) (開發環境程式碼變更自動重啟)

---

## 📁 專案目錄結構 (Project Structure)

```text
.
├── cmd/
│   ├── migrate/            # 資料庫遷移工具執行入口 (go run cmd/migrate/main.go)
│   └── server/             # Web 伺服器啟動入口 (go run cmd/server/main.go)
├── internal/
│   ├── config/             # 設定檔載入、i18n 與模板引擎初始化
│   ├── handler/            # 控制器/處理器 (HTTP Handlers)
│   │   ├── api/            # API 端點處理器 (例如：語系切換、註冊 API)
│   │   └── frontend/       # 前端頁面渲染處理器 (例如：登入、註冊頁面)
│   ├── locales/            # 嵌入式多語系翻譯檔 (TOML)
│   ├── router/             # 路由定義與分組管理
│   └── views/              # HTML 模板檔案 (支援巢狀 Layout 與元件)
├── pkg/                    # 專案內置的共享模組 (i18n, renderer, validator, cast)
├── migrations/             # SQL 資料庫遷移檔案 (.up.sql / .down.sql)
├── databases/              # 本機 SQLite 資料庫儲存目錄 (已設定於 .gitignore 排除)
├── .air.toml               # Air 熱重載設定檔
├── .env                    # 環境變數設定檔
├── go.mod                  # Go 模組定義檔
└── README.md               # 本開發說明檔
```

---

## 🚀 快速開始 (Quick Start)

### 1. 安裝環境與依賴
請確保本機已安裝 Go 1.21+。複製專案後執行以下指令下載所有依賴：
```bash
go mod download
```

### 2. 設定環境變數
專案根目錄下需包含 `.env` 檔案（若無可自行建立）：
```env
USERS_ACCOUNT_MIN_LENGTH=6
USERS_PASSWORD_MIN_LENGTH=8
```

### 3. 執行資料庫遷移
本專案的 SQLite 資料庫會在第一次執行遷移或啟動時自動建立於 `databases/main.db`。執行以下指令套用最新的資料庫結構：
```bash
go run cmd/migrate/main.go
```

### 4. 啟動開發伺服器
*   **一般啟動**：
    ```bash
    go run cmd/server/main.go
    ```
*   **使用 Air 進行熱重載開發**（需先安裝 Air：`go install github.com/air-verse/air@latest`）：
    ```bash
    air
    ```
    伺服器預設會運行在 `http://localhost:1323`。

---

## 🗄 資料庫遷移 (Migrations)

當你需要修改資料庫結構（例如新增資料表或欄位）時，請使用 `golang-migrate` 進行版本控制。

### 建立新的遷移檔案
若你安裝了 `migrate` CLI 工具，可執行：
```bash
migrate create -ext sql -dir migrations -seq <migration_name>
```
*範例：*
```bash
migrate create -ext sql -dir migrations -seq create_messages_table
```
執行後會在 `migrations/` 目錄下產生兩個檔案：
*   `00000X_<migration_name>.up.sql`：寫入結構變更的 SQL (例如 `CREATE TABLE...`)
*   `00000X_<migration_name>.down.sql`：寫入還原該變更的 SQL (例如 `DROP TABLE...`)

### 執行遷移
編輯好 SQL 檔後，執行以下 Go 程式以自動更新資料庫：
```bash
go run cmd/migrate/main.go
```

---

## 🌐 多語系支援 (i18n)

本專案支援多語系，翻譯檔案存放在 `internal/locales/` 下，採用 TOML 格式：
*   語系切換 API：`POST /api/lang`
*   模板中的多語系使用：可在 HTML 模板中直接調用翻譯函式。

---

## 📝 開發規範

1.  **資料庫排除**：請勿提交任何本機的 SQLite 庫檔案 (`databases/*.db`) 至 Git。
2.  **新增頁面**：
    *   在 `internal/views/` 下建立對應的 HTML 模板。
    *   在 `internal/handler/frontend/` 建立對應的渲染處理器。
    *   在 `internal/router/routing/frontend.go` 配置路由。
3.  **新增 API**：
    *   在 `internal/handler/api/` 建立 API 邏輯。
    *   在 `internal/router/routing/api.go` 配置路由。
