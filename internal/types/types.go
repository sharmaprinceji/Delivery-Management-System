package types
//models for each table in the database

import "time"

type Student struct {
	ID        int64 `json:"id"`  
	Name      string `validate:"required" json:"name"`
	Age       int    `validate:"required" json:"age"`
	Email     string `validate:"required" json:"email"`
	City	  string `validate:"required" json:"city"`
}

type Agent struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	WarehouseID int64   `json:"warehouse_id"`
	CheckedIn   bool    `json:"checked_in"`
}

type Warehouse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}

type Order struct {
	ID          int64   `json:"id"`
	Customer    string  `json:"customer"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	WarehouseID int64   `json:"warehouse_id"`
	Assigned    bool    `json:"assigned"`
	AgentID     *int64  `json:"agent_id,omitempty"`
}


type Assignment struct {
	ID        int64 `json:"id"`
	AgentID   int64
	OrderID   int64
	AssignedAt time.Time
}