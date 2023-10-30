package jirarest

import (
	"errors"
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

// CustomField takes a custom value key such as `customfield_12345`.`
func (im *IssueMore) CustomField(customFieldLabel string) (IssueCustomField, error) {
	cf := IssueCustomField{}
	if im.Issue == nil {
		return cf, errors.New("issue not set")
	}
	err := GetUnmarshalCustomValue(*im.Issue, customFieldLabel, &cf)
	return cf, err
}

func (im *IssueMore) EpicKey() string {
	if im.Issue == nil || im.Issue.Fields == nil || im.Issue.Fields.Epic == nil {
		return ""
	} else {
		return im.Issue.Fields.Epic.Key
	}
}

func (im *IssueMore) EpicName() string {
	if im.Issue == nil || im.Issue.Fields == nil || im.Issue.Fields.Epic == nil {
		return ""
	} else if strings.TrimSpace(im.Issue.Fields.Epic.Name) != "" {
		return im.Issue.Fields.Epic.Name
	} else {
		return ""
	}
}

func (im *IssueMore) EpicNameOrSummary() string {
	if im.Issue == nil || im.Issue.Fields == nil || im.Issue.Fields.Epic == nil {
		return ""
	} else if strings.TrimSpace(im.Issue.Fields.Epic.Name) != "" {
		return im.Issue.Fields.Epic.Name
	} else if strings.TrimSpace(im.Issue.Fields.Epic.Summary) != "" {
		return im.Issue.Fields.Epic.Summary
	} else {
		return ""
	}
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
	if im.Issue == nil || im.Issue.Fields == nil || im.Issue.Fields.Resolution == nil {
		return ""
	}
	return im.Issue.Fields.Resolution.Name
}

func (im *IssueMore) Status() string {
	if im.Issue == nil || im.Issue.Fields == nil || im.Issue.Fields.Status == nil {
		return ""
	}
	return im.Issue.Fields.Status.Name
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
