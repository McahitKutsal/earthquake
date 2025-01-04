package performance

import (
	"bytes"
	"context"
	"earthquake/internal/config"
	"earthquake/internal/database"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestResult struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Endpoint    string             `bson:"endpoint"`
	Concurrency int                `bson:"concurrency"`
	Requests    int                `bson:"requests"`
	Method      string             `bson:"method"`
	Headers     map[string]string  `bson:"headers"`
	Body        string             `bson:"body"`
	TestSummary TestSummary        `bson:"test_summary"`
	CreatedAt   time.Time          `bson:"created_at"`
}

type Result struct {
	StatusCode int           `json:"status_code"`
	Duration   time.Duration `json:"duration"`
	Error      string        `json:"error,omitempty"`
}

type TestSummary struct {
	TotalRequests int           `json:"total_requests"`
	Success       int           `json:"success"`
	Failures      int           `json:"failures"`
	SuccessRate   float64       `json:"success_rate"`
	FailureRate   float64       `json:"failure_rate"`
	TotalTime     time.Duration `json:"total_time"`
	AverageTime   time.Duration `json:"average_time"`
	StatusCodes   map[int]int   `json:"status_codes"`
	Results       []Result      `json:"results"`
}

func ExecuteTest(cfg config.Config) TestSummary {
	results := make(chan Result, cfg.Requests)
	var wg sync.WaitGroup

	startTime := time.Now()

	// Calculate base requests per goroutine and the remainder
	baseRequests := cfg.Requests / cfg.Concurrency
	remainder := cfg.Requests % cfg.Concurrency

	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)

		// Determine how many requests this goroutine should handle
		requestsForThisGoroutine := baseRequests
		if i < remainder { // Distribute the remainder
			requestsForThisGoroutine++
		}

		go func(requestCount int) {
			defer wg.Done()
			client := &http.Client{}
			for j := 0; j < requestCount; j++ {
				start := time.Now()
				req, err := http.NewRequest(cfg.Method, cfg.Endpoint, bytes.NewBuffer([]byte(cfg.Body)))
				if err != nil {
					results <- Result{Error: err.Error()}
					continue
				}
				for k, v := range cfg.Headers {
					req.Header.Set(k, v)
				}
				resp, err := client.Do(req)
				duration := time.Since(start)
				if err != nil {
					results <- Result{Error: err.Error(), Duration: duration}
					continue
				}
				results <- Result{StatusCode: resp.StatusCode, Duration: duration}
				resp.Body.Close()
			}
		}(requestsForThisGoroutine)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var success, failures int
	var totalDuration time.Duration
	statusCodes := make(map[int]int)
	var testResults []Result

	for result := range results {
		testResults = append(testResults, result)

		totalDuration += result.Duration
		statusCodes[result.StatusCode]++

		if result.Error != "" || result.StatusCode != 200 {
			failures++
		} else {
			success++
		}
	}

	totalTime := time.Since(startTime)
	averageTime := time.Duration(0)
	if success > 0 {
		averageTime = totalDuration / time.Duration(success)
	}

	successRate := (float64(success) / float64(cfg.Requests)) * 100
	failureRate := (float64(failures) / float64(cfg.Requests)) * 100

	return TestSummary{
		TotalRequests: cfg.Requests,
		Success:       success,
		Failures:      failures,
		SuccessRate:   successRate,
		FailureRate:   failureRate,
		TotalTime:     totalTime,
		AverageTime:   averageTime,
		StatusCodes:   statusCodes,
		Results:       testResults,
	}
}

func RunTest(cfg config.Config, id primitive.ObjectID) TestResult {
	summary := ExecuteTest(cfg)

	result := TestResult{
		ID:          id,
		Endpoint:    cfg.Endpoint,
		Concurrency: cfg.Concurrency,
		Requests:    cfg.Requests,
		Method:      cfg.Method,
		Headers:     cfg.Headers,
		Body:        cfg.Body,
		TestSummary: summary,
		CreatedAt:   time.Now(),
	}

	// Save to MongoDB
	collection := database.GetCollection("performance", "results")
	_, err := collection.InsertOne(context.Background(), result)
	if err != nil {
		panic(err) // Replace with better error handling
	}
	fmt.Println("done")

	return result
}
