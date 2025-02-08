package gojira

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/grokify/mogo/type/stringsutil/join"
)

type JQLMeta struct {
	Name            string
	Key             string
	Description     string
	FilterID        int
	QueryTime       time.Time
	QueryTotalCount int
}

// JQL is a JQL builder. It will create a JQL string using `JQL.String()` from the supplied infomration.
type JQL struct {
	Meta            JQLMeta // Not part of JQL
	CreatedGT       *time.Time
	CreatedGTE      *time.Time
	CreatedLT       *time.Time
	CreatedLTE      *time.Time
	UpdatedGT       *time.Time
	UpdatedGTE      *time.Time
	UpdatedLT       *time.Time
	UpdatedLTE      *time.Time
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
	Any             JQLAndOrStringer
}

type JQLAndOrStringer [][]fmt.Stringer

// JQLOrAndStringer combines `fmt.Stringer` slice of slice to create a JQL.
// Outer slide is "AND", Inner slice is "OR".
func (j JQLAndOrStringer) String() string {
	if len(j) == 0 {
		return ""
	}
	var andClauses []string
	for _, orClausesRaw := range j {
		if len(orClausesRaw) == 0 {
			continue
		}
		var orClausesStr []string
		for _, orClauseRaw := range orClausesRaw {
			orClauseStr := strings.TrimSpace(orClauseRaw.String())
			if orClauseStr == "" {
				continue
			} else {
				orClausesStr = append(orClausesStr, addParen(orClauseStr))
			}
		}
		andClauses = append(andClauses, addParen(strings.Join(orClausesStr, " OR ")))
	}
	return strings.Join(andClauses, " AND ")
}

func addParen(s string) string {
	return addPrefixSuffix("(", ")", s)
}

func addPrefixSuffix(prefix, suffix, s string) string {
	return prefix + s + suffix
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

	if clauses := j.clausesCreated(); len(clauses) > 0 {
		parts = append(parts, clauses...)
	}
	if clauses := j.clausesUpdated(); len(clauses) > 0 {
		parts = append(parts, clauses...)
	}
	if clauses := j.clausesCustomFields(); len(clauses) > 0 {
		parts = append(parts, clauses...)
	}
	if clause := j.Any.String(); clause != "" {
		parts = append(parts, clause)
	}
	parts = append(parts, j.Raw...)

	if len(parts) > 0 {
		return strings.Join(parts, " AND ")
	} else {
		return ""
	}
}

func (j JQL) clausesCreated() []string {
	var clauses []string
	if j.CreatedGT != nil {
		clauses = append(clauses, fmtFieldOperatorDate(FieldCreatedDate, OperatorGT, *j.CreatedGT))
	}
	if j.CreatedGTE != nil {
		clauses = append(clauses, fmtFieldOperatorDate(FieldCreatedDate, OperatorGTE, *j.CreatedGTE))
	}
	if j.CreatedLT != nil {
		clauses = append(clauses, fmtFieldOperatorDate(FieldCreatedDate, OperatorLT, *j.CreatedLT))
	}
	if j.CreatedLTE != nil {
		clauses = append(clauses, fmtFieldOperatorDate(FieldCreatedDate, OperatorLTE, *j.CreatedLTE))
	}
	return clauses
}

func (j JQL) clausesUpdated() []string {
	var clauses []string
	if j.UpdatedGT != nil {
		clauses = append(clauses, fmtFieldOperatorDate(FieldUpdated, OperatorGT, *j.UpdatedGT))
	}
	if j.UpdatedGTE != nil {
		clauses = append(clauses, fmtFieldOperatorDate(FieldUpdated, OperatorGTE, *j.UpdatedGTE))
	}
	if j.UpdatedLT != nil {
		clauses = append(clauses, fmtFieldOperatorDate(FieldUpdated, OperatorLT, *j.UpdatedLT))
	}
	if j.UpdatedLTE != nil {
		clauses = append(clauses, fmtFieldOperatorDate(FieldUpdated, OperatorLTE, *j.UpdatedLTE))
	}
	return clauses
}

func (j JQL) clausesCustomFields() []string {
	var clauses []string
	for cfk, cfv := range j.CustomFieldIncl {
		cfv = stringsutil.SliceCondenseSpace(cfv, true, false)
		if len(cfv) == 0 {
			continue
		}
		if cfkCanonicalID, err := CustomFieldLabelToID(cfk); err == nil {
			cfk = cfkCanonicalID.StringBrackets()
		}
		if clause := inClause(cfk, cfv, false); clause != "" {
			clauses = append(clauses, clause)
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
			clauses = append(clauses, clause)
		}
	}
	return clauses
}

func fmtFieldOperatorDate(field, op string, dt time.Time) string {
	return fmt.Sprintf("%s %s %s", field, op, dt.Format(timeutil.RFC3339FullDate))
}

func (j JQL) QueryString() string {
	return "jql=" + url.QueryEscape(j.String())
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
