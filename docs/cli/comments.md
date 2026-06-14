# comments

Get comments for a Jira issue.

## Usage

```bash
gojira comments <issue-key> [flags]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `issue-key` | The Jira issue key (e.g., `FOO-123`) |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--max` | 50 | Maximum number of comments to return |

Plus [global flags](index.md#global-flags).

## Examples

### Basic Usage

```bash
# Get comments for an issue
gojira comments FOO-123

# Limit to 10 most recent comments
gojira comments FOO-123 --max 10
```

### Scripting

```bash
# Get comment count
gojira comments FOO-123 -q | jq '.total'

# Get latest comment body
gojira comments FOO-123 --max 1 -q | jq -r '.comments[0].body'

# List all comment authors
gojira comments FOO-123 -q | jq -r '.comments[].author'
```

## Output

### JSON Format

```json
{
  "key": "FOO-123",
  "total": 3,
  "comments": [
    {
      "id": "10001",
      "author": "John Doe",
      "body": "I've started working on this issue.",
      "created": "2024-01-15T10:30:00.000+0000",
      "updated": "2024-01-15T10:30:00.000+0000"
    },
    {
      "id": "10002",
      "author": "Jane Smith",
      "body": "Please check the latest commit.",
      "created": "2024-01-16T14:20:00.000+0000",
      "updated": "2024-01-16T14:25:00.000+0000"
    }
  ]
}
```

## Use Cases

### Review Discussion History

```bash
# Get full comment thread for context
gojira comments FOO-123
```

### AI Agent Workflows

```bash
# Get recent discussion for LLM context
gojira comments FOO-123 --max 5 -q
```

### Integration with Other Tools

```bash
# Export comments to a file
gojira comments FOO-123 > comments.json

# Combine with issue details
echo '{"issue":' && gojira get FOO-123 -q && echo ',"comments":' && gojira comments FOO-123 -q && echo '}'
```
