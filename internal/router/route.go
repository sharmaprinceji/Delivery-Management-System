package router

import (
	"log"
	"net/http"

	"github.com/sharmaprinceji/delivery-management-system/db"
	"github.com/sharmaprinceji/delivery-management-system/internal/config"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
)

func StudentRoute() (*http.ServeMux,storage.Storage){
	router:=http.NewServeMux()
	cfg := config.MustLoad()

	storage, err := db.Mydb(cfg)

	if err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}

	
	log.Println("Db connection on..", cfg.HTTPServer.Addr)
    
	err:= storage.InitSchema()
	
    if err != nil {
		log.Fatalf("schema error: %v", err)
	}

	for {
		now := time.Now()
		if now.Hour() == 7 {
			log.Println("Running allocation job...")
			err := jobs.AllocateOrders(storage)
			if err != nil {
				log.Printf("Job failed: %v", err)
			}
		}
		time.Sleep(time.Minute * 10)
	}

	return router,storage;
}