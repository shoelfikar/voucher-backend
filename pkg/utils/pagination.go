package utils

import (
	"math"
	"strconv"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page      int
	Limit     int
	Offset    int
	SortBy    string
	SortOrder string
}

// ParsePaginationParams parses pagination parameters from query strings
func ParsePaginationParams(pageStr, limitStr, sortBy, sortOrder string) PaginationParams {
	page := parsePage(pageStr)
	limit := parseLimit(limitStr)
	offset := calculateOffset(page, limit)

	// Default sort order
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// Validate sort order
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	return PaginationParams{
		Page:      page,
		Limit:     limit,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}

// parsePage parses the page parameter with default value
func parsePage(pageStr string) int {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return DefaultPage
	}
	return page
}

// parseLimit parses the limit parameter with default and max values
func parseLimit(limitStr string) int {
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return DefaultLimit
	}
	if limit > MaxLimit {
		return MaxLimit
	}
	return limit
}

// calculateOffset calculates the offset based on page and limit
func calculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// CalculateTotalPages calculates the total number of pages
func CalculateTotalPages(total int64, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}
