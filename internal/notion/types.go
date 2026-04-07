package notion

type QueryDatabaseRequest struct {
	Filter      map[string]any   `json:"filter,omitempty"`
	Sorts       []map[string]any `json:"sorts,omitempty"`
	PageSize    int              `json:"page_size,omitempty"`
	StartCursor string           `json:"start_cursor,omitempty"`
}

type QueryDatabaseResponse[T any] struct {
	Results    []T    `json:"results"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor,omitempty"`
}

type UpsertPageRequest[T any] struct {
	PageID     string `json:"-"`
	Properties T      `json:"properties"`
}
