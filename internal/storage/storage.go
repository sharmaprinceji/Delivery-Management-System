package storage


 import "github.com/sharmaprinceji/delivery-management-system/internal/types"

type Storage interface {
	Save(data any) error
	GetCheckedInAgents() ([]types.Agent, error)
	GetUnassignedOrders() ([]types.Order, error)
	AssignOrderToAgent(orderID int64, agentID int64) error
	CheckInAgent(agent types.Agent) error
	GetAllAssignments() ([]types.Assignment, error)


	// Optional schema setup (used in SQLite for migration)
	InitSchema() error
	CreateWarehouse(name string, location types.Location) (int64, error)
	CheckInAgents(name string, warehouseID int64) (int64, error)
	CreateOrder(o types.Order) (int64, error)
	CreateBulkOrders(orders []types.Order) (int, error)
	GetAgentSummaryPaginated(page int, limit int) (types.PaginatedAgentSummary, error)
    GetSystemSummary() (types.SystemSummary, error)


}

