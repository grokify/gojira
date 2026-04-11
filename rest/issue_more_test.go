package rest

import (
	"testing"
	"time"

	jira "github.com/andygrunwald/go-jira"
)

func TestIssueMoreKey(t *testing.T) {
	tests := []struct {
		name  string
		issue *jira.Issue
		want  string
	}{
		{
			name:  "nil issue",
			issue: nil,
			want:  "",
		},
		{
			name:  "valid key",
			issue: &jira.Issue{Key: "TEST-123"},
			want:  "TEST-123",
		},
		{
			name:  "key with whitespace",
			issue: &jira.Issue{Key: "  TEST-456  "},
			want:  "TEST-456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIssueMore(tt.issue)
			if got := im.Key(); got != tt.want {
				t.Errorf("IssueMore.Key() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIssueMoreType(t *testing.T) {
	tests := []struct {
		name  string
		issue *jira.Issue
		want  string
	}{
		{
			name:  "nil issue",
			issue: nil,
			want:  "",
		},
		{
			name: "bug type",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Type: jira.IssueType{Name: "Bug"},
				},
			},
			want: "Bug",
		},
		{
			name: "story type",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Type: jira.IssueType{Name: "Story"},
				},
			},
			want: "Story",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIssueMore(tt.issue)
			if got := im.Type(); got != tt.want {
				t.Errorf("IssueMore.Type() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIssueMoreStatus(t *testing.T) {
	tests := []struct {
		name  string
		issue *jira.Issue
		want  string
	}{
		{
			name:  "nil issue",
			issue: nil,
			want:  "",
		},
		{
			name: "nil fields",
			issue: &jira.Issue{
				Fields: nil,
			},
			want: "",
		},
		{
			name: "nil status",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Status: nil,
				},
			},
			want: "",
		},
		{
			name: "open status",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Status: &jira.Status{Name: "Open"},
				},
			},
			want: "Open",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIssueMore(tt.issue)
			if got := im.Status(); got != tt.want {
				t.Errorf("IssueMore.Status() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIssueMoreSummary(t *testing.T) {
	tests := []struct {
		name  string
		issue *jira.Issue
		want  string
	}{
		{
			name:  "nil issue",
			issue: nil,
			want:  "",
		},
		{
			name: "valid summary",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Summary: "Fix login bug",
				},
			},
			want: "Fix login bug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIssueMore(tt.issue)
			if got := im.Summary(); got != tt.want {
				t.Errorf("IssueMore.Summary() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIssueMoreAssigneeName(t *testing.T) {
	tests := []struct {
		name  string
		issue *jira.Issue
		want  string
	}{
		{
			name:  "nil issue",
			issue: nil,
			want:  "",
		},
		{
			name: "nil assignee",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Assignee: nil,
				},
			},
			want: "",
		},
		{
			name: "valid assignee",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Assignee: &jira.User{DisplayName: "John Doe"},
				},
			},
			want: "John Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIssueMore(tt.issue)
			if got := im.AssigneeName(); got != tt.want {
				t.Errorf("IssueMore.AssigneeName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIssueMoreLabels(t *testing.T) {
	tests := []struct {
		name    string
		issue   *jira.Issue
		sortAsc bool
		want    []string
	}{
		{
			name:    "nil issue",
			issue:   nil,
			sortAsc: false,
			want:    []string{},
		},
		{
			name: "no labels",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Labels: []string{},
				},
			},
			sortAsc: false,
			want:    []string{},
		},
		{
			name: "multiple labels unsorted",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Labels: []string{"zebra", "apple", "banana"},
				},
			},
			sortAsc: false,
			want:    []string{"zebra", "apple", "banana"},
		},
		{
			name: "multiple labels sorted",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Labels: []string{"zebra", "apple", "banana"},
				},
			},
			sortAsc: true,
			want:    []string{"apple", "banana", "zebra"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIssueMore(tt.issue)
			got := im.Labels(tt.sortAsc)
			if len(got) != len(tt.want) {
				t.Errorf("IssueMore.Labels() len = %d, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("IssueMore.Labels()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestIssueMoreLabelExists(t *testing.T) {
	issue := &jira.Issue{
		Fields: &jira.IssueFields{
			Labels: []string{"bug", "High-Priority", "needs-review"},
		},
	}
	im := NewIssueMore(issue)

	tests := []struct {
		name       string
		label      string
		ignoreCase bool
		want       bool
	}{
		{
			name:       "exact match",
			label:      "bug",
			ignoreCase: false,
			want:       true,
		},
		{
			name:       "case mismatch strict",
			label:      "BUG",
			ignoreCase: false,
			want:       false,
		},
		{
			name:       "case mismatch ignore case",
			label:      "BUG",
			ignoreCase: true,
			want:       true,
		},
		{
			name:       "not found",
			label:      "nonexistent",
			ignoreCase: false,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := im.LabelExists(tt.label, tt.ignoreCase); got != tt.want {
				t.Errorf("IssueMore.LabelExists(%q, %v) = %v, want %v", tt.label, tt.ignoreCase, got, tt.want)
			}
		})
	}
}

func TestIssueMoreCreateTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		issue *jira.Issue
		want  time.Time
	}{
		{
			name:  "nil issue",
			issue: nil,
			want:  time.Time{},
		},
		{
			name: "nil fields",
			issue: &jira.Issue{
				Fields: nil,
			},
			want: time.Time{},
		},
		{
			name: "valid time",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Created: jira.Time(now),
				},
			},
			want: now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIssueMore(tt.issue)
			got := im.CreateTime()
			if !got.Equal(tt.want) {
				t.Errorf("IssueMore.CreateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssueMoreParentKey(t *testing.T) {
	tests := []struct {
		name  string
		issue *jira.Issue
		want  string
	}{
		{
			name:  "nil issue",
			issue: nil,
			want:  "",
		},
		{
			name: "nil parent",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Parent: nil,
				},
			},
			want: "",
		},
		{
			name: "valid parent",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Parent: &jira.Parent{Key: "PARENT-100"},
				},
			},
			want: "PARENT-100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIssueMore(tt.issue)
			if got := im.ParentKey(); got != tt.want {
				t.Errorf("IssueMore.ParentKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIssueMoreProject(t *testing.T) {
	tests := []struct {
		name  string
		issue *jira.Issue
		want  string
	}{
		{
			name:  "nil issue",
			issue: nil,
			want:  "",
		},
		{
			name: "valid project",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Project: jira.Project{Name: "My Project"},
				},
			},
			want: "My Project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIssueMore(tt.issue)
			if got := im.Project(); got != tt.want {
				t.Errorf("IssueMore.Project() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIssueMoreProjectKey(t *testing.T) {
	tests := []struct {
		name  string
		issue *jira.Issue
		want  string
	}{
		{
			name:  "nil issue",
			issue: nil,
			want:  "",
		},
		{
			name: "valid project key",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Project: jira.Project{Key: "PROJ"},
				},
			},
			want: "PROJ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIssueMore(tt.issue)
			if got := im.ProjectKey(); got != tt.want {
				t.Errorf("IssueMore.ProjectKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIssueMoreResolution(t *testing.T) {
	tests := []struct {
		name  string
		issue *jira.Issue
		want  string
	}{
		{
			name:  "nil issue",
			issue: nil,
			want:  "",
		},
		{
			name: "nil resolution",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Resolution: nil,
				},
			},
			want: "",
		},
		{
			name: "valid resolution",
			issue: &jira.Issue{
				Fields: &jira.IssueFields{
					Resolution: &jira.Resolution{Name: "Fixed"},
				},
			},
			want: "Fixed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewIssueMore(tt.issue)
			if got := im.Resolution(); got != tt.want {
				t.Errorf("IssueMore.Resolution() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewIssueMore(t *testing.T) {
	issue := &jira.Issue{Key: "TEST-1"}
	im := NewIssueMore(issue)

	if im.Issue != issue {
		t.Error("NewIssueMore() did not set Issue correctly")
	}
}
