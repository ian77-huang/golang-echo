package routing

import (
	"github.com/ian77-huang/golang-echo/internal/handler/admin"
)

func (h *Routing) Admin() {
	admin.New(&admin.AdminParameter{DB: h.DB, Echo: h.Echo})
}
