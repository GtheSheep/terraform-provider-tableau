package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetProjectPermissions(projectID string) (*GranteeCapabilities, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%s/permissions", c.ApiUrl, projectID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	projectPermissionsResponse := ProjectPermissionsResponse{}
	err = json.Unmarshal(body, &projectPermissionsResponse)
	if err != nil {
		return nil, err
	}
	return &projectPermissionsResponse.ProjectPermissions, nil
}
