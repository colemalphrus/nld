package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	c.rootCmd.PersistentFlags().StringVar(&c.outputFormat, "output-format", "text", "Output format (text, json)")
	
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
	var force bool
	
	validateCmd := &cobra.Command{
		Use:   "validate [file...]",
		Short: "Validate an NLD document",
		Long:  "Validate one or more NLD documents against their schema",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runValidateFiles(args, schemaPath, force)
		},
	}
	
	// Add validate-specific flags
	validateCmd.Flags().StringVarP(&schemaPath, "schema", "s", "", "Path to schema file (optional)")
	validateCmd.Flags().BoolVar(&force, "force", false, "Continue validation even if some files fail")
	
	c.rootCmd.AddCommand(validateCmd)
}

// addInitCommand adds the init command
func (c *CLI) addInitCommand() {
	var docType string
	var outputPath string
	var force bool
	var interactive bool
	var title string
	
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new NLD document",
		Long:  "Initialize a new NLD document with a specified template",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runInit(docType, outputPath, force, interactive, title)
		},
	}
	
	// Add init-specific flags
	initCmd.Flags().StringVarP(&docType, "type", "t", "contract", "Type of document to initialize (contract, receipt, agreement)")
	initCmd.Flags().StringVarP(&outputPath, "output", "o", "document.json", "Output file path")
	initCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file if it exists")
	initCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode to prompt for metadata")
	initCmd.Flags().StringVar(&title, "title", "", "Document title")
	
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

// runValidateFiles runs the validate command for multiple files
func (c *CLI) runValidateFiles(filePaths []string, schemaPath string, force bool) error {
	validCount := 0
	invalidCount := 0
	
	for _, filePath := range filePaths {
		err := c.runValidate(filePath, schemaPath)
		if err != nil {
			invalidCount++
			if !force {
				return err
			}
		} else {
			validCount++
		}
	}
	
	// Summary output
	if !c.quiet {
		if len(filePaths) > 1 {
			fmt.Printf("\nValidation summary: %d valid, %d invalid\n", validCount, invalidCount)
		}
	}
	
	// Return error if any files were invalid
	if invalidCount > 0 {
		return fmt.Errorf("%d file(s) failed validation", invalidCount)
	}
	
	return nil
}

// runValidate runs the validate command for a single file
func (c *CLI) runValidate(filePath, schemaPath string) error {
	if c.verbose {
		fmt.Printf("Validating file: %s\n", filePath)
		if schemaPath != "" {
			fmt.Printf("Using schema: %s\n", schemaPath)
		}
	}
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if !c.quiet {
			fmt.Printf("✗ %s: file not found\n", filePath)
		}
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
			if !c.quiet {
				fmt.Printf("✗ %s: failed to determine schema: %v\n", filePath, err)
			}
			return fmt.Errorf("failed to determine schema: %w", err)
		}
		
		// Read the document
		docBytes, err := os.ReadFile(filePath)
		if err != nil {
			if !c.quiet {
				fmt.Printf("✗ %s: failed to read document: %v\n", filePath, err)
			}
			return fmt.Errorf("failed to read document: %w", err)
		}
		
		// Validate using the determined schema
		result, err = s.Validate(docBytes)
	}
	
	if err != nil {
		if !c.quiet {
			fmt.Printf("✗ %s: validation error: %v\n", filePath, err)
		}
		return fmt.Errorf("validation error: %w", err)
	}
	
	// Output the result
	if !c.quiet {
		if c.outputFormat == "json" {
			// Output as JSON
			jsonResult, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format result as JSON: %w", err)
			}
			fmt.Println(string(jsonResult))
		} else {
			// Output as text with colors
			if result.Valid {
				fmt.Println(validator.ColoredOutput(true, fmt.Sprintf("✓ %s is valid", filePath)))
			} else {
				fmt.Println(validator.ColoredOutput(false, fmt.Sprintf("✗ %s has %d errors:", filePath, len(result.Errors))))
				for _, err := range result.Errors {
					lineInfo := ""
					if err.Line > 0 {
						lineInfo = fmt.Sprintf("Line %d: ", err.Line)
					}
					fmt.Printf("  - %s%s\n", lineInfo, err.Message)
					if c.verbose && err.Field != "" {
						fmt.Printf("    at %s\n", err.Field)
					}
				}
			}
		}
	}
	
	// Return an error if the document is invalid
	if !result.Valid {
		return fmt.Errorf("document validation failed")
	}
	
	return nil
}

