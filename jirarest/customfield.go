package jirarest

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/net/urlutil"
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

type CustomFieldsService struct {
	JRClient *Client
}

func NewCustomFieldsService(client *Client) *CustomFieldsService {
	return &CustomFieldsService{JRClient: client}
}

func (svc *CustomFieldsService) GetCustomFields() (CustomFields, error) {
	var cfs CustomFields
	if svc.JRClient == nil {
		return cfs, ErrJiraRESTClientCannotBeNil
	}
	apiURL := urlutil.JoinAbsolute(svc.JRClient.Config.ServerURL, APIURL2ListCustomFields)
	hclient := svc.JRClient.HTTPClient
	if hclient == nil {
		hclient = &http.Client{}
	}

	resp, err := hclient.Get(apiURL)
	if err != nil {
		return cfs, err
	}
	if resp.StatusCode >= 300 {
		return cfs, fmt.Errorf("error status code (%d)", resp.StatusCode)
	}
	_, err = jsonutil.UnmarshalReader(resp.Body, &cfs)
	return cfs, err
}

func (svc *CustomFieldsService) GetCustomFieldEpicLink() (CustomField, error) {
	return svc.GetCustomField(CustomFieldNameEpicLink)
}

func (svc *CustomFieldsService) GetCustomField(customFieldName string) (CustomField, error) {
	cfs, err := svc.GetCustomFields()
	if err != nil {
		return CustomField{}, err
	}
	cfsName := cfs.FilterByNames(customFieldName)
	if len(cfsName) != 1 {
		return CustomField{}, errors.New("epic link custom field not found")
	}
	return cfsName[0], nil
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
	tbl.Columns = []string{"ID", "Name", "Clause Names"}
	for _, cf := range cfs {
		row := []string{
			cf.ID, cf.Name, stringsutil.JoinLiteraryQuote(cf.ClauseNames, `"`, `"`, `, `, ""),
		}
		tbl.Rows = append(tbl.Rows, row)
	}
	return tbl
}

func (cfs CustomFields) WriteTable(w io.Writer) {
	cfs.SortByName(true)
	tbl := cfs.Table("")
	tw := tablewriter.NewWriter(w)
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
	} else if key, err := CustomFieldKeyCanonical(customFieldKey); err != nil {
		return err
	} else if unv, ok := iss.Fields.Unknowns[key]; !ok {
		return nil
	} else {
		return jsonutil.UnmarshalAny(unv, v)
	}
}

const customfieldPrefix = "customfield_"

var (
	ErrInvalidCustomFieldFormat = errors.New("invalid customfield format")
	rxCustomFieldBrackets       = regexp.MustCompile(`^cf\[([0-9]+)\]$`)
	rxCustomFieldCanonical      = regexp.MustCompile(`^customfield_[0-9]+$`)
	rxCustomFieldDigits         = regexp.MustCompile(`^[0-9]+$`)
)

// CustomFieldKeyCanonical converts a custom field string to `customfield_12345`.
func CustomFieldKeyCanonical(key string) (string, error) {
	key = strings.ToLower(strings.TrimSpace(key))
	if rxCustomFieldCanonical.MatchString(key) {
		return key, nil
	} else if rxCustomFieldDigits.MatchString(key) {
		return customfieldPrefix + key, nil
	} else if m := rxCustomFieldBrackets.FindAllStringSubmatch(key, -1); len(m) > 0 {
		n := m[0]
		if len(n) > 1 {
			return customfieldPrefix + n[1], nil
		}
	}
	return "", ErrInvalidCustomFieldFormat
}

func IsCustomFieldKey(key string) (string, bool) {
	if can, err := CustomFieldKeyCanonical(key); err != nil {
		return key, false
	} else {
		return can, true
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
