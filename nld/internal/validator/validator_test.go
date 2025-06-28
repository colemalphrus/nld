package validator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateDocument(t *testing.T) {
	// Create a validator
	v := New()

	// Get the current directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Go up to the project root
	projectRoot := filepath.Join(wd, "..", "..")

	// Define test cases
	testCases := []struct {
		name           string
		docPath        string
		schemaPath     string
		expectValid    bool
		expectErrorMsg string
	}{
		{
			name:        "Valid Contract",
			docPath:     filepath.Join(projectRoot, "examples", "valid-contract.json"),
			schemaPath:  filepath.Join(projectRoot, "schemas", "document-v1.json"),
			expectValid: true,
		},
		{
			name:        "Valid Receipt",
			docPath:     filepath.Join(projectRoot, "examples", "valid-receipt.json"),
			schemaPath:  filepath.Join(projectRoot, "schemas", "document-v1.json"),
			expectValid: true,
		},
		{
			name:        "Invalid Missing Fields",
			docPath:     filepath.Join(projectRoot, "examples", "invalid-missing-fields.json"),
			schemaPath:  filepath.Join(projectRoot, "schemas", "document-v1.json"),
			expectValid: false,
			expectErrorMsg: "missing properties",
		},
		{
			name:        "Invalid Type",
			docPath:     filepath.Join(projectRoot, "examples", "invalid-type.json"),
			schemaPath:  filepath.Join(projectRoot, "schemas", "document-v1.json"),
			expectValid: false,
			expectErrorMsg: "value must be one of",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Validate the document
			result, err := v.ValidateDocument(tc.docPath, tc.schemaPath)
			if err != nil {
				t.Fatalf("Validation failed with error: %v", err)
			}

			// Check if the result matches expectations
			if result.Valid != tc.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tc.expectValid, result.Valid)
			}

			// If we expect an error message, check that it's present
			if tc.expectErrorMsg != "" && !result.Valid {
				found := false
				for _, err := range result.Errors {
					if strings.Contains(err.Message, tc.expectErrorMsg) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error message containing '%s', but didn't find it in errors: %v", 
						tc.expectErrorMsg, result.Errors)
				}
			}
		})
	}
}

func TestValidateBytes(t *testing.T) {
	// Create a validator
	v := New()

	// Get the current directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Go up to the project root
	projectRoot := filepath.Join(wd, "..", "..")

	// Load the schema
	schemaPath := filepath.Join(projectRoot, "schemas", "document-v1.json")
	schema, err := v.LoadSchema(schemaPath)
	if err != nil {
		t.Fatalf("Failed to load schema: %v", err)
	}

	// Test valid JSON
	validJSON := []byte(`{
		"metadata": {
			"version": "1.0.0",
			"type": "contract",
			"created": "2025-06-27T12:00:00Z",
			"title": "Test Contract"
		},
		"content": {
			"sections": [
				{
					"id": "test",
					"title": "Test Section",
					"content": "Test content"
				}
			]
		}
	}`)

	result, err := v.ValidateBytes(validJSON, schema)
	if err != nil {
		t.Fatalf("Validation failed with error: %v", err)
	}
	if !result.Valid {
		t.Errorf("Expected valid document, got invalid: %v", result.Errors)
	}

	// Test invalid JSON
	invalidJSON := []byte(`{
		"metadata": {
			"version": "1.0.0",
			"type": "invalid-type",
			"created": "2025-06-27T12:00:00Z",
			"title": "Test Contract"
		},
		"content": {
			"sections": [
				{
					"id": "test",
					"title": "Test Section",
					"content": "Test content"
				}
			]
		}
	}`)

	result, err = v.ValidateBytes(invalidJSON, schema)
	if err != nil {
		t.Fatalf("Validation failed with error: %v", err)
	}
	if result.Valid {
		t.Errorf("Expected invalid document, got valid")
	}

	// Test malformed JSON
	malformedJSON := []byte(`{
		"metadata": {
			"version": "1.0.0",
			"type": "contract",
			"created": "2025-06-27T12:00:00Z",
			"title": "Test Contract"
		},
		"content": {
			"sections": [
				{
					"id": "test",
					"title": "Test Section",
					"content": "Test content"
				}
			]
		`)

	result, err = v.ValidateBytes(malformedJSON, schema)
	if err != nil {
		t.Fatalf("Validation failed with error: %v", err)
	}
	if result.Valid {
		t.Errorf("Expected invalid document due to malformed JSON, got valid")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}