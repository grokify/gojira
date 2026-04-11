# GoJira Tasks

Open items identified from project assessment.

## High Priority

All high priority tasks completed.

## Medium Priority

- [x] Add architecture documentation (added to README.md)
- [ ] Clean up commented-out code blocks

## Low Priority

- [ ] Add issue creation helpers (beyond patching)
- [ ] Add sprint/board support (Agile API)
- [ ] Add comment/attachment management
- [ ] Add webhook/event handling support
- [ ] Improve error messages with more context

## Completed

- [x] Implement gojira CLI (search, get, version commands)
- [x] Add cobra dependency for CLI
- [x] Implement JSON, Table, TOON output formats
- [x] Add flexible authentication (env vars, goauth files, CLI flags)
- [x] Add unit tests for IssueService operations (search, get, patch)
- [x] Fix lint issue in `rest/issue_service.go:88` (ineffectual assignment)
- [x] Unify logging to slog (remove zerolog dependency)
- [x] Rename `SearchIssuesDeprecated` to `SearchIssuesOnPremise` with clear documentation
- [x] Add godoc comments to exported functions in rest/
- [x] Refactor package names: `jirarest/` → `rest/`, `jiraxml/` → `xml/`, `jiraweb/` → `web/`
- [x] Add integration tests with mock Jira server (supports live server via env vars)
