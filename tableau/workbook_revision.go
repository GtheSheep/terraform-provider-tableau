package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type WorkbookRevision struct {
	WorkbookID     string
	Current        bool   `json:"current,omitempty"`
	Deleted        bool   `json:"deleted,omitempty"`
	PublishedAt    string `json:"publishedAt,omitempty"`
	RevisionNumber string `json:"revisionNumber,omitempty"`
	Publisher      struct {
		ID string `json:"id,omitempty"`
		// Name string `json:"name,omitempty"`
	} `json:"publisher,omitempty"`
}

type WorkbookRevisionRequest struct {
	WorkbookRevision WorkbookRevision `json:"workbookRevisions"`
}

type WorkbookRevisionsResponse struct {
	WorkbookRevisions []WorkbookRevision `json:"revision"`
}

type WorkbookRevisionListResponse struct {
	WorkbookRevisionsResponse WorkbookRevisionsResponse `json:"revisions"`
	Pagination                PaginationDetails         `json:"pagination"`
}

func (c *Client) GetWorkbookRevisions(workbookID string) ([]WorkbookRevision, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/workbooks/%s/revisions", c.ApiUrl, workbookID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	workbookRevisionsListResponse := WorkbookRevisionListResponse{}
	err = json.Unmarshal(body, &workbookRevisionsListResponse)
	if err != nil {
		return nil, err
	}
	pageNumber, totalPageCount, totalAvailable, err := GetPaginationNumbers(workbookRevisionsListResponse.Pagination)
	if err != nil {
		return nil, err
	}

	allWorkbookRevisions := make([]WorkbookRevision, 0, totalAvailable)
	allWorkbookRevisions = append(allWorkbookRevisions, workbookRevisionsListResponse.WorkbookRevisionsResponse.WorkbookRevisions...)
	for page := pageNumber + 1; page <= totalPageCount; page++ {
		fmt.Printf("Searching page %d", page)
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/workbooks/%s/revisions?pageNumber=%s", c.ApiUrl, workbookID, strconv.Itoa(page)), nil)
		if err != nil {
			return nil, err
		}
		body, err = c.doRequest(req)
		if err != nil {
			return nil, err
		}
		workbookRevisionsListResponse = WorkbookRevisionListResponse{}
		err = json.Unmarshal(body, &workbookRevisionsListResponse)
		if err != nil {
			return nil, err
		}
		allWorkbookRevisions = append(allWorkbookRevisions, workbookRevisionsListResponse.WorkbookRevisionsResponse.WorkbookRevisions...)
	}
	for idx := range allWorkbookRevisions {
		allWorkbookRevisions[idx].WorkbookID = workbookID
	}
	return allWorkbookRevisions, nil
}
