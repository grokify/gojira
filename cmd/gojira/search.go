package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagSearchJQL    string
	flagSearchMax    int
	flagSearchAll    bool
	flagSearchFields string
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search issues with JQL",
	Long: `Search searches for Jira issues using JQL (Jira Query Language).

Examples:
  # Search for open issues in a project
  gojira search --jql "project = FOO AND status = Open"

  # Search with limit
  gojira search --jql "project = FOO" --max 100

  # Search and retrieve all results
  gojira search --jql "project = FOO" --all

  # Output as table
  gojira search --jql "assignee = currentUser()" --table

  # Token-optimized output for LLMs
  gojira search --jql "project = FOO" --toon`,
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVar(&flagSearchJQL, "jql", "", "JQL query string (required)")
	searchCmd.Flags().IntVarP(&flagSearchMax, "max", "m", 50, "Maximum number of results")
	searchCmd.Flags().BoolVarP(&flagSearchAll, "all", "a", false, "Retrieve all results (paginate automatically)")
	searchCmd.Flags().StringVarP(&flagSearchFields, "fields", "f", "", "Comma-separated list of fields to include")

	if err := searchCmd.MarkFlagRequired("jql"); err != nil {
		panic(err)
	}
}

func runSearch(cmd *cobra.Command, args []string) error {
	if flagSearchJQL == "" {
		return fmt.Errorf("--jql flag is required")
	}

	// Create client
	client, err := NewClientFromOptions(getAuthOptions())
	if err != nil {
		return fmt.Errorf("failed to create Jira client: %w", err)
	}

	// Search issues
	var issues, errSearch = client.IssueAPI.SearchIssues(flagSearchJQL, flagSearchAll || flagSearchMax == 0)
	if errSearch != nil {
		return fmt.Errorf("search failed: %w", errSearch)
	}

	// Apply max limit if not retrieving all
	if !flagSearchAll && flagSearchMax > 0 && len(issues) > flagSearchMax {
		issues = issues[:flagSearchMax]
	}

	if len(issues) == 0 {
		if !flagQuiet {
			fmt.Fprintln(os.Stderr, "No issues found")
		}
		return nil
	}

	if !flagQuiet {
		fmt.Fprintf(os.Stderr, "Found %d issue(s)\n", len(issues))
	}

	// Output results
	cfg := NewOutputConfig(getOutputFormat())
	return WriteIssues(issues, cfg)
}
