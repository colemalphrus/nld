package cli

import (
	"fmt"
	"os"

	"github.com/colemalphrus/nld/internal/validator"
	"github.com/spf13/cobra"
)

const (
	// Version is the current version of the NLD tool
	Version = "0.1.0"
)

// CLI handles command-line interface operations for the NLD tool
type CLI struct {
	rootCmd     *cobra.Command
	validator   *validator.Validator
	verbose     bool
	quiet       bool
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
	validateCmd := &cobra.Command{
		Use:   "validate [file]",
		Short: "Validate an NLD document",
		Long:  "Validate an NLD document against its schema",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runValidate(args[0])
		},
	}
	
	c.rootCmd.AddCommand(validateCmd)
}

// addInitCommand adds the init command
func (c *CLI) addInitCommand() {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new NLD document",
		Long:  "Initialize a new NLD document with a specified template",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runInit()
		},
	}
	
	// Add init-specific flags
	initCmd.Flags().String("type", "default", "Type of document to initialize")
	
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
func (c *CLI) runValidate(filePath string) error {
	if c.verbose {
		fmt.Printf("Validating file: %s\n", filePath)
	}
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}
	
	// Placeholder for validation logic
	fmt.Println("Validation successful (placeholder)")
	return nil
}

// runInit runs the init command
func (c *CLI) runInit() error {
	// Placeholder for init logic
	fmt.Println("Initialized new NLD document (placeholder)")
	return nil
}

// showVersion displays version information
func (c *CLI) showVersion() {
	fmt.Printf("NLD - Next-Gen Layout Document Tool\n")
	fmt.Printf("Version: %s\n", Version)
}