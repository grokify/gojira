package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/grokify/gojira/rest"
	"github.com/spf13/cobra"
)

var fieldsCmd = &cobra.Command{
	Use:   "fields [flags]",
	Short: "List and filter custom fields",
	Long: `List custom fields from Jira, with optional filtering.

Examples:
  # List all custom fields
  gojira fields

  # Filter by field ID
  gojira fields --id customfield_10001

  # Filter by multiple IDs
  gojira fields --id customfield_10001,customfield_10002

  # Filter by name (partial match)
  gojira fields --name "Epic"

  # Filter by exact name
  gojira fields --name-exact "Epic Link"

  # Show only custom fields (exclude system fields)
  gojira fields --custom-only

  # Output as JSON
  gojira fields --json

  # Get Epic Link field specifically
  gojira fields --epic-link`,
	RunE: runFields,
}

var (
	fieldsFilterIDs   string
	fieldsFilterName  string
	fieldsFilterExact string
	fieldsCustomOnly  bool
	fieldsEpicLink    bool
	fieldsOutputJSON  bool
	fieldsOutputTable bool
)

func init() {
	rootCmd.AddCommand(fieldsCmd)

	fieldsCmd.Flags().StringVar(&fieldsFilterIDs, "id", "", "Filter by field ID(s), comma-separated")
	fieldsCmd.Flags().StringVar(&fieldsFilterName, "name", "", "Filter by name (partial match)")
	fieldsCmd.Flags().StringVar(&fieldsFilterExact, "name-exact", "", "Filter by exact name")
	fieldsCmd.Flags().BoolVar(&fieldsCustomOnly, "custom-only", false, "Show only custom fields")
	fieldsCmd.Flags().BoolVar(&fieldsEpicLink, "epic-link", false, "Show Epic Link field")
	fieldsCmd.Flags().BoolVar(&fieldsOutputJSON, "json", false, "Output as JSON")
	fieldsCmd.Flags().BoolVar(&fieldsOutputTable, "table", true, "Output as table (default)")
}

func runFields(cmd *cobra.Command, args []string) error {
	client, err := NewClientFromOptions(getAuthOptions())
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Special case: get Epic Link field
	if fieldsEpicLink {
		return showEpicLinkField(client)
	}

	// Get all fields
	fields, err := client.CustomFieldAPI.GetCustomFields()
	if err != nil {
		return fmt.Errorf("failed to get fields: %w", err)
	}

	// Apply filters
	fields = applyFieldFilters(fields)

	if len(fields) == 0 {
		fmt.Fprintln(os.Stderr, "No fields found matching criteria")
		return nil
	}

	// Output
	if fieldsOutputJSON {
		return outputFieldsJSON(fields)
	}

	return outputFieldsTable(fields)
}

func showEpicLinkField(client *rest.Client) error {
	field, err := client.CustomFieldAPI.GetCustomFieldEpicLink()
	if err != nil {
		return fmt.Errorf("failed to get Epic Link field: %w", err)
	}

	if field.ID == "" {
		fmt.Fprintln(os.Stderr, "Epic Link field not found")
		return nil
	}

	fields := rest.CustomFields{field}

	if fieldsOutputJSON {
		return outputFieldsJSON(fields)
	}

	return outputFieldsTable(fields)
}

func applyFieldFilters(fields rest.CustomFields) rest.CustomFields {
	// Filter by custom only
	if fieldsCustomOnly {
		var filtered rest.CustomFields
		for _, f := range fields {
			if f.Custom {
				filtered = append(filtered, f)
			}
		}
		fields = filtered
	}

	// Filter by IDs
	if fieldsFilterIDs != "" {
		ids := parseCommaSeparated(fieldsFilterIDs)
		fields = fields.FilterByIDs(ids...)
	}

	// Filter by exact name
	if fieldsFilterExact != "" {
		names := parseCommaSeparated(fieldsFilterExact)
		fields = fields.FilterByNames(names...)
	}

	// Filter by partial name match
	if fieldsFilterName != "" {
		var filtered rest.CustomFields
		searchLower := strings.ToLower(fieldsFilterName)
		for _, f := range fields {
			if strings.Contains(strings.ToLower(f.Name), searchLower) {
				filtered = append(filtered, f)
			}
		}
		fields = filtered
	}

	return fields
}

func parseCommaSeparated(s string) []string {
	var result []string
	for _, item := range strings.Split(s, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func outputFieldsJSON(fields rest.CustomFields) error {
	data, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func outputFieldsTable(fields rest.CustomFields) error {
	return fields.WriteTable(os.Stdout)
}
