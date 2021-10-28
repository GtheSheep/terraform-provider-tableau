package tableau

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type User struct {
	ID          *string `json:"id"`
	Email       *string `json:"email"`
	Name        *string `json:"name"`
	FullName    *string `json:"fullName"`
	SiteRole    *string `json:"siteRole"`
	AuthSetting *string `json:"authSetting"`
}

type UserResponse struct {
	User User `json:"user"`
}

func (c *Client) GetUser(userID string) (*User, error) {
	log.Printf(fmt.Sprintf("%s/users/%s/", c.ApiUrl, userID))
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
		Email:       &email,
		Name:        &name,
		FullName:    &fullName,
		SiteRole:    &siteRole,
		AuthSetting: &authSetting,
	}

	newUserJson, err := json.Marshal(newUser)
	if err != nil {
		return nil, err
	}

	log.Printf(string(newUserJson))
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
		Name:        &name,
		SiteRole:    &siteRole,
		AuthSetting: &authSetting,
	}

	newUserJson, err := json.Marshal(newUser)
	if err != nil {
		return nil, err
	}

	log.Printf(string(newUserJson))
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

func (c *Client) DeleteUser(userID string) (*User, error) {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%s", c.ApiUrl, userID), nil)
	if err != nil {
		return nil, err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
