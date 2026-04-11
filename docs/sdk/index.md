# SDK Guide

GoJira provides a Go SDK for interacting with the Jira REST API. The SDK is organized into packages with clear dependency boundaries.

## Package Structure

| Package | Import | Description |
|---------|--------|-------------|
| `gojira` | `github.com/grokify/gojira` | JQL builder, config, constants (no external deps) |
| `rest` | `github.com/grokify/gojira/rest` | REST API client (requires go-jira) |
| `xml` | `github.com/grokify/gojira/xml` | XML export parser (no external deps) |
| `web` | `github.com/grokify/gojira/web` | URL helpers (no external deps) |

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/grokify/gojira/rest"
)

func main() {
    // Create client with basic auth
    client, err := rest.NewClientFromBasicAuth(
        "https://your-instance.atlassian.net",
        "your-email@example.com",
        "your-api-token",
        false, // don't load custom fields on init
    )
    if err != nil {
        log.Fatal(err)
    }

    // Search issues
    issues, err := client.IssueAPI.SearchIssues("project = FOO AND status = Open", false)
    if err != nil {
        log.Fatal(err)
    }

    // Print results
    for _, issue := range issues {
        fmt.Printf("%s: %s\n", issue.Key, issue.Fields.Summary)
    }
}
```

## Key Concepts

### Client

The `rest.Client` is your entry point to the API. It provides access to service APIs:

- `client.IssueAPI` - Issue operations (search, get, update)
- `client.CustomFieldAPI` - Custom field operations
- `client.BacklogAPI` - Backlog operations

See [Client](client.md) for details.

### Issues

Issues are returned as `rest.Issues` (a slice of `jira.Issue`). Use `IssuesSet` for advanced operations:

```go
issuesSet, err := issues.IssuesSet(nil)
if err != nil {
    log.Fatal(err)
}

// Get counts by status
counts := issuesSet.CountsByStatus()
```

See [Issues](issues.md) for details.

### JQL Builder

The root package provides a JQL builder for constructing queries programmatically:

```go
import "github.com/grokify/gojira"

jql := gojira.JQL{
    ProjectsIncl: [][]string{{"FOO", "BAR"}},
    StatusesIncl: [][]string{{"Open", "In Progress"}},
}

query := jql.String()
// Result: "project IN ('FOO', 'BAR') AND status IN ('Open', 'In Progress')"
```

See [JQL Builder](jql.md) for details.

## Authentication Options

```go
// Option 1: Direct credentials
client, err := rest.NewClientFromBasicAuth(url, user, token, false)

// Option 2: From goauth credentials file
client, err := rest.NewClientGoauthBasicAuthFile(
    "~/.config/goauth/credentials.json",
    "jira-prod",
    false,
)

// Option 3: Interactive selection from goauth CLI
client, err := rest.NewClientFromGoauthCLI(true, false)
```

## Error Handling

All operations return errors that should be checked:

```go
issues, err := client.IssueAPI.SearchIssues(jql, false)
if err != nil {
    // Handle error (authentication, network, API errors)
    return err
}
```

## Next Steps

- [Client](client.md) - Creating and configuring clients
- [Issues](issues.md) - Working with issues
- [JQL Builder](jql.md) - Building JQL queries
