package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	//"log/slog"
	"net/http"
	"strconv"

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
		id := r.PathValue("id")
		if id == "" {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("id is required")))
			return
		}
        intid,err:=strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id format: %v", err)))
			return
		}

		student, err := storage.GetStudentById(intid)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to get student: %v", err)))
			return
		}

		response.WriteJSON(w, http.StatusOK, student)
	}
}

