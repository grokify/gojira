package jirarest

import (
	"github.com/grokify/mogo/errors/errorsutil"
	"golang.org/x/exp/slices"
)

// RetrieveChildrenOfType retrieves all children of supplied parent types. If the child matches a base type,
// it is inserted into the current `IssuesSet`. If it is not a baseType, it is inserted into `Parents`. Of
// note, this will only load children of parent types that are already in the `IssuesSet`.
func (is *IssuesSet) RetrieveChildrenOfType(client *Client, parentTypes, baseTypes []string) error {
	if len(parentTypes) == 0 {
		return nil
	} else if client == nil {
		return errorsutil.Wrap(ErrClientCannotBeNil, "called in IssuesSet.RetrieveChildrenOfType")
	}
	parentKeys := is.KeysForTypes(parentTypes, true, true)

	if len(parentKeys) == 0 {
		return nil
	}

	children, err := client.IssueAPI.SearchChildrenIssues(parentKeys...)
	if err != nil {
		return err
	} else if is.Parents == nil {
		is.Parents = NewIssuesSet(is.Config)
	}

	for len(children) > 0 {
		var unknownChildrenKeys []string
		for _, c := range children {
			c := c
			im := NewIssueMore(&c)
			childKey := im.Key()
			if childKey == "" || is.KeyExists(childKey, true) {
				continue
			} else if slices.Index(baseTypes, im.Type()) > -1 {
				is.IssuesMap[childKey] = c
				continue
			} else { // not base type.
				is.Parents.IssuesMap[childKey] = c
				unknownChildrenKeys = append(unknownChildrenKeys, childKey)
			}
		}

		if len(unknownChildrenKeys) == 0 {
			break
		} else {
			children, err = client.IssueAPI.SearchChildrenIssues(unknownChildrenKeys...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
