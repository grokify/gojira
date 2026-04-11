# Client

The `rest.Client` provides access to the Jira REST API. It wraps the [go-jira](https://github.com/andygrunwald/go-jira) library with additional utilities.

## Creating a Client

### From Direct Credentials

```go
import "github.com/grokify/gojira/rest"

client, err := rest.NewClientFromBasicAuth(
    "https://your-instance.atlassian.net", // Server URL
    "your-email@example.com",              // Username
    "your-api-token",                       // API token/password
    false,                                   // Load custom fields on init
)
if err != nil {
    log.Fatal(err)
}
```

### From goauth Credentials File

```go
client, err := rest.NewClientGoauthBasicAuthFile(
    "~/.config/goauth/credentials.json", // Credentials file path
    "jira-prod",                          // Account key
    false,                                 // Load custom fields
)
```

### From Environment Variables

Create a helper to read from environment:

```go
import (
    "os"
    "github.com/grokify/gojira/rest"
)

func NewClientFromEnv() (*rest.Client, error) {
    return rest.NewClientFromBasicAuth(
        os.Getenv("JIRA_URL"),
        os.Getenv("JIRA_USER"),
        os.Getenv("JIRA_TOKEN"),
        false,
    )
}
```

### Interactive Selection

Use goauth CLI for interactive account selection:

```go
client, err := rest.NewClientFromGoauthCLI(
    true,  // Include accounts on error
    false, // Load custom fields
)
```

## Client Structure

```go
type Client struct {
    Config         *gojira.Config
    HTTPClient     *http.Client
    JiraClient     *jira.Client
    Logger         *slog.Logger
    BacklogAPI     *BacklogService
    CustomFieldAPI *CustomFieldService
    IssueAPI       *IssueService
    CustomFieldSet *CustomFieldSet
}
```

## Service APIs

### IssueAPI

The most commonly used service for issue operations:

```go
// Search issues
issues, err := client.IssueAPI.SearchIssues("project = FOO", false)

// Get single issue
ctx := context.Background()
issue, err := client.IssueAPI.Issue(ctx, "FOO-123", nil)

// Get multiple issues
issues, err := client.IssueAPI.Issues(ctx, []string{"FOO-1", "FOO-2"}, nil)

// Update issue
resp, err := client.IssueAPI.IssuePatch(ctx, "FOO-123", patchBody)
```

### CustomFieldAPI

Access custom field definitions:

```go
// Get all custom fields
fields, err := client.CustomFieldAPI.GetCustomFields()

// Get Epic Link field
epicField, err := client.CustomFieldAPI.GetCustomFieldEpicLink()
```

### BacklogAPI

Backlog operations:

```go
// Get backlog for a board
backlog, err := client.BacklogAPI.Backlog(boardID)
```

## Loading Custom Fields

Custom fields can be loaded during client initialization or later:

```go
// Load during init
client, err := rest.NewClientFromBasicAuth(url, user, token, true)

// Or load later
err := client.LoadCustomFields()

// Access loaded fields
if client.CustomFieldSet != nil {
    field, ok := client.CustomFieldSet.ByID("customfield_10001")
}
```

## Configuration

The client includes configuration from the root package:

```go
// Access config
serverURL := client.Config.ServerURL

// Modify config
client.Config.ServerURL = "https://new-instance.atlassian.net"
```

## Logging

Set a logger for debug output:

```go
import "log/slog"

client.Logger = slog.Default()
```

## Low-Level Access

For operations not covered by service APIs, access the underlying clients:

```go
// go-jira client
jiraClient := client.JiraClient

// Raw HTTP client
httpClient := client.HTTPClient
```

## Example: Complete Workflow

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/grokify/gojira/rest"
)

func main() {
    // Create client
    client, err := rest.NewClientFromBasicAuth(
        "https://your-instance.atlassian.net",
        "your-email@example.com",
        "your-api-token",
        true, // Load custom fields
    )
    if err != nil {
        log.Fatal(err)
    }

    // Search issues
    issues, err := client.IssueAPI.SearchIssues("project = FOO AND status = Open", false)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d issues\n", len(issues))

    // Get details for first issue
    if len(issues) > 0 {
        ctx := context.Background()
        issue, err := client.IssueAPI.Issue(ctx, issues[0].Key, nil)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Issue: %s - %s\n", issue.Key, issue.Fields.Summary)
    }
}
```
