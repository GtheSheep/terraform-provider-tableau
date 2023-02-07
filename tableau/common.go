package tableau

import (
	"math"
	"strconv"
)

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
