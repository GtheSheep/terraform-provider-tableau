package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type VirtualConnection struct {
	ID      string
	Project struct {
		ID string `json:"id,omitempty"`
	} `json:"project,omitempty"`
	Owner struct {
		ID string `json:"id,omitempty"`
	} `json:"owner,omitempty"`
	Content string `json:"content,omitempty"`
	Name    string `json:"name,omitempty"`
}

type ListedVirtualConnection struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
	HasExtracts bool   `json:"hasExtracts,omitempty"`
	IsCertified bool   `json:"isCertified,omitempty"`
	WebPageURL  string `json:"webpageUrl,omitempty"`
}

type VirtualConnectionsRequest struct {
	VirtualConnection ListedVirtualConnection `json:"virtualConnection"`
}

type VirtualConnectionResponse struct {
	VirtualConnection VirtualConnection `json:"virtualConnection"`
}

type VirtualConnectionsResponse struct {
	VirtualConnections []ListedVirtualConnection `json:"virtualConnection"`
}

type VirtualConnectionsListResponse struct {
	VirtualConnectionsResponse VirtualConnectionsResponse `json:"virtualConnections"`
	Pagination                 PaginationDetails          `json:"pagination"`
}

func (c *Client) GetVirtualConnection(ID string) (*VirtualConnection, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/virtualconnections/%s", c.ApiUrl, ID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	virtualConnectionResponse := VirtualConnectionResponse{}
	err = json.Unmarshal(body, &virtualConnectionResponse)
	if err != nil {
		return nil, err
	}
	return &virtualConnectionResponse.VirtualConnection, nil
}

func (c *Client) GetVirtualConnections() ([]ListedVirtualConnection, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/virtualconnections", c.ApiUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	virtualConnectionListResponse := VirtualConnectionsListResponse{}
	err = json.Unmarshal(body, &virtualConnectionListResponse)
	if err != nil {
		return nil, err
	}

	pageNumber, totalPageCount, totalAvailable, err := GetPaginationNumbers(virtualConnectionListResponse.Pagination)
	if err != nil {
		return nil, err
	}

	allVirtualConnections := make([]ListedVirtualConnection, 0, totalAvailable)
	allVirtualConnections = append(allVirtualConnections, virtualConnectionListResponse.VirtualConnectionsResponse.VirtualConnections...)

	for page := pageNumber + 1; page <= totalPageCount; page++ {
		fmt.Printf("Searching page %d", page)
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/virtualconnections?pageNumber=%d", c.ApiUrl, page), nil)
		if err != nil {
			return nil, err
		}
		body, err = c.doRequest(req)
		if err != nil {
			return nil, err
		}
		virtualConnectionListResponse = VirtualConnectionsListResponse{}
		err = json.Unmarshal(body, &virtualConnectionListResponse)
		if err != nil {
			return nil, err
		}
		allVirtualConnections = append(allVirtualConnections, virtualConnectionListResponse.VirtualConnectionsResponse.VirtualConnections...)
	}

	return allVirtualConnections, nil
}
