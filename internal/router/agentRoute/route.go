package agentroute

import (
	"net/http"

	"github.com/sharmaprinceji/delivery-management-system/internal/http/handlers/agent"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
)


func RegisterAgentRoutes(router *http.ServeMux, storage storage.Storage) {
	router.HandleFunc("POST /api/warehouse", agent.CreateWareHouse(storage))
	router.HandleFunc("POST /api/checkin", agent.CheckedInAgents(storage))  
	router.HandleFunc("GET /api/checkin", agent.CheckInAgent(storage))
	router.HandleFunc("GET /api/assignments", agent.GetAssignments(storage))

}