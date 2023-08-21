# JiraXML

[![Build Status][build-status-svg]][build-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

 [build-status-svg]: https://github.com/grokify/go-jiraxml/workflows/build/badge.svg
 [build-status-url]: https://github.com/grokify/go-jiraxml/actions
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/go-jiraxml
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/go-jiraxml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/go-jiraxml
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/go-jiraxml
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/go-jiraxml/blob/master/LICENSE

This module contains code to accss Jira, parsing both the JSON API via [`github.com/andygrunwald/go-jira`](https://github.com/andygrunwald/go-jira) in `jirarest` and code to parse a Jira XML file in `jiraxml`.

Various aggregate staticstics and reports are calculated/generated.


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

# Backlog

Downloading the project backlog view appears to be challenging. Here is some info on attempts to do this.

* [What is the equivalent JQL query for a scrum board backlog](https://community.atlassian.com/t5/Jira-Software-questions/What-is-the-equivalent-JQL-query-for-a-scrum-board-backlog/qaq-p/868136)
* [How can I export the Backlog to csv or Excel or TXT?](https://community.atlassian.com/t5/Jira-Software-questions/How-can-I-export-the-Backlog-to-csv-or-Excel-or-TXT/qaq-p/1322548)

Roughly:

`project = <project_name> AND resolution = Unresolved AND status!=Closed AND (Sprint not in openSprints() OR Sprint is EMPTY) AND type not in (Epic, Sub-Task) ORDER BY Rank ASC`
