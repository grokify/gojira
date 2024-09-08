package jirarest

import (
	"errors"
	"io"
	"sort"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/olekukonko/tablewriter"
)

const (
	CustomFieldNameEpicLink = "Epic Link"
)

var ErrJiraRESTClientCannotBeNil = errors.New("jirarest.Client cannot be nil")

type CustomFields []CustomField

type CustomField struct {
	ID               string            `json:"id"` // "customfield_12345"
	Key              string            `json:"key"`
	Name             string            `json:"name"`
	UntranslatedName string            `json:"untranslatedName"`
	Custom           bool              `json:"custom"`
	Orderable        bool              `json:"orderable"`
	Navigable        bool              `json:"navigable"`
	Searchable       bool              `json:"searchable"`
	ClauseNames      []string          `json:"clauseNames"`
	Schema           CustomFieldSchema `json:"schema"`
}

type CustomFieldSchema struct {
	Type     string `json:"type"`
	Custom   string `json:"custom"`
	CustomID int    `json:"customId"`
}

func (cfs CustomFields) SortByName(asc bool) CustomFields {
	if asc {
		sort.Slice(cfs, func(i, j int) bool {
			return cfs[i].Name < cfs[j].Name
		})
	} else {
		sort.Slice(cfs, func(i, j int) bool {
			return cfs[i].Name > cfs[j].Name
		})
	}
	return cfs
}

func (cfs CustomFields) FilterByIDs(ids ...string) CustomFields {
	filtered := CustomFields{}
	if len(ids) == 0 {
		return filtered
	}
	ids = slicesutil.Dedupe(ids)
	idsMap := map[string]int{}
	for _, id := range ids {
		idsMap[id] = 1
	}
	for _, cf := range cfs {
		if _, ok := idsMap[cf.ID]; ok {
			filtered = append(filtered, cf)
			if len(filtered) == len(ids) {
				return filtered
			}
		}
	}
	return filtered
}

func (cfs CustomFields) FilterByNames(names ...string) CustomFields {
	filtered := CustomFields{}
	if len(names) == 0 {
		return filtered
	}
	names = slicesutil.Dedupe(names)
	namesMap := map[string]int{}
	for _, name := range names {
		namesMap[name] = 1
	}
	for _, cf := range cfs {
		if _, ok := namesMap[cf.Name]; ok {
			filtered = append(filtered, cf)
		}
	}
	return filtered
}

func (cfs CustomFields) Table(name string) table.Table {
	if strings.TrimSpace(name) == "" {
		name = "Custom Fields"
	}
	tbl := table.NewTable(name)
	tbl.Columns = []string{"Name", "ID", "Clause Names"}
	for _, cf := range cfs {
		row := []string{
			cf.Name, cf.ID, stringsutil.JoinLiteraryQuote(cf.ClauseNames, `"`, `"`, `, `, ""),
		}
		tbl.Rows = append(tbl.Rows, row)
	}
	return tbl
}

func (cfs CustomFields) WriteTable(w io.Writer) {
	cfs.SortByName(true)
	tbl := cfs.Table("")
	tw := tablewriter.NewWriter(w)
	tw.SetRowLine(true)
	tw.SetRowSeparator("-")
	tw.SetHeader(tbl.Columns)
	tw.AppendBulk(tbl.Rows)
	tw.Render()
}

// IssueFieldsCustomFieldString returns a string custom field, e.g "Epic Link"
func IssueFieldsCustomFieldString(fields *jira.IssueFields, id string) string {
	if fields == nil {
		return ""
	}
	val, err := fields.Unknowns.String(id)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(val)
}

// GetCustomValueString attempts to return a string if either the custom value is a simple string
// or is an `IssueCustomField`, in which case it returns the `value` property.
func GetCustomValueString(iss jira.Issue, customFieldKey string) (string, error) {
	if iss.Fields == nil {
		return "", nil
	}
	any, ok := iss.Fields.Unknowns[customFieldKey]
	if !ok {
		return "", nil
	}
	if strval, ok := any.(string); ok {
		return strval, nil
	}
	icf := &IssueCustomField{}
	err := GetUnmarshalCustomValue(iss, customFieldKey, icf)
	if err != nil {
		return "", err
	}
	return icf.Value, nil
}

// GetUnmarshalCustomValue can be used to unmarshal a value to `IssueCustomField{}`.
func GetUnmarshalCustomValue(iss jira.Issue, customFieldKey string, v *IssueCustomField) error {
	if iss.Fields == nil {
		return nil
	} else if key, err := gojira.CustomFieldKeyCanonical(customFieldKey); err != nil {
		return err
	} else if unv, ok := iss.Fields.Unknowns[key]; !ok {
		return nil
	} else {
		return jsonutil.UnmarshalAny(unv, v)
	}
}

type IssueCustomField struct {
	ID    string `json:"id"`
	Self  string `json:"self"`
	Value string `json:"value"`
}

/*
func UnmarshalAny(data, v any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
*/
