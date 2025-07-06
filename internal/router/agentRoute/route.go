package agentroute

import (
	"net/http"

	"github.com/sharmaprinceji/delivery-management-system/internal/http/handlers/agent"
	"github.com/sharmaprinceji/delivery-management-system/internal/router"
)


func AgentRouter() *http.ServeMux {
	route,storage :=  router.StudentRoute()

	route.HandleFunc("GET /api/student/{id}", agent.GetById(storage))
	route.HandleFunc("POST /api/student", agent.Create(storage))

	return route;
}