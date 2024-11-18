package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type VirtualConnectionPermission struct {
	VirtualConnectionID string
	EntityID            string
	EntityType          string
	CapabilityName      string
	CapabilityMode      string
}

type VirtualConnectionPermissions struct {
	GranteeCapabilities []GranteeCapability `json:"granteeCapabilities"`
}

type VirtualConnectionPermissionsRequest struct {
	VirtualConnectionPermissions VirtualConnectionPermissions `json:"permissions"`
}

type VirtualConnectionPermissionsResponse struct {
	VirtualConnectionPermissions VirtualConnectionPermissions `json:"permissions"`
}

func (c *Client) GetVirtualConnectionPermission(virtualConnectionID, entityID, entityType, capabilityName, capabilityMode string) (*VirtualConnectionPermission, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/virtualconnections/%s/permissions", c.ApiUrl, virtualConnectionID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	virtualConnectionPermissionsResponse := VirtualConnectionPermissionsResponse{}
	err = json.Unmarshal(body, &virtualConnectionPermissionsResponse)
	if err != nil {
		return nil, err
	}
	for _, granteeCapabilitie := range virtualConnectionPermissionsResponse.VirtualConnectionPermissions.GranteeCapabilities {
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
				return &VirtualConnectionPermission{
					VirtualConnectionID: virtualConnectionID,
					EntityID:            permissionEntityID,
					EntityType:          "users",
					CapabilityName:      capabilities.Name,
					CapabilityMode:      capabilities.Mode,
				}, nil
			}
			if entityType == "groups" && permissionEntityID == entityID && capabilityName == capabilities.Name && capabilities.Mode == capabilityMode {
				return &VirtualConnectionPermission{
					VirtualConnectionID: virtualConnectionID,
					EntityID:            permissionEntityID,
					EntityType:          "groups",
					CapabilityName:      capabilities.Name,
					CapabilityMode:      capabilities.Mode,
				}, nil
			}
		}
	}
	return nil, nil
}

func (c *Client) CreateVirtualConnectionPermissions(virtualConnectionID string, virtualConnectionPermissions VirtualConnectionPermissions) (*VirtualConnectionPermissions, error) {

	virtualConnectionPermissionsRequest := VirtualConnectionPermissionsRequest{
		VirtualConnectionPermissions: virtualConnectionPermissions,
	}

	newVirtualConnectionPermissionsJson, err := json.Marshal(virtualConnectionPermissionsRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/virtualconnections/%s/permissions", c.ApiUrl, virtualConnectionID), strings.NewReader(string(newVirtualConnectionPermissionsJson)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	virtualConnectionPermissionsResponse := VirtualConnectionPermissionsResponse{}
	err = json.Unmarshal(body, &virtualConnectionPermissionsResponse)
	if err != nil {
		return nil, err
	}

	return &virtualConnectionPermissionsResponse.VirtualConnectionPermissions, nil
}

func (c *Client) DeleteVirtualConnectionPermission(userID, groupID *string, virtualConnectionID, capabilityName, capabilityMode string) error {
	var entityID string
	entityType := "users"
	if userID != nil {
		entityID = *userID
	} else {
		entityType = "groups"
		entityID = *groupID
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/virtualconnections/%s/permissions/%s/%s/%s/%s", c.ApiUrl, virtualConnectionID, entityType, entityID, capabilityName, capabilityMode), nil)

	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
