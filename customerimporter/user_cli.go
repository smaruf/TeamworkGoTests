package customerimporter

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CLI() {
	// Get the input file path
	inputFile := getInputFilePath()

	// Get the output file path
	outputFile := getOutputFilePath()

	// Get the processing mode
	mode := getProcessingMode()

	// Open the input file
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Error: Unable to open input file '%s': %v", inputFile, err)
	}
	defer file.Close()

	// Process the file based on the chosen mode
	var domainCounts map[string]int
	switch mode {
	case "2":
		log.Println("Running in concurrent-streaming mode...")
		domainCounts, err = ProcessWithConcurrentStreaming(file)
		if err != nil {
			log.Fatalf("Error in concurrent-streaming processing: %v", err)
		}
	case "1":
		log.Println("Running in single-threaded mode...")
		domainCounts, err = Process(inputFile, outputFile)
		if err != nil {
			log.Fatalf("Error in single-threaded processing: %v", err)
		}
	default:
		// This should never happen due to prior validation
		log.Fatalf("Error: Invalid processing mode '%s'", mode)
	}

	// Handle output
	handleOutput(domainCounts, outputFile)
}

// getInputFilePath prompts the user for the input file path and validates it
func getInputFilePath() string {
	for {
		fmt.Print("Enter the path to the input CSV file: ")
		var inputFile string
		_, err := fmt.Scanln(&inputFile)
		if err != nil || strings.TrimSpace(inputFile) == "" {
			fmt.Println("Error: Invalid input. Please provide a valid file path.")
			continue
		}
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			fmt.Printf("Error: File '%s' does not exist. Please provide a valid file path.\n", inputFile)
			continue
		}
		return inputFile
	}
}

// getOutputFilePath prompts the user for the output file path and validates it
func getOutputFilePath() string {
	for {
		fmt.Print("Enter the path to the output file (or 'console' to print to terminal): ")
		var outputFile string
		_, err := fmt.Scanln(&outputFile)
		if err != nil || strings.TrimSpace(outputFile) == "" {
			fmt.Println("Error: Invalid input. Please provide a valid file path or 'console'.")
			continue
		}
		if strings.ToLower(outputFile) == "console" {
			return "console"
		}
		if err := validateOutputFilePath(outputFile); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		return outputFile
	}
}

// validateOutputFilePath checks if the output file path is valid and writable
func validateOutputFilePath(outputFile string) error {
	absPath, err := filepath.Abs(outputFile)
	if err != nil {
		return errors.New("unable to resolve absolute path for the output file")
	}
	dir := filepath.Dir(absPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' does not exist", dir)
	}
	testFile := filepath.Join(dir, ".test_writable")
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("directory '%s' is not writable", dir)
	}
	file.Close()
	os.Remove(testFile)
	return nil
}

// getProcessingMode prompts the user for the processing mode and validates it
func getProcessingMode() string {
	for {
		fmt.Println("Choose processing mode:")
		fmt.Println("1: Single-threaded")
		fmt.Println("2: Concurrent-streaming")
		fmt.Print("Enter mode (1 or 2): ")
		var mode string
		_, err := fmt.Scanln(&mode)
		if err != nil || (mode != "1" && mode != "2") {
			fmt.Println("Error: Invalid mode. Please enter '1' or '2'.")
			continue
		}
		return mode
	}
}

// handleOutput processes the output based on the user's choice
func handleOutput(domainCounts map[string]int, outputFile string) {
	if outputFile == "console" {
		// Print results to the console
		fmt.Println("Processing completed. Results:")
		for domain, count := range domainCounts {
			fmt.Printf("%s: %d\n", domain, count)
		}
	} else {
		// Write results to the output file
		err := writeOutput(domainCounts, outputFile)
		if err != nil {
			log.Fatalf("Error: Unable to write to output file '%s': %v", outputFile, err)
		}
		log.Printf("Processing completed successfully. Output written to '%s'", outputFile)
	}
}
