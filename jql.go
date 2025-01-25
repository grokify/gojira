package gojira

import (
	"fmt"
	"strings"
	"time"

	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/mogo/type/stringsutil/join"
)

// JQL is a JQL builder. It will create a JQL string using `JQL.String()` from the supplied infomration.
type JQL struct {
	Name            string
	FilterID        int
	CreatedGTE      time.Time
	CreatedLT       time.Time
	FiltersIncl     [][]string // outer level is `AND`, inner level is `IN`.
	FiltersExcl     [][]string
	IssuesIncl      [][]string
	IssuesExcl      [][]string
	KeysIncl        [][]string
	KeysExcl        [][]string
	LabelsIncl      [][]string
	LabelsExcl      [][]string
	ParentsIncl     [][]string
	ParentsExcl     [][]string
	ProjectsIncl    [][]string
	ProjectsExcl    [][]string
	ResolutionIncl  [][]string
	ResolutionExcl  [][]string
	StatusesIncl    [][]string
	StatusesExcl    [][]string
	TypesIncl       [][]string
	TypesExcl       [][]string
	Raw             []string
	CustomFieldIncl map[string][]string // slice is `IN`
	CustomFieldExcl map[string][]string
}

func (j JQL) String() string {
	var parts []string

	type inclExclProc struct {
		Field   string
		Values  [][]string
		Exclude bool
	}

	procs := []inclExclProc{
		{Field: FieldFilter, Values: j.FiltersIncl, Exclude: false},
		{Field: FieldFilter, Values: j.FiltersExcl, Exclude: true},
		{Field: FieldIssue, Values: j.IssuesIncl, Exclude: false},
		{Field: FieldIssue, Values: j.IssuesExcl, Exclude: true},
		{Field: FieldKey, Values: j.KeysIncl, Exclude: false},
		{Field: FieldKey, Values: j.KeysExcl, Exclude: true},
		{Field: FieldLabels, Values: j.LabelsIncl, Exclude: false},
		{Field: FieldLabels, Values: j.LabelsExcl, Exclude: true},
		{Field: FieldParent, Values: j.ParentsIncl, Exclude: false},
		{Field: FieldParent, Values: j.ParentsExcl, Exclude: true},
		{Field: FieldProject, Values: j.ProjectsIncl, Exclude: false},
		{Field: FieldProject, Values: j.ProjectsExcl, Exclude: true},
		{Field: FieldResolution, Values: j.ResolutionIncl, Exclude: false},
		{Field: FieldResolution, Values: j.ResolutionExcl, Exclude: true},
		{Field: FieldStatus, Values: j.StatusesIncl, Exclude: false},
		{Field: FieldStatus, Values: j.StatusesExcl, Exclude: true},
		{Field: FieldType, Values: j.TypesIncl, Exclude: false},
		{Field: FieldType, Values: j.TypesExcl, Exclude: true},
	}
	for _, proc := range procs {
		if field := strings.TrimSpace(proc.Field); field == "" {
			panic("field is empty")
		} else if len(proc.Values) > 0 {
			for _, inClauseVals := range proc.Values {
				if clause := inClause(proc.Field, inClauseVals, proc.Exclude); clause != "" {
					parts = append(parts, clause)
				}
			}
		}
		/*
			if len(proc.Values) == 0 {
				continue
			} else if field := strings.TrimSpace(proc.Field); field == "" {
				panic("field is empty")
				// } else if clause := inClause(proc.Field, proc.Values, proc.Exclude); clause != "" {
				// parts = append(parts, clause)
			}
		*/
	}

	if !j.CreatedGTE.IsZero() {
		parts = append(parts, fmt.Sprintf("%s >= %s", FieldCreatedDate, j.CreatedGTE.Format(timeutil.RFC3339FullDate)))
	}
	if !j.CreatedLT.IsZero() {
		parts = append(parts, fmt.Sprintf("%s < %s", FieldCreatedDate, j.CreatedLT.Format(timeutil.RFC3339FullDate)))
	}
	for cfk, cfv := range j.CustomFieldIncl {
		cfv = stringsutil.SliceCondenseSpace(cfv, true, false)
		if len(cfv) == 0 {
			continue
		}
		if cfkCanonicalID, err := CustomFieldLabelToID(cfk); err == nil {
			cfk = cfkCanonicalID.StringBrackets()
		}
		if clause := inClause(cfk, cfv, false); clause != "" {
			parts = append(parts, clause)
		}
	}
	for cfk, cfv := range j.CustomFieldExcl {
		cfv = stringsutil.SliceCondenseSpace(cfv, true, false)
		if len(cfv) == 0 {
			continue
		}
		if cfkCanonicalID, err := CustomFieldLabelToID(cfk); err == nil {
			cfk = cfkCanonicalID.StringBrackets()
		}
		if clause := inClause(cfk, cfv, true); clause != "" {
			parts = append(parts, clause)
		}
	}

	parts = append(parts, j.Raw...)

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
		return fmt.Sprintf("%s %s (%s)", field, operator, join.JoinQuote(values, "'", "'", JQLInSep))
	} else {
		return ""
	}
}

// JQLStringsSimple provides a set of JQLs for a single field and values. The purpose of this function
// is to split very long lists of values so that each JQL is under a certain length limit.
func JQLStringsSimple(field string, exclude bool, vals []string, jqlMaxLength int) []string {
	if jqlMaxLength < 0 {
		jqlMaxLength *= -1
	}
	field = strings.TrimSpace(field)
	if field == "" {
		return []string{}
	}
	vals = stringsutil.SliceCondenseSpace(vals, true, true)
	if len(vals) == 0 {
		return []string{}
	}
	var jqls []string
	operator := "IN"
	if exclude {
		operator = "NOT IN"
	}
	baseString := fmt.Sprintf("%s %s ()", field, operator)
	baseStringLen := len(baseString)
	quoter := stringsutil.Quoter{
		Beg:         "'",
		End:         "'",
		SkipNesting: true,
	}
	if jqlMaxLength == 0 {
		jqlMaxLength = JQLMaxLength
	}
	valsMaxLength := jqlMaxLength - baseStringLen
	valsString := ""
	for i, val := range vals {
		valQuoted := quoter.Quote(val)
		if len(valsString)+len(valQuoted) > valsMaxLength {
			jqls = append(jqls, fmt.Sprintf("%s %s (%s)", field, operator, valsString))
			valsString = ""
		}
		valsString += valQuoted
		if i < len(vals)-1 {
			valsString += JQLInSep
		}
	}
	if len(valsString) > 0 {
		jqls = append(jqls, fmt.Sprintf("%s %s (%s)", field, operator, valsString))
	}
	return jqls
}

type JQLs []JQL

func (j JQLs) JoinString(keyword string) string {
	var parts []string
	for _, jql := range j {
		parts = append(parts, "("+jql.String()+")")
	}
	return strings.Join(parts, " "+keyword+" ")
}
