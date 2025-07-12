package agentRoute

import (
	"github.com/gorilla/mux"
	"github.com/sharmaprinceji/delivery-management-system/internal/http/handlers/agent"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
)

func RegisterAgentRoutes(router *mux.Router, storage storage.Storage) {
	router.HandleFunc("/api/warehouse", agent.CreateWareHouse(storage)).Methods("POST")
	router.HandleFunc("/api/checkin", agent.CheckedInAgents(storage)).Methods("POST")
	router.HandleFunc("/api/checkin", agent.CheckInAgent(storage)).Methods("GET")
	router.HandleFunc("/api/assignments", agent.GetAssignments(storage)).Methods("GET")
}
