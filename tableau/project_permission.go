package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Capability struct {
	Name string `json:"name"`
	Mode string `json:"mode"`
}

type Capabilities struct {
	Capabilities []Capability `json:"capability"`
}

type GranteeCapability struct {
	User         *User        `json:"user,omitempty"`
	Group        *Group       `json:"group,omitempty"`
	Capabilities Capabilities `json:"capabilities"`
}

type ProjectPermission struct {
	ProjectID      string
	EntityID       string
	EntityType     string
	CapabilityName string
	CapabilityMode string
}

type GranteeCapabilities struct {
	GranteeCapabilities []GranteeCapability `json:"granteeCapabilities"`
}

type ProjectPermissionsRequest struct {
	ProjectPermissions GranteeCapabilities `json:"permissions"`
}

type ProjectPermissionsResponse struct {
	ProjectPermissions GranteeCapabilities `json:"permissions"`
}

func (c *Client) GetProjectPermission(projectID, entityID, entityType, capabilityName, capabilityMode string) (*ProjectPermission, error) {
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
	for _, granteeCapabilitie := range projectPermissionsResponse.ProjectPermissions.GranteeCapabilities {
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
				return &ProjectPermission{
					ProjectID:      projectID,
					EntityID:       permissionEntityID,
					EntityType:     "users",
					CapabilityName: capabilities.Name,
					CapabilityMode: capabilities.Mode,
				}, nil
			}
			if entityType == "groups" && permissionEntityID == entityID && capabilityName == capabilities.Name && capabilities.Mode == capabilityMode {
				return &ProjectPermission{
					ProjectID:      projectID,
					EntityID:       permissionEntityID,
					EntityType:     "groups",
					CapabilityName: capabilities.Name,
					CapabilityMode: capabilities.Mode,
				}, nil
			}
		}
	}
	return nil, nil
}

func (c *Client) CreateProjectPermissions(projectID string, projectPermissions GranteeCapabilities) (*GranteeCapabilities, error) {

	projectPermissionsRequest := ProjectPermissionsRequest{
		ProjectPermissions: projectPermissions,
	}

	newProjectPermissionsJson, err := json.Marshal(projectPermissionsRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/projects/%s/permissions", c.ApiUrl, projectID), strings.NewReader(string(newProjectPermissionsJson)))
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

func (c *Client) DeleteProjectPermission(userID, groupID *string, projectID, capabilityName, capabilityMode string) error {
	var entityID string
	entityType := "users"
	if userID != nil {
		entityID = *userID
	} else {
		entityType = "groups"
		entityID = *groupID
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s/permissions/%s/%s/%s/%s", c.ApiUrl, projectID, entityType, entityID, capabilityName, capabilityMode), nil)

	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
