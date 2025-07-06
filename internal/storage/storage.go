package storage


 import "github.com/sharmaprinceji/delivery-management-system/internal/types"

type Storage interface {
	CreateStudent(name string,email string,age int,city string) (int64,error)
	GetStudentById(id int64) (types.Student, error)
	Save(data any) error
	GetCheckedInAgents() ([]types.Agent, error)
	GetUnassignedOrders() ([]types.Order, error)
	AssignOrderToAgent(orderID int64, agentID int64) error
	CheckInAgent(agent types.Agent) error

	// Optional schema setup (used in SQLite for migration)
	InitSchema() error
}

