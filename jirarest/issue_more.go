package jirarest

import (
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/time/timeutil"
	"golang.org/x/exp/slices"
)

type IssueMore struct {
	issue *jira.Issue
}

func NewIssueMore(iss *jira.Issue) IssueMore {
	return IssueMore{issue: iss}
}

func (im *IssueMore) AdditionalFields(additionalFieldNames []string) map[string]*string {
	out := map[string]*string{}
	for _, name := range additionalFieldNames {
		cfStr, err := im.CustomFieldString(name)
		if err != nil {
			out[name] = nil
		} else {
			out[name] = pointer.Pointer(cfStr)
		}
	}
	return out
}

func (im *IssueMore) AssigneeName() string {
	if im.issue == nil || im.issue.Fields == nil || im.issue.Fields.Assignee == nil {
		return ""
	}
	return im.issue.Fields.Assignee.DisplayName
}

func (im *IssueMore) CreateTime() time.Time {
	if im.issue == nil || im.issue.Fields == nil {
		return time.Time{}
	}
	return time.Time(im.issue.Fields.Created)
}

func (im *IssueMore) CreatorName() string {
	if im.issue == nil || im.issue.Fields == nil || im.issue.Fields.Creator == nil {
		return ""
	}
	return im.issue.Fields.Creator.DisplayName
}

// CustomField takes a custom value key such as `customfield_12345`.`
func (im *IssueMore) CustomField(customFieldLabel string) (IssueCustomField, error) {
	cf := IssueCustomField{}
	if im.issue == nil {
		return cf, errors.New("issue not set")
	}
	err := GetUnmarshalCustomValue(*im.issue, customFieldLabel, &cf)
	return cf, err
}

// CustomFieldString takes a custom value key such as `customfield_12345`.`
func (im *IssueMore) CustomFieldString(customFieldLabel string) (string, error) {
	cf, err := im.CustomField(customFieldLabel)
	return cf.Value, err
}

// CustomFieldStringOrEmpty takes a custom value key such as `customfield_12345`.`
func (im *IssueMore) CustomFieldStringOrDefault(customFieldLabel, def string) string {
	if cf, err := im.CustomField(customFieldLabel); err != nil {
		return def
	} else {
		return cf.Value
	}
}

func (im *IssueMore) EpicKey() string {
	if im.issue == nil || im.issue.Fields == nil || im.issue.Fields.Epic == nil {
		return ""
	} else {
		return im.issue.Fields.Epic.Key
	}
}

func (im *IssueMore) EpicName() string {
	if im.issue == nil || im.issue.Fields == nil || im.issue.Fields.Epic == nil {
		return ""
	} else if strings.TrimSpace(im.issue.Fields.Epic.Name) != "" {
		return im.issue.Fields.Epic.Name
	} else {
		return ""
	}
}

func (im *IssueMore) EpicNameOrSummary() string {
	if im.issue == nil || im.issue.Fields == nil || im.issue.Fields.Epic == nil {
		return ""
	} else if strings.TrimSpace(im.issue.Fields.Epic.Name) != "" {
		return im.issue.Fields.Epic.Name
	} else if strings.TrimSpace(im.issue.Fields.Epic.Summary) != "" {
		return im.issue.Fields.Epic.Summary
	} else {
		return ""
	}
}

