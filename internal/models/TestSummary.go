package models

type TestSummary struct {
	TotalRequests int             `bson:"total_requests"`
	Success       int             `bson:"success"`
	Failures      int             `bson:"failures"`
	SuccessRate   float64         `bson:"success_rate"`
	FailureRate   float64         `bson:"failure_rate"`
	TotalTime     string          `bson:"total_time"`
	AverageTime   string          `bson:"average_time"`
	StatusCodes   map[int]int     `bson:"status_codes"`
	Results       []RequestResult `bson:"results"`
}
