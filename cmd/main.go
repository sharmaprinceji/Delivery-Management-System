package main

import (
	"context"
	"fmt"

	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sharmaprinceji/delivery-management-system/internal/config"
	 "github.com/sharmaprinceji/delivery-management-system/internal/router"
	// "github.com/sharmaprinceji/delivery-management-system/internal/router/agentRoute"
	//"github.com/sharmaprinceji/delivery-management-system/internal/router/orderRoute"
)

func main() {
	//loadConfig()
	cfg := config.MustLoad() /// 👈 call only once here

	//seprate db setup ..
	//route:=agentroute.AgentRouter();
	//route := orderroute.OrderRouter() // 👈 call only once here
	route,_ := router.StudentRoute() // 👈 call only once here

	//setup server.
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: route,
	}

	slog.Info("Starting server...", slog.String("address", cfg.HTTPServer.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			return
		}
	}()

	<-done // Wait for a signal to stop the server

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	er := server.Shutdown(ctx)
	if er != nil {
		slog.Error("failed to shutting down server", slog.String("error", er.Error()))
	}

	slog.Info("Server stopped gracefully")

}