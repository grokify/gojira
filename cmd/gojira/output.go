package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gojira/rest"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// OutputFormat represents the output format type.
type OutputFormat int

const (
	OutputJSON OutputFormat = iota
	OutputTable
	OutputTOON
)

// OutputConfig holds output configuration.
type OutputConfig struct {
	Format OutputFormat
	Writer io.Writer
}

// NewOutputConfig creates a new OutputConfig with default writer (stdout).
func NewOutputConfig(format OutputFormat) *OutputConfig {
	return &OutputConfig{
		Format: format,
		Writer: os.Stdout,
	}
}

// WriteIssues writes issues in the specified format.
func WriteIssues(issues rest.Issues, cfg *OutputConfig) error {
	if cfg == nil {
		cfg = NewOutputConfig(OutputJSON)
	}
	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}

	switch cfg.Format {
	case OutputJSON:
		return writeIssuesJSON(issues, cfg.Writer)
	case OutputTable:
		return writeIssuesTable(issues, cfg.Writer)
	case OutputTOON:
		return writeIssuesToon(issues, cfg.Writer)
	default:
		return writeIssuesJSON(issues, cfg.Writer)
	}
}

// WriteIssue writes a single issue in the specified format.
func WriteIssue(issue *jira.Issue, cfg *OutputConfig) error {
	if issue == nil {
		return nil
	}
	issues := rest.Issues{*issue}
	return WriteIssues(issues, cfg)
}

// writeIssuesJSON outputs issues as JSON.
func writeIssuesJSON(issues rest.Issues, w io.Writer) error {
	metas := issuesToMetas(issues)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(metas)
}

// writeIssuesTable outputs issues as an ASCII table.
func writeIssuesTable(issues rest.Issues, w io.Writer) error {
	tw := tablewriter.NewWriter(w)
	tw.Header([]string{"Key", "Type", "Status", "Assignee", "Summary"})

	var rows [][]string
	for _, iss := range issues {
		im := rest.NewIssueMore(&iss)
		summary := truncateString(im.Summary(), 50)
		rows = append(rows, []string{
			im.Key(),
			im.Type(),
			im.Status(),
			im.AssigneeName(),
			summary,
		})
	}

	if err := tw.Bulk(rows); err != nil {
		return err
	}
	return tw.Render()
}

// writeIssuesToon outputs issues in TOON (Token-Optimized Object Notation) format.
// TOON is a compact key-value format designed for minimal token usage with LLMs.
// Format: K:KEY|T:Type|S:Status|A:Assignee|Su:Summary
func writeIssuesToon(issues rest.Issues, w io.Writer) error {
	for _, iss := range issues {
		im := rest.NewIssueMore(&iss)
		line := formatTOON(im)
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

// formatTOON formats a single issue in TOON format.
func formatTOON(im rest.IssueMore) string {
	parts := []string{
		"K:" + im.Key(),
		"T:" + im.Type(),
		"S:" + im.Status(),
	}

	if assignee := im.AssigneeName(); assignee != "" {
		parts = append(parts, "A:"+assignee)
	}

	if resolution := im.Resolution(); resolution != "" {
		parts = append(parts, "R:"+resolution)
	}

	if project := im.ProjectKey(); project != "" {
		parts = append(parts, "P:"+project)
	}

	// Summary - truncate for readability
	if summary := im.Summary(); summary != "" {
		parts = append(parts, "Su:"+truncateString(summary, 80))
	}

	return strings.Join(parts, "|")
}

// IssueMeta represents simplified issue metadata for JSON output.
type IssueMeta struct {
	Key        string     `json:"key"`
	Type       string     `json:"type"`
	Status     string     `json:"status"`
	Resolution string     `json:"resolution,omitempty"`
	Summary    string     `json:"summary"`
	Project    string     `json:"project"`
	ProjectKey string     `json:"projectKey"`
	Assignee   string     `json:"assignee,omitempty"`
	Creator    string     `json:"creator,omitempty"`
	Created    *time.Time `json:"created,omitempty"`
	Updated    *time.Time `json:"updated,omitempty"`
	ParentKey  string     `json:"parentKey,omitempty"`
	EpicKey    string     `json:"epicKey,omitempty"`
	Labels     []string   `json:"labels,omitempty"`
}

// issuesToMetas converts Issues to a slice of IssueMeta.
func issuesToMetas(issues rest.Issues) []IssueMeta {
	metas := make([]IssueMeta, 0, len(issues))
	for _, iss := range issues {
		im := rest.NewIssueMore(&iss)
		meta := IssueMeta{
			Key:        im.Key(),
			Type:       im.Type(),
			Status:     im.Status(),
			Resolution: im.Resolution(),
			Summary:    im.Summary(),
			Project:    im.Project(),
			ProjectKey: im.ProjectKey(),
			Assignee:   im.AssigneeName(),
			Creator:    im.CreatorName(),
			ParentKey:  im.ParentKey(),
			EpicKey:    im.EpicKey(),
			Labels:     im.Labels(true),
		}

		created := im.CreateTime()
		if !created.IsZero() {
			meta.Created = &created
		}

		updated := im.UpdateTime()
		if !updated.IsZero() {
			meta.Updated = &updated
		}

		metas = append(metas, meta)
	}
	return metas
}

// truncateString truncates a string to maxLen and adds ellipsis if needed.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// outputResult writes any JSON-serializable result to stdout.
// Used for commands that return non-Issue results like create, update, etc.
func outputResult(_ *cobra.Command, result any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
