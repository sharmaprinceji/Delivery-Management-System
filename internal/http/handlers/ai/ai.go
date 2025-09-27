package ai

import (
	"fmt"
	"net/http"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
)

func Ai_input(storage storage.Storage) http.HandlerFunc{
    return func(w http.ResponseWriter, r *http.Request) {
		 fmt.Println("ai handler working !")
	}
}

func Ai_output(storage storage.Storage) http.HandlerFunc{
    return func(w http.ResponseWriter, r *http.Request) {
		 fmt.Println("ai handler working !")
	}
}