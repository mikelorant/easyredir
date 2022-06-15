package option

type Metadata struct {
	HasMore bool `json:"has_more,omitempty"`
}

type Links struct {
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

type Pagination struct {
	StartingAfter string
	EndingBefore  string
}
