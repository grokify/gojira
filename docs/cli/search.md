# search

Search for Jira issues using JQL (Jira Query Language).

## Usage

```bash
gojira search --jql <query> [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--jql` | | (required) | JQL query string |
| `--max` | `-m` | 50 | Maximum number of results |
| `--all` | `-a` | false | Retrieve all results (paginate automatically) |
| `--fields` | `-f` | | Comma-separated list of fields to include |

Plus [global flags](index.md#global-flags).

## Examples

### Basic Search

```bash
# Search for open issues in a project
gojira search --jql "project = FOO AND status = Open"

# Search assigned to current user
gojira search --jql "assignee = currentUser()"
```

### Limiting Results

```bash
# Get first 10 results
gojira search --jql "project = FOO" --max 10

# Get all results (may be slow for large datasets)
gojira search --jql "project = FOO" --all
```

### Output Formats

```bash
# JSON output (default)
gojira search --jql "project = FOO" --json

# Human-readable table
gojira search --jql "project = FOO" --table

# Token-optimized for LLMs
gojira search --jql "project = FOO" --toon
```

### Piping to jq

```bash
# Extract just the keys
gojira search --jql "project = FOO" | jq -r '.[].key'

# Get summaries
gojira search --jql "project = FOO" | jq -r '.[] | "\(.key): \(.fields.summary)"'
```

### Quiet Mode

```bash
# Suppress progress messages (useful for scripting)
gojira search --jql "project = FOO" -q
```

## Output

### JSON Format

```json
[
  {
    "key": "FOO-123",
    "fields": {
      "summary": "Fix login bug",
      "status": {
        "name": "Open"
      },
      "issuetype": {
        "name": "Bug"
      }
    }
  }
]
```

### Table Format

```
KEY        TYPE    STATUS    SUMMARY
FOO-123    Bug     Open      Fix login bug
FOO-124    Story   Done      Add user profile
```

## JQL Tips

See [JQL Examples](../guides/jql-examples.md) for common query patterns.

```bash
# Complex queries
gojira search --jql "project = FOO AND status in (Open, 'In Progress') AND priority = High"

# Date filtering
gojira search --jql "project = FOO AND created >= -7d"

# Text search
gojira search --jql "project = FOO AND text ~ 'error'"
```
