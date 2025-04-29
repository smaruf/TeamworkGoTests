// Package customerimporter reads from a CSV file and returns a sorted (data
// structure of your choice) of email domains along with the number of customers
// with e-mail addresses for each domain. This should be able to be ran from the
// CLI and output the sorted domains to the terminal or to a file. Any errors
// should be logged (or handled). Performance matters (this is only ~3k lines,
// but could be 1m lines or run on a small machine).
package customerimporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

// Main Process Function
func Process(inputFileName string, outputFileName string) (map[string]int, error) {
	records, err := readCSV(inputFileName)
	if err != nil {
		return nil, fmt.Errorf("error processing CSV: %v", err)
	}

	domainCounts := countEmailDomains(records)
	sortedDomains := sortDomains(domainCounts)
	err = writeOutput(sortedDomains, outputFileName)
	if err != nil {
		return nil, fmt.Errorf("error writing output: %v", err)
	}

	return domainCounts, nil
}

// CSV Reading and Validation
func readCSV(fileName string) ([]Record, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	records, err := parseCSVRecords(file)
	if err != nil {
		return nil, err
	}

	if err := validateRecordCounts(records, MinRecords, MaxRecords); err != nil {
		return nil, err
	}
	return records, nil
}

func parseCSVRecords(file *os.File) ([]Record, error) {
	reader := csv.NewReader(file)
	var records []Record

	// Skip the first line (header row)
	_, err := reader.Read()
	if err == io.EOF {
		return nil, fmt.Errorf("file is empty or contains only headers")
	}
	if err != nil {
		return nil, fmt.Errorf("error reading header row: %v", err)
	}

	// Process the remaining rows
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %v", err)
		}
		if len(fields) < FieldCount {
			log.Printf("Skipping malformed row: %v", fields) // Log malformed rows
			continue
		}
		record, err := createRecord(fields)
		if err != nil {
			log.Printf("Skipping row due to error: %v", err) // Log the error and skip the row
			continue
		}
		records = append(records, record)
	}
	return records, nil
}
