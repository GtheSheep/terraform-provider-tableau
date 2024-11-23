package tableau

import (
	"math"
	"strconv"
	"strings"
)

type Owner struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type PaginationDetails struct {
	PageNumber     string `json:"pageNumber"`
	PageSize       string `json:"pageSize"`
	TotalAvailable string `json:"totalAvailable"`
}

func GetPaginationNumbers(paginationDetails PaginationDetails) (int, int, error) {
	pageNumber, err := strconv.Atoi(paginationDetails.PageNumber)
	if err != nil {
		return 0, 0, err
	}
	pageSize, err := strconv.Atoi(paginationDetails.PageSize)
	if err != nil {
		return 0, 0, err
	}
	totalAvailable, err := strconv.Atoi(paginationDetails.TotalAvailable)
	if err != nil {
		return 0, 0, err
	}
	totalPageCount := int(math.Ceil(float64(totalAvailable) / float64(pageSize)))

	return pageNumber, totalPageCount, nil
}

func GetCombinedID(id1, id2 string) string {
	combined := strings.Join([]string{id1, id2}, ":")
	return combined
}

func GetIDsFromCombinedID(id string) (string, string) {
	split := strings.Split(id, ":")
	return split[0], split[1]
}
