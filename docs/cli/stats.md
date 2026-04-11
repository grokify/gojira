# stats

Show aggregate statistics for issues grouped by a field.

## Usage

```bash
gojira stats --jql <query> --by <field> [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--jql` | (required) | JQL query to search issues |
| `--by` | (required) | Field to group by (see [Grouping Fields](#grouping-fields)) |
| `--format` | `toon` | Output format: `toon`, `json`, or `table` |

Plus [global flags](index.md#global-flags).

## Grouping Fields

| Field | Description |
|-------|-------------|
| `status` | Issue status (Open, In Progress, Done, etc.) |
| `type` | Issue type (Bug, Story, Task, etc.) |
| `priority` | Issue priority (High, Medium, Low, etc.) |
| `assignee` | Assigned user |
| `project` | Project key |
| `resolution` | Resolution status |
| `customfield_XXXXX` | Any custom field by ID |

## Examples

### Count by Status

```bash
gojira stats --jql "project = FOO" --by status --format table
```

Output:

```
VALUE          COUNT      %
Open              45   30.0%
In Progress       30   20.0%
Done              75   50.0%
------         ------  ------
TOTAL            150  100.0%
```

### Count by Type

```bash
gojira stats --jql "project = FOO" --by type --format table
```

### Count by Priority

```bash
gojira stats --jql "project = FOO AND status != Done" --by priority --format table
```

### Count by Assignee

```bash
gojira stats --jql "project = FOO AND status = 'In Progress'" --by assignee --format table
```

### Count by Custom Field

```bash
# Find field ID first
gojira fields --name "Team"

# Then use it
gojira stats --jql "project = FOO" --by customfield_10005 --format table
```

## Output Formats

### TOON (default)

Token-optimized format for LLMs:

```bash
gojira stats --jql "project = FOO" --by status
```

### JSON

```bash
gojira stats --jql "project = FOO" --by status --format json
```

Output:

```json
{
  "field": "status",
  "total": 150,
  "results": [
    {"value": "Done", "count": 75},
    {"value": "Open", "count": 45},
    {"value": "In Progress", "count": 30}
  ]
}
```

### Table

```bash
gojira stats --jql "project = FOO" --by status --format table
```

## Use Cases

### Sprint Planning

```bash
# Status breakdown for sprint
gojira stats --jql "sprint = 'Sprint 42'" --by status --format table

# Type breakdown
gojira stats --jql "sprint = 'Sprint 42'" --by type --format table
```

### Workload Analysis

```bash
# Issues per assignee
gojira stats --jql "project = FOO AND status != Done" --by assignee --format table
```

### Project Health

```bash
# Unresolved issues by priority
gojira stats --jql "project = FOO AND resolution is EMPTY" --by priority --format table
```

### Script Integration

```bash
# Get count of open bugs
OPEN_BUGS=$(gojira stats --jql "project = FOO AND type = Bug AND status = Open" --by status --format json -q | jq '.total')
echo "Open bugs: $OPEN_BUGS"
```
