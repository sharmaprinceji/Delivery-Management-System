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

	"github.com/gorilla/mux"
	"github.com/sharmaprinceji/delivery-management-system/internal/router"

	"github.com/sharmaprinceji/delivery-management-system/internal/router/agentRoute"
	"github.com/sharmaprinceji/delivery-management-system/internal/router/orderRoute"

	_ "github.com/sharmaprinceji/delivery-management-system/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)
   // delivery-management-system-h5nh.onrender.com  // localhost:5002

// @title Delivery Management System API
// @version 1.0
// @description API for managing delivery agents, orders and allocation system
// @contact.name Prince Raj
// @contact.email princesh1411@gmail.com
// @host delivery-management-system-h5nh.onrender.com
// @BasePath /
func main() {
	//cfg := config.MustLoad()

	route, storage := router.SetupRouter()

	// Enable CORS
	route.Use(mux.CORSMethodMiddleware(route))
	route.Use(corsMiddleware)

	agentRoute.RegisterAgentRoutes(route, storage)
	orderroute.RegisterOrderRoutes(route, storage)

	// Swagger route
	route.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5002"
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: route,
	}

	slog.Info("Starting server...", slog.String("address",  ":" + port))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	<-done
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutting down server", slog.String("error", err.Error()))
	}

	slog.Info("Server stopped gracefully")
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

	
		allowedOrigins := map[string]bool{
			"https://delivery-management-system-h5nh.onrender.com": true,
			"http://localhost:5002":                                 true,
			"https://delivery-management-system-h5nh.onrender.com/swagger": true,
		}

		// If the origin is allowed, set appropriate headers
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Set headers required for preflight and CORS
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}


