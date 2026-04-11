# export

Export Jira issues to JSON or XLSX format.

## Usage

```bash
gojira export [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `--jql` | JQL query to search issues |
| `--keys` | Comma-separated issue keys to export |
| `--json` | Output JSON file path |
| `--xlsx` | Output XLSX file path |
| `--from-json` | Read issues from existing JSON file instead of querying |
| `--include-parents` | Include parent issues in export |
| `--sheet` | Sheet name for XLSX export (default: "issues") |

Plus [global flags](index.md#global-flags).

## Examples

### Export to JSON

```bash
# Export search results to JSON
gojira export --jql "project = FOO" --json issues.json

# Export specific issues
gojira export --keys FOO-1,FOO-2,FOO-3 --json output.json
```

### Export to Excel

```bash
# Export to XLSX
gojira export --jql "project = FOO AND status = Open" --xlsx report.xlsx

# Custom sheet name
gojira export --jql "project = FOO" --xlsx report.xlsx --sheet "Open Issues"
```

### Export Both Formats

```bash
# Export to both JSON and XLSX
gojira export --jql "project = FOO" --json backup.json --xlsx report.xlsx
```

### Include Parent Issues

When exporting sub-tasks or stories linked to epics, include the parent issues:

```bash
gojira export --jql "project = FOO AND issuetype = Sub-task" --include-parents --xlsx subtasks.xlsx
```

### Convert JSON to XLSX

If you have an existing JSON export, convert it to XLSX:

```bash
gojira export --from-json issues.json --xlsx report.xlsx
```

## Output

### JSON Output

The JSON file contains the full issue data including all fields:

```json
{
  "issues": [
    {
      "key": "FOO-123",
      "fields": {
        "summary": "Fix login bug",
        "status": {"name": "Open"},
        ...
      }
    }
  ]
}
```

### XLSX Output

The Excel file contains a table with columns:

- Key
- Type
- Status
- Priority
- Summary
- Assignee
- Created
- Updated
- And more...

## Use Cases

### Backup Issues

```bash
# Regular backup of project issues
gojira export --jql "project = FOO" --json "backup-$(date +%Y%m%d).json"
```

### Generate Reports

```bash
# Sprint report
gojira export --jql "project = FOO AND sprint = 'Sprint 42'" --xlsx sprint-42.xlsx

# Unresolved bugs report
gojira export --jql "project = FOO AND type = Bug AND resolution is EMPTY" --xlsx bugs.xlsx
```

### Data Analysis

```bash
# Export for external analysis
gojira export --jql "project = FOO AND created >= -30d" --json recent.json

# Then use with jq, Python, etc.
cat recent.json | jq '.issues | length'
```
