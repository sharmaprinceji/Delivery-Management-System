package sqlite

import (
	"database/sql"
	"fmt"
	"math"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sharmaprinceji/delivery-management-system/internal/config"
	"github.com/sharmaprinceji/delivery-management-system/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

func (s *Sqlite) Save(data any) error {
	switch v := data.(type) {
	case types.Student:
		_, err := s.Db.Exec(`INSERT INTO Users (name, age, email, city) VALUES (?, ?, ?, ?)`,
			v.Name, v.Age, v.Email, v.City)
		return err
	default:
		return fmt.Errorf("unsupported type in Save: %T", v)
	}
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		age INTEGER,
		email TEXT UNIQUE,
		city Text
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}




func (s *Sqlite) InitSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS warehouses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			lat REAL NOT NULL,
			lng REAL NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS agents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			warehouse_id INTEGER NOT NULL,
			checked_in BOOLEAN NOT NULL,
			FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
		);`,

		`CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			customer TEXT NOT NULL,
			lat REAL NOT NULL,
			lng REAL NOT NULL,
			warehouse_id INTEGER NOT NULL,
			assigned BOOLEAN DEFAULT 0,
			agent_id INTEGER,
			FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
			FOREIGN KEY (agent_id) REFERENCES agents(id)
        );`,
	}

	for _, q := range queries {
		_, err := s.Db.Exec(q)
		if err != nil {
			return fmt.Errorf("schema error: %w in query: %q", err, q)
		}
	}
	return nil
}


func (s *Sqlite) GetCheckedInAgents() ([]types.Agent, error) {
	rows, err := s.Db.Query("SELECT id, name, warehouse_id, checked_in FROM agents WHERE checked_in = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []types.Agent
	for rows.Next() {
		var a types.Agent
		if err := rows.Scan(&a.ID, &a.Name, &a.WarehouseID, &a.CheckedIn); err != nil {
			return nil, err
		}
		agents = append(agents, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return agents, nil
}

func (s *Sqlite) GetUnassignedOrders() ([]types.Order, error) {
	rows, err := s.Db.Query("SELECT id, customer, lat, lng, warehouse_id, assigned, agent_id FROM orders WHERE assigned = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []types.Order
	for rows.Next() {
		var o types.Order
		var agentID sql.NullInt64

		if err := rows.Scan(&o.ID, &o.Customer, &o.Lat, &o.Lng, &o.WarehouseID, &o.Assigned, &agentID); err != nil {
			return nil, err
		}

		if agentID.Valid {
			o.AgentID = &agentID.Int64
		} else {
			o.AgentID = nil
		}

		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Sqlite) AssignOrderToAgent(orderID int64, agentID int64) error {
	_, err := s.Db.Exec(`UPDATE orders SET assigned = 1, agent_id = ? WHERE id = ?`, agentID, orderID)
	return err
}

func (s *Sqlite) CheckInAgent(agent types.Agent) error {
	_, err := s.Db.Exec(`
		INSERT INTO agents (name, warehouse_id, checked_in)
		VALUES (?, ?, 1)
		ON CONFLICT(name, warehouse_id) DO UPDATE SET checked_in = 1
	`, agent.Name, agent.WarehouseID)
	return err
}

func (s *Sqlite) GetAllAssignments() ([]types.Assignment, error) {
	rows, err := s.Db.Query(`
		SELECT id, agent_id, order_id, assigned_at
		FROM assignments
		ORDER BY assigned_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []types.Assignment
	for rows.Next() {
		var a types.Assignment
		if err := rows.Scan(&a.ID, &a.AgentID, &a.OrderID, &a.AssignedAt); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, nil
}

