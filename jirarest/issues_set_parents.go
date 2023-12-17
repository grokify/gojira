package jirarest

import (
	"errors"
	"strings"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/type/stringsutil"
)

/*
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
*/

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
		im := IssueMore{Issue: pointer.Pointer(iss)}
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

// Lineage returns a slice of `IssueMeta` where the suppied key is in index 0 and the most senior
// parent is the last element of the slice. If a parent is not found in the set, an error is returned.
func (is *IssuesSet) Lineage(key string) (IssueMetas, error) {
	if key == "Epic" {
		panic("Lineage Epic")
	}
	ims := IssueMetas{}
	iss, err := is.Get(key)
	if err != nil {
		return ims, errorsutil.Wrapf(err, "key not found (%s)", key)
	}
	im := IssueMore{Issue: &iss}
	imeta := im.Meta(is.Config.ServerURL)
	ims = append(ims, imeta)
	parKey := im.ParentKey()

	if parKey != "" && is.Parents == nil {
		return ims, errors.New("parents not set")
	}

	for parKey != "" {
		parIss, err := is.Get(parKey)
		if err != nil {
			return ims, errorsutil.Wrap(err, "parent not found")
		}
		parIM := IssueMore{Issue: &parIss}
		parM := parIM.Meta(is.Config.ServerURL)
		ims = append(ims, parM)
		parKey = parIM.ParentKey()
	}

	return ims, nil
}

/*
func (is *IssuesSet) GetLineage(key string) (Issues, error) {
	key = strings.TrimSpace(key)
	lineage := Issues{}
	if key == "" {
		return lineage, nil
	}
	iss, err := is.Get(key)
	if err != nil {
		return lineage, err
	}

	im := IssueMore{Issue: &iss}
	parKey := im.ParentKey()
	for {
		if parKey == "" {
			return lineage, nil
		}
		if parIss, err := is.Get(parKey); err != nil {
			return lineage, err
		} else {
			lineage = append(lineage, parIss)
			parIM := IssueMore{Issue: &parIss}
			parKey = parIM.ParentKey()
		}
	}
	return lineage, nil
}
*/
