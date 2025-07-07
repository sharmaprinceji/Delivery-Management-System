package types

import "time"

type Student struct {
	ID    int64  `json:"id"`
	Name  string `json:"name" validate:"required"`
	Age   int    `json:"age" validate:"required"`
	Email string `json:"email" validate:"required,email"` // optional: email format
	City  string `json:"city" validate:"required"`
}

// Location used in Warehouse and geo fields
type Location struct {
	Lat float64 `json:"lat" validate:"required"`
	Lng float64 `json:"lng" validate:"required"`
}

// Agent model
type Agent struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" validate:"required"`
	WarehouseID int64  `json:"warehouse_id" validate:"required"`
	CheckedIn   bool   `json:"checked_in,omitempty"` 
}

// Warehouse model
type Warehouse struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name" validate:"required"`
	Location Location `json:"location" validate:"required"` 
}

// Order model
type Order struct {
	ID          int64   `json:"id"`
	Customer    string  `json:"customer" validate:"required"`
	Lat         float64 `json:"lat" validate:"required"`
	Lng         float64 `json:"lng" validate:"required"`
	WarehouseID int64   `json:"warehouse_id" validate:"required"`
	Assigned    bool    `json:"assigned"`              
	AgentID     *int64  `json:"agent_id,omitempty"`    
}

// Assignment model
type Assignment struct {
	ID         int64     `json:"id"`
	AgentID    int64     `json:"agent_id"`
	OrderID    int64     `json:"order_id"`
	AssignedAt time.Time `json:"assigned_at"`
}

type AgentSummary struct {
	AgentID      int64   `json:"agent_id"`
	TotalOrders  int     `json:"total_orders"`
	TotalKm      float64 `json:"total_km"`
	TotalMinutes float64 `json:"total_minutes"`
	Profit       float64 `json:"profit"`
}

type PaginatedAgentSummary struct {
	CurrentPage int              `json:"current_page"`
	TotalPages  int              `json:"total_pages"`
	Data        []AgentSummary  `json:"data"`
}

type SystemSummary struct {
	TotalOrders     int                  `json:"total_orders"`
	AssignedOrders  int                  `json:"assigned_orders"`
	DeferredOrders  int                  `json:"deferred_orders"`
	AgentUtilization []AgentSummary      `json:"agent_utilization"`
}

