package pagination

const MaxPageSize = 100

// Clamp bounds list page size to sensible defaults.
func Clamp(pageSize, defaultSize int) int {
	if pageSize < 1 {
		return defaultSize
	}
	if pageSize > MaxPageSize {
		return MaxPageSize
	}
	return pageSize
}
