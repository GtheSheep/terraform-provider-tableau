package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var defaultPermissionTargetTypes = []string{
	"databases",
	"dataroles",
	"datasources",
	"flows",
	"lenses",
	"metrics",
	"tables",
	"virtualconnections",
	"workbooks",
}

func (c *Client) GetDefaultPermissions(projectID, targetType string) (*ProjectPermissions, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%s/default-permissions/%s", c.ApiUrl, projectID, targetType), nil)
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
