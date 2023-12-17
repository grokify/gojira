package jirarest

import (
	"errors"
	"strconv"
	"strings"

	jira "github.com/andygrunwald/go-jira"
)

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
	c := Client{JiraClient: jclient}
	newEpics, err := c.GetIssuesSetForKeys(epicKeys)
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
