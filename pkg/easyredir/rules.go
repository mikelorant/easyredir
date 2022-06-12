package easyredir

import (
	"fmt"
	"net/http"
	"strings"
)

type Rules struct {
	options  RulesOptions
	Data     []RuleData `json:"data"`
	Metadata Metadata   `json:"meta"`
	Links    Links      `json:"links"`
}

type RuleData struct {
	ID         *string        `json:"id"`
	Type       *string        `json:"type"`
	Attributes RuleAttributes `json:"attributes"`
}

type RuleAttributes struct {
	ForwardParams *bool     `json:"forward_params"`
	ForwardPath   *bool     `json:"forward_path"`
	ResponseType  *string   `json:"response_type"`
	SourceURLs    []*string `json:"source_urls"`
	TargetURL     *string   `json:"target_url"`
}

type Metadata struct {
	HasMore bool `json:"has_more"`
}

type Links struct {
	Next string `json:"next"`
	Prev string `json:"prev"`
}

type Pagination struct {
	startingAfter string
	endingBefore  string
}

type RulesOptions struct {
	sourceFilter string
	targetFilter string
	limit        int
	pagination   Pagination
}

func WithSourceFilter(url string) func(*RulesOptions) {
	return func(o *RulesOptions) {
		o.sourceFilter = url
	}
}

func WithTargetFilter(url string) func(*RulesOptions) {
	return func(o *RulesOptions) {
		o.targetFilter = url
	}
}

func WithLimit(limit int) func(*RulesOptions) {
	return func(o *RulesOptions) {
		o.limit = limit
	}
}

func (e *Easyredir) ListRules(opts ...func(*RulesOptions)) (r Rules, err error) {
	options := &RulesOptions{}
	for _, o := range opts {
		o(options)
	}

	pathQuery := listRulesPathQuery(options)
	reader, err := e.client.sendRequest(e.config.baseURL, pathQuery, http.MethodGet, nil)
	if err != nil {
		return r, fmt.Errorf("unable to send request: %w", err)
	}

	if err := decodeJSON(reader, &r); err != nil {
		return r, fmt.Errorf("unable to get json: %w", err)
	}

	return r, nil
}

func listRulesPathQuery(opts *RulesOptions) string {
	var sb strings.Builder
	var params []string

	fmt.Fprint(&sb, "/rules")

	if opts.pagination.startingAfter != "" {
		params = append(params, fmt.Sprintf("starting_after=%v", opts.pagination.startingAfter))
	}

	if opts.pagination.endingBefore != "" {
		params = append(params, fmt.Sprintf("ending_before=%v", opts.pagination.endingBefore))
	}

	if opts.sourceFilter != "" {
		params = append(params, fmt.Sprintf("sq=%v", opts.sourceFilter))
	}

	if opts.targetFilter != "" {
		params = append(params, fmt.Sprintf("tq=%v", opts.targetFilter))
	}

	if opts.limit != 0 {
		params = append(params, fmt.Sprintf("limit=%v", opts.limit))
	}

	if len(params) != 0 {
		fmt.Fprintf(&sb, "?%v", strings.Join(params, "&"))
	}

	return sb.String()
}
