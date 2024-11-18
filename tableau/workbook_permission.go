package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type WorkbookPermission struct {
	WorkbookID     string
	EntityID       string
	EntityType     string
	CapabilityName string
	CapabilityMode string
}

type WorkbookPermissions struct {
	GranteeCapabilities []GranteeCapability `json:"granteeCapabilities"`
}

type WorkbookPermissionsRequest struct {
	WorkbookPermissions WorkbookPermissions `json:"permissions"`
}

type WorkbookPermissionsResponse struct {
	WorkbookPermissions WorkbookPermissions `json:"permissions"`
}

func (c *Client) GetWorkbookPermission(workbookID, entityID, entityType, capabilityName, capabilityMode string) (*WorkbookPermission, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/workbooks/%s/permissions", c.ApiUrl, workbookID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workbookPermissionsResponse := WorkbookPermissionsResponse{}
	err = json.Unmarshal(body, &workbookPermissionsResponse)
	if err != nil {
		return nil, err
	}
	for _, granteeCapabilitie := range workbookPermissionsResponse.WorkbookPermissions.GranteeCapabilities {
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
				return &WorkbookPermission{
					WorkbookID:     workbookID,
					EntityID:       permissionEntityID,
					EntityType:     "users",
					CapabilityName: capabilities.Name,
					CapabilityMode: capabilities.Mode,
				}, nil
			}
			if entityType == "groups" && permissionEntityID == entityID && capabilityName == capabilities.Name && capabilities.Mode == capabilityMode {
				return &WorkbookPermission{
					WorkbookID:     workbookID,
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

func (c *Client) CreateWorkbookPermissions(workbookID string, workbookPermissions WorkbookPermissions) (*WorkbookPermissions, error) {

	workbookPermissionsRequest := WorkbookPermissionsRequest{
		WorkbookPermissions: workbookPermissions,
	}

	newWorkbookPermissionsJson, err := json.Marshal(workbookPermissionsRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/workbooks/%s/permissions", c.ApiUrl, workbookID), strings.NewReader(string(newWorkbookPermissionsJson)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workbookPermissionsResponse := WorkbookPermissionsResponse{}
	err = json.Unmarshal(body, &workbookPermissionsResponse)
	if err != nil {
		return nil, err
	}

	return &workbookPermissionsResponse.WorkbookPermissions, nil
}

func (c *Client) DeleteWorkbookPermission(userID, groupID *string, workbookID, capabilityName, capabilityMode string) error {
	var entityID string
	entityType := "users"
	if userID != nil {
		entityID = *userID
	} else {
		entityType = "groups"
		entityID = *groupID
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/workbooks/%s/permissions/%s/%s/%s/%s", c.ApiUrl, workbookID, entityType, entityID, capabilityName, capabilityMode), nil)

	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
