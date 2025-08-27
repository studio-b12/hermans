package api

import (
	"net/http"
	"time"

	"github.com/zekrotja/hermans/pkg/model"
)

type API struct {
	ctl    Controller
	server *http.Server
}

func New(ctl Controller, addr string) *API {
	mux := http.NewServeMux()
	t := API{
		ctl:    ctl,
		server: &http.Server{Addr: addr, Handler: mux},
	}

	mux.Handle("/", http.FileServer(http.Dir("webapp")))

	// API Routen
	mux.HandleFunc("OPTIONS /", t.handleOptions)
	mux.HandleFunc("GET /api/items", multiHandler(t.setCORSHeader, t.handleGetStoreItems))
	mux.HandleFunc("POST /api/lists", multiHandler(t.setCORSHeader, t.handleCreateOrderList))
	mux.HandleFunc("GET /api/lists/{id}", multiHandler(t.setCORSHeader, t.handleGetOrderList))
	mux.HandleFunc("DELETE /api/lists/{id}", multiHandler(t.setCORSHeader, t.handleDeleteOrderList))
	mux.HandleFunc("POST /api/lists/{id}/orders", multiHandler(t.setCORSHeader, t.handleCreateOrder))
	mux.HandleFunc("PUT /api/lists/{listId}/orders/{orderId}", multiHandler(t.setCORSHeader, t.handleUpdateOrder))
	mux.HandleFunc("DELETE /api/lists/{listId}/orders/{orderId}", multiHandler(t.setCORSHeader, t.handleDeleteOrder))
	mux.HandleFunc("GET /api/lists/{listId}/orders/{orderId}", multiHandler(t.setCORSHeader, t.handleGetOrder))
	mux.HandleFunc("POST /api/feedback", multiHandler(t.setCORSHeader, t.handleCreateFeedback))
	mux.HandleFunc("GET /api/dev/clearall", multiHandler(t.setCORSHeader, t.handleClearAll))

	return &t
}

func (t *API) Start() error {
	return t.server.ListenAndServe()
}

func (t *API) handleOptions(w http.ResponseWriter, r *http.Request) {
	t.setCORSHeader(w, r)
	w.WriteHeader(http.StatusNoContent)
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
	if err != nil && err.Error() != "EOF" {
		respondErr(w, err)
		return
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
	if err := t.ctl.DeleteOrderList(orderListId); err != nil {
		respondErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (t *API) handleCreateFeedback(w http.ResponseWriter, r *http.Request) {
	feedback, err := readJsonBody[model.Feedback](r)
	if err != nil {
		respondErr(w, err)
		return
	}
	newFeedback, err := t.ctl.CreateFeedback(&feedback)
	if err != nil {
		respondErr(w, err)
		return
	}
	respondJson(w, http.StatusCreated, newFeedback)
}

func (t *API) setCORSHeader(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-control-allow-headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}
