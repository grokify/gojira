package jirarest

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
)

func (is *IssuesSet) RetrieveParentsIssuesSet(client *Client) (*IssuesSet, error) {
	parIssuesSet := NewIssuesSet(is.Config)
	parIDs := is.UnknownParents()
	if len(parIDs) == 0 {
		return parIssuesSet, nil
	}

	err := parIssuesSet.RetrieveIssues(client, parIDs)
	if err != nil {
		return nil, err
	}

	err = parIssuesSet.RetrieveParents(client)

	return parIssuesSet, err
}

func (is *IssuesSet) RetrieveParents(client *Client) error {
	parIDs := is.UnknownParents()
	i := 0
	for len(parIDs) > 0 {
		err := is.RetrieveIssues(client, parIDs)
		if err != nil {
			return err
		}
		parIDs = is.UnknownParents()
		i++
		if i > 10 {
			return errors.New("exceeded max retrieve parent iterations")
		}
	}
	return nil
}

func (is *IssuesSet) RetrieveIssues(client *Client, ids []string) error {
	ids = stringsutil.SliceCondenseSpace(ids, true, true)
	if len(ids) == 0 {
		return nil
	}
	jql := "key in (" + strings.Join(ids, ",") + ")"
	iss, err := client.SearchIssues(jql)
	if err != nil {
		return err
	}
	return is.Add(iss...)
}

func (is *IssuesSet) UnknownParents() []string {
	parKeys := []string{}
	for _, iss := range is.IssuesMap {
		im := IssueMore{Issue: &iss}
		parKey := im.ParentKey()
		if parKey == "" {
			continue
		}
		if _, ok := is.IssuesMap[parKey]; ok {
			continue
		}
		parKeys = append(parKeys, parKey)
	}
	return stringsutil.SliceCondenseSpace(parKeys, true, true)
}

func (is *IssuesSet) Lineage(key string) (IssueMetas, error) {
	ims := IssueMetas{}
	iss, ok := is.IssuesMap[key]
	if !ok {
		return ims, fmt.Errorf("key not found (%s)", key)
	}
	im := IssueMore{Issue: &iss}
	imeta := im.Meta(is.Config.BaseURL)
	ims = append(ims, imeta)

	parKey := im.ParentKey()
	if parKey == "" {
		return ims, nil
	}

	if is.Parents == nil {
		return ims, errors.New("parents not set")
	}

	for parKey != "" {
		parIss, ok := is.Parents.IssuesMap[parKey]
		if !ok {
			return ims, errors.New("parent not found")
		}
		parIM := IssueMore{Issue: &parIss}
		parM := parIM.Meta(is.Config.BaseURL)
		ims = append(ims, parM)
		parKey = parIM.ParentKey()
	}

	return ims, nil
}
