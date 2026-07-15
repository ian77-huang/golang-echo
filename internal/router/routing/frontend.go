package routing

import (
	"github.com/ian77-huang/golang-echo/internal/handler/frontend"
)

func (h *Routing) Frontend() {
	frontend.New(&frontend.FrontendParameter{DB: h.DB, Echo: h.Echo})
}
