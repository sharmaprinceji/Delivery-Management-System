# 🚚 Delivery Management System

A simplified last-mile delivery management system to assign orders to agents based on location, time, and profitability constraints — inspired by logistics platforms like Amazon, Delhivery, and Bluedart.

---

## 📦 Tech Stack

| Layer      | Tech              |
|------------|-------------------|
| Language   | Go (Golang)       |
| DB         | SQLite            |
| Routing    | `http.ServeMux`   |
| Scheduler  | Custom via `goroutine` |
| Validator  | `go-playground/validator` |
| Logging    | `slog` / `log`    |

---

## 🛠 Setup Instructions

### 🔁 Clone & Run

```bash
git clone https://github.com/your-username/delivery-management-system.git
cd delivery-management-system
go mod tidy
go run cmd/main.go --config=config/local.yaml


Required Project Structure:
delivery-management-system/
│
├── cmd/
│   └── main.go
├── config/
│   └── local.yaml
├── internal/
│   ├── types/
│   ├── storage/
│   ├── http/
│   ├── router/
│   ├── schedular/
│   └── jobs/
└── db/
    └── sqlite.go


#System Setup Flow
Follow the below route flow in sequence to fully simulate the system:
1.Agent Check-In:
POST /api/checkin
payload:
{
  "name": "Ravi Sharma",
  "warehouse_id": 1
}

2.Check Agent Assignments:
GET /api/assignments?page=1&limit=3

3.Create Warehouse:
POST /api/warehouse
payload:
{
  "name": "Bangalore North Hub",
  "location": {
    "lat": 12.9716,
    "lng": 77.5946
  }
}


4.Check-in Agent Again (After Warehouse):
POST /api/checkin

5.Create a Single Order:
POST /api/order
payload:
{
  "customer": "John Doe",
  "lat": 12.9721,
  "lng": 77.5940,
  "warehouse_id": 1
}

6. Bulk Create Orders:
POST /api/orders
payload:
[
  {
    "customer": "A",
    "lat": 12.9730,
    "lng": 77.5935,
    "warehouse_id": 2
  },
  {
    "customer": "B",
    "lat": 12.9735,
    "lng": 77.5931,
    "warehouse_id": 2
  }
]

7. Trigger Manual Allocation:
GET /api/allocate

8.Get Agent Utilization Summary (with pagination):
GET /api/agent-summary?page=1
response:
[
  {
    "agent_id": 1,
    "total_orders": 25,
    "total_km": 50,
    "total_minutes": 250,
    "profit": 875
  }
]

9. System Summary (with pagination):
GET /api/system-summary?page=2
response:
{
  "total_orders": 100,
  "assigned_orders": 85,
  "deferred_orders": 15,
  "agent_utilization": [ ... ]
}


***Business Rules Implemented***
Rule	Value
Max Agent Distance	100 km
Max Agent Time/Day	10 hours
Travel Time	1 km = 5 min
Min Agent Profit	₹500
25+ orders	₹35/order
50+ orders	₹42/order

#end.....


  

