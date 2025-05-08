package customerimporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sync"
)

const ChunkSize = 1000

// Regex pattern for validating email domains
var domainRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`)

// ProcessWithConcurrentStreaming processes a CSV file concurrently with streaming
func ProcessWithConcurrentStreaming(file *os.File) (map[string]int, error) {
	reader := csv.NewReader(file)
	ch := make(chan map[string]int)
	var domainCounts sync.Map
	var wg sync.WaitGroup

	// Process chunks and collect results concurrently
	processChunks(reader, ch, &wg)
	collectResults(ch, &domainCounts)

	// Convert sync.Map to regular map and return
	return convertSyncMapToRegularMap(&domainCounts), nil
}

// processChunks reads CSV records in chunks and processes them concurrently
func processChunks(reader *csv.Reader, ch chan map[string]int, wg *sync.WaitGroup) {
	var chunk [][]string
	rowNumber := 1

	for {
		fields, err := reader.Read()
		rowNumber++
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading row %d: %v", rowNumber, err)
			continue
		}
		chunk = append(chunk, fields)
		if len(chunk) >= ChunkSize {
			wg.Add(1)
			go processChunk(chunk, ch, wg)
			chunk = nil
		}
	}

	// Process any remaining records in the last chunk
	if len(chunk) > 0 {
		wg.Add(1)
		go processChunk(chunk, ch, wg)
	}

	// Close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(ch)
	}()
}

// processChunk processes a single chunk of records and sends results to the channel
func processChunk(chunk [][]string, ch chan map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	localCounts := make(map[string]int)

	for _, fields := range chunk {
		// Validate row length
		if len(fields) < 3 {
			log.Printf("Skipping malformed row: %v", fields)
			continue
		}

		// Extract and validate email domain
		domain := extractDomain(fields[2])
		if domain == "" || !domainRegex.MatchString(domain) {
			log.Printf("Skipping row due to invalid domain: %s", fields[2])
			continue
		}

		// Increment domain count
		localCounts[domain]++
	}

	// Send local counts to the channel
	ch <- localCounts
}

// collectResults aggregates results from multiple goroutines into a sync.Map
func collectResults(ch chan map[string]int, domainCounts *sync.Map) {
	for localCounts := range ch {
		for domain, count := range localCounts {
			// Atomically update the sync.Map
			actual, _ := domainCounts.LoadOrStore(domain, count)
			if actual != nil {
				domainCounts.Store(domain, actual.(int)+count)
			}
		}
	}
}

// convertSyncMapToRegularMap converts sync.Map to a regular map
func convertSyncMapToRegularMap(domainCounts *sync.Map) map[string]int {
	finalCounts := make(map[string]int)
	domainCounts.Range(func(key, value interface{}) bool {
		finalCounts[key.(string)] = value.(int)
		return true
	})
	return finalCounts
}

// extractDomain extracts the domain from an email address
func extractDomain(email string) string {
	at := strings.Index(email, "@")
	if at == -1 {
		return ""
	}
	return email[at+1:]
}
