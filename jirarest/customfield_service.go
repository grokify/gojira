package jirarest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/net/urlutil"
)

type CustomFieldService struct {
	JRClient *Client
}

func NewCustomFieldService(client *Client) *CustomFieldService {
	return &CustomFieldService{JRClient: client}
}

func (svc *CustomFieldService) GetCustomFields() (CustomFields, error) {
	var cfs CustomFields
	if svc.JRClient == nil {
		return cfs, ErrJiraRESTClientCannotBeNil
	}
	apiURL := urlutil.JoinAbsolute(svc.JRClient.Config.ServerURL, APIV2URLListCustomFields)
	hclient := svc.JRClient.HTTPClient
	if hclient == nil {
		hclient = &http.Client{}
	}

	resp, err := hclient.Get(apiURL)
	if err != nil {
		return cfs, err
	}
	if resp.StatusCode >= 300 {
		return cfs, fmt.Errorf("error status code (%d)", resp.StatusCode)
	}
	_, err = jsonutil.UnmarshalReader(resp.Body, &cfs)
	return cfs, err
}

func (svc *CustomFieldService) GetCustomFieldEpicLink() (CustomField, error) {
	return svc.GetCustomField(CustomFieldNameEpicLink)
}

func (svc *CustomFieldService) GetCustomField(customFieldName string) (CustomField, error) {
	cfs, err := svc.GetCustomFields()
	if err != nil {
		return CustomField{}, err
	}
	cfsName := cfs.FilterByNames(customFieldName)
	if len(cfsName) != 1 {
		return CustomField{}, errors.New("epic link custom field not found")
	}
	return cfsName[0], nil
}
