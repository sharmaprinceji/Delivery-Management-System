package router

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/sharmaprinceji/delivery-management-system/db"
	"github.com/sharmaprinceji/delivery-management-system/internal/config"
	"github.com/sharmaprinceji/delivery-management-system/internal/schedular"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
)

func SetupRouter() (*mux.Router, storage.Storage) {
	router := mux.NewRouter()
	cfg := config.MustLoad()

	st, err := db.Mydb(cfg)
	if err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}
	log.Println("DB connection on..", cfg.HTTPServer.Addr)

	if err := st.InitSchema(); err != nil {
		log.Fatalf("schema error: %v", err)
	}

	schedular.SchedularJob(st)

	return router, st
}
