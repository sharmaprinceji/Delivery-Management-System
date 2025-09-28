package ai

import (
	"context"
	"encoding/json"

	"fmt"
	// "log/slog"
	"net/http"

	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
	"github.com/spf13/viper"
	// "google.golang.org/api/apikeys/v2"
	"google.golang.org/api/option"
)

type InputReq struct {
	Prompt string `json:"prompt"`
}

// AiInput godoc
// @Summary Generate AI response using Gemini LLM
// @Description Accepts a prompt and returns AI-generated response using Google Gemini API
// @Tags AI
// @Accept json
// @Produce json
// @Param input body InputReq true "AI prompt input"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/inputai [post]
func Ai_input(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body InputReq
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		ctx := context.Background()

		// 1. Try to get key from environment variable
		apiKey := os.Getenv("GEMINI_API_KEY")
		//fmt.Println("API Key from env:", apiKey)
		// 2. Agar env me nahi mila to config.yaml se le lo (via Viper)
		if apiKey == "" {
			apiKey = viper.GetString("variables.gemini.GEMINI_API_KEY")
		}

		// 3. Agar dono jagah nahi mila to error
		if apiKey == "" {
			http.Error(w, "Gemini API key not found", http.StatusInternalServerError)
			return
		}

		// 4. Create Gemini client
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			http.Error(w, "failed to init gemini client: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer client.Close()
		
         //list of models available
		// it := client.ListModels(ctx)
		// for {
		// 	model, err := it.Next()
		// 	if err != nil {
		// 		break
		// 	}
		// 	fmt.Printf("Model: %s, Display Name: %s\n", model.Name, model.DisplayName)
		// }

		// 5. Select Model
		model := client.GenerativeModel("models/gemini-2.5-flash")

		// 6. Generate content
		resp, err := model.GenerateContent(ctx, genai.Text(body.Prompt))

		// if resp != nil {
		// 	fmt.Println("Gemini response received", resp)
		// }

		if resp == nil {
			fmt.Println("Gemini response is nil", err)
		}

		if err != nil {
			http.Error(w, "gemini request failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 7. Extract text response
		output := ""
		for _, cand := range resp.Candidates {
			for _, part := range cand.Content.Parts {
				if text, ok := part.(genai.Text); ok {
					output += string(text)
				}
			}
		}

		// (Optional) Agar storage me save karna hai
		// _ = storage.Save(map[string]string{"prompt": body.Prompt, "response": output})

		// 8. Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"answer": output,
		})
	}
}

// ---- Output API ----
func Ai_output(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Example: agar aapne DB me save kiya hoga to fetch kar sakte ho
		// yaha abhi demo static response bhej raha hu
		resp := map[string]string{
			"message": "Gemini Output API working",
			"result":  "Fetch from DB or use Ai_input directly",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
