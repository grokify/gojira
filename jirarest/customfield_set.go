package jirarest

import (
	"errors"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gojira"
)

type CustomFieldSet struct {
	Data map[string]CustomField
}

func NewCustomFieldSet() *CustomFieldSet {
	return &CustomFieldSet{Data: map[string]CustomField{}}
}

func (set *CustomFieldSet) Init() {
	if set.Data == nil {
		set.Data = map[string]CustomField{}
	}
}

func (set *CustomFieldSet) Add(fields ...CustomField) error {
	set.Init()
	for _, ci := range fields {
		id := strings.TrimSpace(ci.ID)
		if id == "" {
			return errors.New("custom field cannot have empty id")
		}
		set.Data[ci.ID] = ci
	}
	return nil
}

func (set *CustomFieldSet) IDToName(id string) (string, error) {
	if can, ok := gojira.IsCustomFieldKey(id); ok {
		if cf, ok := set.Data[can]; ok {
			if name := strings.TrimSpace(cf.Name); name != "" {
				return name, nil
			} else {
				return "", errors.New("customfield has no name")
			}
		} else {
			return "", errors.New("customfieldid not found")
		}
	} else {
		return "", errors.New("cannot parse custom field key")
	}
}

// CreateFuncIssueToMap creates a function to use with `IssuesSet.HistogramMapFunc`.
func (set *CustomFieldSet) CreateFuncIssueToMap(fieldsWithDefaults map[string]string, useCustomFieldDisplayNames bool) FuncIssueToMap {
	return func(iss *jira.Issue) (map[string]string, error) {
		out := map[string]string{}
		im := NewIssueMore(iss)
		for k, def := range fieldsWithDefaults {
			v := im.ValueOrDefault(k, def)
			if useCustomFieldDisplayNames {
				if can, ok := gojira.IsCustomFieldKey(k); !ok {
					out[k] = v
				} else if canName, err := set.IDToName(can); err != nil {
					return out, err
				} else {
					out[canName] = v
				}
			} else {
				out[k] = v
			}
		}
		return out, nil
	}
}