func (s *Sqlite) CreateWarehouse(name string, location types.Location) (int64, error) {
	stmt, err := s.Db.Prepare(`
		INSERT INTO warehouses (name, lat, lng)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, location.Lat, location.Lng)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Sqlite) CheckInAgents(name string, warehouseID int64) (int64, error) {
	stmt, err := s.Db.Prepare(`
		INSERT INTO agents (name, warehouse_id, checked_in)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, warehouseID, true)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Sqlite) CreateOrder(o types.Order) (int64, error) {
	stmt, err := s.Db.Prepare(`
		INSERT INTO orders (customer, lat, lng, warehouse_id, assigned)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(o.Customer, o.Lat, o.Lng, o.WarehouseID, o.Assigned)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (s *Sqlite) CreateBulkOrders(orders []types.Order) (int, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO orders (customer, lat, lng, warehouse_id, assigned)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	count := 0
	for _, order := range orders {
		_, err := stmt.Exec(order.Customer, order.Lat, order.Lng, order.WarehouseID, order.Assigned)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		count++
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Sqlite) GetAgentSummary() ([]types.AgentSummary, error) {
	rows, err := s.Db.Query(`
		SELECT agent_id, COUNT(*) as total_orders
		FROM orders
		WHERE assigned = 1
		GROUP BY agent_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []types.AgentSummary
	for rows.Next() {
		var summary types.AgentSummary
		err := rows.Scan(&summary.AgentID, &summary.TotalOrders)
		if err != nil {
			return nil, err
		}

		// mock distance/time = 2.0 km/order = 10 minutes/order
		summary.TotalKm = float64(summary.TotalOrders) * 2.0
		summary.TotalMinutes = float64(summary.TotalOrders) * 10.0

		// Profit tiers
		if summary.TotalOrders >= 50 {
			summary.Profit = float64(summary.TotalOrders) * 42
		} else if summary.TotalOrders >= 25 {
			summary.Profit = float64(summary.TotalOrders) * 35
		} else {
			summary.Profit = float64(summary.TotalOrders) * 20 // base rate
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (s *Sqlite) GetAgentSummaryPaginated(page int, limit int) (types.PaginatedAgentSummary, error) {
	offset := (page - 1) * limit

	// 1. Get total agent count
	var totalCount int
	err := s.Db.QueryRow("SELECT COUNT(DISTINCT agent_id) FROM orders WHERE assigned = 1").Scan(&totalCount)
	if err != nil {
		return types.PaginatedAgentSummary{}, err
	}

	// 2. Get actual data with pagination
	rows, err := s.Db.Query(`
		SELECT agent_id, COUNT(*) as total_orders
		FROM orders
		WHERE assigned = 1
		GROUP BY agent_id
		ORDER BY agent_id
		LIMIT ? OFFSET ?
	`, limit, offset)

	if err != nil {
		return types.PaginatedAgentSummary{}, err
	}
	defer rows.Close()

	var summaries []types.AgentSummary
	for rows.Next() {
		var summary types.AgentSummary
		err := rows.Scan(&summary.AgentID, &summary.TotalOrders)
		if err != nil {
			return types.PaginatedAgentSummary{}, err
		}

		summary.TotalKm = float64(summary.TotalOrders) * 2.0
		summary.TotalMinutes = float64(summary.TotalOrders) * 10.0

		if summary.TotalOrders >= 50 {
			summary.Profit = float64(summary.TotalOrders) * 42
		} else if summary.TotalOrders >= 25 {
			summary.Profit = float64(summary.TotalOrders) * 35
		} else {
			summary.Profit = float64(summary.TotalOrders) * 20
		}

		summaries = append(summaries, summary)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return types.PaginatedAgentSummary{
		CurrentPage: page,
		TotalPages:  totalPages,
		Data:        summaries,
	}, nil
}

func (s *Sqlite) GetSystemSummary() (types.SystemSummary, error) {
	var summary types.SystemSummary

	err := s.Db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&summary.TotalOrders)
	if err != nil {
		return summary, err
	}

	err = s.Db.QueryRow("SELECT COUNT(*) FROM orders WHERE assigned = 1").Scan(&summary.AssignedOrders)
	if err != nil {
		return summary, err
	}

	summary.DeferredOrders = summary.TotalOrders - summary.AssignedOrders

	agentUtil, err := s.GetAgentSummary()
	if err != nil {
		return summary, err
	}

	summary.AgentUtilization = agentUtil
	return summary, nil
}
