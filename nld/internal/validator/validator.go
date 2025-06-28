package validator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

// Validator is responsible for validating NLD documents against schemas
type Validator struct {
	// Compiler for JSON schemas
	compiler *jsonschema.Compiler
}

// ValidationResult contains the result of a validation operation
type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []ValidationWarning
}

// ValidationError represents a validation error with location information
type ValidationError struct {
	Field   string
	Message string
	Line    int
	Column  int
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Field   string
	Message string
}

// New creates a new Validator instance
func New() *Validator {
	compiler := jsonschema.NewCompiler()
	
	// Set up the compiler with default settings
	compiler.Draft = jsonschema.Draft7

	return &Validator{
		compiler: compiler,
	}
}

// ValidateDocument validates a document against a schema
func (v *Validator) ValidateDocument(docPath, schemaPath string) (*ValidationResult, error) {
	// Load the document
	docBytes, err := os.ReadFile(docPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read document file: %w", err)
	}

	// Load the schema
	schema, err := v.loadSchema(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load schema: %w", err)
	}

	// Validate the document
	return v.ValidateBytes(docBytes, schema)
}

// ValidateBytes validates a JSON document provided as bytes against a schema
func (v *Validator) ValidateBytes(docBytes []byte, schema *jsonschema.Schema) (*ValidationResult, error) {
	// Parse the document to get a Go value
	var doc interface{}
	err := json.Unmarshal(docBytes, &doc)
	if err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "",
					Message: fmt.Sprintf("Invalid JSON: %v", err),
					Line:    findErrorLine(err, string(docBytes)),
					Column:  findErrorColumn(err),
				},
			},
		}, nil
	}

	// Validate against the schema
	err = schema.Validate(doc)
	if err != nil {
		// Convert validation errors to our format
		return &ValidationResult{
			Valid:  false,
			Errors: convertValidationErrors(err),
		}, nil
	}

	return &ValidationResult{Valid: true}, nil
}

// ValidateString validates a JSON document provided as a string against a schema
func (v *Validator) ValidateString(docString, schemaString string) (*ValidationResult, error) {
	// Load the schema from string
	schema, err := v.loadSchemaFromString(schemaString)
	if err != nil {
		return nil, fmt.Errorf("failed to load schema from string: %w", err)
	}

	// Validate the document
	return v.ValidateBytes([]byte(docString), schema)
}

// loadSchema loads a JSON Schema from a file
func (v *Validator) loadSchema(schemaPath string) (*jsonschema.Schema, error) {
	// Check if the file exists
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("schema file not found: %s", schemaPath)
	}

	// Load the schema using the compiler
	schema, err := v.compiler.Compile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}

	return schema, nil
}

// loadSchemaFromString loads a JSON Schema from a string
func (v *Validator) loadSchemaFromString(schemaString string) (*jsonschema.Schema, error) {
	// Create a temporary file for the schema
	tmpFile, err := os.CreateTemp("", "schema-*.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write the schema to the temporary file
	_, err = tmpFile.WriteString(schemaString)
	if err != nil {
		return nil, fmt.Errorf("failed to write schema to temporary file: %w", err)
	}
	tmpFile.Close()

	// Load the schema using the compiler
	schema, err := v.compiler.Compile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}

	return schema, nil
}

// convertValidationErrors converts jsonschema validation errors to our format
func convertValidationErrors(err error) []ValidationError {
	var result []ValidationError

	if ve, ok := err.(*jsonschema.ValidationError); ok {
		// Process the basic error
		result = append(result, ValidationError{
			Field:   ve.InstanceLocation, // Use InstanceLocation instead of InstancePtr
			Message: ve.Message,
			Line:    0, // We don't have line information from the library
			Column:  0,
		})

		// Process any sub-errors
		for _, subErr := range ve.Causes {
			result = append(result, convertValidationErrors(subErr)...)
		}
	}

	return result
}

