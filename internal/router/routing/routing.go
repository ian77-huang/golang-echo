package routing

func New(rp *RoutingParameter) {
	h := &Routing{DB: rp.DB, Echo: rp.Echo}

	h.Frontend()
	h.Api()
	h.Admin()
}
