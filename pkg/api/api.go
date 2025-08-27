package api

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/zekrotja/hermans/pkg/model"
)

type API struct {
	ctl    Controller
	server *http.Server
}

func New(ctl Controller, addr string) *API {
	apiMux := http.NewServeMux()
	t := API{ctl: ctl}

	apiMux.HandleFunc("GET /items", t.handleGetStoreItems)
	apiMux.HandleFunc("POST /lists", t.handleCreateOrderList)
	apiMux.HandleFunc("GET /lists/{id}", t.handleGetOrderList)
	apiMux.HandleFunc("DELETE /lists/{id}", t.handleDeleteOrderList)
	apiMux.HandleFunc("POST /lists/{id}/orders", t.handleCreateOrder)
	apiMux.HandleFunc("GET /lists/{listId}/orders/{orderId}", t.handleGetOrder)
	apiMux.HandleFunc("PUT /lists/{listId}/orders/{orderId}", t.handleUpdateOrder)
	apiMux.HandleFunc("DELETE /lists/{listId}/orders/{orderId}", t.handleDeleteOrder)
	apiMux.HandleFunc("GET /dev/clearall", t.handleClearAll)

	mainMux := http.NewServeMux()
	mainMux.Handle("/api/", http.StripPrefix("/api", apiMux))
	mainMux.Handle("/", http.FileServer(http.Dir("webapp")))

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		mainMux.ServeHTTP(w, r)
	})

	t.server = &http.Server{Addr: addr, Handler: finalHandler}
	return &t
}

func (t *API) Start() error {
	return t.server.ListenAndServe()
}

func (t *API) handleClearAll(w http.ResponseWriter, r *http.Request) {
	if err := t.ctl.ClearAllData(); err != nil {
		respondErr(w, err)
		return
	}
	w.Write([]byte("Alle Daten wurden gelöscht."))
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
	payload, err := readJsonBody[struct {
		Deadline *time.Time `json:"deadline"`
	}](r)
	if err != nil && err != io.EOF {
		slog.Warn("failed reading optional deadline body", "err", err)
	}

	list, err := t.ctl.CreateOrderList(payload.Deadline)
	if err != nil {
		respondErr(w, err)
		return
	}
	respondJson(w, http.StatusCreated, list)
}

func (t *API) handleGetOrderList(w http.ResponseWriter, r *http.Request) {
	orderListId := r.PathValue("id")
	list, err := t.ctl.GetOrderList(orderListId)
	if err != nil {
		respondErr(w, err)
		return
	}
	orders, err := t.ctl.GetOrders(orderListId)
	if err != nil {
		respondErr(w, err)
		return
	}
	respondJson(w, http.StatusOK, map[string]interface{}{
		"id":       list.Id,
		"created":  list.Created,
		"deadline": list.Deadline,
		"orders":   orders,
	})
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
	respondJson(w, http.StatusCreated, map[string]interface{}{
		"id":          newOrder.Id,
		"created":     newOrder.Created,
		"creator":     newOrder.Creator,
		"store_items": newOrder.StoreItems,
		"drink":       newOrder.Drink,
		"editKey":     newOrder.EditKey,
	})
}

func (t *API) handleUpdateOrder(w http.ResponseWriter, r *http.Request) {
	listId := r.PathValue("listId")
	orderId := r.PathValue("orderId")

	payload, err := readJsonBody[struct {
		model.Order
		EditKey string `json:"editKey"`
	}](r)
	if err != nil {
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

	payload, err := readJsonBody[struct {
		EditKey string `json:"editKey"`
	}](r)
	if err != nil {
		respondErr(w, err)
		return
	}

	if err := t.ctl.DeleteOrder(listId, orderId, payload.EditKey); err != nil {
		respondErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
