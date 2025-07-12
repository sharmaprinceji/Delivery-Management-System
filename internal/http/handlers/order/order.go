package order

import (
	"encoding/json"
	"fmt"
	"strconv"

	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sharmaprinceji/delivery-management-system/internal/jobs"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
	"github.com/sharmaprinceji/delivery-management-system/internal/types"
	"github.com/sharmaprinceji/delivery-management-system/internal/utils/response"
)

// CreateOrder godoc
// @Summary Create a new order
// @Description Creates a new customer order and stores it in the database
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body types.OrderRequest true "Order details"
// @Success 201 {object} map[string]int64
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/order [post]
func CreateOrder(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderRequest 

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request: %v", err)))
			return
		}

	
		if err := validator.New().Struct(req); err != nil {
			validationErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validationErrs))
			return
		}

		order := types.Order{
			Customer:    req.Customer,
			Lat:         req.Lat,
			Lng:         req.Lng,
			WarehouseID: req.WarehouseID,
			Assigned:    false,
		}

		id, err := storage.CreateOrder(order)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to create order: %v", err)))
			return
		}

	
		response.WriteJSON(w, http.StatusCreated, map[string]int64{
			"Order has been created successfully with id": id,
		})
	}
}



// CreateBulkOrders godoc
// @Summary Create multiple orders in bulk
// @Description Accepts a list of customer orders and stores them in the database
// @Tags Orders
// @Accept json
// @Produce json
// @Param orders body types.BulkOrderRequest true "List of order requests"
// @Success 201 {object} map[string]int
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/orders/bulk [post]
func CreateBulkOrders(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BulkOrderRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body: %v", err)))
			return
		}

		validate := validator.New()
		for i, order := range req.Orders {
			if err := validate.Struct(order); err != nil {
				response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("validation failed at index %d: %v", i, err)))
				return
			}
		}

	
		var orders []types.Order
		for _, o := range req.Orders {
			orders = append(orders, types.Order{
				Customer:    o.Customer,
				Lat:         o.Lat,
				Lng:         o.Lng,
				WarehouseID: o.WarehouseID,
				Assigned:    false,
			})
		}

		count, err := storage.CreateBulkOrders(orders)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to insert orders: %v", err)))
			return
		}

		response.WriteJSON(w, http.StatusCreated, map[string]int{"Ordered inserted": count})
	}
}




// ManualAllocation godoc
// @Summary Trigger manual allocation of orders
// @Description Runs the allocation algorithm to assign orders to agents
// @Tags Orders
// @Produce plain
// @Success 200 {string} string "Allocation successful"
// @Failure 500 {string} string "Allocation failed"
// @Router /api/allocate [get]
func ManualAllocation(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := jobs.AllocateOrders(s)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(" Allocation failed: " + err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Allocation successful"))
	}
}


// GetAgentSummary godoc
// @Summary Get agent summary with pagination
// @Description Returns a paginated summary of agents, including total orders, distance, time, and profit
// @Tags Summary
// @Accept json
// @Produce json
// @Param page query int false "Page number (default is 1)"
// @Success 200 {object} types.PaginatedAgentSummary
// @Failure 500 {object} response.Response
// @Router /api/agent-summary [get]
func GetAgentSummary(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		page := 1
		if pageStr != "" {
			p, err := strconv.Atoi(pageStr)
			if err == nil && p > 0 {
				page = p
			}
		}

		limit := 10 
		summaries, err := storage.GetAgentSummaryPaginated(page, limit)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to fetch summary: %v", err)))
			return
		}

		response.WriteJSON(w, http.StatusOK, summaries)
	}
}


// GetSystemSummary godoc
// @Summary Get system summary with paginated agent utilization
// @Description Returns a system-wide summary including total, assigned, and deferred orders, along with agent utilization
// @Tags Summary
// @Accept json
// @Produce json
// @Param page query int false "Page number (default is 1)"
// @Success 200 {object} types.SystemSummary
// @Failure 500 {object} response.Response
// @Router /api/system-summary [get]
func GetSystemSummary(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		page := 1
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
		limit := 10

		summary, err := storage.GetSystemSummaryPaginated(page, limit)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to get system summary: %v", err)))
			return
		}
		response.WriteJSON(w, http.StatusOK, summary)
	}
}

