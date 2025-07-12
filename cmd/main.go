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
	"github.com/sharmaprinceji/delivery-management-system/internal/router/agentRoute"
	"github.com/sharmaprinceji/delivery-management-system/internal/router/orderRoute"

	"github.com/swaggo/http-swagger"
	_ "github.com/sharmaprinceji/delivery-management-system/docs"                     
)

// @title Delivery Management System API
// @version 1.0
// @description API for managing delivery agents, orders and allocation system
// @contact.name Prince Raj
// @contact.email prince@example.com
// @host localhost:5002
// @BasePath /

func main() {
	cfg := config.MustLoad()

	route, storage := router.SetupRouter()

	agentroute.RegisterAgentRoutes(route, storage)
	orderroute.RegisterOrderRoutes(route, storage)

	// Swagger route
	route.Handle("/", httpSwagger.WrapHandler)


	port := os.Getenv("PORT")
	if port == "" {
		port = "5002"
	}

	server := http.Server{
		Addr:    ":" + port,
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

	<-done
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	er := server.Shutdown(ctx)
	if er != nil {
		slog.Error("failed to shutting down server", slog.String("error", er.Error()))
	}

	slog.Info("Server stopped gracefully")
}
