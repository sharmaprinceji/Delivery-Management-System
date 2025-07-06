package sqlite

import (
	"database/sql"
	 "fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sharmaprinceji/delivery-management-system/internal/config"
	 "github.com/sharmaprinceji/delivery-management-system/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	
	_,err=db.Exec(`CREATE TABLE IF NOT EXISTS students (
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
	stmt, err := s.Db.Prepare("INSERT INTO students(name, email, age, city) VALUES(?, ?, ?, ?)")
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
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 2")
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

	return student, nil
}

