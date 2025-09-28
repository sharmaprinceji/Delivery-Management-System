package openairoute
import (
	"github.com/gorilla/mux"
	"github.com/sharmaprinceji/delivery-management-system/internal/http/handlers/ai"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
)

func OpenAiRoutes(router *mux.Router, storage storage.Storage) {
	router.HandleFunc("/api/inputai",ai.Ai_input(storage)).Methods("POST")
	router.HandleFunc("/api/outputai",ai.Ai_output(storage)).Methods("GET")
}
