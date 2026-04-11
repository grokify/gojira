# fields

List and filter Jira custom fields.

## Usage

```bash
gojira fields [flags]
```

## Flags

| Flag | Description |
|------|-------------|
| `--id` | Filter by field ID(s), comma-separated |
| `--name` | Filter by name (partial match) |
| `--name-exact` | Filter by exact name |
| `--custom-only` | Show only custom fields (exclude system fields) |
| `--epic-link` | Show Epic Link field |
| `--json` | Output as JSON |
| `--table` | Output as table (default) |

Plus [global flags](index.md#global-flags).

## Examples

### List All Fields

```bash
# List all custom fields as table
gojira fields

# List as JSON
gojira fields --json
```

### Filter by ID

```bash
# Get specific field
gojira fields --id customfield_10001

# Get multiple fields
gojira fields --id customfield_10001,customfield_10002
```

### Search by Name

```bash
# Partial match search
gojira fields --name "Epic"

# Exact name match
gojira fields --name-exact "Epic Link"
```

### Custom Fields Only

```bash
# Exclude system fields
gojira fields --custom-only
```

### Get Epic Link Field

Find the Epic Link field (useful for linking issues to epics):

```bash
gojira fields --epic-link
```

## Output

### Table Format

```
ID                  NAME              TYPE         CUSTOM
customfield_10001   Epic Link         string       true
customfield_10002   Story Points      number       true
customfield_10003   Sprint            array        true
```

### JSON Format

```json
[
  {
    "id": "customfield_10001",
    "name": "Epic Link",
    "custom": true,
    "schema": {
      "type": "string"
    }
  }
]
```

## Use Cases

### Find Field IDs for Queries

```bash
# Find the Sprint field ID
gojira fields --name Sprint

# Use in search
gojira search --jql "project = FOO AND cf[10003] is not EMPTY"
```

### Discover Available Fields

```bash
# See all custom fields your Jira has
gojira fields --custom-only | head -20
```

### Script Integration

```bash
# Get Epic Link field ID
EPIC_FIELD=$(gojira fields --epic-link --json | jq -r '.[0].id')
echo "Epic Link field: $EPIC_FIELD"
```
