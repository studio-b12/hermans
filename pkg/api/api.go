package api

import (
	"net/http"

	"github.com/zekrotja/hermans/pkg/model"
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
	mux.HandleFunc("GET /api/lists/{id}", t.handleGetOrderList)
	mux.HandleFunc("DELETE /api/lists/{id}", t.handleDeleteOrderList)

	mux.HandleFunc("POST /api/lists/{id}/orders", t.handleCreateOrder)

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

func (t *API) handleGetOrderList(w http.ResponseWriter, r *http.Request) {
	orderListId := r.PathValue("id")

	list, err := t.ctl.GetOrders(orderListId)
	if err != nil {
		respondErr(w, err)
		return
	}

	respondJson(w, http.StatusOK, list)
}

func (t *API) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	orderListId := r.PathValue("id")

	order, err := readJsonBody[model.Order](r)
	if err != nil {
		respondErr(w, err)
		return
	}

	newOrder, err := t.ctl.CreateOrder(orderListId, &order)
	if err != nil {
		respondErr(w, err)
		return
	}

	respondJson(w, http.StatusCreated, newOrder)
}

func (t *API) handleDeleteOrderList(w http.ResponseWriter, r *http.Request) {
	orderListId := r.PathValue("id")

	err := t.ctl.DeleteOrderList(orderListId)
	if err != nil {
		respondErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
