package agent

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	// "strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
	"github.com/sharmaprinceji/delivery-management-system/internal/types"
	"github.com/sharmaprinceji/delivery-management-system/internal/utils/response"
)

//var validate = validator.New()

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
		var req types.WarehouseRequest // ✅ Correct type for Swagger match

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("failed to decode request body: %v", err)))
			return
		}

		// Validate the struct 
		if err := validator.New().Struct(req); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		id, err := storage.CreateWarehouse(req.Name, req.Location)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to create warehouse: %v", err)))
			return
		}

		slog.Info("warehouse created successfully", slog.Int64("id", id))
		response.WriteJSON(w, http.StatusCreated, map[string]int64{
			"warehouse created successfully with Id": id,
		})
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
		var req types.AgentCheckInRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request: %v", err)))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		id, err := storage.CheckInAgents(req.Name, req.WarehouseID)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("check-in failed: %v", err)))
			return
		}

		slog.Info("Agent checked in successfully :", slog.Int64("id", id))
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"Agent checked successfully with Id": id})
	}
}


// In handler/agent.go
// GetAgentDetails godoc
// @Summary Get Agent Details
// @Description Returns full summary of agent including total orders, profit, etc.
// @Tags Agent
// @Produce json
// @Param agent_id path int true "Agent ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/agent/{agent_id} [get]
func GetAgentDetails(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		agentIDStr := vars["agent_id"]
		agentID, err := strconv.ParseInt(agentIDStr, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid agent ID")))
			return
		}

		agentData, err := storage.GetAgentDetails(agentID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response.WriteJSON(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("agent not found")))
				return
			}
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, agentData)
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

