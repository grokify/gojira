package gojira

import (
	"fmt"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/mogo/type/stringsutil/join"
)

type JQL struct {
	IssuesIncl   []string
	IssuesExcl   []string
	ProjectsIncl []string
	ProjectsExcl []string
	StatusesIncl []string
	StatusesExcl []string
	TypesIncl    []string
	TypesExcl    []string
}

func (j JQL) String() string {
	var parts []string

	type inclExclProc struct {
		Field   string
		Values  []string
		Exclude bool
	}

	procs := []inclExclProc{
		{Field: FieldIssue, Values: j.IssuesIncl, Exclude: false},
		{Field: FieldIssue, Values: j.IssuesExcl, Exclude: true},
		{Field: FieldProject, Values: j.ProjectsIncl, Exclude: false},
		{Field: FieldProject, Values: j.ProjectsExcl, Exclude: true},
		{Field: FieldStatus, Values: j.StatusesIncl, Exclude: false},
		{Field: FieldStatus, Values: j.StatusesExcl, Exclude: true},
		{Field: FieldType, Values: j.TypesIncl, Exclude: false},
		{Field: FieldType, Values: j.TypesExcl, Exclude: true},
	}
	for _, proc := range procs {
		if clause := inClause(proc.Field, proc.Values, proc.Exclude); clause != "" {
			parts = append(parts, clause)
		}
	}

	if len(parts) > 0 {
		return strings.Join(parts, " AND ")
	} else {
		return ""
	}
}

func inClause(field string, values []string, exclude bool) string {
	field = strings.TrimSpace(field)
	values = stringsutil.SliceCondenseSpace(values, true, true)
	if field == "" || len(values) == 0 {
		return ""
	} else if len(values) == 1 {
		operator := "="
		if exclude {
			operator = "!="
		}
		qtr := stringsutil.Quoter{
			Beg:         "'",
			End:         "'",
			SkipNesting: true,
		}
		return fmt.Sprintf("%s %s %s", field, operator, qtr.Quote(values[0]))
	} else if len(values) > 1 {
		operator := "IN"
		if exclude {
			operator = "NOT IN"
		}
		return fmt.Sprintf("%s %s (%s)", field, operator, join.JoinQuote(values, "'", "'", ","))
	} else {
		return ""
	}
}
