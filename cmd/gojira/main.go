// gojira is a command-line interface for the Jira REST API.
//
// It is designed primarily for AI agents (with JSON output) but also
// supports human-readable table output and token-optimized TOON format
// for LLMs.
//
// Usage:
//
//	gojira search --jql "project = FOO"
//	gojira get ISSUE-123
//	gojira version
package main

func main() {
	Execute()
}
