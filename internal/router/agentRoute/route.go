package agentroute

import (
	"net/http"

	"github.com/sharmaprinceji/delivery-management-system/internal/http/handlers/agent"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
)


func RegisterAgentRoutes(router *http.ServeMux, storage storage.Storage) {
	router.HandleFunc("POST /api/warehouse", agent.CreateWareHouse(storage))//agent
	router.HandleFunc("POST /api/checkin", agent.CheckedInAgents(storage))  //order
	router.HandleFunc("GET /api/checkin", agent.CheckInAgent(storage))
	router.HandleFunc("GET /api/assignments", agent.GetAssignments(storage))
}