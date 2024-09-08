package jirarest

type API struct {
	Client       *Client
	Backlog      *BacklogService
	CustomFields *CustomFieldService
}

func NewAPI(client *Client) API {
	return API{
		Client:       client,
		Backlog:      NewBacklogService(client),
		CustomFields: NewCustomFieldService(client),
	}
}
