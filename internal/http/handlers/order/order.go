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

func CreateOrder(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var o types.Order

		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request: %v", err)))
			return
		}

		if err := validator.New().Struct(o); err != nil {
			validationErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validationErrs))
			return
		}

		id, err := storage.CreateOrder(o)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to create order: %v", err)))
			return
		}

		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": id})
	}
}

func CreateBulkOrders(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var orders []types.Order

		if err := json.NewDecoder(r.Body).Decode(&orders); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid request body: %v", err)))
			return
		}

		validate := validator.New()

		for i, order := range orders {
			if err := validate.Struct(order); err != nil {
				response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("validation failed at index %d: %v", i, err)))
				return
			}
		}

		count, err := storage.CreateBulkOrders(orders)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to insert orders: %v", err)))
			return
		}

		response.WriteJSON(w, http.StatusCreated, map[string]int{"inserted": count})
	}
}

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

		limit := 10 // fixed limit per page
		summaries, err := storage.GetAgentSummaryPaginated(page, limit)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to fetch summary: %v", err)))
			return
		}

		response.WriteJSON(w, http.StatusOK, summaries)
	}
}

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
