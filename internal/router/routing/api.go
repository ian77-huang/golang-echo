package routing

import (
	apiHandler "github.com/ian77-huang/golang-echo/internal/handler/api"
)

func (h *Routing) Api() {
	apiHandler.New(&apiHandler.ApiParameter{DB: h.DB, Echo: h.Echo})
}
