# get

Get one or more Jira issues by their keys.

## Usage

```bash
gojira get <issue-key> [issue-key...] [flags]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `issue-key` | One or more Jira issue keys (e.g., `FOO-123`) |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--expand` | false | Expand changelog and other fields |
| `--fields` | | Comma-separated list of fields to include |
| `--raw` | false | Output full API JSON instead of simplified metadata |

Plus [global flags](index.md#global-flags).

## Examples

### Single Issue

```bash
# Get a single issue
gojira get FOO-123

# Human-readable output
gojira get FOO-123 --table
```

### Multiple Issues

```bash
# Get multiple issues
gojira get FOO-123 FOO-456 FOO-789

# Output as JSON array
gojira get FOO-123 FOO-456 --json
```

### With Changelog

```bash
# Include issue history
gojira get FOO-123 --expand
```

### Full API JSON

```bash
# Output complete API response (all fields)
gojira get FOO-123 --raw
```

The `--raw` flag outputs the full `jira.Issue` JSON from the API instead of the simplified metadata format. Use this when you need access to all fields including custom fields, comments, or other data not included in the default output.

## Output

### JSON Format (Default)

The default output is a simplified metadata format:

```json
[
  {
    "key": "FOO-123",
    "type": "Bug",
    "status": "Open",
    "summary": "Fix login bug",
    "description": "Users cannot log in when...",
    "project": "Foo Project",
    "projectKey": "FOO",
    "assignee": "John Doe",
    "creator": "Jane Smith",
    "created": "2024-01-15T10:30:00Z",
    "updated": "2024-01-16T14:20:00Z",
    "labels": ["urgent", "backend"]
  }
]
```

### Raw JSON Format (`--raw`)

The `--raw` flag outputs the complete API response:

```json
[
  {
    "key": "FOO-123",
    "fields": {
      "summary": "Fix login bug",
      "description": "Users cannot log in when...",
      "status": { "name": "Open" },
      "issuetype": { "name": "Bug" },
      "priority": { "name": "High" },
      "assignee": { "displayName": "John Doe" },
      "created": "2024-01-15T10:30:00.000+0000",
      "updated": "2024-01-16T14:20:00.000+0000"
    }
  }
]
```

### Table Format

```
KEY        TYPE    STATUS    PRIORITY    ASSIGNEE     SUMMARY
FOO-123    Bug     Open      High        John Doe     Fix login bug
```

## Use Cases

### Scripting

```bash
# Get issue and extract status
STATUS=$(gojira get FOO-123 -q | jq -r '.[0].fields.status.name')
echo "Issue status: $STATUS"

# Check if issue is closed
if [[ $(gojira get FOO-123 -q | jq -r '.[0].fields.status.name') == "Closed" ]]; then
  echo "Issue is closed"
fi
```

### AI Agent Workflows

```bash
# Get issue details for LLM context
gojira get FOO-123 --toon
```
