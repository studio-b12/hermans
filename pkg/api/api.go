package api

import (
	"net/http"
)

type API struct {
	ctl    Controller
	server *http.Server
}

func New(ctl Controller, addr string) *API {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	t := API{
		ctl:    ctl,
		server: server,
	}

	mux.HandleFunc("POST /api/lists", t.handleCreateOrderList)

	return &t
}

func (t *API) Start() error {
	return t.server.ListenAndServe()
}

func (t *API) handleCreateOrderList(w http.ResponseWriter, r *http.Request) {
	list, err := t.ctl.CreateOrderList()
	if err != nil {
		respondErr(w, err)
		return
	}

	respondJson(w, http.StatusCreated, list)
}
