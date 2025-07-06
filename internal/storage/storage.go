package storage


 import "github.com/sharmaprinceji/delivery-management-system/internal/types"

type Storage interface {
	CreateStudent(name string,email string,age int,city string) (int64,error)
	GetStudentById(id int64) (types.Student, error)
	Save(data any) error
	GetCheckedInAgents()
	GetUnassignedOrders()
	AssignOrderToAgent()
	InitSchema()
}

