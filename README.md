# Go Echo Web 專案

以 Go 與 Echo v5 建置的 Web 應用程式練習專案，提供 SQLite 資料庫遷移、HTML 模板、多語系與使用者註冊、登入／登出、session 驗證等基礎功能。

## 目前進度

| 功能                 | 狀態   | 說明                                                |
| -------------------- | ------ | --------------------------------------------------- |
| Echo Web 伺服器      | 已完成 | 預設監聽 `http://localhost:1323`。                  |
| HTML 模板            | 已完成 | 包含首頁、登入頁與註冊頁，以及可巢狀使用的 layout。 |
| 多語系               | 已完成 | 支援繁體中文（`zh-TW`）與英文（`en`）。             |
| SQLite 與 migrations | 已完成 | 已建立 users、messages、session 三組遷移。          |
| 使用者註冊 API       | 已完成 | 建立帳號、以 Argon2 雜湊密碼並建立 session。        |
| Session 驗證         | 已完成 | 透過 cookie 與 JWT 處理 session 的建立與刷新。      |
| 登入 API             | 已完成 | 驗證帳號密碼，成功後建立 session 並寫入 cookie。    |
| 登出 API             | 已完成 | 清除 session 並刪除 cookie，重新導向至登入頁。      |

## 技術棧

- Go 1.26.1
- Echo v5
- GORM + SQLite
- golang-migrate
- go-i18n
- go-playground/validator
- Argon2、JWT 與 cookie session

## 專案結構

```text
.
├── cmd/
│   ├── migrate/       # 資料庫遷移入口
│   └── server/        # Web 伺服器入口
├── internal/
│   ├── config/        # 環境設定、i18n、模板設定
│   ├── handler/       # 頁面與 API handlers
│   ├── locales/       # TOML 翻譯檔
│   ├── models/        # users、session 等資料模型
│   ├── router/        # 路由設定
│   └── views/         # HTML 模板
├── migrations/        # SQLite migration SQL
├── pkg/               # auth、database、i18n、renderer 等共用套件
├── .env.example       # 可提交的環境變數範本
└── .env               # 本機設定（不納入 Git）
```

## 快速開始

### 1. 安裝依賴

請先安裝 Go 1.26.1，接著下載模組：

```bash
go mod download
```

### 2. 建立 `.env`

從範本建立本機設定檔：

```bash
cp .env.example .env
```

請將 `SECRET_KEY` 換成自己的高強度隨機字串。`.env` 已由 Git 忽略，請勿提交其中的機密值。

| 變數                        | 必填 | 用途                                    | 開發環境範例                        |
| --------------------------- | ---- | --------------------------------------- | ----------------------------------- |
| `SECRET_KEY`                | 是   | JWT 簽章與 session 驗證金鑰；不可為空。 | `replace-with-a-long-random-secret` |
| `DATABASE_PATH`             | 是   | SQLite 資料庫檔案路徑。                 | `databases/main.db`                 |
| `USERS_ACCOUNT_MIN_LENGTH`  | 是   | 註冊帳號最小長度。                      | `6`                                 |
| `USERS_PASSWORD_MIN_LENGTH` | 是   | 註冊密碼最小長度。                      | `8`                                 |

### 3. 套用資料庫遷移

```bash
go run cmd/migrate/main.go
```

此指令會建立 `DATABASE_PATH` 的父目錄，並套用 `migrations/` 中尚未執行的 migration。

### 4. 啟動伺服器

```bash
go run cmd/server/main.go
```

開發時也可使用 Air 熱重載：

```bash
go install github.com/air-verse/air@latest
air
```

## 路由與 API

| 方法   | 路徑                 | 說明                     |
| ------ | -------------------- | ------------------------ |
| `GET`  | `/`                  | 首頁                     |
| `GET`  | `/user`              | 使用者首頁               |
| `GET`  | `/user/login`        | 登入頁面                 |
| `GET`  | `/user/register`     | 註冊頁面                 |
| `GET`  | `/user/logout`       | 登出並重新導向至登入頁   |
| `GET`  | `/api/ping`          | 健康檢查                 |
| `POST` | `/api/lang`          | 切換語系                 |
| `POST` | `/api/user/register` | 註冊使用者並建立 session |
| `POST` | `/api/user/login`    | 登入並建立 session       |

註冊 API（`POST /api/user/register`）請求內容：

```json
{
  "account": "example_user",
  "password": "your-secure-password",
  "confirmPassword": "your-secure-password"
}
```

登入 API（`POST /api/user/login`）請求內容：

```json
{
  "account": "example_user",
  "password": "your-secure-password"
}
```

## 執行測試

```bash
go test ./...
```

## 開發注意事項

- `.env`、本機資料庫與日誌都不應提交至 Git。
- 變更資料表時，請新增一組 `.up.sql` 與 `.down.sql` migration，再執行遷移指令。
- 新增前端頁面時，請同步建立模板、handler 與路由；新增 API 時，請建立 handler 並在 `internal/router/routing/api.go` 註冊。
