package iostream_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/baepo-cloud/baepo-cli/pkg/iostream"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Country string `json:"country"`
}

func TestMessagePlainText(t *testing.T) {
	var stdout bytes.Buffer
	stream := iostream.New(false)
	stream.Stdout = &stdout

	stream.Message("Hello %s", "World")

	expected := "Hello World\n"
	fmt.Printf("TestMessagePlainText: %q\n", stdout.String())
	if stdout.String() != expected {
		t.Errorf("Expected %q, got %q", expected, stdout.String())
	}
}

func TestMessageJSON(t *testing.T) {
	var stdout bytes.Buffer
	stream := iostream.New(true)
	stream.Stdout = &stdout

	stream.Message("Hello %s", "World")

	fmt.Printf("TestMessageJSON: %s", stdout.String())
	var result map[string]string
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["message"] != "Hello World" {
		t.Errorf("Expected message to be %q, got %q", "Hello World", result["message"])
	}
}

func TestErrorPlainText(t *testing.T) {
	var stderr bytes.Buffer
	stream := iostream.New(false)
	stream.Stderr = &stderr

	stream.Error("Failed with code %d", 404)

	expected := "Failed with code 404\n"
	fmt.Printf("TestErrorPlainText: %q\n", stderr.String())
	if stderr.String() != expected {
		t.Errorf("Expected %q, got %q", expected, stderr.String())
	}
}

func TestErrorWithDetailsPlainText(t *testing.T) {
	var stderr bytes.Buffer
	stream := iostream.New(false)
	stream.Stderr = &stderr

	opts := iostream.ErrorOptions{
		Error:   "Failed with code %d",
		Details: "Resource users not found",
		Code:    "NOT_FOUND",
	}
	stream.ErrorWithDetails(opts, 404)

	expected := "Failed with code 404 (Resource users not found) [code: NOT_FOUND]\n"
	fmt.Printf("TestErrorWithDetailsPlainText: %q\n", stderr.String())
	if stderr.String() != expected {
		t.Errorf("Expected %q, got %q", expected, stderr.String())
	}
}

func TestErrorWithDetailsJSON(t *testing.T) {
	var stderr bytes.Buffer
	stream := iostream.New(true)
	stream.Stderr = &stderr

	opts := iostream.ErrorOptions{
		Error:   "Failed with code %d",
		Details: "Resource users not found",
		Code:    "NOT_FOUND",
	}
	stream.ErrorWithDetails(opts, 404)

	fmt.Printf("TestErrorWithDetailsJSON: %s", stderr.String())
	var result map[string]string
	if err := json.Unmarshal(stderr.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["error"] != "Failed with code 404" {
		t.Errorf("Expected error to be %q, got %q", "Failed with code %d", result["error"])
	}

	if result["details"] != "Resource users not found" {
		t.Errorf("Expected details to be %q, got %q", "Resource users not found", result["details"])
	}

	if result["code"] != "NOT_FOUND" {
		t.Errorf("Expected code to be %q, got %q", "NOT_FOUND", result["code"])
	}
}

func TestObjectPlainText(t *testing.T) {
	var stdout bytes.Buffer
	stream := iostream.New(false)
	stream.Stdout = &stdout

	person := Person{
		Name:    "John Doe",
		Age:     30,
		Country: "USA",
	}

	stream.Object(person)

	output := stdout.String()
	fmt.Printf("TestObjectPlainText:\n%s", output)

	// Check that all fields are present
	if !strings.Contains(output, "name:") || !strings.Contains(output, "John Doe") {
		t.Error("Output missing name field")
	}

	if !strings.Contains(output, "age:") || !strings.Contains(output, "30") {
		t.Error("Output missing age field")
	}

	if !strings.Contains(output, "country:") || !strings.Contains(output, "USA") {
		t.Error("Output missing country field")
	}
}

func TestObjectJSON(t *testing.T) {
	var stdout bytes.Buffer
	stream := iostream.New(true)
	stream.Stdout = &stdout

	person := Person{
		Name:    "John Doe",
		Age:     30,
		Country: "USA",
	}

	stream.Object(person)

	fmt.Printf("TestObjectJSON: %s", stdout.String())
	var result Person
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result.Name != "John Doe" {
		t.Errorf("Expected Name to be %q, got %q", "John Doe", result.Name)
	}

	if result.Age != 30 {
		t.Errorf("Expected Age to be %d, got %d", 30, result.Age)
	}

	if result.Country != "USA" {
		t.Errorf("Expected Country to be %q, got %q", "USA", result.Country)
	}
}

func TestArrayPlainText(t *testing.T) {
	var stdout bytes.Buffer
	stream := iostream.New(false)
	stream.Stdout = &stdout

	people := []Person{
		{Name: "John Doe", Age: 30, Country: "USA"},
		{Name: "Jane Smith", Age: 28, Country: "Canada"},
	}

	stream.Array(people, func(key string, value any) (string, string) {
		switch key {
		case "name":
			return "NAME", value.(string)
		case "age":
			return "AGE", fmt.Sprintf("%d", value)
		case "country":
			return "COUNTRY", fmt.Sprintf("%v", value)
		default:
			return "", ""
		}
	})

	output := stdout.String()
	fmt.Printf("TestArrayPlainText:\n%s", output)
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines (header + 2 data rows), got %d", len(lines))
	}

	// Check header
	if !strings.Contains(lines[0], "NAME") ||
		!strings.Contains(lines[0], "AGE") ||
		!strings.Contains(lines[0], "COUNTRY") {
		t.Errorf("Header missing expected columns: %s", lines[0])
	}

	// Check data rows
	if !strings.Contains(lines[1], "John Doe") ||
		!strings.Contains(lines[1], "30") ||
		!strings.Contains(lines[1], "USA") {
		t.Errorf("First data row missing expected values: %s", lines[1])
	}

	if !strings.Contains(lines[2], "Jane Smith") ||
		!strings.Contains(lines[2], "28") ||
		!strings.Contains(lines[2], "Canada") {
		t.Errorf("Second data row missing expected values: %s", lines[2])
	}
}

func TestArrayJSON(t *testing.T) {
	var stdout bytes.Buffer
	stream := iostream.New(true)
	stream.Stdout = &stdout

	people := []Person{
		{Name: "John Doe", Age: 30, Country: "USA"},
		{Name: "Jane Smith", Age: 28, Country: "Canada"},
	}

	stream.Array(people, func(key string, value any) (string, string) {
		// This function is not used for JSON output
		return key, fmt.Sprintf("%v", value)
	})

	fmt.Printf("TestArrayJSON: %s", stdout.String())
	var result []Person
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(result))
	}

	if result[0].Name != "John Doe" || result[0].Age != 30 || result[0].Country != "USA" {
		t.Errorf("First item doesn't match: %+v", result[0])
	}

	if result[1].Name != "Jane Smith" || result[1].Age != 28 || result[1].Country != "Canada" {
		t.Errorf("Second item doesn't match: %+v", result[1])
	}
}
