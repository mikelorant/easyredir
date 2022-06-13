package rule

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir"
	"github.com/mikelorant/easyredir-cli/pkg/structutil"
)

type Rules struct {
	options  Options
	Data     []Data   `json:"data"`
	Metadata Metadata `json:"meta"`
	Links    Links    `json:"links"`
}

type Data struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Attributes    Attributes    `json:"attributes,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"`
}

type Attributes struct {
	ForwardParams *bool    `json:"forward_params,omitempty"`
	ForwardPath   *bool    `json:"forward_path,omitempty"`
	ResponseType  *string  `json:"response_type,omitempty"`
	SourceURLs    []string `json:"source_urls,omitempty"`
	TargetURL     *string  `json:"target_url,omitempty"`
}

type ResponseType string

const (
	ResponseMovedPermanently ResponseType = "moved permanently"
	ResponseFound            ResponseType = "found"
)

type Relationships struct {
	SourceHosts SourceHosts `json:"source_hosts,omitempty"`
}

type SourceHosts struct {
	Data  []SourceHostData `json:"data,omitempty"`
	Links SourceHostsLinks `json:"links,omitempty"`
}

type SourceHostData struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type SourceHostsLinks struct {
	Related string `json:"related,omitempty"`
}

type Metadata struct {
	HasMore bool `json:"has_more,omitempty"`
}

type Links struct {
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

type Pagination struct {
	startingAfter string
	endingBefore  string
}

type Options struct {
	sourceFilter string
	targetFilter string
	limit        int
	pagination   Pagination
}

func WithSourceFilter(url string) func(*Options) {
	return func(o *Options) {
		o.sourceFilter = url
	}
}

func WithTargetFilter(url string) func(*Options) {
	return func(o *Options) {
		o.targetFilter = url
	}
}

func WithLimit(limit int) func(*Options) {
	return func(o *Options) {
		o.limit = limit
	}
}

func ListRulesPaginator(e *easyredir.Easyredir, opts ...func(*Options)) (r Rules, err error) {
	r = Rules{
		Data: []Data{},
	}

	rules := Rules{}
	for {
		optsWithPage := opts
		if rules.HasMore() {
			optsWithPage = append(optsWithPage, rules.NextPage())
		}

		rules, err = ListRules(e, optsWithPage...)
		if err != nil {
			return r, fmt.Errorf("unable to get a rules page: %w", err)
		}
		r.Data = append(r.Data, rules.Data...)
		if !rules.HasMore() {
			break
		}
	}

	return r, nil
}

func ListRules(e *easyredir.Easyredir, opts ...func(*Options)) (r Rules, err error) {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}

	pathQuery := buildListRules(options)
	reader, err := e.Client.SendRequest(e.Config.BaseURL, pathQuery, http.MethodGet, nil)
	if err != nil {
		return r, fmt.Errorf("unable to send request: %w", err)
	}

	if err := decodeJSON(reader, &r); err != nil {
		return r, fmt.Errorf("unable to get json: %w", err)
	}

	return r, nil
}

func (r *Rules) NextPage() func(o *Options) {
	return func(o *Options) {
		o.pagination.startingAfter = strings.Split(r.Links.Next, "=")[1]
	}
}

func (r *Rules) HasMore() bool {
	return r.Metadata.HasMore
}

func (r Data) String() string {
	str, _ := structutil.Sprint(r)

	return str
}

func (r Rules) String() string {
	ss := []string{}
	i := 0
	for _, v := range r.Data {
		ss = append(ss, fmt.Sprint(v))
		i++
	}
	ss = append(ss, fmt.Sprintf("Total: %v\n", i))
	return strings.Join(ss, "\n")
}

func buildListRules(opts *Options) string {
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

func decodeJSON(r io.ReadCloser, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return fmt.Errorf("unable to json decode: %w", err)
	}
	r.Close()

	return nil
}
