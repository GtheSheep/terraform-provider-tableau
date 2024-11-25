package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Tag struct {
	Label string `json:"label"`
}

type Tags struct {
	Tags []Tag `json:"tag,omitempty"`
}

type Datasource struct {
	ID                  string  `json:"id,omitempty"`
	Name                string  `json:"name,omitempty"`
	Description         string  `json:"description,omitempty"`
	CertificationNote   string  `json:"certificationNote,omitempty"`
	Type                string  `json:"type,omitempty"`
	ContentURL          string  `json:"contentUrl,omitempty"`
	CreatedAt           string  `json:"createdAt,omitempty"`
	UpdatedAt           string  `json:"updatedAt,omitempty"`
	EncryptExtracts     string  `json:"encryptExtracts,omitempty"`
	HasExtracts         bool    `json:"hasExtracts,omitempty"`
	IsCertified         bool    `json:"isCertified,omitempty"`
	UseRemoteQueryAgent bool    `json:"useRemoteQueryAgent,omitempty"`
	WebPageURL          string  `json:"webpageUrl,omitempty"`
	Owner               Owner   `json:"owner,omitempty"`
	Project             Project `json:"project,omitempty"`
	Tags                Tags    `json:"tags,omitempty"`
}

type DatasourceRequest struct {
	Datasource Datasource `json:"datasource"`
}

type DatasourceResponse struct {
	Datasource Datasource `json:"datasource"`
}

type DatasourcesResponse struct {
	Datasources []Datasource `json:"datasource"`
}

type DatasourceListResponse struct {
	DatasourcesResponse DatasourcesResponse `json:"datasources"`
	Pagination          PaginationDetails   `json:"pagination"`
}

func (c *Client) GetDatasource(datasourceID, name string) (*Datasource, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/datasources", c.ApiUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	datasourceListResponse := DatasourceListResponse{}
	err = json.Unmarshal(body, &datasourceListResponse)
	if err != nil {
		return nil, err
	}

	// TODO: Generalise pagination handling and use elsewhere
	pageNumber, totalPageCount, _, err := GetPaginationNumbers(datasourceListResponse.Pagination)
	if err != nil {
		return nil, err
	}
	for i, datasource := range datasourceListResponse.DatasourcesResponse.Datasources {
		if (datasource.ID == datasourceID) || (datasource.Name == name) {
			return &datasourceListResponse.DatasourcesResponse.Datasources[i], nil
		}
	}

	for page := pageNumber + 1; page <= totalPageCount; page++ {
		fmt.Printf("Searching page %d", page)
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/datasources?pageNumber=%s", c.ApiUrl, strconv.Itoa(page)), nil)
		if err != nil {
			return nil, err
		}
		body, err = c.doRequest(req)
		if err != nil {
			return nil, err
		}
		datasourceListResponse = DatasourceListResponse{}
		err = json.Unmarshal(body, &datasourceListResponse)
		if err != nil {
			return nil, err
		}
		// check if we found the datasource in this page
		for i, datasource := range datasourceListResponse.DatasourcesResponse.Datasources {
			if (datasource.ID == datasourceID) || (datasource.Name == name) {
				return &datasourceListResponse.DatasourcesResponse.Datasources[i], nil
			}
		}
	}

	return nil, fmt.Errorf("Did not find datasource ID %s", datasourceID)
}
