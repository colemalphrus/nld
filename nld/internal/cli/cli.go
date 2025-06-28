package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/colemalphrus/nld/internal/schema"
	"github.com/colemalphrus/nld/internal/validator"
	"github.com/spf13/cobra"
)

const (
	// Version is the current version of the NLD tool
	Version = "0.1.0"
)

// CLI handles command-line interface operations for the NLD tool
type CLI struct {
	rootCmd      *cobra.Command
	validator    *validator.Validator
	verbose      bool
	quiet        bool
	outputFormat string
}

// New creates a new CLI instance
func New() *CLI {
	cli := &CLI{
		validator: validator.New(),
	}
	
	cli.setupCommands()
	return cli
}

// Execute executes the CLI with the given arguments
func (c *CLI) Execute(args []string) error {
	c.rootCmd.SetArgs(args)
	return c.rootCmd.Execute()
}

// setupCommands initializes all CLI commands
func (c *CLI) setupCommands() {
	// Root command
	c.rootCmd = &cobra.Command{
		Use:   "nld",
		Short: "NLD - Next-Gen Layout Document Tool",
		Long: `NLD is a tool for working with Next-Gen Layout Documents.
It provides functionality for creating, validating, and managing NLD documents.`,
		SilenceUsage: true,
		SilenceErrors: true,
	}

	// Global flags
	c.rootCmd.PersistentFlags().BoolVarP(&c.verbose, "verbose", "v", false, "Enable verbose output")
	c.rootCmd.PersistentFlags().BoolVarP(&c.quiet, "quiet", "q", false, "Suppress all output except errors")
	c.rootCmd.PersistentFlags().StringVarP(&c.outputFormat, "output-format", "o", "text", "Output format (text, json)")
	
	// Version flag on root command
	c.rootCmd.Flags().BoolP("version", "V", false, "Display version information")
	c.rootCmd.Run = func(cmd *cobra.Command, args []string) {
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			c.showVersion()
			return
		}
		
		// If no subcommand is provided, show help
		if len(args) == 0 {
			cmd.Help()
		}
	}

	// Add subcommands
	c.addValidateCommand()
	c.addInitCommand()
	c.addVersionCommand()
}

// addValidateCommand adds the validate command
func (c *CLI) addValidateCommand() {
	var schemaPath string
	
	validateCmd := &cobra.Command{
		Use:   "validate [file]",
		Short: "Validate an NLD document",
		Long:  "Validate an NLD document against its schema",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runValidate(args[0], schemaPath)
		},
	}
	
	// Add validate-specific flags
	validateCmd.Flags().StringVarP(&schemaPath, "schema", "s", "", "Path to schema file (optional)")
	
	c.rootCmd.AddCommand(validateCmd)
}

// addInitCommand adds the init command
func (c *CLI) addInitCommand() {
	var docType string
	var outputPath string
	
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new NLD document",
		Long:  "Initialize a new NLD document with a specified template",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runInit(docType, outputPath)
		},
	}
	
	// Add init-specific flags
	initCmd.Flags().StringVarP(&docType, "type", "t", "contract", "Type of document to initialize (contract, receipt, agreement)")
	initCmd.Flags().StringVarP(&outputPath, "output", "o", "document.json", "Output file path")
	
	c.rootCmd.AddCommand(initCmd)
}

// addVersionCommand adds the version command
func (c *CLI) addVersionCommand() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  "Display detailed version information about the NLD tool",
		Run: func(cmd *cobra.Command, args []string) {
			c.showVersion()
		},
	}
	
	c.rootCmd.AddCommand(versionCmd)
}

// runValidate runs the validate command
func (c *CLI) runValidate(filePath, schemaPath string) error {
	if c.verbose {
		fmt.Printf("Validating file: %s\n", filePath)
		if schemaPath != "" {
			fmt.Printf("Using schema: %s\n", schemaPath)
		}
	}
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}
	
	var result *validator.ValidationResult
	var err error
	
	if schemaPath != "" {
		// Use the specified schema
		result, err = c.validator.ValidateDocument(filePath, schemaPath)
	} else {
		// Determine the schema based on the document type
		s, err := schema.GetDocumentSchema(filePath)
		if err != nil {
			return fmt.Errorf("failed to determine schema: %w", err)
		}
		
		// Read the document
		docBytes, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read document: %w", err)
		}
		
		// Validate using the determined schema
		result, err = s.Validate(docBytes)
	}
	
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	
	// Output the result
	if c.outputFormat == "json" {
		// Output as JSON
		jsonResult, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format result as JSON: %w", err)
		}
		fmt.Println(string(jsonResult))
	} else {
		// Output as text
		if result.Valid {
			fmt.Println("Document is valid.")
		} else {
			fmt.Println(validator.FormatValidationResult(result))
		}
	}
	
	// Return an error if the document is invalid
	if !result.Valid {
		return fmt.Errorf("document validation failed")
	}
	
	return nil
}

// runInit runs the init command
func (c *CLI) runInit(docType, outputPath string) error {
	if c.verbose {
		fmt.Printf("Initializing new %s document: %s\n", docType, outputPath)
	}
	
	// Create a basic document template based on the type
	doc := map[string]interface{}{
		"metadata": map[string]interface{}{
			"version": "1.0.0",
			"type":    docType,
			"created": "2025-06-27T00:00:00Z",
			"title":   fmt.Sprintf("New %s", docType),
		},
		"content": map[string]interface{}{
			"sections": []map[string]interface{}{
				{
					"id":      "section1",
					"title":   "Section 1",
					"content": "Enter your content here.",
				},
			},
		},
	}
	
	// Convert to JSON
	jsonData, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	
	// Write to file
	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}
	
	fmt.Printf("Created new %s document: %s\n", docType, outputPath)
	return nil
}

// showVersion displays version information
func (c *CLI) showVersion() {
	fmt.Printf("NLD - Next-Gen Layout Document Tool\n")
	fmt.Printf("Version: %s\n", Version)
}