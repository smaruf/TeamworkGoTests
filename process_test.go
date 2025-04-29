package customerimporter

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestProcess(t *testing.T) {
	// Use input_test.csv as the base testing file
	inputFile := "input_test.csv"
	outputFile := "output_test.txt"
	defer os.Remove(outputFile) // Clean up the output file after the test

	// Run the Process function
	_, err := Process(inputFile, outputFile)
	if err != nil {
		t.Errorf("Process() returned an error: %v", err)
	}

	// Verify the output file contents
	outputData, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedOutput := "example.com"
	if !strings.Contains(string(outputData), expectedOutput) {
		t.Errorf("Output file does not contain expected output. Got = %q; want it to contain %q", string(outputData), expectedOutput)
	}
}

func TestReadCSV(t *testing.T) {
	// Use input_test.csv as the base testing file
	inputFile := "input_test.csv"

	// Run the readCSV function
	records, err := readCSV(inputFile)
	if err != nil {
		t.Errorf("readCSV() returned an error: %v", err)
	}

	// Verify the records
	expected := []Record{
		{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Gender: "male", IPAddress: "205.102.45.112"},
		{FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com", Gender: "female", IPAddress: "102.205.22.123"},
		{FirstName: "Bonnie", LastName: "Ortiz", Email: "bortiz1@cyberchimps.com", Gender: "Female", IPAddress: "197.54.209.129"},
		{FirstName: "Dennis", LastName: "Henry", Email: "dhenry2@hubpages.com", Gender: "Male", IPAddress: "155.75.186.217"},
	}
	if len(records) < len(expected) {
		t.Fatalf("readCSV() returned %d records; want at least %d", len(records), len(expected))
	}
	if !reflect.DeepEqual(records[:len(expected)], expected) {
		t.Errorf("readCSV() = %v; want %v", records[:len(expected)], expected)
	}
}

func TestParseCSVRecords(t *testing.T) {
	// Open input_test.csv for reading
	file, err := os.Open("input_test.csv")
	if err != nil {
		t.Fatalf("Failed to open input_test.csv: %v", err)
	}
	defer file.Close()

	// Run the parseCSVRecords function
	records, err := parseCSVRecords(file)
	if err != nil {
		t.Errorf("parseCSVRecords() returned an error: %v", err)
	}

	// Verify the records
	expectedEmails := []string{
		"john.doe@example.com",
		"jane.smith@example.com",
		"bortiz1@cyberchimps.com",
		"dhenry2@hubpages.com",
	}
	if len(records) < len(expectedEmails) {
		t.Errorf("parseCSVRecords() returned %d records; want at least %d", len(records), len(expectedEmails))
	}
	for i, email := range expectedEmails {
		if records[i].Email != email {
			t.Errorf("parseCSVRecords() returned record with Email = %q; want %q", records[i].Email, email)
		}
	}
}

func TestWriteOutput(t *testing.T) {
	// Create a temporary output file
	outputFile := "output_test.txt"
	defer os.Remove(outputFile) // Clean up the output file after the test

	// Run the writeOutput function
	sortedDomains := []string{"example.com", "another.com"}
	err := writeOutput(sortedDomains, outputFile)
	if err != nil {
		t.Errorf("writeOutput() returned an error: %v", err)
	}

	// Verify the file contents
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expected := "example.com\n\ranother.com\n\r"
	if string(data) != expected {
		t.Errorf("writeOutput() = %q; want %q", string(data), expected)
	}
}
