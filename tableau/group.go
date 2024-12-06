package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type GroupImport struct {
	DomainName       *string `json:"domainName"`
	MinimumSiteRole  *string `json:"siteRole"`
	GrantLicenseMode *string `json:"grantLicenseMode"`
}

type Group struct {
	ID              string       `json:"id,omitempty"`
	Name            string       `json:"name"`
	MinimumSiteRole string       `json:"minimumSiteRole,omitempty"`
	OnDemandAccess  *bool        `json:"externalUserEnabled,omitempty"`
	Import          *GroupImport `json:"import,omitempty"`
}

type GroupRequest struct {
	Group Group `json:"group"`
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

func (c *Client) GetGroups() ([]Group, error) {
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

	// TODO: Generalise pagination handling and use elsewhere
	pageNumber, totalPageCount, totalAvailable, err := GetPaginationNumbers(groupListResponse.Pagination)
	if err != nil {
		return nil, err
	}

	allGroups := make([]Group, 0, totalAvailable)
	for _, group := range groupListResponse.GroupsResponse.Groups {
		allGroups = append(allGroups, group)
	}

	for page := pageNumber + 1; page <= totalPageCount; page++ {
		fmt.Printf("Searching page %d", page)
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/groups?pageNumber=%s", c.ApiUrl, strconv.Itoa(page)), nil)
		if err != nil {
			return nil, err
		}
		body, err = c.doRequest(req)
		if err != nil {
			return nil, err
		}
		groupListResponse = GroupListResponse{}
		err = json.Unmarshal(body, &groupListResponse)
		if err != nil {
			return nil, err
		}
		for _, group := range groupListResponse.GroupsResponse.Groups {
			allGroups = append(allGroups, group)
		}
	}

	return allGroups, nil
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

	// TODO: Generalise pagination handling and use elsewhere
	pageNumber, totalPageCount, _, err := GetPaginationNumbers(groupListResponse.Pagination)
	if err != nil {
		return nil, err
	}
	for i, group := range groupListResponse.GroupsResponse.Groups {
		if group.ID == groupID {
			return &groupListResponse.GroupsResponse.Groups[i], nil
		}
	}

	for page := pageNumber + 1; page <= totalPageCount; page++ {
		fmt.Printf("Searching page %d", page)
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/groups?pageNumber=%s", c.ApiUrl, strconv.Itoa(page)), nil)
		if err != nil {
			return nil, err
		}
		body, err = c.doRequest(req)
		if err != nil {
			return nil, err
		}
		groupListResponse = GroupListResponse{}
		err = json.Unmarshal(body, &groupListResponse)
		if err != nil {
			return nil, err
		}

		for i, group := range groupListResponse.GroupsResponse.Groups {
			if group.ID == groupID {
				return &groupListResponse.GroupsResponse.Groups[i], nil
			}
		}
	}

	return nil, fmt.Errorf("Did not find group ID %s", groupID)
}

func (c *Client) CreateGroup(name string, minimumSiteRole *string, onDemandAccess *bool) (*Group, error) {
	group := Group{
		Name: name,
	}

	if minimumSiteRole != nil && *minimumSiteRole != "" {
		group.MinimumSiteRole = *minimumSiteRole
	}

	if onDemandAccess != nil {
		group.OnDemandAccess = onDemandAccess
	}

	groupRequest := GroupRequest{Group: group}
	groupJSON, err := json.Marshal(groupRequest)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/groups", c.ApiUrl), strings.NewReader(string(groupJSON)))
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

func (c *Client) UpdateGroup(groupID, name string, minimumSiteRole *string, onDemandAccess *bool) (*Group, error) {
	group := Group{
		Name: name,
	}

	if minimumSiteRole != nil && *minimumSiteRole != "" {
		group.MinimumSiteRole = *minimumSiteRole
	}

	if onDemandAccess != nil {
		group.OnDemandAccess = onDemandAccess
	}

	groupRequest := GroupRequest{
		Group: group,
	}

	updateGroupJson, err := json.Marshal(groupRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/groups/%s", c.ApiUrl, groupID), strings.NewReader(string(updateGroupJson)))
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

func (c *Client) DeleteGroup(groupID string) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/groups/%s", c.ApiUrl, groupID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
