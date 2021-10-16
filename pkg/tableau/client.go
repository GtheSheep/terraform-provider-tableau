package tableau

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	ApiUrl     string
	HTTPClient *http.Client
	AuthToken  string
}

type Site struct {
	ContentUrl string `json:"contentUrl`
}

type Credentials struct {
	Name     string `json:"name`
	Password string `json:"password`
	Site     Site   `json:"site`
}

type SignInRequest struct {
	Credentials Credentials `json:"credentials`
}

type SignInResponse struct {
}

func NewClient(server, username, password, site, server_version *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	if (server != nil) && (username != nil) && (site != nil) && (server_version != nil) {
		url := fmt.Sprintf("%s/api/%s/auth/signin", server, server_version)

		// authenticate
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return nil, err
		}

		body, err := c.doRequest(req)

		// parse response body
		ar := SignInResponse{}
		err = json.Unmarshal(body, &ar)
		if err != nil {
			return nil, err
		}

		c.ApiUrl = url
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-tableau-auth", c.AuthToken)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != 201) {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
