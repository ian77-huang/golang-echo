package api

type ApiHandler struct {
}

type ChangeLangRequest struct {
	Code string `json:"code" validate:"required,oneof=zh-TW en"`
}
