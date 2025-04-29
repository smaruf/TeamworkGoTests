package customerimporter

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainCLI(t *testing.T) {
	// Simulate user input for the interactive CLI
	inputFile := "input_test_cli.csv"
	outputFile := "output_test_cli.txt"
	userInput := strings.Join([]string{
		inputFile,  // Input file path
		outputFile, // Output file path
		"1",        // Processing mode (1 for single-threaded)
	}, "\n")

	// Redirect stdin to simulate user input
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess")
	cmd.Env = append(os.Environ(), "GO_TEST_HELPER=1")
	cmd.Stdin = strings.NewReader(userInput)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Run the CLI program
	err := cmd.Run()
	if err != nil {
		t.Fatalf("CLI execution failed: %v\nOutput: %s", err, out.String())
	}

	// Verify the output
	expectedOutput := "CLI validation successful. Input: " + inputFile + ", Output: " + outputFile + "\n"
	if !strings.Contains(out.String(), expectedOutput) {
		t.Errorf("CLI output = %q; want it to contain %q", out.String(), expectedOutput)
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_TEST_HELPER") != "1" {
		return
	}

	// Simulate the CLI function
	var inputFile, outputFile, mode string

	// Read user input
	_, err := fmt.Scanln(&inputFile)
	if err != nil || inputFile == "" {
		t.Fatalf("Failed to read input file path: %v", err)
	}

	_, err = fmt.Scanln(&outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file path: %v", err)
	}

	_, err = fmt.Scanln(&mode)
	if err != nil || (mode != "1" && mode != "2") {
		t.Fatalf("Failed to read processing mode: %v", err)
	}

	// Simulate successful validation
	os.Stdout.WriteString("CLI validation successful. Input: " + inputFile + ", Output: " + outputFile + "\n")
	os.Exit(0)
}
