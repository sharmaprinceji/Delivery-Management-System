package openairoute
import (
	"github.com/gorilla/mux"
	"github.com/sharmaprinceji/delivery-management-system/internal/http/handlers/ai"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
)

func OpenAiRoutes(router *mux.Router, storage storage.Storage) {
	router.HandleFunc("/api/input-Ai",ai.Ai_input(storage)).Methods("POST")
	router.HandleFunc("/api/output-Ai",ai.Ai_output(storage)).Methods("GET")
}
