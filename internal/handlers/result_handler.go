package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"earthquake/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTestResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	testID := r.URL.Query().Get("id")
	if testID == "" {
		http.Error(w, "Missing test ID", http.StatusBadRequest)
		return
	}

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(testID)
	if err != nil {
		http.Error(w, "Invalid test ID", http.StatusBadRequest)
		return
	}

	// Fetch from MongoDB
	collection := database.GetCollection("performance", "results")
	var result map[string]interface{}
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		http.Error(w, "Test result not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
