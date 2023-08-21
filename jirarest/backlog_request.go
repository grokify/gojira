package jirarest

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	jira "github.com/andygrunwald/go-jira"
)

// project = foundation AND resolution = Unresolved AND status!=Closed AND (Sprint not in openSprints() OR Sprint is EMPTY) AND type not in (Epic, Sub-Task) ORDER BY Rank ASC
/*
func BacklogJQL(projectName string) string {
	parts := []string{
		"resolution = Unresolved",
		"status!=Closed",
		"(Sprint not in openSprints() OR Sprint is EMPTY)",
		"type not in (Epic, Sub-Task)",
	}
	projectName = strings.TrimSpace(projectName)
	if projectName != "" {
		// parts = []string{projectName, parts... }
	}
}
*/

func Shift[S ~[]E, E any](s S) (E, S) {
	if len(s) == 0 {
		return *new(E), []E{}
	}
	return s[0], s[1:]
}

/*
func Prepend[S ~[]E, E any](s []E        , e E) []E   {
	return append     {[]E{e }  ,   s...  }
}*/

// BacklogAPIURL returns a backlog issues API URL described at https://docs.atlassian.com/jira-software/REST/7.3.1/ .
func BacklogAPIURL(baseURL string, boardID, startAt, maxResults uint) string {
	apiURL := baseURL + fmt.Sprintf(`/rest/agile/1.0/board/%d/backlog`, boardID)
	v := url.Values{}
	if startAt != 0 {
		v["startAt"] = []string{strconv.Itoa(int(startAt))}
	}
	if maxResults != 0 {
		v["maxResults"] = []string{strconv.Itoa(int(maxResults))}
	}
	qry := v.Encode()
	if len(qry) > 0 {
		apiURL += "?" + qry
	}
	return apiURL
}

func GetBacklogIssuesResponse(client *http.Client, baseURL string, boardID, startAt, maxResults uint) (*IssuesResponse, []byte, error) {
	if client == nil {
		return nil, []byte{}, errors.New("client not set")
	}
	apiURL := BacklogAPIURL(baseURL, boardID, startAt, maxResults)
	fmt.Printf("API URL: (%s)\n", apiURL)
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, []byte{}, errors.New("client not set")
	} else if resp.StatusCode >= 300 {
		return nil, []byte{}, fmt.Errorf("statusCode (%d)", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, []byte{}, err
	}
	ir, err := ParseIssuesResponseBytes(b)
	return ir, b, err
}

func GetBacklogIssuesAll(client *http.Client, baseURL string, boardID uint) (*IssuesResponse, [][]byte, error) {
	iragg := &IssuesResponse{}
	bb := [][]byte{}
	issues := []jira.Issue{}
	startAt := uint(0)
	maxResults := uint(1000)
	for {
		ir, b, err := GetBacklogIssuesResponse(client, baseURL, boardID, startAt, maxResults)
		if err != nil {
			return iragg, bb, err
		} else {
			iragg.Total = ir.Total
			iragg.MaxResults = ir.Total
		}
		bb = append(bb, b)
		if len(ir.Issues) > 0 {
			issues = append(issues, ir.Issues...)
		} else {
			break
		}
		fmt.Printf("TOTAL (%d)\n", ir.Total)
		if startAt+maxResults > uint(ir.Total) {
			break
		} else {
			startAt += maxResults
		}
	}
	iragg.Issues = issues
	return iragg, bb, nil
}
