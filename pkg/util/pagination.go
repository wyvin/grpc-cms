package util

// GetPage get page parameters
func FormatPageAndPerPage(page, perPage *uint32) {
	if *perPage <= 0 {
		*perPage = 10
	}
	if *page > 0 {
		*page = (*page - 1) * *perPage
	}
}
