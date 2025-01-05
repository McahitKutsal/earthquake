package handlers

import (
	"encoding/json"
	"net/http"

	"earthquake/internal/models"
	"earthquake/internal/performance"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleTestRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var testRequests []models.TestRequest
	err := json.NewDecoder(r.Body).Decode(&testRequests)
	if err != nil {
		http.Error(w, "Invalid Test Request", http.StatusBadRequest)
		return
	}

	// Generate a unique test ID
	testRequestID := primitive.NewObjectID()

	// Run the test asynchronously
	for _, testRequest := range testRequests {
		go func() {
			performance.RunTest(testRequest, testRequestID)
		}()
	}

	// Respond with the test ID immediately
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"test_id": testRequestID.Hex()})
}
