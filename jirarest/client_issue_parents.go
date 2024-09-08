package jirarest

import "errors"

func (c *IssueAPI) IssuesSetAddParents(is *IssuesSet) error {
	if is == nil {
		return errors.New("issues set is nil")
	} else if parents, err := c.SearchIssuesSetParents(is); err != nil {
		return err
	} else {
		is.Parents = parents
		return nil
	}
}
