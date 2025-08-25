package api

import (
	"net/http"
	"strings"

	"github.com/zekrotja/hermans/pkg/model"
)

type API struct {
	ctl    Controller
	server *http.Server
}

func New(ctl Controller, addr string) *API {
	mux := http.NewServeMux()
	t := API{ctl: ctl}

	mux.HandleFunc("GET /items", t.handleGetStoreItems)
	mux.HandleFunc("POST /lists", t.handleCreateOrderList)
	mux.HandleFunc("GET /lists/{id}", t.handleGetOrderList)
	mux.HandleFunc("DELETE /lists/{id}", t.handleDeleteOrderList)
	mux.HandleFunc("POST /lists/{id}/orders", t.handleCreateOrder)
	mux.HandleFunc("PUT /lists/{listId}/orders/{orderId}", t.handleUpdateOrder)
	mux.HandleFunc("DELETE /lists/{listId}/orders/{orderId}", t.handleDeleteOrder)
	mux.HandleFunc("GET /lists/{listId}/orders/{orderId}", t.handleGetOrder)

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.StripPrefix("/api", mux).ServeHTTP(w, r)
			return
		}
		http.FileServer(http.Dir("webapp")).ServeHTTP(w, r)
	})

	t.server = &http.Server{Addr: addr, Handler: finalHandler}
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
		http.Error(w, "Liste nicht gefunden", http.StatusNotFound)
		return
	}
	respondJson(w, http.StatusOK, list)
}

func (t *API) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	orderListId := r.PathValue("id")
	var order model.Order
	if err := readJsonBody(r, &order); err != nil {
		respondErr(w, err)
		return
	}
	newOrder, err := t.ctl.CreateOrder(orderListId, &order)
	if err != nil {
		respondErr(w, err)
		return
	}
	respondJson(w, http.StatusCreated, map[string]interface{}{
		"id":         newOrder.Id,
		"created":    newOrder.Created,
		"creator":    newOrder.Creator,
		"store_item": newOrder.StoreItem,
		"drink":      newOrder.Drink,
		"editKey":    newOrder.EditKey,
	})
}

func (t *API) handleDeleteOrderList(w http.ResponseWriter, r *http.Request) {
	orderListId := r.PathValue("id")
	if err := t.ctl.DeleteOrderList(orderListId); err != nil {
		respondErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (t *API) handleUpdateOrder(w http.ResponseWriter, r *http.Request) {
	listId := r.PathValue("listId")
	orderId := r.PathValue("orderId")
	var payload struct {
		model.Order
		EditKey string `json:"editKey"`
	}
	if err := readJsonBody(r, &payload); err != nil {
		respondErr(w, err)
		return
	}
	updatedOrder, err := t.ctl.UpdateOrder(listId, orderId, payload.EditKey, &payload.Order)
	if err != nil {
		respondErr(w, err)
		return
	}
	respondJson(w, http.StatusOK, updatedOrder)
}

func (t *API) handleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	listId := r.PathValue("listId")
	orderId := r.PathValue("orderId")
	var payload struct {
		EditKey string `json:"editKey"`
	}
	if err := readJsonBody(r, &payload); err != nil {
		respondErr(w, err)
		return
	}
	if err := t.ctl.DeleteOrder(listId, orderId, payload.EditKey); err != nil {
		respondErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (t *API) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	listId := r.PathValue("listId")
	orderId := r.PathValue("orderId")
	order, err := t.ctl.GetOrder(listId, orderId)
	if err != nil {
		respondErr(w, err)
		return
	}
	respondJson(w, http.StatusOK, order)
}
