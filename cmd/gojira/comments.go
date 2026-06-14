package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagCommentsMaxResults int
)

var commentsCmd = &cobra.Command{
	Use:   "comments <issue-key>",
	Short: "Get comments for an issue",
	Long: `Retrieves the comment thread for a Jira issue.

Examples:
  # Get comments for an issue
  gojira comments ISSUE-123

  # Limit the number of comments returned
  gojira comments ISSUE-123 --max 10`,
	Args: cobra.ExactArgs(1),
	RunE: runComments,
}

func init() {
	rootCmd.AddCommand(commentsCmd)

	commentsCmd.Flags().IntVar(&flagCommentsMaxResults, "max", 50, "Maximum number of comments to return")
}

func runComments(cmd *cobra.Command, args []string) error {
	key := args[0]

	// Create client
	client, err := NewClientFromOptions(getAuthOptions())
	if err != nil {
		return fmt.Errorf("failed to create Jira client: %w", err)
	}

	ctx := context.Background()

	// Use shared GetComments method
	response, err := client.GetComments(ctx, key, flagCommentsMaxResults)
	if err != nil {
		return fmt.Errorf("failed to get comments: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(response)
}
