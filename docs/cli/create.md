# create

Create a Jira issue from a YAML file.

## Usage

```bash
gojira create -f <file> [flags]
```

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--file` | `-f` | YAML file containing issue data (required) |
| `--dry-run` | | Validate and preview without creating |
| `--project` | | Override project key from file |
| `--parent` | | Override parent issue key from file |
| `--type` | | Override issue type from file |
| `--json` | `-j` | Output full result as JSON |

## YAML Format

The YAML file supports all standard Jira fields plus custom fields.

### Standard Fields

| Field | Required | Description |
|-------|----------|-------------|
| `project` | Yes | Project key (e.g., PROJ) |
| `type` | Yes | Issue type (Story, Bug, Task, Epic) |
| `summary` | Yes | Issue title |
| `description` | No | Issue description (supports Jira markdown) |
| `parent` | No | Parent issue key for subtasks or stories under epics |
| `labels` | No | List of labels |
| `priority` | No | Priority name (High, Medium, Low) |
| `assignee` | No | Username or email |
| `reporter` | No | Username or email |
| `components` | No | List of component names |
| `fix_versions` | No | List of fix version names |

### Custom Fields

Any field starting with `customfield_` is passed directly to the Jira API:

```yaml
customfield_12345: "Value for custom field"
customfield_10001: "Q2-2024"
```

## Examples

### Basic Story

```yaml
project: PROJ
type: Story
summary: Add user authentication
description: |
  Implement OAuth2 login flow for the application.
labels:
  - auth
  - mvp
```

```bash
gojira create -f story.yaml
# Output: Created PROJ-123
```

### Full Example with Custom Fields

```yaml
project: PROJ
type: Story
summary: Add user authentication via OAuth2
description: |
  ## Background
  Users currently have no way to securely authenticate. We need to implement
  OAuth2 login flow to support SSO with corporate identity providers.

  ## Requirements
  - Support Google and Microsoft identity providers
  - Store refresh tokens securely
  - Handle token expiration gracefully

  ## Technical Notes
  Use the existing session middleware for token validation.
parent: PROJ-100
labels:
  - auth
  - mvp
  - security
priority: High
assignee: john.doe
components:
  - backend
  - security
customfield_12345: |
  Given a user is on the login page
  When they click "Sign in with Google"
  Then they should be redirected to Google OAuth
  And upon success, redirected back to the dashboard
customfield_10001: Q2-2024
```

### Dry Run

Preview what would be created without actually creating:

```bash
gojira create -f story.yaml --dry-run
```

Output:

```json
{
  "valid": true,
  "project": "PROJ",
  "type": "Story",
  "summary": "Add user authentication via OAuth2",
  "description": "## Background\n...",
  "parent": "PROJ-100",
  "labels": ["auth", "mvp", "security"],
  "priority": "High",
  "assignee": "john.doe",
  "custom_fields": {
    "customfield_12345": "Given a user is on the login page..."
  }
}
```

### Override Project or Parent

```bash
# Create in a different project
gojira create -f story.yaml --project OTHER

# Create under a different epic
gojira create -f story.yaml --parent EPIC-456

# Override issue type
gojira create -f story.yaml --type Task
```

### JSON Output

Get full response including issue ID and self URL:

```bash
gojira create -f story.yaml --json
```

Output:

```
Created PROJ-123
{
  "key": "PROJ-123",
  "id": "10001",
  "self": "https://company.atlassian.net/rest/api/2/issue/10001",
  "summary": "Add user authentication via OAuth2"
}
```

## Multiline Strings

Use YAML block scalars for multiline content:

| Syntax | Behavior |
|--------|----------|
| `\|` | Preserves newlines (literal) |
| `>` | Folds newlines to spaces |

```yaml
# Literal block - preserves newlines
description: |
  Line 1
  Line 2

  Line 4 after blank

# Folded block - joins lines
description: >
  This becomes one long
  paragraph joined together.
```

For markdown code blocks, use the literal style (`|`):

```yaml
description: |
  ## Example

  ```go
  func main() {
      fmt.Println("Hello")
  }
  ```
```