// runInit runs the init command
func (c *CLI) runInit(docType, outputPath string, force, interactive bool, title string) error {
	if c.verbose {
		fmt.Printf("Initializing new %s document: %s\n", docType, outputPath)
	}
	
	// Check if file exists and force flag is not set
	if _, err := os.Stat(outputPath); err == nil && !force {
		return fmt.Errorf("file already exists: %s (use --force to overwrite)", outputPath)
	}
	
	// Create metadata based on document type
	metadata := map[string]interface{}{
		"version": "1.0.0",
		"type":    docType,
		"created": time.Now().Format(time.RFC3339),
	}
	
	// Set title from flag or default
	if title != "" {
		metadata["title"] = title
	} else {
		metadata["title"] = fmt.Sprintf("New %s", docType)
	}
	
	// Interactive mode to prompt for additional metadata
	if interactive {
		if !c.quiet {
			fmt.Println("Enter document metadata (press Enter to use default):")
		}
		
		// Prompt for title if not provided via flag
		if title == "" {
			fmt.Printf("Title [%s]: ", metadata["title"])
			var input string
			fmt.Scanln(&input)
			if input != "" {
				metadata["title"] = input
			}
		}
		
		// Prompt for author
		fmt.Print("Author: ")
		var author string
		fmt.Scanln(&author)
		if author != "" {
			metadata["author"] = author
		}
		
		// Prompt for jurisdiction if it's a contract or agreement
		if docType == "contract" || docType == "agreement" {
			fmt.Print("Jurisdiction: ")
			var jurisdiction string
			fmt.Scanln(&jurisdiction)
			if jurisdiction != "" {
				metadata["jurisdiction"] = jurisdiction
			}
		}
	}
	
	// Create document content based on type
	var sections []map[string]interface{}
	
	switch docType {
	case "contract":
		sections = []map[string]interface{}{
			{
				"id":      "parties",
				"title":   "Parties",
				"content": "This agreement is between the following parties:",
			},
			{
				"id":      "scope",
				"title":   "Scope of Work",
				"content": "The scope of work includes the following:",
			},
			{
				"id":      "terms",
				"title":   "Terms and Conditions",
				"content": "The following terms and conditions apply:",
			},
		}
	case "receipt":
		sections = []map[string]interface{}{
			{
				"id":      "transaction",
				"title":   "Transaction Details",
				"content": "Transaction details go here.",
			},
			{
				"id":      "items",
				"title":   "Items",
				"content": "List of items purchased.",
			},
			{
				"id":      "payment",
				"title":   "Payment Information",
				"content": "Payment details go here.",
			},
		}
	case "agreement":
		sections = []map[string]interface{}{
			{
				"id":      "introduction",
				"title":   "Introduction",
				"content": "This agreement is made on the date specified above.",
			},
			{
				"id":      "terms",
				"title":   "Terms of Agreement",
				"content": "The parties agree to the following terms:",
			},
			{
				"id":      "signatures",
				"title":   "Signatures",
				"content": "The parties have executed this agreement as follows:",
			},
		}
	default:
		sections = []map[string]interface{}{
			{
				"id":      "section1",
				"title":   "Section 1",
				"content": "Enter your content here.",
			},
		}
	}
	
	// Create the document structure
	doc := map[string]interface{}{
		"metadata": metadata,
		"content": map[string]interface{}{
			"sections": sections,
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
	
	if !c.quiet {
		fmt.Println(validator.ColoredOutput(true, fmt.Sprintf("Created new %s document: %s", docType, outputPath)))
	}
	return nil
}

// showVersion displays version information
func (c *CLI) showVersion() {
	fmt.Printf("NLD - Next-Gen Layout Document Tool\n")
	fmt.Printf("Version: %s\n", Version)
}