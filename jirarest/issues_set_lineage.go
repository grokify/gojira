package jirarest

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/type/stringsutil"
)

var ErrLineageNotFound = errors.New("lineage not found")

func (is *IssuesSet) LineageValidateSet() (popLineage []string, unpopLineage []string, allValid bool) {
	issKeys := is.Keys()
	for _, issKey := range issKeys {
		_, err := is.LineageValidateKey(issKey)
		if err != nil {
			unpopLineage = append(unpopLineage, issKey)
		} else {
			popLineage = append(popLineage, issKey)
		}
	}
	popLineage = stringsutil.SliceCondenseSpace(popLineage, true, true)
	unpopLineage = stringsutil.SliceCondenseSpace(unpopLineage, true, true)
	if len(popLineage) == len(is.IssuesMap) && len(unpopLineage) == 0 {
		allValid = true
		return popLineage, unpopLineage, allValid
	}
	return popLineage, unpopLineage, allValid
}

func (is *IssuesSet) LineageTopKeysPopulated() ([]string, error) {
	var linPop []string
	issKeys := is.Keys()
	for _, issKey := range issKeys {
		issKey = strings.TrimSpace(issKey)
		if issKey == "" {
			return linPop, errors.New("issue map key is empty string")
		}
		lin, err := is.LineageValidateKey(issKey)
		if err != nil {
			if errors.Is(err, ErrLineageNotFound) {
				continue
			} else {
				panic(err)
				// return linUnpop, err
			}
		} else {
			if len(lin) > 0 {
				linPop = append(linPop, lin[len(lin)-1])
			} else {
				panic("lineage empty - LineageTopKeysPopulated")
			}
		}
	}
	return stringsutil.SliceCondenseSpace(linPop, true, true), nil
}

func (is *IssuesSet) LineageTopKeysUnpopulated() ([]string, error) {
	var linUnpop []string
	issKeys := is.Keys()
	for _, issKey := range issKeys {
		issKey = strings.TrimSpace(issKey)
		if issKey == "" {
			return linUnpop, errors.New("issue map key is empty string")
		}
		lin, err := is.LineageValidateKey(issKey)
		if err != nil {
			if errors.Is(err, ErrLineageNotFound) {
				if len(lin) > 0 {
					linUnpop = append(linUnpop, lin[len(lin)-1])
				} else {
					panic("linage empty - LineageTopKeysUnpopulated")
				}
			} else {
				panic(err)
				// return linUnpop, err
			}
		}
	}
	return stringsutil.SliceCondenseSpace(linUnpop, true, true), nil
}

// LineageValidateKey returns a lineage slice where the leaf key is in index position 0 (little-endian).
// This is done in case a parent cannot be found in which case the boolean returned is false.
func (is *IssuesSet) LineageValidateKey(key string) ([]string, error) {
	key = strings.TrimSpace(key)
	var lineage []string
	if key == "" {
		return lineage, errors.New("key not provided")
	}
	iss, ok := is.IssueOrParent(key)
	if !ok {
		return lineage, fmt.Errorf("key not found for (%s)", key)
	}
	im := IssueMore{Issue: iss}
	issKey := im.Key()
	if issKey != key {
		return lineage, fmt.Errorf("found key (%s) did not match request (%s)", issKey, key)
	} else {
		lineage = append(lineage, issKey)
	}
	parKey := im.ParentKey()
	for {
		if parKey == "" {
			return lineage, nil
		}
		lineage = append(lineage, parKey)
		parIss, ok := is.IssueOrParent(parKey)
		if !ok {
			return lineage, errorsutil.Wrapf(ErrLineageNotFound, "parent key not found (%s)", parKey)
		}
		parIssMore := IssueMore{Issue: parIss}
		parIssKey := parIssMore.Key()
		if parIssKey != parKey {
			return lineage, fmt.Errorf("found key (%s) did not match request (%s)", parIssKey, parKey)
		}
		parKey = parIssMore.ParentKey()
	}
	return lineage, nil
}
