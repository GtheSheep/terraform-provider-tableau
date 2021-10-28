package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type GroupUsersResponse struct {
	Users []User `json:"user"`
}

type GroupUsersListResponse struct {
	GroupUsersResponse GroupUsersResponse `json:"users"`
	Pagination         PaginationDetails  `json:"pagination"`
}

func (c *Client) GetGroupUser(groupID, userID string) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/groups/%s/users", c.ApiUrl, groupID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	groupUsersListResponse := GroupUsersListResponse{}
	err = json.Unmarshal(body, &groupUsersListResponse)
	if err != nil {
		return nil, err
	}

	for i, user := range groupUsersListResponse.GroupUsersResponse.Users {
		if *user.ID == userID {
			return &groupUsersListResponse.GroupUsersResponse.Users[i], nil
		}
	}

	return nil, fmt.Errorf("Did not find user ID %s in group ID %s", userID, groupID)
}

func (c *Client) CreateGroupUser(groupID, userID string) (*User, error) {

	newGroupUser := User{
		ID: &userID,
	}

	newGroupUserJson, err := json.Marshal(newGroupUser)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/groups/%s/users", c.ApiUrl, groupID), strings.NewReader(string(newGroupUserJson)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	groupUserResponse := UserResponse{}
	err = json.Unmarshal(body, &groupUserResponse)
	if err != nil {
		return nil, err
	}

	return &groupUserResponse.User, nil
}

func (c *Client) DeleteGroupUser(groupID, userID string) (*User, error) {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/groups/%s/users", c.ApiUrl, groupID), nil)
	if err != nil {
		return nil, err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
