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

JiraXML currently parses a Jira XML file consisting of multiple issues.

In addition to parsing the Jira XML into a Go struct, various aggregate staticstics are calculated.

## Note on Hours Per Day and Days Per Week

This module supports custom `hoursPerDay` and `daysPerWeek` settings per Jira.

This is described here and set in the UI via the screenshot below,

https://community.atlassian.com/t5/Jira-Software-questions/What-it-JIRA-counting-as-a-quot-day-quot-in-Time-Tracking/qaq-p/1703409

![](ss_jira_time-tracking.png)