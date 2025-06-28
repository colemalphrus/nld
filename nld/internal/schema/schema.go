package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/colemalphrus/nld/internal/validator"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

// Schema represents an NLD document schema
type Schema struct {
	Path     string
	RawData  []byte
	Compiled *jsonschema.Schema
}

// Load loads a schema from a file
func Load(path string) (*Schema, error) {
	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("schema file not found: %s", path)
	}

	// Read the schema file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	// Validate that it's valid JSON
	var jsonObj interface{}
	if err := json.Unmarshal(data, &jsonObj); err != nil {
		return nil, fmt.Errorf("invalid JSON in schema file: %w", err)
	}

	// Compile the schema
	v := validator.New()
	compiled, err := v.LoadSchema(path)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}

	return &Schema{
		Path:     path,
		RawData:  data,
		Compiled: compiled,
	}, nil
}

// Validate validates a document against this schema
func (s *Schema) Validate(document []byte) (*validator.ValidationResult, error) {
	v := validator.New()
	return v.ValidateBytes(document, s.Compiled)
}

// GetDocumentSchema returns the appropriate schema for a document
func GetDocumentSchema(docPath string) (*Schema, error) {
	// Read the document to determine its type
	data, err := os.ReadFile(docPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read document: %w", err)
	}

	// Parse the document to extract the type
	var doc struct {
		Metadata struct {
			Type string `json:"type"`
		} `json:"metadata"`
	}

	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("invalid JSON in document: %w", err)
	}

	// Get the schema path for this document type
	v := validator.New()
	schemaPath, err := v.GetSchemaForDocumentType(doc.Metadata.Type)
	if err != nil {
		// If we can't determine the type, use the default schema
		execPath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed to get executable path: %w", err)
		}
		
		baseDir := filepath.Dir(execPath)
		schemaPath = filepath.Join(baseDir, "schemas/document-v1.json")
	}

	// Load the schema
	return Load(schemaPath)
}

// GetSchemaVersion extracts the version from a schema file
func GetSchemaVersion(schemaPath string) (string, error) {
	// Read the schema file
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return "", fmt.Errorf("failed to read schema file: %w", err)
	}

	// Parse the schema to extract the title which contains the version
	var schema struct {
		Title string `json:"title"`
	}

	if err := json.Unmarshal(data, &schema); err != nil {
		return "", fmt.Errorf("invalid JSON in schema file: %w", err)
	}

	return schema.Title, nil
}