package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ViewPermission struct {
	ViewID         string
	EntityID       string
	EntityType     string
	CapabilityName string
	CapabilityMode string
}

type ViewPermissions struct {
	GranteeCapabilities []GranteeCapability `json:"granteeCapabilities"`
}

type ViewPermissionsRequest struct {
	ViewPermissions ViewPermissions `json:"permissions"`
}

type ViewPermissionsResponse struct {
	ViewPermissions ViewPermissions `json:"permissions"`
}

func (c *Client) GetViewPermission(viewID, entityID, entityType, capabilityName, capabilityMode string) (*ViewPermission, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/views/%s/permissions", c.ApiUrl, viewID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	viewPermissionsResponse := ViewPermissionsResponse{}
	err = json.Unmarshal(body, &viewPermissionsResponse)
	if err != nil {
		return nil, err
	}
	for _, granteeCapabilitie := range viewPermissionsResponse.ViewPermissions.GranteeCapabilities {
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
				return &ViewPermission{
					ViewID:         viewID,
					EntityID:       permissionEntityID,
					EntityType:     "users",
					CapabilityName: capabilities.Name,
					CapabilityMode: capabilities.Mode,
				}, nil
			}
			if entityType == "groups" && permissionEntityID == entityID && capabilityName == capabilities.Name && capabilities.Mode == capabilityMode {
				return &ViewPermission{
					ViewID:         viewID,
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

func (c *Client) CreateViewPermissions(viewID string, viewPermissions ViewPermissions) (*ViewPermissions, error) {

	viewPermissionsRequest := ViewPermissionsRequest{
		ViewPermissions: viewPermissions,
	}

	newViewPermissionsJson, err := json.Marshal(viewPermissionsRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/views/%s/permissions", c.ApiUrl, viewID), strings.NewReader(string(newViewPermissionsJson)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	viewPermissionsResponse := ViewPermissionsResponse{}
	err = json.Unmarshal(body, &viewPermissionsResponse)
	if err != nil {
		return nil, err
	}

	return &viewPermissionsResponse.ViewPermissions, nil
}

func (c *Client) DeleteViewPermission(userID, groupID *string, viewID, capabilityName, capabilityMode string) error {
	var entityID string
	entityType := "users"
	if userID != nil {
		entityID = *userID
	} else {
		entityType = "groups"
		entityID = *groupID
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/views/%s/permissions/%s/%s/%s/%s", c.ApiUrl, viewID, entityType, entityID, capabilityName, capabilityMode), nil)

	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
