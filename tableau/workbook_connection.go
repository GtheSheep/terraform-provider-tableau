package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WorkbookConnection struct {
	WorkbookID              string
	ID                      string `json:"id,omitempty"`
	Type                    string `json:"type,omitempty"`
	ServerAddress           string `json:"serverAddress,omitempty"`
	ServerPort              string `json:"serverPort,omitempty"`
	UserName                string `json:"userName,omitempty"`
	QueryTaggingEnabled     bool   `json:"query_tagging_enabled,omitempty"`
	AuthenticationType      string `json:"authenticationType,omitempty"`
	EmbedPassword           bool   `json:"embedPassword,omitempty"`
	UseOAuthManagedKeychain bool   `json:"useOauthManagedKeychain,omitempty"`
	DataSourceID            struct {
		ID string `json:"id,omitempty"`
		// Name string `json:"name,omitempty"`
	} `json:"datasource,omitempty"`
}

type WorkbookConnectionRequest struct {
	WorkbookConnection WorkbookConnection `json:"workbookConnections"`
}

type WorkbookConnectionsResponse struct {
	WorkbookConnections []WorkbookConnection `json:"connection"`
}

type WorkbookConnectionListResponse struct {
	WorkbookConnectionsResponse WorkbookConnectionsResponse `json:"connections"`
}

func (c *Client) GetWorkbookConnections(workbookID string) ([]WorkbookConnection, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/workbooks/%s/connections", c.ApiUrl, workbookID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	workbookConnectionsListResponse := WorkbookConnectionListResponse{}
	err = json.Unmarshal(body, &workbookConnectionsListResponse)
	if err != nil {
		return nil, err
	}
	// workbook connections don't seem to have pagination
	allWorkbookConnections := workbookConnectionsListResponse.WorkbookConnectionsResponse.WorkbookConnections
	for idx := range allWorkbookConnections {
		allWorkbookConnections[idx].WorkbookID = workbookID
	}
	return allWorkbookConnections, nil
}
