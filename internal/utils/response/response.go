package response


import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  string  `json:"status"`      
	Error   string   `json:"error"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOk="OK"
	StatusError="Error"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return  json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) ErrorResponse {
	return ErrorResponse{
		Status: StatusError,
		Error: err.Error(),
	}
}

func ValidationError(errs  validator.ValidationErrors) ErrorResponse {
	var errMsgs []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s failed validation for tag\n", err.Field()))
		}
	}

	return ErrorResponse{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "	),
	}
}