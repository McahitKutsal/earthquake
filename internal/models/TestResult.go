package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestResult struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	TestRequestID primitive.ObjectID `bson:"test_Request_ID"`
	Endpoint      string             `bson:"endpoint"`
	Concurrency   int                `bson:"concurrency"`
	Requests      int                `bson:"requests"`
	Method        string             `bson:"method"`
	Headers       map[string]string  `bson:"headers"`
	Body          string             `bson:"body"`
	TestSummary   TestSummary        `bson:"test_summary"`
	CreatedAt     time.Time          `bson:"created_at"`
}
