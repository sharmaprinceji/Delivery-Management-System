package orderroute

import (
	"net/http"

	"github.com/sharmaprinceji/delivery-management-system/internal/http/handlers/order"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
)


func RegisterOrderRoutes(router *http.ServeMux, storage storage.Storage) {
	router.HandleFunc("POST /api/order", order.CreateOrder(storage))  
	router.HandleFunc("POST /api/orders", order.CreateBulkOrders(storage))
	router.HandleFunc("GET /api/allocate", order.ManualAllocation(storage))  
	router.HandleFunc("GET /api/agent-summary", order.GetAgentSummary(storage))
    router.HandleFunc("GET /api/system-summary", order.GetSystemSummary(storage)) 
}