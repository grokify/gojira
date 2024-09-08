package jirarest

import "errors"

func (svc *IssueService) IssuesSetAddParents(is *IssuesSet) error {
	if is == nil {
		return errors.New("issues set is nil")
	} else if parents, err := svc.SearchIssuesSetParents(is); err != nil {
		return err
	} else {
		is.Parents = parents
		return nil
	}
}
