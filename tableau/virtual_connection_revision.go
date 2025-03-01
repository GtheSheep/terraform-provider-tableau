package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type VirtualConnectionRevision struct {
	VirtualConnectionID string
	Publisher           struct {
		ID string `json:"id,omitempty"`
		// Name string `json:"name,omitempty"`
	} `json:"publisher,omitempty"`
	Current        bool   `json:"current,omitempty"`
	Deleted        bool   `json:"deleted,omitempty"`
	PublishedAt    string `json:"publishedAt,omitempty"`
	RevisionNumber string `json:"revisionNumber,omitempty"`
}

type VirtualConnectionRevisionRequest struct {
	VirtualConnectionRevision VirtualConnectionRevision `json:"virtualConnectionRevisions"`
}

type VirtualConnectionRevisionsResponse struct {
	VirtualConnectionRevisions []VirtualConnectionRevision `json:"revision"`
}

type VirtualConnectionRevisionListResponse struct {
	VirtualConnectionRevisionsResponse VirtualConnectionRevisionsResponse `json:"revisions"`
	Pagination                         PaginationDetails                  `json:"pagination"`
}

func (c *Client) GetVirtualConnectionRevisions(virtualConnectionID string) ([]VirtualConnectionRevision, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/virtualconnections/%s/revisions", c.ApiUrl, virtualConnectionID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	virtualConnectionRevisionsListResponse := VirtualConnectionRevisionListResponse{}
	err = json.Unmarshal(body, &virtualConnectionRevisionsListResponse)
	if err != nil {
		return nil, err
	}
	pageNumber, totalPageCount, totalAvailable, err := GetPaginationNumbers(virtualConnectionRevisionsListResponse.Pagination)
	if err != nil {
		return nil, err
	}

	allVirtualConnectionRevisions := make([]VirtualConnectionRevision, 0, totalAvailable)
	allVirtualConnectionRevisions = append(allVirtualConnectionRevisions, virtualConnectionRevisionsListResponse.VirtualConnectionRevisionsResponse.VirtualConnectionRevisions...)
	for page := pageNumber + 1; page <= totalPageCount; page++ {
		fmt.Printf("Searching page %d", page)
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/virtualconnections/%s/revisions?pageNumber=%d", c.ApiUrl, virtualConnectionID, page), nil)
		if err != nil {
			return nil, err
		}
		body, err = c.doRequest(req)
		if err != nil {
			return nil, err
		}
		virtualConnectionRevisionsListResponse = VirtualConnectionRevisionListResponse{}
		err = json.Unmarshal(body, &virtualConnectionRevisionsListResponse)
		if err != nil {
			return nil, err
		}
		allVirtualConnectionRevisions = append(allVirtualConnectionRevisions, virtualConnectionRevisionsListResponse.VirtualConnectionRevisionsResponse.VirtualConnectionRevisions...)
	}
	for idx := range allVirtualConnectionRevisions {
		allVirtualConnectionRevisions[idx].VirtualConnectionID = virtualConnectionID
	}
	return allVirtualConnectionRevisions, nil
}
