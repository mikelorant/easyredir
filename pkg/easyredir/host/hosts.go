package host

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir"
)

type Options struct {
	limit      int
	pagination Pagination
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

func WithHostsLimit(limit int) func(*Options) {
	return func(o *Options) {
		o.limit = limit
	}
}

func ListHostsPaginator(e *easyredir.Easyredir, opts ...func(*Options)) (h Hosts, err error) {
	h = Hosts{
		Data: []Data{},
	}

	hosts := Hosts{}
	for {
		optsWithPage := opts
		if hosts.HasMore() {
			optsWithPage = append(optsWithPage, hosts.NextPage())
		}

		hosts, err = ListHosts(e, optsWithPage...)
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

func (h *Hosts) NextPage() func(o *Options) {
	return func(o *Options) {
		o.pagination.startingAfter = strings.Split(h.Links.Next, "=")[1]
	}
}

func (h *Hosts) HasMore() bool {
	return h.Metadata.HasMore
}

func ListHosts(e *easyredir.Easyredir, opts ...func(*Options)) (h Hosts, err error) {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}

	pathQuery := buildListHosts(options)
	reader, err := e.Client.SendRequest(e.Config.BaseURL, pathQuery, http.MethodGet, nil)
	if err != nil {
		return h, fmt.Errorf("unable to send request: %w", err)
	}

	if err := decodeJSON(reader, &h); err != nil {
		return h, fmt.Errorf("unable to get json: %w", err)
	}

	return h, nil
}

func buildListHosts(opts *Options) string {
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

func GetHost(e *easyredir.Easyredir, id string) (h Host, err error) {
	pathQuery := buildGetHost(id)
	reader, err := e.Client.SendRequest(e.Config.BaseURL, pathQuery, http.MethodGet, nil)
	if err != nil {
		return h, fmt.Errorf("unable to send request: %w", err)
	}

	if err := decodeJSON(reader, &h); err != nil {
		return h, fmt.Errorf("unable to get json: %w", err)
	}

	if ok := (h.Data.ID == id); !ok {
		return h, fmt.Errorf("received incorrect host: %v", h.Data.ID)
	}

	return h, nil
}

func buildGetHost(id string) string {
	return fmt.Sprintf("/hosts/%v", id)
}

func decodeJSON(r io.ReadCloser, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return fmt.Errorf("unable to json decode: %w", err)
	}
	r.Close()

	return nil
}
