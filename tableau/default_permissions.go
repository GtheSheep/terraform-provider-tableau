package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DefaultPermissions struct {
	ProjectID           string
	TargetType          string
	GranteeCapabilities []GranteeCapability `json:"granteeCapabilities"`
}

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

func (c *Client) GetDefaultPermissions(projectID, targetType string) (*DefaultPermissions, error) {
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
	defaultPermissions := &DefaultPermissions{
		ProjectID:           projectID,
		TargetType:          targetType,
		GranteeCapabilities: projectPermissionsResponse.ProjectPermissions.GranteeCapabilities,
	}
	return defaultPermissions, nil
}

func (c *Client) CreateDefaultPermissions(projectID, targetType string, defaultPermissions []GranteeCapability) (*GranteeCapabilities, error) {

	defaultPermissionsRequest := ProjectPermissionsRequest{
		ProjectPermissions: GranteeCapabilities{
			GranteeCapabilities: defaultPermissions,
		},
	}

	newDefaultPermissionsJson, err := json.Marshal(defaultPermissionsRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/projects/%s/default-permissions/%s", c.ApiUrl, projectID, targetType),
		strings.NewReader(string(newDefaultPermissionsJson)),
	)
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

func (c *Client) DeleteDefaultPermission(userID, groupID *string, projectID, targetType, capabilityName, capabilityMode string) error {
	var entityID string
	entityType := "users"
	if userID != nil && *userID != "" {
		entityID = *userID
	} else if groupID != nil && *groupID != "" {
		entityType = "groups"
		entityID = *groupID
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s/default-permissions/%s/%s/%s/%s/%s", c.ApiUrl, projectID, targetType, entityType, entityID, capabilityName, capabilityMode), nil)
	if err != nil {
		return err
	}
	_, err = c.doRequest(req)
	if err != nil {
		return err
	}
	return nil
}
