package tableau

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

type WorkbookOwner struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type WorkbookProject struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type XMLWorkbookOwner struct {
	XMLName xml.Name `xml:"owner"`
	ID      string   `xml:"id,attr"`
}

type XMLWorkbookProject struct {
	XMLName xml.Name `xml:"project"`
	ID      string   `xml:"id,attr"`
}

type XMLWorkbook struct {
	XMLName          xml.Name           `xml:"workbook"`
	Name             string             `xml:"name,attr"`
	Description      string             `xml:"description,attr"`
	Project          XMLWorkbookProject `xml:"project"`
	ShowTabs         string             `xml:"showTabs,attr"`
	ThumbnailsUserID string             `xml:"thumbnailsUserId,attr"`
	EncryptExtracts  string             `xml:"encryptExtracts,attr"`
	Owner            XMLWorkbookOwner   `xml:"owner"`
}

type XMLTsRequest struct {
	XMLName  xml.Name    `xml:"tsRequest"`
	Workbook XMLWorkbook `xml:"workbook"`
}

type Workbook struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	EncryptExtracts string `json:"encryptExtracts,omitempty"`
	ShowTabs        string `json:"showTabs,omitempty"`
	Size            string `json:"size,omitempty"` // size in megabytes
	ContentURL      string `json:"contentUrl,omitempty"`
	WebPageURL      string `json:"webpageUrl,omitempty"`
	Location        struct {
		ID   string `json:"id,omitempty"`
		Type string `json:"type,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"location,omitempty"`
	Project   WorkbookProject `json:"project,omitempty"`
	Owner     WorkbookOwner   `json:"owner,omitempty"`
	CreatedAt string          `json:"createdAt,omitempty"`
	UpdatedAt string          `json:"updatedAt,omitempty"`
	// Tags
}

type WorkbookRequest struct {
	Workbook Workbook `json:"workbook"`
}

type WorkbooksResponse struct {
	Workbooks []Workbook `json:"workbook"`
}

type WorkbookResponse struct {
	Workbook Workbook `json:"workbook"`
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

func (c *Client) GetWorkbook(id string) (*Workbook, error) {
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
	for idx := range allWorkbooks {
		if allWorkbooks[idx].ID == id {
			return &allWorkbooks[idx], nil
		}
	}
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
		for idx := range allWorkbooks {
			if allWorkbooks[idx].ID == id {
				return &allWorkbooks[idx], nil
			}
		}
	}
	return nil, fmt.Errorf("failed to find workbook with ID %s", id)
}

func (c *Client) CreateWorkbook(ctx context.Context, name, projectID, showTabs, thumbnailsUserID, wbFilename string, wbContent []byte) (string, error) {

	workbookRequest := XMLTsRequest{
		Workbook: XMLWorkbook{
			Name:             name,
			Project:          XMLWorkbookProject{ID: projectID},
			ShowTabs:         showTabs,
			ThumbnailsUserID: thumbnailsUserID,
		},
	}
	newWorkbookXml, err := xml.Marshal(workbookRequest)
	if err != nil {
		return "", err
	}
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	part, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": {"name=\"request_payload\""},
		"Content-Type":        {"application/xml"},
	})
	if err != nil {
		return "", err
	}
	if _, err = part.Write(newWorkbookXml); err != nil {
		return "", err
	}
	part, err = writer.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": {fmt.Sprintf("name=\"tableau_workbook\"; filename=\"%s\"", wbFilename)},
		"Content-Type":        {"application/octet-stream"},
	})
	if err != nil {
		return "", err
	}
	if _, err = part.Write(wbContent); err != nil {
		return "", err
	}
	if err = writer.Close(); err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/workbooks", c.ApiUrl), reqBody)
	if err != nil {
		return "", err
	}
	respBody, err := c.doRequest(req, func(o *requestOptions) {
		o.contentType = "multipart/mixed; boundary=" + writer.Boundary()
	})
	if err != nil {
		return "", err
	}
	workbookResponse := WorkbookResponse{}
	if err = json.Unmarshal(respBody, &workbookResponse); err != nil {
		return "", err
	}
	return workbookResponse.Workbook.ID, nil
}

func (c *Client) UpdateWorkbook(id, name, projectID, showTabs, description, encryptExtracts, ownerID string) (*Workbook, error) {

	workbookRequest := WorkbookRequest{
		Workbook: Workbook{
			ID:              id,
			Name:            name,
			Project:         WorkbookProject{ID: projectID},
			ShowTabs:        showTabs,
			Description:     description,
			EncryptExtracts: encryptExtracts,
			Owner:           WorkbookOwner{ID: ownerID},
		},
	}

	newWorkbookJson, err := json.Marshal(workbookRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/workbooks/%s", c.ApiUrl, id), strings.NewReader(string(newWorkbookJson)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	workbookResponse := WorkbookResponse{}
	err = json.Unmarshal(body, &workbookResponse)
	if err != nil {
		return nil, err
	}

	return &workbookResponse.Workbook, nil
}

func (c *Client) DeleteWorkbook(id string) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/workbooks/%s", c.ApiUrl, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
