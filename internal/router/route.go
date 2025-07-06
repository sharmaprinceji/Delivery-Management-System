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

func StudentRoute() (*http.ServeMux,storage.Storage){
	router:=http.NewServeMux()
	cfg := config.MustLoad()

	storage, err := db.Mydb(cfg)

	if err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}

	
	log.Println("Db connection on..", cfg.HTTPServer.Addr)
    
	er:= storage.InitSchema()

    if er != nil {
		log.Fatalf("schema error: %v", er)
	}

    schedular.SchedularJob(storage);

	//insert all specific router here..
	//  _:=agentroute.AgentRouter()
	//  _:=orderroute.OrderRouter()

	
	return router,storage;
}