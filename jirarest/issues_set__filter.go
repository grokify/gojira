package jirarest

import (
	"github.com/grokify/mogo/pointer"
	"golang.org/x/exp/slices"
)

func (set *IssuesSet) FilterByStatus(inclStatuses, exclStatuses []string) (*IssuesSet, error) {
	filteredIssuesSet := NewIssuesSet(set.Config)
	inclStatusesMap := map[string]int{}
	for _, s := range inclStatuses {
		inclStatusesMap[s]++
	}
	exclStatusesMap := map[string]int{}
	for _, s := range exclStatuses {
		exclStatusesMap[s]++
	}
	for _, iss := range set.IssuesMap {
		im := NewIssueMore(pointer.Pointer(iss))
		// ifs := IssueFieldsSimple{Fields: iss.Fields}
		statusName := im.Status()
		_, inclStatusOk := inclStatusesMap[statusName]
		_, exclStatusOk := exclStatusesMap[statusName]
		if len(inclStatusesMap) > 0 && !inclStatusOk {
			continue
		} else if len(exclStatuses) > 0 && exclStatusOk {
			continue
		}
		err := filteredIssuesSet.Add(iss)
		if err != nil {
			return nil, err
		}
	}
	return filteredIssuesSet, nil
}

/*
func (set *IssuesSet) FilterStatus(inclStatuses ...string) (*IssuesSet, error) {
	n := NewIssuesSet(set.Config)
	if len(inclStatuses) == 0 {
		return n, nil
	}
	for _, iss := range set.IssuesMap {
		im := NewIssueMore(pointer.Pointer(iss))
		if slices.Index(inclStatuses, im.Status()) >= 0 {
			err := n.Add(iss)
			if err != nil {
				return nil, err
			}
		}
	}
	return n, nil
}
*/

func (set *IssuesSet) FilterByType(inclTypes ...string) (*IssuesSet, error) {
	n := NewIssuesSet(set.Config)
	if len(inclTypes) == 0 {
		return n, nil
	}
	for _, iss := range set.IssuesMap {
		im := NewIssueMore(pointer.Pointer(iss))
		if slices.Index(inclTypes, im.Type()) >= 0 {
			err := n.Add(iss)
			if err != nil {
				return nil, err
			}
		}
	}
	return n, nil
}
