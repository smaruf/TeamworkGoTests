package customerimporter

import (
	"fmt"
	"log"
	"os"
)

func CLI() {
	// Ask the user for the input file path
	fmt.Print("Enter the path to the input CSV file: ")
	var inputFile string
	_, err := fmt.Scanln(&inputFile)
	if err != nil || inputFile == "" {
		log.Fatalf("Invalid input file path: %v", err)
	}

	// Ask the user for the output file path
	fmt.Print("Enter the path to the output file (leave empty to print to console): ")
	var outputFile string
	_, err = fmt.Scanln(&outputFile)
	if err != nil {
		log.Fatalf("Invalid output file path: %v", err)
	}

	// Ask the user for the processing mode
	fmt.Println("Choose processing mode:")
	fmt.Println("1: Single-threaded")
	fmt.Println("2: Concurrent-streaming")
	var mode string
	fmt.Print("Enter mode (1 or 2): ")
	_, err = fmt.Scanln(&mode)
	if err != nil {
		log.Fatalf("Failed to read mode: %v", err)
	}

	// Open the input file
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer file.Close()

	// Choose the processing mode
	var domainCounts map[string]int
	switch mode {
	case "2":
		log.Println("Running in concurrent-streaming mode...")
		domainCounts, err = ProcessWithConcurrentStreaming(file)
		if err != nil {
			log.Fatalf("Error in concurrent-streaming processing: %v", err)
		}
		outputData := sortDomains(domainCounts)
		writeOutput(outputData, outputFile)
	default:
		log.Println("Running in single-threaded mode...")
		domainCounts, err = Process(inputFile, outputFile)
		if err != nil {
			log.Fatalf("Error in single-threaded processing: %v", err)
		}
	}

	// Handle output
	if outputFile == "" {
		// Print results to the console
		fmt.Println("Processing completed. Results:")
		for domain, count := range domainCounts {
			fmt.Printf("%s: %d\n", domain, count)
		}
	} else {
		// Write results to the output file
		log.Printf("Processing completed successfully. Output written to %s", outputFile)
	}
}
