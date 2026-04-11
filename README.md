# GoJira

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/gojira/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/gojira/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/gojira/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/gojira/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/gojira/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/gojira/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/gojira
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/gojira
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/gojira
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/gojira
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fgojira
 [loc-svg]: https://tokei.rs/b1/github/grokify/gojira
 [repo-url]: https://github.com/grokify/gojira
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/gojira/blob/master/LICENSE

GoJira is a Go SDK and CLI for Jira that provides:

- **REST API client** (`rest/`) - wrapper around [`go-jira`](https://github.com/andygrunwald/go-jira) with additional utilities
- **XML parser** (`xml/`) - parse Jira XML exports when API access is unavailable
- **JQL builder** (root package) - programmatically construct JQL queries
- **CLI tool** (`cmd/gojira/`) - command-line interface optimized for AI agents and humans

## Installation

```bash
# Install the CLI
go install github.com/grokify/gojira/cmd/gojira@latest

# Use as a library
go get github.com/grokify/gojira
```

## Quick Start

### CLI Usage

```bash
# Set credentials via environment variables
export JIRA_URL=https://your-instance.atlassian.net
export JIRA_USER=your-email@example.com
export JIRA_TOKEN=your-api-token

# Search issues with JQL
gojira search --jql "project = FOO AND status = Open"

# Get a specific issue
gojira get ISSUE-123

# Output formats
gojira search --jql "assignee = currentUser()" --table  # Human-readable
gojira search --jql "project = FOO" --json              # Machine-readable (default)
gojira search --jql "project = FOO" --toon              # Token-optimized for LLMs
```

### Library Usage

```go
import "github.com/grokify/gojira/rest"

// Create client with basic auth
client, err := rest.NewClientFromBasicAuth(
    "https://your-instance.atlassian.net",
    "your-email@example.com",
    "your-api-token",
    false,
)
if err != nil {
    log.Fatal(err)
}

// Search issues
issues, err := client.IssueAPI.SearchIssues("project = FOO", false)
if err != nil {
    log.Fatal(err)
}

for _, issue := range issues {
    im := rest.NewIssueMore(&issue)
    fmt.Printf("%s: %s [%s]\n", im.Key(), im.Summary(), im.Status())
}
```

### JQL Builder

```go
import "github.com/grokify/gojira"

// Build JQL programmatically
jql := gojira.JQL{
    Project: "FOO",
    Status:  "Open",
}
query := jql.String() // "project = FOO AND status = Open"
```

## Package Structure

| Package | Description | Dependencies |
|---------|-------------|--------------|
| `gojira` | JQL builder, config, constants | None (lightweight) |
| `gojira/rest` | REST API client | go-jira SDK |
| `gojira/xml` | XML export parser | None |
| `gojira/web` | URL helpers | None |

## Use Cases

1. Programmatically construct JQL queries
1. Generate Markdown reports from JQL results
1. Parse Jira XML exports when API access is unavailable
1. Automate Jira operations via CLI (ideal for AI agents)

## URL Formats

Accessing a list of issues by JQL is avialable via the UI and API:

* UI: `https://{jira_host}/issues/?jql=`
* API: `https://{jira_host}/rest/api/2/search?jql=`

### REST API Authentication: Basic Auth

The API auth can be provided by Basic Auth using an personsal API Token.

* Docs: https://developer.atlassian.com/cloud/jira/platform/basic-auth-for-rest-apis/
* Token Page: https://id.atlassian.com/manage-profile/security/api-tokens

## Note on Hours Per Day and Days Per Week

This module supports custom `hoursPerDay` and `daysPerWeek` settings per Jira.

This is described here and set in the UI via the screenshot below,

Ref: https://community.atlassian.com/t5/Jira-Software-questions/What-it-JIRA-counting-as-a-quot-day-quot-in-Time-Tracking/qaq-p/1703409

Also of note is that the hours per day can be set to a decimal value, such as `8.5`, but the UI may not show it:

Ref: https://community.atlassian.com/t5/Jira-questions/change-quot-Working-hours-per-day-quot-by-a-decimal-value/qaq-p/583095

![](ss_jira_time-tracking.png)

## Additional Discussion on Jira XML

### General Discussion

General discussion including using Jira XML to:

1. export comments and issue link types
1. create CSV for flexible reporting and import

Ref: https://community.atlassian.com/t5/Jira-questions/JIRA-Issue-XML-Export-What-is-it-good-for/qaq-p/603308

### Global Config

Working Hours Per Day and Working Days Per Week are global values and cannot be set on a per-project basis.

Ref: https://community.atlassian.com/t5/Jira-Software-questions/Time-Tracking-Hours-Is-it-still-a-global-change/qaq-p/1337399

# JQL Examples

| Goal | Example |
|------|---------|
| Query by key | key = ABC-123 |
| Query by parent | parent = ABC-123 |
| Query by linked issue | issue in linkedIssues (ABC-123) |
| Query by reporter | reporter = "foo@bar.com" |

https://community.atlassian.com/t5/Jira-questions/How-to-search-all-linked-issues-with-issues-from-specific/qaq-p/1027269

## Backlog

Downloading the project backlog view appears to be challenging. Here is some info on attempts to do this.

* [What is the equivalent JQL query for a scrum board backlog](https://community.atlassian.com/t5/Jira-Software-questions/What-is-the-equivalent-JQL-query-for-a-scrum-board-backlog/qaq-p/868136)
* [How can I export the Backlog to csv or Excel or TXT?](https://community.atlassian.com/t5/Jira-Software-questions/How-can-I-export-the-Backlog-to-csv-or-Excel-or-TXT/qaq-p/1322548)

Roughly:

`project = <project_name> AND resolution = Unresolved AND status!=Closed AND (Sprint not in openSprints() OR Sprint is EMPTY) AND type not in (Epic, Sub-Task) ORDER BY Rank ASC`

## Agile Life Cycle Stages

Stage is a common way to understand the stages of development given that many companies can use different statuses and workflows.

Stage the following stages (from [Asana](https://asana.com/resources/agile-methodology)) and will provide a grouping capability for implementation-specific workflows into these standard meta statuses for consistent and canonical understanding of the process.

Not all the stages need to be used, however, the stages that are used can be understood in a canonical way.

![](docs/agile-methodology_asana.png)
*image courtesy of Asana*

An alternate, but related agile methodology is avialable from [eSparkBiz](https://www.esparkinfo.com/blog/difference-between-product-backlog-grooming-and-sprint-planning.html):

![](docs/agile-methodology_esparkbiz.png)
*image courtesy of eSparkBiz*

As well as from [BISS](https://handbook.bosphorusiss.com/way-of-working/process-design/agile-software-development-lifecycle/):

![](docs/agile-methodology_biss.png)
*image courtesy of BISS*

This is described as ["The 5 Stages of the Agile Software Development Lifecycle"](https://www.mendix.com/blog/agile-software-development-lifecycle-stages/) by Mendix.

![](docs/agile-methodology_mendix.png)
*image courtesy of Mendix*