package jirarest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/http/httpsimple"

	"github.com/grokify/gojira"
)

type TransitionPayload struct {
	Transition jira.TransitionPayload  `json:"transition,omitempty"`
	Fields     map[string]any          `json:"fields,omitempty"`
	Update     TransitionPayloadUpdate `json:"update,omitempty"`
}

type TransitionPayloadUpdate struct {
	Worklog []WorklogOperation `json:"worklog,omitempty"`
}

type WorklogOperation struct {
	Add jira.WorklogRecord `json:"add"`
}

type TranstionFieldIDOrValue struct {
	ID    any `json:"id,omitempty"`
	Value any `json:"value,omitempty"`
}

type TransitionsAPIReponse struct {
	Transitions []Transition `json:"transitions"`
}

func (svc *IssueService) GetTransitionsWithContext(ctx context.Context, id string, expandTransitionsFields bool) (Transitions, *jira.Response, error) {
	if expandTransitionsFields {
		if !gojira.KeyIsValid(id) {
			return nil, nil, fmt.Errorf("id is not a valid Jira key (%s)", id)
		}
		sr := httpsimple.Request{
			Method: http.MethodGet,
			URL:    fmt.Sprintf("/rest/api/2/issue/%s/transitions", id),
			Query:  map[string][]string{"expand": {"transitions.fields"}}}
		var result TransitionsAPIReponse
		if resp, err := svc.Client.simpleClient.Do(ctx, sr); err != nil {
			return nil, nil, err
		} else if resp.StatusCode > 299 {
			return nil, nil, fmt.Errorf("bad api response status code (%d)", resp.StatusCode)
		} else if b, err := io.ReadAll(resp.Body); err != nil {
			return nil, nil, err
		} else {
			err = json.Unmarshal(b, &result)
			return result.Transitions, &jira.Response{Response: resp}, err
		}
	} else {
		txnsSDK, resp, err := svc.Client.JiraClient.Issue.GetTransitionsWithContext(ctx, id)
		if err != nil {
			return nil, resp, err
		}
		txns := Transitions{}
		txns.AddTransitionsSDK(txnsSDK)
		return txns, resp, nil
	}
}

func (svc *IssueService) DoTransitionWithNameAndPayload(ctx context.Context, issueID, updateTransitionName string, payload *TransitionPayload) error {
	issue, resp, err := svc.Client.JiraClient.Issue.Get(issueID, nil)
	if err != nil {
		return err
	} else if resp.StatusCode > 299 {
		return fmt.Errorf("bad api response status code (%d)", resp.StatusCode)
	}
	if issue.Fields.Status.Name == updateTransitionName {
		return nil
	}
	possibleTxns, resp, err := svc.GetTransitionsWithContext(ctx, issueID, true)
	if err != nil {
		return err
	} else if resp.StatusCode > 299 {
		return fmt.Errorf("bad api response status code (%d)", resp.StatusCode)
	}
	txns := Transitions(possibleTxns)
	wantTxn, err := txns.GetByName(updateTransitionName)
	if err != nil {
		return errorsutil.Wrapf(err, "jira id (%s)", issueID)
	}
	if payload != nil {
		payload.Transition.ID = wantTxn.ID
		if resp, err = svc.Client.JiraClient.Issue.DoTransitionWithPayloadWithContext(ctx, issueID, *payload); err != nil {
			return err
		} else if resp.StatusCode > 299 {
			return fmt.Errorf("bad api response status code (%d)", resp.StatusCode)
		}
	} else {
		if resp, err = svc.Client.JiraClient.Issue.DoTransitionWithContext(ctx, issueID, wantTxn.ID); err != nil {
			return err
		} else if resp.StatusCode > 299 {
			return fmt.Errorf("bad api response status code (%d)", resp.StatusCode)
		}
	}

	return nil
}