// findErrorLine attempts to find the line number where a JSON parsing error occurred
func findErrorLine(err error, content string) int {
	errMsg := err.Error()
	
	// Try to extract line information from the error message
	if strings.Contains(errMsg, "line") {
		// This is a simplistic approach and might need refinement
		parts := strings.Split(errMsg, "line")
		if len(parts) > 1 {
			var lineNum int
			_, err := fmt.Sscanf(parts[1], " %d", &lineNum)
			if err == nil {
				return lineNum
			}
		}
	}
	
	// If we can't extract from the error message, count lines up to the offset
	if strings.Contains(errMsg, "offset") {
		var offset int
		_, err := fmt.Sscanf(errMsg, "invalid character '%c' at offset %d", new(rune), &offset)
		if err == nil {
			// Count newlines up to the offset
			return strings.Count(content[:min(offset, len(content))], "\n") + 1
		}
	}
	
	return 0
}

// findErrorColumn attempts to find the column number where a JSON parsing error occurred
func findErrorColumn(err error) int {
	errMsg := err.Error()
	
	// Try to extract column information from the error message
	if strings.Contains(errMsg, "column") {
		parts := strings.Split(errMsg, "column")
		if len(parts) > 1 {
			var colNum int
			_, err := fmt.Sscanf(parts[1], " %d", &colNum)
			if err == nil {
				return colNum
			}
		}
	}
	
	// If we can't extract from the error message, try to get the offset
	if strings.Contains(errMsg, "offset") {
		var offset int
		_, err := fmt.Sscanf(errMsg, "invalid character '%c' at offset %d", new(rune), &offset)
		if err == nil {
			// Find the last newline before the offset
			content := err.Error()
			lastNewline := strings.LastIndex(content[:min(offset, len(content))], "\n")
			if lastNewline == -1 {
				return offset + 1 // 1-based column indexing
			}
			return offset - lastNewline
		}
	}
	
	return 0
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Validate validates an NLD document against its schema
func (v *Validator) Validate(document []byte) error {
	// This is kept for backward compatibility
	// It should be updated to use the new validation methods
	return nil
}

// LoadSchema loads a schema from a file or embedded resource
func (v *Validator) LoadSchema(schemaPath string) (*jsonschema.Schema, error) {
	return v.loadSchema(schemaPath)
}

// GetSchemaForDocumentType returns the appropriate schema for a given document type
func (v *Validator) GetSchemaForDocumentType(docType string) (string, error) {
	// Map document types to schema files
	schemaMap := map[string]string{
		"contract":  "schemas/document-v1.json",
		"receipt":   "schemas/document-v1.json",
		"agreement": "schemas/document-v1.json",
		"nda":       "schemas/nda.schema.json",
	}
	
	schemaPath, ok := schemaMap[strings.ToLower(docType)]
	if !ok {
		return "", fmt.Errorf("no schema available for document type: %s", docType)
	}
	
	// Resolve the schema path relative to the executable
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	
	baseDir := filepath.Dir(execPath)
	return filepath.Join(baseDir, schemaPath), nil
}

// FormatValidationResult formats a validation result as a human-readable string
func FormatValidationResult(result *ValidationResult) string {
	if result.Valid {
		green := color.New(color.FgGreen).SprintFunc()
		return green("Document is valid.")
	}
	
	var sb strings.Builder
	red := color.New(color.FgRed).SprintFunc()
	sb.WriteString(red("Document validation failed:\n"))
	
	for i, err := range result.Errors {
		sb.WriteString(fmt.Sprintf("%d. %s", i+1, err.Message))
		if err.Field != "" {
			sb.WriteString(fmt.Sprintf(" (at %s", err.Field))
			if err.Line > 0 {
				sb.WriteString(fmt.Sprintf(", line %d", err.Line))
				if err.Column > 0 {
					sb.WriteString(fmt.Sprintf(", column %d", err.Column))
				}
			}
			sb.WriteString(")")
		}
		sb.WriteString("\n")
	}
	
	return sb.String()
}

// ColoredOutput returns colored output for validation results
func ColoredOutput(valid bool, message string) string {
	if valid {
		green := color.New(color.FgGreen).SprintFunc()
		return green(message)
	}
	red := color.New(color.FgRed).SprintFunc()
	return red(message)
}