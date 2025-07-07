package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	// "strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
	"github.com/sharmaprinceji/delivery-management-system/internal/types"
	"github.com/sharmaprinceji/delivery-management-system/internal/utils/response"
)


var validate = validator.New()
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

func GetAssignments(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assignments, err := storage.GetAllAssignments()
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to fetch assignments: %v", err)))
			return
		}
		response.WriteJSON(w, http.StatusOK, assignments)
	}
}

func CreateWareHouse(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var warehouse types.Warehouse

		err := json.NewDecoder(r.Body).Decode(&warehouse)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("failed to decode request body: %v", err)))
			return
		}

		// Validate the struct (assuming Lat/Lng is inside `Location`)
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
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": id})
	}
}

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

		slog.Info("Agent checked in", slog.Int64("id", id))
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": id})
	}
}






 