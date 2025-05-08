package customerimporter

import (
	"os"
	"reflect"
	"sync"
	"testing"
)

func TestProcessWithConcurrentStreaming(t *testing.T) {
	// Create a temporary CSV file
	file, err := os.CreateTemp("", "input_test.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write sample data to the file
	data := `FirstName,LastName,Email,Phone,Address,City,State,Zip,Country,Company,JobTitle,Website,Notes
John,Doe,john.doe@example.com,1234567890,123 Main St,City,State,12345,Country,Company,Job Title,www.example.com,Notes
Jane,Smith,jane.smith@example.com,9876543210,456 Elm St,City,State,67890,Country,Company,Job Title,www.example.com,Notes
Invalid,User,invalid-email,0000000000,789 Pine St,City,State,11111,Country,Company,Job Title,www.example.com,Notes
`
	if _, err := file.WriteString(data); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	// Open the file for reading
	file.Seek(0, 0)

	// Run the ProcessWithConcurrentStreaming function
	result, err := ProcessWithConcurrentStreaming(file)
	if err != nil {
		t.Errorf("ProcessWithConcurrentStreaming() returned an error: %v", err)
	}

	// Verify the results
	expected := map[string]int{
		"example.com": 2,
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ProcessWithConcurrentStreaming() = %v; want %v", result, expected)
	}
}

func TestProcessWithConcurrentStreaming_EmptyFile(t *testing.T) {
	// Create an empty temporary CSV file
	file, err := os.CreateTemp("", "empty_test.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Run the ProcessWithConcurrentStreaming function
	result, err := ProcessWithConcurrentStreaming(file)
	if err != nil {
		t.Errorf("ProcessWithConcurrentStreaming() returned an error for empty file: %v", err)
	}

	// Verify the results
	expected := map[string]int{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ProcessWithConcurrentStreaming() = %v; want %v", result, expected)
	}
}

func TestProcessChunk(t *testing.T) {
	chunk := [][]string{
		{"John", "Doe", "john.doe@example.com", "1234567890"},
		{"Jane", "Smith", "jane.smith@example.com", "9876543210"},
		{"Invalid", "User", "invalid-email", "0000000000"},
	}

	ch := make(chan map[string]int, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go processChunk(chunk, ch, &wg)

	wg.Wait()
	close(ch)

	// Verify the results
	result := <-ch
	expected := map[string]int{
		"example.com": 2,
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("processChunk() = %v; want %v", result, expected)
	}
}

func TestCollectResults(t *testing.T) {
	ch := make(chan map[string]int, 2)
	var domainCounts sync.Map

	// Simulate two chunks of results
	ch <- map[string]int{"example.com": 2}
	ch <- map[string]int{"another.com": 1}
	close(ch)

	collectResults(ch, &domainCounts)

	// Convert sync.Map to regular map for verification
	finalCounts := make(map[string]int)
	domainCounts.Range(func(key, value interface{}) bool {
		finalCounts[key.(string)] = value.(int)
		return true
	})

	// Verify the results
	expected := map[string]int{
		"example.com": 2,
		"another.com": 1,
	}
	if !reflect.DeepEqual(finalCounts, expected) {
		t.Errorf("collectResults() = %v; want %v", finalCounts, expected)
	}
}

func TestConvertSyncMapToRegularMap(t *testing.T) {
	var domainCounts sync.Map
	domainCounts.Store("example.com", 2)
	domainCounts.Store("another.com", 1)

	result := convertSyncMapToRegularMap(&domainCounts)

	// Verify the results
	expected := map[string]int{
		"example.com": 2,
		"another.com": 1,
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("convertSyncMapToRegularMap() = %v; want %v", result, expected)
	}
}
