package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DatasourcePermission struct {
	DatasourceID   string
	EntityID       string
	EntityType     string
	CapabilityName string
	CapabilityMode string
}

type DatasourcePermissions struct {
	GranteeCapabilities []GranteeCapability `json:"granteeCapabilities"`
}

type DatasourcePermissionsRequest struct {
	DatasourcePermissions DatasourcePermissions `json:"permissions"`
}

type DatasourcePermissionsResponse struct {
	DatasourcePermissions DatasourcePermissions `json:"permissions"`
}

func (c *Client) GetDatasourcePermission(datasourceID, entityID, entityType, capabilityName, capabilityMode string) (*DatasourcePermission, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/datasources/%s/permissions", c.ApiUrl, datasourceID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	datasourcePermissionsResponse := DatasourcePermissionsResponse{}
	err = json.Unmarshal(body, &datasourcePermissionsResponse)
	if err != nil {
		return nil, err
	}
	for _, granteeCapabilitie := range datasourcePermissionsResponse.DatasourcePermissions.GranteeCapabilities {
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
				return &DatasourcePermission{
					DatasourceID:   datasourceID,
					EntityID:       permissionEntityID,
					EntityType:     "users",
					CapabilityName: capabilities.Name,
					CapabilityMode: capabilities.Mode,
				}, nil
			}
			if entityType == "groups" && permissionEntityID == entityID && capabilityName == capabilities.Name && capabilities.Mode == capabilityMode {
				return &DatasourcePermission{
					DatasourceID:   datasourceID,
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

func (c *Client) CreateDatasourcePermissions(datasourceID string, datasourcePermissions DatasourcePermissions) (*DatasourcePermissions, error) {

	datasourcePermissionsRequest := DatasourcePermissionsRequest{
		DatasourcePermissions: datasourcePermissions,
	}

	newDatasourcePermissionsJson, err := json.Marshal(datasourcePermissionsRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/datasources/%s/permissions", c.ApiUrl, datasourceID), strings.NewReader(string(newDatasourcePermissionsJson)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	datasourcePermissionsResponse := DatasourcePermissionsResponse{}
	err = json.Unmarshal(body, &datasourcePermissionsResponse)
	if err != nil {
		return nil, err
	}

	return &datasourcePermissionsResponse.DatasourcePermissions, nil
}

func (c *Client) DeleteDatasourcePermission(userID, groupID *string, datasourceID, capabilityName, capabilityMode string) error {
	var entityID string
	entityType := "users"
	if userID != nil {
		entityID = *userID
	} else {
		entityType = "groups"
		entityID = *groupID
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/datasources/%s/permissions/%s/%s/%s/%s", c.ApiUrl, datasourceID, entityType, entityID, capabilityName, capabilityMode), nil)

	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