func (im *IssueMore) Key() string {
	if im.issue == nil {
		return ""
	}
	return strings.TrimSpace(im.issue.Key)
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

func (im *IssueMore) Labels(sortAsc bool) []string {
	if im.issue == nil || im.issue.Fields == nil || len(im.issue.Fields.Labels) == 0 {
		return []string{}
	} else if !sortAsc || len(im.issue.Fields.Labels) == 1 {
		return im.issue.Fields.Labels
	} else {
		labels := im.issue.Fields.Labels
		sort.Strings(labels)
		return labels
	}
}

func (im *IssueMore) LabelExists(label string) bool {
	labels := im.Labels(true)
	return slices.Contains(labels, label)
}

func (im *IssueMore) ParentKey() string {
	if im.issue == nil || im.issue.Fields == nil || im.issue.Fields.Parent == nil {
		return ""
	}
	return strings.TrimSpace(im.issue.Fields.Parent.Key)
}

func (im *IssueMore) Project() string {
	if im.issue == nil || im.issue.Fields == nil {
		return ""
	}
	return im.issue.Fields.Project.Name
}

func (im *IssueMore) ProjectKey() string {
	if im.issue == nil || im.issue.Fields == nil {
		return ""
	}
	return im.issue.Fields.Project.Key
}

func (im *IssueMore) Resolution() string {
	if im.issue == nil || im.issue.Fields == nil || im.issue.Fields.Resolution == nil {
		return ""
	}
	return im.issue.Fields.Resolution.Name
}

func (im *IssueMore) Status() string {
	if im.issue == nil || im.issue.Fields == nil || im.issue.Fields.Status == nil {
		return ""
	}
	return im.issue.Fields.Status.Name
}

func (im *IssueMore) Summary() string {
	if im.issue == nil {
		return ""
	}
	return im.issue.Fields.Summary
}

func (im *IssueMore) Type() string {
	if im.issue == nil {
		return ""
	}
	return im.issue.Fields.Type.Name
}

func (im *IssueMore) UpdateTime() time.Time {
	if im.issue == nil || im.issue.Fields == nil {
		return time.Time{}
	}
	return time.Time(im.issue.Fields.Updated)
}

func (im *IssueMore) Value(fieldSlug string) (string, bool) {
	fieldSlug = strings.ToLower(strings.TrimSpace(fieldSlug))
	switch fieldSlug {
	case gojira.FieldKey:
		return im.Key(), true
	case gojira.AliasIssueKey:
		return im.Key(), true
	case gojira.FieldProjectKey:
		return im.ProjectKey(), true
	case gojira.CalcCreatedAgeDays:
		t := im.CreateTime()
		tm := timeutil.NewTimeMore(t, 0)
		if tm.IsZeroAny() {
			return "0", true
		} else {
			days := timeutil.DurationDays(time.Since(t))
			return strconv.Itoa(int(days)), true
		}
	case gojira.FieldCreatedDate:
		t := im.CreateTime()
		return t.Format(time.RFC3339), true
	case gojira.CalcCreatedMonth:
		tm := timeutil.NewTimeMore(im.CreateTime().UTC(), 0)
		return tm.MonthStart().Format(time.RFC3339), true
	case gojira.FieldResolution:
		return im.Resolution(), true
	case gojira.FieldStatus:
		return im.Status(), true
	case gojira.FieldSummary:
		return im.Summary(), true
	case gojira.FieldType:
		return im.Type(), true
	default:
		if canonicalCustomKey, ok := gojira.IsCustomFieldKey(fieldSlug); ok {
			return im.CustomFieldStringOrDefault(canonicalCustomKey, ""), true
		}
	}
	return "", false
}

func (im *IssueMore) Meta(serverURL string, additionalFieldNames []string) IssueMeta {
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
		AdditionalFields: im.AdditionalFields(additionalFieldNames),
		AssigneeName:     im.AssigneeName(),
		CreateTime:       createdPtr,
		CreatorName:      im.CreatorName(),
		EpicName:         im.EpicName(),
		Key:              im.Key(),
		KeyURL:           im.KeyURL(serverURL),
		Labels:           im.Labels(true),
		ParentKey:        im.ParentKey(),
		Project:          im.Project(),
		ProjectKey:       im.ProjectKey(),
		Resolution:       im.Resolution(),
		Status:           im.Status(),
		Summary:          im.Summary(),
		Type:             im.Type(),
		UpdateTime:       updatedPtr,
	}
}

func (im *IssueMore) WriteFileJSON(filename string, perm os.FileMode, prefix, indent string) error {
	return jsonutil.WriteFile(filename, im.issue, prefix, indent, perm)
}
