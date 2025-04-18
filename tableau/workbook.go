package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Workbook struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	DefaultViewID   string `json:"defaultViewId,omitempty"`
	EncryptExtracts string `json:"encryptExtracts,omitempty"`
	ShowTabs        string `json:"showTabs,omitempty"`
	Size            string `json:"size,omitempty"` // size in megabytes
	ContentURL      string `json:"contentUrl,omitempty"`
	WebPageURL      string `json:"webpageUrl,omitempty"`
	Location        struct {
		ID string `json:"id,omitempty"`
		// Type string `json:"type,omitempty"` // for example "Project"
		// Name string `json:"name,omitempty"`
	} `json:"location,omitempty"`
	Owner struct {
		ID string `json:"id,omitempty"`
		// Name string `json:"name,omitempty"`
	} `json:"owner,omitempty"`
	Project struct {
		ID string `json:"id,omitempty"`
		// Name string `json:"name,omitempty"`
	} `json:"project,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	// Tags
}

type WorkbookRequest struct {
	Workbook Workbook `json:"workbook"`
}

type WorkbooksResponse struct {
	Workbooks []Workbook `json:"workbook"`
}

type WorkbookListResponse struct {
	WorkbooksResponse WorkbooksResponse `json:"workbooks"`
	Pagination        PaginationDetails `json:"pagination"`
}

func (c *Client) GetWorkbooks() ([]Workbook, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/workbooks", c.ApiUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	workbookListResponse := WorkbookListResponse{}
	err = json.Unmarshal(body, &workbookListResponse)
	if err != nil {
		return nil, err
	}

	pageNumber, totalPageCount, totalAvailable, err := GetPaginationNumbers(workbookListResponse.Pagination)
	if err != nil {
		return nil, err
	}

	allWorkbooks := make([]Workbook, 0, totalAvailable)
	allWorkbooks = append(allWorkbooks, workbookListResponse.WorkbooksResponse.Workbooks...)

	for page := pageNumber + 1; page <= totalPageCount; page++ {
		fmt.Printf("Searching page %d", page)
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/workbooks?pageNumber=%d", c.ApiUrl, page), nil)
		if err != nil {
			return nil, err
		}
		body, err = c.doRequest(req)
		if err != nil {
			return nil, err
		}
		workbookListResponse = WorkbookListResponse{}
		err = json.Unmarshal(body, &workbookListResponse)
		if err != nil {
			return nil, err
		}
		allWorkbooks = append(allWorkbooks, workbookListResponse.WorkbooksResponse.Workbooks...)
	}

	return allWorkbooks, nil
}
