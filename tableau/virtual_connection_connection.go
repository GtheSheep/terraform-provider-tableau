package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type VirtualConnectionConnection struct {
	VirtualConnectionID string
	ID                  string `json:"connectionId,omitempty"`
	DBClass             string `json:"dbClass"`
	ServerAddress       string `json:"server"`
	ServerPort          string `json:"port"`
	UserName            string `json:"username"`
}

type VirtualConnectionConnectionRequest struct {
	VirtualConnectionConnection VirtualConnectionConnection `json:"virtualConnectionConnections"`
}

type VirtualConnectionConnectionsResponse struct {
	VirtualConnectionConnections []VirtualConnectionConnection `json:"connection"`
}

type VirtualConnectionConnectionListResponse struct {
	VirtualConnectionConnectionsResponse VirtualConnectionConnectionsResponse `json:"virtualConnectionConnections"`
	Pagination                           PaginationDetails                    `json:"pagination"`
}

func (c *Client) GetVirtualConnectionConnections(virtualConnectionID string) ([]VirtualConnectionConnection, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/virtualconnections/%s/connections", c.ApiUrl, virtualConnectionID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	virtualConnectionConnectionsListResponse := VirtualConnectionConnectionListResponse{}
	err = json.Unmarshal(body, &virtualConnectionConnectionsListResponse)
	if err != nil {
		return nil, err
	}

	pageNumber, totalPageCount, totalAvailable, err := GetPaginationNumbers(virtualConnectionConnectionsListResponse.Pagination)
	if err != nil {
		return nil, err
	}

	allVirtualConnectionConnections := make([]VirtualConnectionConnection, 0, totalAvailable)
	allVirtualConnectionConnections = append(allVirtualConnectionConnections, virtualConnectionConnectionsListResponse.VirtualConnectionConnectionsResponse.VirtualConnectionConnections...)
	for page := pageNumber + 1; page <= totalPageCount; page++ {
		fmt.Printf("Searching page %d", page)
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/virtualconnections/%s/connections?pageNumber=%d", c.ApiUrl, virtualConnectionID, page), nil)
		if err != nil {
			return nil, err
		}
		body, err = c.doRequest(req)
		if err != nil {
			return nil, err
		}
		virtualConnectionConnectionsListResponse = VirtualConnectionConnectionListResponse{}
		err = json.Unmarshal(body, &virtualConnectionConnectionsListResponse)
		if err != nil {
			return nil, err
		}
		allVirtualConnectionConnections = append(allVirtualConnectionConnections, virtualConnectionConnectionsListResponse.VirtualConnectionConnectionsResponse.VirtualConnectionConnections...)
	}
	for idx := range allVirtualConnectionConnections {
		allVirtualConnectionConnections[idx].VirtualConnectionID = virtualConnectionID
	}
	return allVirtualConnectionConnections, nil
}
