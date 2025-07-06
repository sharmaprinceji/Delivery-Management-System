package sqlite

import (
	"database/sql"
	"fmt"
	"log"

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
	
	_,err=db.Exec(`CREATE TABLE IF NOT EXISTS Users (
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
		Db:db,
	},nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int, city string) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO Users (name, email, age, city) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, email, age, city)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM Users WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Age, &student.Email, &student.City)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s",fmt.Sprint(id)) // No student found with the given ID
		}

		return types.Student{}, fmt.Errorf("failed to query student: %w", err) // Other error occurred while querying
	}
	log.Printf("Student found: %v\n", student)

	return student, nil
}

func (s *Sqlite) InitSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS agents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			warehouse_id INTEGER,
			checked_in BOOLEAN
		);`,
		`CREATE TABLE IF NOT EXISTS warehouses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			lat REAL,
			lng REAL
		);`,
		`CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			customer TEXT,
			lat REAL,
			lng REAL,
			warehouse_id INTEGER,
			assigned BOOLEAN,
			agent_id INTEGER
		);`,
	}

	for _, q := range queries {
		_, err := s.Db.Exec(q)
		if err != nil {
			return err
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
		err := rows.Scan(&a.ID, &a.Name, &a.WarehouseID, &a.CheckedIn)
		if err != nil {
			return nil, err
		}
		agents = append(agents, a)
	}
	return agents, nil
}

func (s *Sqlite) GetUnassignedOrders() ([]types.Order, error) {
	rows, err := s.Db.Query("SELECT id, customer, lat, lng, warehouse_id FROM orders WHERE assigned = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []types.Order
	for rows.Next() {
		var o types.Order
		err := rows.Scan(&o.ID, &o.Customer, &o.Lat, &o.Lng, &o.WarehouseID)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (s *Sqlite) AssignOrderToAgent(orderID int64, agentID int64) error {
	_, err := s.Db.Exec(`UPDATE orders SET assigned = 1, agent_id = ? WHERE id = ?`, agentID, orderID)
	return err
}



