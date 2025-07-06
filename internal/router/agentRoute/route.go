package agentroute

import (
	"log"
	"net/http"

	"github.com/sharmaprinceji/delivery-management-system/db"
	"github.com/sharmaprinceji/delivery-management-system/internal/config"
	"github.com/sharmaprinceji/delivery-management-system/internal/http/handlers/agent"
	"github.com/sharmaprinceji/delivery-management-system/internal/router"
)


func StudentRouter() *http.ServeMux {
	cfg := config.MustLoad()

	storage, err := db.Mydb(cfg)
	if err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}

	route :=  router.StudentRoute()

	route.HandleFunc("GET /api/student/{id}", agent.GetById(storage))
	

	return route;
}