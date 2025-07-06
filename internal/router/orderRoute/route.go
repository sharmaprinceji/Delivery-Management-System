package orderroute

import (
	"net/http"

	"github.com/sharmaprinceji/delivery-management-system/internal/http/handlers/order"
	"github.com/sharmaprinceji/delivery-management-system/internal/router"
)


func StudentRouter() *http.ServeMux {
	route ,storage:=  router.StudentRoute()

	route.HandleFunc("GET /api/assignments", order.GetAssignments(storage))
	route.HandleFunc("GET /api/checkin", order.CheckInAgent(storage))
	

	return route;
}