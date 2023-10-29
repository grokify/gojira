package jirarest

type API struct {
	Client       *Client
	Backlog      *BacklogService
	CustomFields *CustomFieldsService
}

func NewAPI(client *Client) API {
	return API{
		Client:       client,
		Backlog:      NewBacklogService(client),
		CustomFields: NewCustomFieldsService(client),
	}
}
