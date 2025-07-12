package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	// "strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
	"github.com/sharmaprinceji/delivery-management-system/internal/types"
	"github.com/sharmaprinceji/delivery-management-system/internal/utils/response"
)

var validate = validator.New()

// CreateWareHouse godoc
// @Summary Create a new warehouse
// @Description Accepts warehouse details and stores them in the system
// @Tags Warehouse
// @Accept json
// @Produce json
// @Param warehouse body types.WarehouseRequest true "Warehouse Details"
// @Success 201 {object} map[string]int64
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/warehouse [post]
func CreateWareHouse(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var warehouse types.Warehouse

		err := json.NewDecoder(r.Body).Decode(&warehouse)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("failed to decode request body: %v", err)))
			return
		}

		// Validate the struct 
		if err := validator.New().Struct(warehouse); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		id, err := storage.CreateWarehouse(warehouse.Name, warehouse.Location)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to create warehouse: %v", err)))
			return
		}

		slog.Info("warehouse created successfully", slog.Int64("id", id))
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"warehouse created successfully with Id": id})
	}
}

// CheckedInAgents godoc
// @Summary Check-in an agent
// @Description Allows an agent to check in to a warehouse
// @Tags Agent
// @Accept json
// @Produce json
// @Param agent body types.AgentCheckInRequest true "Agent Check-In Info"
// @Success 201 {object} map[string]int64
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/agent/checkin [post]
func CheckedInAgents(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var agent types.Agent

		if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request: %v", err)))
			return
		}

		if err := validator.New().Struct(agent); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		id, err := storage.CheckInAgents(agent.Name, agent.WarehouseID)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("check-in failed: %v", err)))
			return
		}

		slog.Info("Agent checked in successfully :", slog.Int64("id", id))
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"Agent checked successfully with Id": id})
	}
}


// CheckInAgent godoc
// @Summary Check in an agent (second version)
// @Description Marks the agent as checked-in at the warehouse
// @Tags Agent
// @Accept json
// @Produce json
// @Param agent body types.AgentCheckInRequest true "Agent Check-In Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/checkin [get]
func CheckInAgent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var agent types.Agent

		err := json.NewDecoder(r.Body).Decode(&agent)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is invalid")))
			return
		}
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("failed to decode request body: %v", err)))
			return
		}
		//debugging log
		log.Printf("Received agent data: %+v\n", agent)

		// Validation
		if err := validate.Struct(agent); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// Call storage layer to mark agent as checked in
		err = storage.CheckInAgent(agent)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to check-in agent: %v", err)))
			return
		}

		log.Printf("Agent checked in successfully: %+v", agent)
		response.WriteJSON(w, http.StatusCreated, map[string]string{"status": "agent checked in"})
	}
}


// GetAssignments godoc
// @Summary Get paginated assignments
// @Description Returns paginated list of assignments with formatted date
// @Tags Assignments
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{} "List of assignments"
// @Failure 500 {object} response.Response
// @Router /api/assignments [get]
func GetAssignments(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limitStr := r.URL.Query().Get("limit")
		limit, err := strconv.Atoi(limitStr)

		if err != nil || limit <= 0 {
			limit = 10
		}

		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		assignments, total, err := storage.GetPaginatedAssignments(limit, offset)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to fetch assignments: %v", err)))
			return
		}

		// Convert to formatted response
		var formatted []types.AssignmentResponse
		for _, a := range assignments {
			formatted = append(formatted, types.AssignmentResponse{
				ID:         a.ID,
				AgentID:    a.AgentID,
				OrderID:    a.OrderID,
				AssignedAt: a.AssignedAt.Format("02/01/2006 03:04 PM"),
			})
		}

		totalPages := int(math.Ceil(float64(total) / float64(limit)))

		response.WriteJSON(w, http.StatusOK, map[string]any{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"data":         formatted,
		})
	}
}

