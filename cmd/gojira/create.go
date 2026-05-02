package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grokify/gojira/core"
	"github.com/spf13/cobra"
)

var (
	createFile    string
	createDryRun  bool
	createProject string
	createParent  string
	createType    string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Jira issue from a YAML file",
	Long: `Create a Jira issue from a YAML file.

The YAML file should contain issue fields:

  project: PROJ
  type: Story
  summary: Add user authentication
  description: |
    Implement OAuth2 login flow.
  parent: PROJ-100
  labels:
    - auth
    - mvp
  customfield_12345: |
    Given A, When B, Then C

Standard fields: project, type, summary, description, parent, labels,
priority, assignee, reporter, components, fix_versions

Custom fields: Any key starting with customfield_ is passed directly
to the Jira API.

Examples:
  # Create issue from file
  gojira create -f story.yaml

  # Dry run to preview
  gojira create -f story.yaml --dry-run

  # Override project
  gojira create -f story.yaml --project PROJ

  # Override parent (for subtasks or stories under epics)
  gojira create -f story.yaml --parent EPIC-123`,
	RunE: runCreate,
}

func init() {
	createCmd.Flags().StringVarP(&createFile, "file", "f", "", "YAML file containing issue data (required)")
	createCmd.Flags().BoolVar(&createDryRun, "dry-run", false, "Validate and show what would be created without actually creating")
	createCmd.Flags().StringVar(&createProject, "project", "", "Override project key from file")
	createCmd.Flags().StringVar(&createParent, "parent", "", "Override parent issue key from file")
	createCmd.Flags().StringVar(&createType, "type", "", "Override issue type from file")

	if err := createCmd.MarkFlagRequired("file"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, _ []string) error {
	// Read and parse the file
	data, err := os.ReadFile(createFile)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	input, err := core.ParseIssueYAML(data)
	if err != nil {
		return fmt.Errorf("parse YAML: %w", err)
	}

	// Apply overrides
	if createProject != "" {
		input.Project = createProject
	}
	if createParent != "" {
		input.Parent = createParent
	}
	if createType != "" {
		input.Type = createType
	}

	// Dry run mode
	if createDryRun {
		result, err := core.DryRunCreate(input)
		if err != nil {
			return err
		}
		return outputResult(cmd, result)
	}

	// Create the issue
	client, err := NewClientFromOptions(getAuthOptions())
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	result, err := core.CreateIssue(context.Background(), client, input)
	if err != nil {
		return err
	}

	// Display the created issue key prominently
	fmt.Printf("Created %s\n", result.Key)

	// Also output full result as JSON if verbose/json output requested
	if flagJSON {
		return outputResult(cmd, result)
	}
	return nil
}
