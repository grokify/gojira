package jirarest

import (
	"errors"
	"fmt"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/fmt/fmtutil"
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
	parIDs := is.KeysParentsUnpopulated()
	i := 0
	for len(parIDs) > 0 {
		err := is.RetrieveIssues(client, parIDs)
		if err != nil {
			return err
		}
		parIDs = is.KeysParentsUnpopulated()
		i++
		if i > 10 {
			return errors.New("exceeded max retrieve parent iterations")
		}
	}
	return nil
}

func (is *IssuesSet) RetrieveIssues(client *Client, ids []string) error {
	if client == nil {
		return errors.New("client cannot be nil")
	}
	ids = stringsutil.SliceCondenseSpace(ids, true, true)
	if len(ids) == 0 {
		return nil
	}
	fmt.Printf("LEN(%d)\n", len(ids))

	idsSlicesMaxResults := SplitMaxLength(ids, gojira.JQLMaxResults)
	fmtutil.PrintJSON(idsSlicesMaxResults)

	for _, idsSlice := range idsSlicesMaxResults {
		jqls := gojira.JQLStringsSimple(gojira.FieldKey, false, idsSlice, 0)
		//fmtutil.PrintJSON(jqls)
		//panic("Z")

		for _, jql := range jqls {
			if iss, err := client.SearchIssuesPages(jql, 0, 0, 0); err != nil {
				return err
			} else {
				return is.Add(iss...)
			}
			/*
				// jql := "key in (" + strings.Join(ids, ",") + ")"
				if iss, err := client.SearchIssues(jql); err != nil {
					return err
				} else {
					return is.Add(iss...)
				}
			*/
		}
	}
	return nil
}

func (is *IssuesSet) IssueOrParent(key string) (*jira.Issue, bool) {
	if iss, ok := is.IssuesMap[key]; ok {
		return &iss, true
	} else if is.Parents == nil {
		return nil, false
	} else if iss, ok := is.Parents.IssuesMap[key]; ok {
		return &iss, true
	} else {
		return nil, false
	}
}

func (is *IssuesSet) KeysParents() []string {
	var parKeys []string
	for _, iss := range is.IssuesMap {
		im := IssueMore{Issue: pointer.Pointer(iss)}
		if parKey := im.ParentKey(); parKey != "" {
			parKeys = append(parKeys, parKey)
		}
	}
	return stringsutil.SliceCondenseSpace(parKeys, true, true)
}

// ParentsPopulated returns issue ids that are in the current set or current parent set.
func (is *IssuesSet) KeysParentsPopulated() []string {
	var parKeysPop []string
	parKeysAll := is.KeysParents()
	for _, parKey := range parKeysAll {
		parIss, ok := is.IssueOrParent(parKey)
		if ok && parIss != nil {
			parKeysPop = append(parKeysPop, parKey)
		}
	}

	return stringsutil.SliceCondenseSpace(parKeysPop, true, true)
}

// ParentsUnpopulated returns issue ids that are not in the current set or current parent set.
func (is *IssuesSet) KeysParentsUnpopulated() []string {
	var parKeysUnpop []string
	parKeysAll := is.KeysParents()
	for _, parKey := range parKeysAll {
		parIss, ok := is.IssueOrParent(parKey)
		if !ok || parIss == nil {
			parKeysUnpop = append(parKeysUnpop, parKey)
		}
	}

	return stringsutil.SliceCondenseSpace(parKeysUnpop, true, true)
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
