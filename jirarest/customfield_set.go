package jirarest

import (
	"errors"
	"strings"
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
