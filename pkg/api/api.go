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

	mux.HandleFunc("GET /api/items", multiHandler(t.setCORSHeader, t.handleGetStoreItems))
	mux.HandleFunc("POST /api/lists", multiHandler(t.setCORSHeader, t.handleCreateOrderList))
	mux.HandleFunc("GET /api/lists/{id}", multiHandler(t.setCORSHeader, t.handleGetOrderList))
	mux.HandleFunc("DELETE /api/lists/{id}", multiHandler(t.setCORSHeader, t.handleDeleteOrderList))
	mux.HandleFunc("POST /api/lists/{id}/orders", multiHandler(t.setCORSHeader, t.handleCreateOrder))

	return &t
}

func (t *API) Start() error {
	return t.server.ListenAndServe()
}

func (t *API) handleGetStoreItems(w http.ResponseWriter, r *http.Request) {
	data, err := t.ctl.GetScrapedData()
	if err != nil {
		respondErr(w, err)
		return
	}

	respondJson(w, http.StatusOK, data)
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

func (t *API) setCORSHeader(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
}
