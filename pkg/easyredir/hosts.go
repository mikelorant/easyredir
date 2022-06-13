package easyredir

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mikelorant/easyredir-cli/pkg/structutil"
)

type Hosts struct {
	Data     []HostData `json:"data"`
	Metadata Metadata   `json:"meta"`
	Links    Links      `json:"links"`
}

type HostData struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes HostAttributes `json:"attributes"`
	Links      HostLinks      `json:"links"`
}

type HostAttributes struct {
	Name              string `json:"name"`
	DNSStatus         string `json:"dns_status"`
	CertificateStatus string `json:"certificate_status"`
}

type HostLinks struct {
	Self string `json:"self"`
}

type HostsOptions struct {
	limit      int
	pagination Pagination
}

func WithHostsLimit(limit int) func(*HostsOptions) {
	return func(o *HostsOptions) {
		o.limit = limit
	}
}

func (e *Easyredir) ListHostsPaginator(opts ...func(*HostsOptions)) (h Hosts, err error) {
	h = Hosts{
		Data: []HostData{},
	}

	hosts := Hosts{}
	for {
		optsWithPage := opts
		if hosts.HasMore() {
			optsWithPage = append(optsWithPage, hosts.NextPage())
		}

		hosts, err = e.ListHosts(optsWithPage...)
		if err != nil {
			return h, fmt.Errorf("unable to get a hosts page: %w", err)
		}
		h.Data = append(h.Data, hosts.Data...)
		if !hosts.HasMore() {
			break
		}
	}

	return h, nil
}

func (h *Hosts) NextPage() func(o *HostsOptions) {
	return func(o *HostsOptions) {
		o.pagination.startingAfter = strings.Split(h.Links.Next, "=")[1]
	}
}

func (h *Hosts) HasMore() bool {
	return h.Metadata.HasMore
}

func (e *Easyredir) ListHosts(opts ...func(*HostsOptions)) (h Hosts, err error) {
	options := &HostsOptions{}
	for _, o := range opts {
		o(options)
	}

	pathQuery := buildListHosts(options)
	reader, err := e.client.sendRequest(e.config.baseURL, pathQuery, http.MethodGet, nil)
	if err != nil {
		return h, fmt.Errorf("unable to send request: %w", err)
	}

	if err := decodeJSON(reader, &h); err != nil {
		return h, fmt.Errorf("unable to get json: %w", err)
	}

	return h, nil
}

func (h HostData) String() string {
	str, _ := structutil.Sprint(h)

	return str
}

func (h Hosts) String() string {
	ss := []string{}
	i := 0
	for _, v := range h.Data {
		ss = append(ss, fmt.Sprint(v))
		i++
	}
	ss = append(ss, fmt.Sprintf("Total: %v\n", i))
	return strings.Join(ss, "\n")
}

func buildListHosts(opts *HostsOptions) string {
	var sb strings.Builder
	var params []string

	fmt.Fprint(&sb, "/hosts")

	if opts.pagination.startingAfter != "" {
		params = append(params, fmt.Sprintf("starting_after=%v", opts.pagination.startingAfter))
	}

	if opts.pagination.endingBefore != "" {
		params = append(params, fmt.Sprintf("ending_before=%v", opts.pagination.endingBefore))
	}

	if opts.limit != 0 {
		params = append(params, fmt.Sprintf("limit=%v", opts.limit))
	}

	if len(params) != 0 {
		fmt.Fprintf(&sb, "?%v", strings.Join(params, "&"))
	}

	return sb.String()
}
