package customerimporter

import (
	"reflect"
	"testing"
)

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		email    string
		expected string
	}{
		{"user@example.com", "example.com"},
		{"invalid-email", ""},
		{"", ""},
		{"user@subdomain.example.com", "subdomain.example.com"},
	}

	for _, test := range tests {
		result := extractDomain(test.email)
		if result != test.expected {
			t.Errorf("extractDomain(%q) = %q; want %q", test.email, result, test.expected)
		}
	}
}

func TestCreateRecord(t *testing.T) {
	validFields := []string{"John", "Doe", "john.doe@example.com", "1234567890", "123 Main St", "City", "State", "12345", "Country", "Company", "Job Title", "www.example.com", "Notes"}
	invalidFields := []string{"John", "Doe", "john.doe@example.com"} // Less than FieldCount

	// Test valid record creation
	record, err := createRecord(validFields)
	if err != nil {
		t.Errorf("createRecord(validFields) returned an error: %v", err)
	}
	if record.Email != "john.doe@example.com" {
		t.Errorf("createRecord(validFields) = %v; want Email = %q", record, "john.doe@example.com")
	}

	// Test invalid record creation
	_, err = createRecord(invalidFields)
	if err == nil {
		t.Errorf("createRecord(invalidFields) did not return an error")
	}
}

func TestValidateRecordCounts(t *testing.T) {
	records := []Record{
		{Email: "user1@example.com"},
		{Email: "user2@example.com"},
	}

	// Test valid record count
	err := validateRecordCounts(records, 1, 10)
	if err != nil {
		t.Errorf("validateRecordCounts(records, 1, 10) returned an error: %v", err)
	}

	// Test too few records
	err = validateRecordCounts(records, 3, 10)
	if err == nil {
		t.Errorf("validateRecordCounts(records, 3, 10) did not return an error")
	}

	// Test too many records
	err = validateRecordCounts(records, 1, 1)
	if err == nil {
		t.Errorf("validateRecordCounts(records, 1, 1) did not return an error")
	}
}

func TestCountEmailDomains(t *testing.T) {
	records := []Record{
		{Email: "user1@example.com"},
		{Email: "user2@example.com"},
		{Email: "user3@another.com"},
		{Email: "invalid-email"},
	}

	expected := map[string]int{
		"example.com": 2,
		"another.com": 1,
	}

	result := countEmailDomains(records)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("countEmailDomains(records) = %v; want %v", result, expected)
	}
}

func TestSortDomains(t *testing.T) {
	domainCounts := map[string]int{
		"example.com": 3,
		"another.com": 3,
		"test.com":    1,
	}

	expected := []string{"another.com: 3", "example.com: 3", "test.com: 1"} // Alphabetical order for ties
	result := sortDomains(domainCounts)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("sortDomains(domainCounts) = %v; want %v", result, expected)
	}
}
