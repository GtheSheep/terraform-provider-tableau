package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type GroupImport struct {
	MinimumSiteRole  *string `json:"minimumSiteRole"`
	GrantLicenseMode *string `json:"grantLicenseMode"`
}

type NewGroup struct {
	Name            string  `json:"name"`
	MinimumSiteRole *string `json:"minimumSiteRole"`
}

type Group struct {
	ID     *string     `json:"id"`
	Name   string      `json:"name"`
	Import GroupImport `json:"import"`
}

type PaginationDetails struct {
	PageNumber     int `json:"pageNumber"`
	PageSize       int `json:"pageSize"`
	TotalAvailable int `json:"totalAvailable"`
}

type GroupResponse struct {
	Group Group `json:"group"`
}

type GroupsResponse struct {
	Groups []Group `json:"group"`
}

type GroupListResponse struct {
	GroupsResponse GroupsResponse    `json:"groups"`
	Pagination     PaginationDetails `json:"pagination"`
}

func (c *Client) GetGroup(groupID string) (*Group, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/groups", c.ApiUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	groupListResponse := GroupListResponse{}
	err = json.Unmarshal(body, &groupListResponse)
	if err != nil {
		return nil, err
	}

	for i, group := range groupListResponse.GroupsResponse.Groups {
		if *group.ID == groupID {
			return &groupListResponse.GroupsResponse.Groups[i], nil
		}
	}

	return nil, fmt.Errorf("Did not find group ID %s", groupID)
}

func (c *Client) CreateGroup(name, minimumSiteRole string) (*Group, error) {

	newGroup := NewGroup{
		Name:            name,
		MinimumSiteRole: &minimumSiteRole,
	}

	newGroupJson, err := json.Marshal(newGroup)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/groups", c.ApiUrl), strings.NewReader(string(newGroupJson)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	groupResponse := GroupResponse{}
	err = json.Unmarshal(body, &groupResponse)
	if err != nil {
		return nil, err
	}

	return &groupResponse.Group, nil
}

func (c *Client) UpdateGroup(groupID, name, minimumSiteRole string) (*Group, error) {

	newGroup := NewGroup{
		Name:            name,
		MinimumSiteRole: &minimumSiteRole,
	}

	newGroupJson, err := json.Marshal(newGroup)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/groups/%s", c.ApiUrl, groupID), strings.NewReader(string(newGroupJson)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	groupResponse := GroupResponse{}
	err = json.Unmarshal(body, &groupResponse)
	if err != nil {
		return nil, err
	}

	return &groupResponse.Group, nil
}

func (c *Client) DeleteGroup(groupID string) (*Group, error) {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/groups/%s", c.ApiUrl, groupID), nil)
	if err != nil {
		return nil, err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
