package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type User struct {
	ID          string `json:"id,omitempty"`
	Email       string `json:"email,omitempty"`
	Name        string `json:"name,omitempty"`
	FullName    string `json:"fullName,omitempty"`
	SiteRole    string `json:"siteRole,omitempty"`
	AuthSetting string `json:"authSetting,omitempty"`
}

type UserRequest struct {
	User User `json:"user"`
}

type UserResponse struct {
	User User `json:"user"`
}

func (c *Client) GetUser(userID string) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/%s/", c.ApiUrl, userID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	userResponse := UserResponse{}
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return nil, err
	}

	return &userResponse.User, nil
}

func (c *Client) CreateUser(email, name, fullName, siteRole, authSetting string) (*User, error) {

	newUser := User{
		Email:       email,
		Name:        name,
		FullName:    fullName,
		SiteRole:    siteRole,
		AuthSetting: authSetting,
	}
	userRequest := UserRequest{
		User: newUser,
	}

	newUserJson, err := json.Marshal(userRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users", c.ApiUrl), strings.NewReader(string(newUserJson)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	userResponse := UserResponse{}
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return nil, err
	}

	return &userResponse.User, nil
}

func (c *Client) UpdateUser(userID, name, siteRole, authSetting string) (*User, error) {

	newUser := User{
		Name:        name,
		SiteRole:    siteRole,
		AuthSetting: authSetting,
	}
	userRequest := UserRequest{
		User: newUser,
	}

	newUserJson, err := json.Marshal(userRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/users/%s", c.ApiUrl, userID), strings.NewReader(string(newUserJson)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	userResponse := UserResponse{}
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return nil, err
	}

	return &userResponse.User, nil
}

func (c *Client) DeleteUser(userID string) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%s", c.ApiUrl, userID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
