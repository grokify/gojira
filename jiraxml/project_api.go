package jiraxml

import (
	"encoding/json"
	"os"

	jira "github.com/andygrunwald/go-jira"
)

// /rest/api/3/project/search?jql=&maxResults=200

type ProjectsAPIResponse struct {
	Self       string   `json:"self"`
	NextPage   string   `json:"nextPage"`
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	Total      int      `json:"total"`
	IsLast     bool     `json:"isLast"`
	Values     Projects `json:"values"`
}

type Projects []jira.Project

func (pp Projects) ProjectsMetasMap() ProjectsMetasMap {
	m := ProjectsMetasMap{Data: map[string]ProjectMeta{}}
	for _, p := range pp {
		m.Data[p.Key] = ProjectMeta{Key: p.Key, Name: p.Name}
	}
	return m
}

func NewProjectsAPIResponseFile(filename string) (ProjectsAPIResponse, error) {
	r := ProjectsAPIResponse{}
	b, err := os.ReadFile(filename)
	if err != nil {
		return r, err
	}
	return r, json.Unmarshal(b, &r)
}

type ProjectMeta struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type ProjectsMetasMap struct {
	Data map[string]ProjectMeta
}
