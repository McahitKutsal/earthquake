package handlers

import (
	"encoding/json"
	"net/http"

	"earthquake/internal/config"
	"earthquake/internal/performance"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleTestRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cfg config.Config
	err := json.NewDecoder(r.Body).Decode(&cfg)
	if err != nil {
		http.Error(w, "Invalid configuration", http.StatusBadRequest)
		return
	}

	// Generate a unique test ID
	testID := primitive.NewObjectID()

	// Run the test asynchronously
	go func() {
		performance.RunTest(cfg, testID)
	}()

	// Respond with the test ID immediately
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"test_id": testID.Hex()})
}
