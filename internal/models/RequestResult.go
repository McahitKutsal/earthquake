package models

import "time"

type RequestResult struct {
	StatusCode int           `bson:"status_code"`
	Duration   time.Duration `bson:"duration"`
	Error      string        `bson:"error,omitempty"`
}
