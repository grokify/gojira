package jirarest

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/mogo/type/stringsutil"
)

func KeysJQL(keys []string) string {
	keys = stringsutil.SliceCondenseSpace(keys, true, false)
	if len(keys) == 0 {
		return ""
	}
	return fmt.Sprintf("key in (%s)", strings.Join(keys, ","))
}

/*
func GetEpics(epicKeys []string) IssuesSet {
	jql := KeysJQL(epicKeys)
	if jql == "" {
		return NewIssuesSet(nil)
	}

}
*/

type EpicsSet struct {
	EpicsMap map[string]jira.Epic
}

func NewEpicsSet() EpicsSet {
	return EpicsSet{EpicsMap: map[string]jira.Epic{}}
}

func (es *EpicsSet) GetKeys(jclient *jira.Client, epicKeys []string) error {
	if jclient == nil {
		return errors.New("jclient cannot be nil")
	}
	newEpics, err := GetIssuesSetForKeys(jclient, epicKeys)
	if err != nil {
		return err
	}
	err = es.AddIssues(newEpics.Issues())
	return err
}

func (es *EpicsSet) AddIssues(issues []jira.Issue) error {
	if es.EpicsMap == nil {
		es.EpicsMap = map[string]jira.Epic{}
	}
	for _, iss := range issues {
		epic, err := IssueToEpic(iss)
		if err != nil {
			return err
		}
		es.EpicsMap[epic.Key] = *epic
	}
	return nil
}

func IssueToEpic(iss jira.Issue) (*jira.Epic, error) {
	idInt, err := strconv.Atoi(iss.ID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(iss.Key) == "" {
		return nil, errors.New("key is empty")
	}
	epic := &jira.Epic{
		ID:  idInt,
		Key: iss.Key,
	}
	if iss.Fields != nil {
		epic.Summary = iss.Fields.Summary
	}
	return epic, nil
}

/*
type Epic struct {
	ID      int    `json:"id" structs:"id"`
	Key     string `json:"key" structs:"key"`
	Self    string `json:"self" structs:"self"`
	Name    string `json:"name" structs:"name"`
	Summary string `json:"summary" structs:"summary"`
	Done    bool   `json:"done" structs:"done"`
}
*/
