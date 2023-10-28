package jirarest

import (
	"strings"
	"time"

	jira "github.com/andygrunwald/go-jira"
)

type IssueMore struct {
	Issue *jira.Issue
}

func (im *IssueMore) AsigneeName() string {
	if im.Issue == nil || im.Issue.Fields == nil || im.Issue.Fields.Assignee == nil {
		return ""
	}
	return im.Issue.Fields.Assignee.DisplayName
}

func (im *IssueMore) CreateTime() time.Time {
	if im.Issue == nil || im.Issue.Fields == nil {
		return time.Time{}
	}
	return time.Time(im.Issue.Fields.Created)
}

func (im *IssueMore) CreatorName() string {
	if im.Issue == nil || im.Issue.Fields == nil || im.Issue.Fields.Creator == nil {
		return ""
	}
	return im.Issue.Fields.Creator.DisplayName
}

func (im *IssueMore) EpicName() string {
	if im.Issue == nil {
		return ""
	}
	ifs := IssueFieldsSimple{Fields: im.Issue.Fields}
	return ifs.EpicName()
}

func (im *IssueMore) Key() string {
	if im.Issue == nil {
		return ""
	}
	return strings.TrimSpace(im.Issue.Key)
}

func (im *IssueMore) KeyURL(baseURL string) string {
	key := im.Key()
	if key == "" {
		return ""
	}
	if strings.TrimSpace(baseURL) == "" {
		return ""
	}
	return BuildJiraIssueURL(baseURL, key)
}

func (im *IssueMore) ParentKey() string {
	if im.Issue == nil || im.Issue.Fields == nil || im.Issue.Fields.Parent == nil {
		return ""
	}
	return strings.TrimSpace(im.Issue.Fields.Parent.Key)
}

func (im *IssueMore) Project() string {
	if im.Issue == nil || im.Issue.Fields == nil {
		return ""
	}
	return im.Issue.Fields.Project.Name
}

func (im *IssueMore) ProjectKey() string {
	if im.Issue == nil || im.Issue.Fields == nil {
		return ""
	}
	return im.Issue.Fields.Project.Key
}

func (im *IssueMore) Resolution() string {
	if im.Issue == nil {
		return ""
	}
	ifs := IssueFieldsSimple{Fields: im.Issue.Fields}
	return ifs.ResolutionName()
}

func (im *IssueMore) Status() string {
	if im.Issue == nil {
		return ""
	}
	ifs := IssueFieldsSimple{Fields: im.Issue.Fields}
	return ifs.StatusName()
}

func (im *IssueMore) Summary() string {
	if im.Issue == nil {
		return ""
	}
	return im.Issue.Fields.Summary
}

func (im *IssueMore) Type() string {
	if im.Issue == nil {
		return ""
	}
	return im.Issue.Fields.Type.Name
}

func (im *IssueMore) UpdateTime() time.Time {
	if im.Issue == nil || im.Issue.Fields == nil {
		return time.Time{}
	}
	return time.Time(im.Issue.Fields.Updated)
}

func (im *IssueMore) Meta(baseURL string) IssueMeta {
	created := im.CreateTime().UTC()
	var createdPtr *time.Time
	if !created.IsZero() {
		createdPtr = &created
	}
	updated := im.UpdateTime().UTC()
	var updatedPtr *time.Time
	if !updated.IsZero() {
		updatedPtr = &updated
	}

	return IssueMeta{
		AssigneeName: im.AsigneeName(),
		CreateTime:   createdPtr,
		CreatorName:  im.CreatorName(),
		EpicName:     im.EpicName(),
		Key:          im.Key(),
		KeyURL:       im.KeyURL(baseURL),
		ParentKey:    im.ParentKey(),
		Project:      im.Project(),
		ProjectKey:   im.ProjectKey(),
		Resolution:   im.Resolution(),
		Status:       im.Status(),
		Summary:      im.Summary(),
		Type:         im.Type(),
		UpdateTime:   updatedPtr,
	}
}

type IssueMetas []IssueMeta

type IssueMeta struct {
	AssigneeName string
	CreateTime   *time.Time
	CreatorName  string
	EpicName     string
	Key          string
	KeyURL       string
	ParentKey    string
	Project      string
	ProjectKey   string
	Resolution   string
	Status       string
	Summary      string
	Type         string
	UpdateTime   *time.Time
}

func (im IssueMeta) String() string {
	k := strings.TrimSpace(im.Key)
	s := strings.TrimSpace(im.Summary)
	if k == "" && s == "" {
		return ""
	}
	parts := []string{}
	if len(k) > 0 {
		parts = append(parts, k)
	}
	if len(s) > 0 {
		parts = append(parts, s)
	}
	return strings.Join(parts, ": ")
}
