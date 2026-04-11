# patch

Update one or more fields on a Jira issue.

## Usage

```bash
gojira patch <issue-key> [flags]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `issue-key` | The Jira issue key to update (e.g., `FOO-123`) |

## Flags

| Flag | Description |
|------|-------------|
| `--set` | Set field value (format: `field=value`). Can be repeated. |
| `--add-label` | Add a label to the issue. Can be repeated. |
| `--remove-label` | Remove a label from the issue. Can be repeated. |
| `--json` | Raw JSON body for complex updates |
| `--dry-run` | Show request body without executing |
| `--show-after` | Show issue after update |

Plus [global flags](index.md#global-flags).

## Examples

### Setting Fields

```bash
# Set a simple field
gojira patch FOO-123 --set summary="New summary"

# Set multiple fields
gojira patch FOO-123 --set summary="New title" --set priority=High

# Set a custom field
gojira patch FOO-123 --set customfield_10001="custom value"
```

### Managing Labels

```bash
# Add a label
gojira patch FOO-123 --add-label bug

# Add multiple labels
gojira patch FOO-123 --add-label bug --add-label urgent

# Remove a label
gojira patch FOO-123 --remove-label obsolete

# Add and remove in one command
gojira patch FOO-123 --add-label reviewed --remove-label needs-review
```

### Complex Updates with JSON

```bash
# Use JSON for complex field structures
gojira patch FOO-123 --json '{"fields":{"summary":"New title","labels":["bug","urgent"]}}'

# Update nested fields
gojira patch FOO-123 --json '{"fields":{"assignee":{"name":"jdoe"}}}'
```

### Dry Run

Preview what would be sent without making changes:

```bash
gojira patch FOO-123 --set summary="Test" --dry-run
```

Output:

```
Would PATCH FOO-123 with:
{
  "fields": {
    "summary": {
      "value": "Test"
    }
  }
}
```

### Show Updated Issue

```bash
# Show issue after update
gojira patch FOO-123 --set summary="Updated" --show-after
```

## Field Names

### Standard Fields

| Field | Example |
|-------|---------|
| `summary` | `--set summary="New summary"` |
| `description` | `--set description="New description"` |
| `priority` | `--set priority=High` |

### Custom Fields

Use the field ID (e.g., `customfield_10001`). Use `gojira fields` to find field IDs:

```bash
# List custom fields
gojira fields --custom-only

# Use in patch
gojira patch FOO-123 --set customfield_10001="value"
```

## Response

On success:

```
Successfully updated FOO-123 (status: 204)
```

On failure, an error message with the HTTP status code is displayed.
