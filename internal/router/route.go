package router

import (
	"log"
	"net/http"

	// "time"

	"github.com/sharmaprinceji/delivery-management-system/db"
	"github.com/sharmaprinceji/delivery-management-system/internal/config"
	//"github.com/sharmaprinceji/delivery-management-system/internal/router/agentRoute"
	//"github.com/sharmaprinceji/delivery-management-system/internal/router/orderRoute"
	"github.com/sharmaprinceji/delivery-management-system/internal/schedular"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
)

// router/router.go
func SetupRouter() (*http.ServeMux, storage.Storage) {
	router := http.NewServeMux()
	cfg := config.MustLoad()

	st, err := db.Mydb(cfg)
	if err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}
	
	log.Println("Db connection on..", cfg.HTTPServer.Addr)

	if err := st.InitSchema(); err != nil {
		log.Fatalf("schema error: %v", err)
	}

	schedular.SchedularJob(st)

	return router, st
}

