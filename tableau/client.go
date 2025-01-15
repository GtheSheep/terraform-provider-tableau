package tableau

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	ApiUrl     string
	HTTPClient *http.Client
	AuthToken  string
}

type SiteDetails struct {
	ID         *string `json:"id"`
	ContentUrl string  `json:"contentUrl"`
}

type Credentials struct {
	Name        *string     `json:"name"`
	Password    *string     `json:"password"`
	TokenName   *string     `json:"personalAccessTokenName"`
	TokenSecret *string     `json:"personalAccessTokenSecret"`
	SiteDetails SiteDetails `json:"site"`
}

type SignInRequest struct {
	Credentials Credentials `json:"credentials"`
}

type SignInResponseData struct {
	SiteDetails               SiteDetails `json:"site"`
	User                      User        `json:"user"`
	Token                     string      `json:"token"`
	EstimatedTimeToExpiration string      `json:"estimatedTimeToExpiration"`
}

type SignInResponse struct {
	SignInResponseData SignInResponseData `json:"credentials"`
}

func NewClient(server, username, password, personalAccessTokenName, personalAccessTokenSecret, site, serverVersion *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	if (server != nil) && (username != nil) && (site != nil) && (serverVersion != nil) {
		baseUrl := fmt.Sprintf("%s/api/%s", *server, *serverVersion)
		url := fmt.Sprintf("%s/auth/signin", baseUrl)

		siteStruct := SiteDetails{ContentUrl: *site}
		credentials := Credentials{
			Name:        username,
			Password:    password,
			TokenName:   personalAccessTokenName,
			TokenSecret: personalAccessTokenSecret,
			SiteDetails: siteStruct,
		}
		authRequest := SignInRequest{
			Credentials: credentials,
		}
		authRequestJson, err := json.Marshal(authRequest)
		if err != nil {
			return nil, err
		}

		// authenticate
		req, err := http.NewRequest("POST", url, strings.NewReader(string(authRequestJson)))
		if err != nil {
			return nil, err
		}

		body, err := c.doRequest(req)
		if err != nil {
			return nil, err
		}

		// parse response body
		ar := SignInResponse{}
		err = json.Unmarshal(body, &ar)
		if err != nil {
			return nil, err
		}

		c.ApiUrl = fmt.Sprintf("%s/sites/%s", baseUrl, *ar.SignInResponseData.SiteDetails.ID)
		c.AuthToken = ar.SignInResponseData.Token
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Tableau-Auth", c.AuthToken)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != 201) && (res.StatusCode != 204) {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
