package customerimporter

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

const (
	NoRecords  = 0
	MinRecords = 1000
	MaxRecords = 1000000
	ChunkSize  = 1000
	FieldCount = 5
)

type Record struct {
	FirstName string
	LastName  string
	Email     string
	Gender    string
	IPAddress string
}

// Regex pattern for validating email domains
var domainRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`)

// Enhanced extractDomain with regex validation
func extractDomain(email string) string {
	at := strings.Index(email, "@")
	if at == -1 {
		log.Printf("Invalid email address (missing '@'): %s", email)
		return ""
	}
	domain := email[at+1:]
	if !domainRegex.MatchString(domain) {
		log.Printf("Invalid email domain: %s", domain)
		return ""
	}
	return domain
}

func createRecord(fields []string) (Record, error) {
	if len(fields) < FieldCount { // Ensure there are enough fields
		return Record{}, fmt.Errorf("invalid number of fields: expected %d, got %d", FieldCount, len(fields))
	}
	return Record{
		FirstName: fields[0],
		LastName:  fields[1],
		Email:     fields[2],
		Gender:    fields[3],
		IPAddress: fields[4],
	}, nil
}

func validateRecordCounts(records []Record, minRecords, maxRecords int) error {
	if len(records) == NoRecords {
		return fmt.Errorf("no records found in file")
	}
	if len(records) > maxRecords {
		return fmt.Errorf("too many records in file: %d", len(records))
	}
	if len(records) < minRecords {
		return fmt.Errorf("too few records in file: %d", len(records))
	}
	return nil
}

// Email Domain Processing
func countEmailDomains(records []Record) map[string]int {
	domainCounts := make(map[string]int, len(records)/2)

	for _, record := range records {
		domain := extractDomain(record.Email)
		if domain == "" {
			log.Printf("Skipping record with invalid email: %s", record.Email)
			continue
		}
		domainCounts[domain]++
	}

	return domainCounts
}

// Refactored sortDomains for clarity
func sortDomains(domainCounts map[string]int) []string {
	domains := make([]string, 0, len(domainCounts))
	for domain := range domainCounts {
		domains = append(domains, domain)
	}
	// Sort by count (descending), and alphabetically for ties
	sort.Slice(domains, func(i, j int) bool {
		if domainCounts[domains[i]] == domainCounts[domains[j]] {
			return domains[i] < domains[j] // Alphabetical order for ties
		}
		return domainCounts[domains[i]] > domainCounts[domains[j]] // Descending order by count
	})

	// Format domains with their counts
	for i, domain := range domains {
		domains[i] = fmt.Sprintf("%s: %d", domain, domainCounts[domain])
	}

	return domains
}

// Improved writeOutput with standardized newline
func writeOutput(sortedDomains []string, outputFileName string) error {
	// Create or overwrite the output file
	file, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	// Write sorted domains to the file
	for _, domain := range sortedDomains {
		if _, err := file.WriteString(fmt.Sprintf("%s\n", domain)); err != nil {
			return fmt.Errorf("error writing to output file: %v", err)
		}
	}

	return nil
}
