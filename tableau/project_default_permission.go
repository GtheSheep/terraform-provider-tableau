package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ProjectDefaultPermission struct {
	ProjectID      string
	TargetType     string
	EntityID       string
	EntityType     string
	CapabilityName string
	CapabilityMode string
}

type ProjectDefaultPermissions struct {
	GranteeCapabilities []GranteeCapability `json:"granteeCapabilities"`
}

type ProjectDefaultPermissionsRequest struct {
	ProjectDefaultPermissions ProjectDefaultPermissions `json:"permissions"`
}

type ProjectDefaultPermissionsResponse struct {
	ProjectDefaultPermissions ProjectDefaultPermissions `json:"permissions"`
}

func (c *Client) GetProjectDefaultPermission(projectID, entityID, entityType, targetType, capabilityName, capabilityMode string) (*ProjectDefaultPermission, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%s/default-permissions/%s", c.ApiUrl, projectID, targetType), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	projectDefaultPermissionsResponse := ProjectDefaultPermissionsResponse{}
	err = json.Unmarshal(body, &projectDefaultPermissionsResponse)
	if err != nil {
		return nil, err
	}
	for _, granteeCapabilitie := range projectDefaultPermissionsResponse.ProjectDefaultPermissions.GranteeCapabilities {
		for _, capabilities := range granteeCapabilitie.Capabilities.Capabilities {
			var permissionEntityID string
			if granteeCapabilitie.User != nil {
				entity := granteeCapabilitie.User
				permissionEntityID = entity.ID
			} else {
				entity := granteeCapabilitie.Group
				permissionEntityID = entity.ID
			}
			if entityType == "users" && permissionEntityID == entityID && capabilityName == capabilities.Name && capabilities.Mode == capabilityMode {
				return &ProjectDefaultPermission{
					ProjectID:      projectID,
					EntityID:       permissionEntityID,
					EntityType:     "users",
					TargetType:     targetType,
					CapabilityName: capabilities.Name,
					CapabilityMode: capabilities.Mode,
				}, nil
			}
			if entityType == "groups" && permissionEntityID == entityID && capabilityName == capabilities.Name && capabilities.Mode == capabilityMode {
				return &ProjectDefaultPermission{
					ProjectID:      projectID,
					EntityID:       permissionEntityID,
					EntityType:     "groups",
					TargetType:     targetType,
					CapabilityName: capabilities.Name,
					CapabilityMode: capabilities.Mode,
				}, nil
			}
		}
	}
	return nil, nil
}

func (c *Client) CreateProjectDefaultPermissions(projectID string, targetType string, projectDefaultPermissions ProjectDefaultPermissions) (*ProjectDefaultPermissions, error) {

	projectDefaultPermissionsRequest := ProjectDefaultPermissionsRequest{
		ProjectDefaultPermissions: projectDefaultPermissions,
	}

	newProjectDefaultPermissionsJson, err := json.Marshal(projectDefaultPermissionsRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/projects/%s/default-permissions/%s", c.ApiUrl, projectID, targetType), strings.NewReader(string(newProjectDefaultPermissionsJson)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	projectDefaultPermissionsResponse := ProjectDefaultPermissionsResponse{}
	err = json.Unmarshal(body, &projectDefaultPermissionsResponse)
	if err != nil {
		return nil, err
	}

	return &projectDefaultPermissionsResponse.ProjectDefaultPermissions, nil
}

func (c *Client) DeleteProjectDefaultPermission(userID, groupID *string, projectID, targetType, capabilityName, capabilityMode string) error {
	var entityID string
	entityType := "users"
	if userID != nil {
		entityID = *userID
	} else {
		entityType = "groups"
		entityID = *groupID
	}

	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf(
			"%s/projects/%s/default-permissions/%s",
			c.ApiUrl, projectID,
			strings.Join([]string{entityType, entityID, targetType, capabilityName, capabilityMode}, "/"),
		),
		nil,
	)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
