package orderroute

import (
	"github.com/gorilla/mux"
	"github.com/sharmaprinceji/delivery-management-system/internal/http/handlers/order"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
)

func RegisterOrderRoutes(router *mux.Router, storage storage.Storage) {
	router.HandleFunc("/api/order", order.CreateOrder(storage)).Methods("POST")
	router.HandleFunc("/api/orders/bulk", order.CreateBulkOrders(storage)).Methods("POST")
	router.HandleFunc("/api/allocate", order.ManualAllocation(storage)).Methods("GET")
	router.HandleFunc("/api/agent-summary", order.GetAgentSummary(storage)).Methods("GET")
	router.HandleFunc("/api/system-summary", order.GetSystemSummary(storage)).Methods("GET")
}
