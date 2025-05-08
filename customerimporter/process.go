package customerimporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

const (
	NoRecords  = 0
	MinRecords = 1000
	MaxRecords = 1000000
	ChunkSize  = 1000
	FieldCount = 5
)

// Regex pattern for validating email addresses
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Main Process Function
func Process(inputFileName string, outputFileName string) (map[string]int, error) {
	records, skipped, err := readCSV(inputFileName)
	if err != nil {
		return nil, fmt.Errorf("error processing CSV: %v", err)
	}

	domainCounts := countEmailDomains(records)
	sortedDomains := sortDomains(domainCounts)
	err = writeOutput(sortedDomains, outputFileName)
	if err != nil {
		return nil, fmt.Errorf("error writing output: %v", err)
	}

	log.Printf("Summary: Processed %d records, Skipped %d malformed rows", len(records), skipped)
	return domainCounts, nil
}

// CSV Reading and Validation
func readCSV(fileName string) ([]Record, int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, 0, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	records, skipped, err := parseCSVRecords(file)
	if err != nil {
		return nil, 0, err
	}

	if err := validateRecordCounts(records, MinRecords, MaxRecords); err != nil {
		return nil, 0, err
	}
	return records, skipped, nil
}

func parseCSVRecords(file *os.File) ([]Record, int, error) {
	reader := csv.NewReader(file)
	var records []Record
	var skipped int

	// Skip the first line (header row)
	_, err := reader.Read()
	if err == io.EOF {
		return nil, 0, fmt.Errorf("file is empty or contains only headers")
	}
	if err != nil {
		return nil, 0, fmt.Errorf("error reading header row: %v", err)
	}

	// Process the remaining rows
	rowNumber := 1
	for {
		fields, err := reader.Read()
		rowNumber++
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading row %d: %v", rowNumber, err)
			skipped++
			continue
		}
		if len(fields) < FieldCount {
			log.Printf("Skipping malformed row %d: %v", rowNumber, fields)
			skipped++
			continue
		}
		if !emailRegex.MatchString(fields[2]) {
			log.Printf("Skipping row %d due to invalid email: %s", rowNumber, fields[2])
			skipped++
			continue
		}
		record, err := createRecord(fields)
		if err != nil {
			log.Printf("Skipping row %d due to error: %v", rowNumber, err)
			skipped++
			continue
		}
		records = append(records, record)
	}
	return records, skipped, nil
}
