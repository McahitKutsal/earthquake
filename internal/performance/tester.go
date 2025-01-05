package performance

import (
	"bytes"
	"context"
	"earthquake/internal/database"
	"earthquake/internal/models"
	"earthquake/pkg/logger"
	"earthquake/pkg/utils"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ExecuteTest(testRequest models.TestRequest) models.TestSummary {

	results := make(chan models.RequestResult, testRequest.Requests)
	var wg sync.WaitGroup

	startTime := time.Now()

	if testRequest.Requests < testRequest.Concurrency {
		testRequest.Concurrency = testRequest.Requests
	}
	// Calculate base requests per goroutine and the remainder
	baseRequests := testRequest.Requests / testRequest.Concurrency
	remainder := testRequest.Requests % testRequest.Concurrency

	logger.LogInfo(fmt.Sprintf("Test Started: Requests, %d Threads, %d", testRequest.Requests, testRequest.Concurrency))

	for i := 0; i < testRequest.Concurrency; i++ {
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
				req, err := http.NewRequest(testRequest.Method, testRequest.Endpoint, bytes.NewBuffer([]byte(testRequest.Body)))
				if err != nil {
					results <- models.RequestResult{Error: err.Error()}
					continue
				}
				for k, v := range testRequest.Headers {
					req.Header.Set(k, v)
				}
				resp, err := client.Do(req)
				duration := time.Since(start)
				if err != nil {
					results <- models.RequestResult{Error: err.Error(), Duration: duration}
					continue
				}
				results <- models.RequestResult{StatusCode: resp.StatusCode, Duration: duration}
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
	var testResults []models.RequestResult

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

	successRate := (float64(success) / float64(testRequest.Requests)) * 100
	failureRate := (float64(failures) / float64(testRequest.Requests)) * 100

	return models.TestSummary{
		TotalRequests: testRequest.Requests,
		Success:       success,
		Failures:      failures,
		SuccessRate:   successRate,
		FailureRate:   failureRate,
		TotalTime:     utils.FormatDuration(totalTime),   // Total time formatted as "hh:mm:ss"
		AverageTime:   utils.FormatDuration(averageTime), // Average time formatted as "hh:mm:ss"
		StatusCodes:   statusCodes,
		Results:       testResults,
	}
}

func RunTest(testRequest models.TestRequest, id primitive.ObjectID) models.TestResult {
	summary := ExecuteTest(testRequest)

	result := models.TestResult{
		ID:            primitive.NewObjectID(),
		TestRequestID: id,
		Endpoint:      testRequest.Endpoint,
		Concurrency:   testRequest.Concurrency,
		Requests:      testRequest.Requests,
		Method:        testRequest.Method,
		Headers:       testRequest.Headers,
		Body:          testRequest.Body,
		TestSummary:   summary,
		CreatedAt:     time.Now(),
	}

	// Save to MongoDB
	collection := database.GetCollection("performance", "results")
	_, err := collection.InsertOne(context.Background(), result)
	if err != nil {
		panic(err) // Replace with better error handling
	}

	logger.LogInfo(fmt.Sprintf("Done %s", id.Hex()))

	return result
}
