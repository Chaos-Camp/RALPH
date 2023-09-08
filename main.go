package main

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

func sanitizeURL(url string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(url, "_")
}

func makeRequest(url string, results *[]map[string]interface{}) {
	startTime := time.Now()
	resp, err := http.Get(url)
	var status string
	if err != nil {
		status = fmt.Sprintf("Error: %v", err) // Enhanced error message
		fmt.Println(status)
	} else {
		status = resp.Status
		resp.Body.Close()
	}
	endTime := time.Now()
	result := map[string]interface{}{
		"url":           url,
		"status":        status,
		"response_time": endTime.Sub(startTime).Seconds(),
		"timestamp":     startTime.Format("2006-01-02 15:04:05"),
	}
	*results = append(*results, result)
}

func simpleTest(url string, iterations int, results *[]map[string]interface{}) {
	for i := 0; i < iterations; i++ {
		makeRequest(url, results)
	}
}

func stressTest(url string, concurrentRequests int, results *[]map[string]interface{}) {
	var wg sync.WaitGroup
	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			makeRequest(url, results)
		}()
	}
	wg.Wait()
}

func spikeTest(url string, spikes int, spikeInterval float64, results *[]map[string]interface{}) {
	for i := 0; i < spikes; i++ {
		stressTest(url, 5, results) // assuming 5 concurrent requests as default
		time.Sleep(time.Duration(spikeInterval) * time.Second)
	}
}

func enduranceTest(url string, duration float64, results *[]map[string]interface{}) {
	endTime := time.Now().Add(time.Duration(duration) * time.Second)
	for time.Now().Before(endTime) {
		makeRequest(url, results)
	}
}

func rampUpTest(url string, maxUsers int, rampUpPeriod int, results *[]map[string]interface{}) {
	for i := 1; i <= maxUsers; i++ {
		stressTest(url, i, results)
		time.Sleep(time.Duration(rampUpPeriod/maxUsers) * time.Second)
	}
}

func saveToCSV(results []map[string]interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"url", "status", "response_time", "timestamp"}
	writer.Write(header)
	for _, result := range results {
		writer.Write([]string{
			result["url"].(string),
			result["status"].(string),
			fmt.Sprintf("%f", result["response_time"].(float64)),
			result["timestamp"].(string),
		})
	}
	return nil
}

func uploadToGCS(filename string, bucketName string) error {
	ctx := context.Background()
	saPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if saPath == "" {
		return fmt.Errorf("the environment variable GOOGLE_APPLICATION_CREDENTIALS must be set")
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(filepath.Base(filename))

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return fmt.Errorf("failed to write file to bucket: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close bucket writer: %v", err)
	}
	return nil
}

func main() {
	var url string
	var maxUsers, rampUpPeriod, iterations, concurrentRequests, spikes int
	var spikeInterval, duration float64
	var bucketName, testType string
	var noCloud bool

	flag.StringVar(&url, "url", "https://www.example.com", "URL to test")
	flag.StringVar(&testType, "type", "ramp_up", "Type of load test: simple, stress, spike, endurance, or ramp_up")
	flag.IntVar(&iterations, "iterations", 1, "Number of iterations for simple test")
	flag.IntVar(&concurrentRequests, "concurrentRequests", 5, "Number of concurrent requests for stress test")
	flag.IntVar(&spikes, "spikes", 3, "Number of spikes for spike test")
	flag.Float64Var(&spikeInterval, "spikeInterval", 1.0, "Interval between spikes in seconds")
	flag.Float64Var(&duration, "duration", 300.0, "Duration of the endurance test in seconds")
	flag.IntVar(&maxUsers, "maxUsers", 10, "Maximum number of users for ramp up test")
	flag.IntVar(&rampUpPeriod, "rampUpPeriod", 10, "Ramp up period in seconds")
	flag.StringVar(&bucketName, "bucket", "my-bucket", "Google Cloud Storage bucket name")
	flag.BoolVar(&noCloud, "nocloud", false, "If set, will not upload to Google Cloud Storage and only save the CSV locally.")

	flag.Parse()

	results := make([]map[string]interface{}, 0)

	switch testType {
	case "simple":
		simpleTest(url, iterations, &results)
	case "stress":
		stressTest(url, concurrentRequests, &results)
	case "spike":
		spikeTest(url, spikes, spikeInterval, &results)
	case "endurance":
		enduranceTest(url, duration, &results)
	case "ramp_up":
		rampUpTest(url, maxUsers, rampUpPeriod, &results)
	default:
		fmt.Println("Unknown test type:", testType)
		return
	}

	sanitizedURL := sanitizeURL(url)
	filename := sanitizedURL + "_results.csv"

	err := saveToCSV(results, filename)
	if err != nil {
		fmt.Println("Error saving results to CSV:", err)
		return
	}

	if !noCloud {
		err = uploadToGCS(filename, bucketName)
		if err != nil {
			fmt.Println("Error uploading to GCS:", err)
			return
		}

		// Deleting the local CSV after successful upload to GCS
		err = os.Remove(filename)
		if err != nil {
			fmt.Println("Error deleting local CSV file:", err)
		}
	}

	fmt.Println("Done!")
}
