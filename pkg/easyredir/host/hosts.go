package host

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir/option"
	"github.com/mikelorant/easyredir-cli/pkg/jsonutil"
)

func ListHostsPaginator(cl ClientAPI, opts ...Option) (h Hosts, err error) {
	h = Hosts{
		Data: []Data{},
	}

	hosts := Hosts{}
	for {
		optsWithPage := opts
		if hosts.HasMore() {
			optsWithPage = append(optsWithPage, hosts.NextPage())
		}

		hosts, err = ListHosts(cl, optsWithPage...)
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

func ListHosts(cl ClientAPI, opts ...Option) (h Hosts, err error) {
	o := &option.Options{}
	for _, opt := range opts {
		opt.Apply(o)
	}

	pathQuery := buildListHosts(o)
	reader, err := cl.SendRequest(pathQuery, http.MethodGet, nil)
	if err != nil {
		return h, fmt.Errorf("unable to send request: %w", err)
	}

	if err := jsonutil.DecodeJSON(reader, &h); err != nil {
		return h, fmt.Errorf("unable to get json: %w", err)
	}

	return h, nil
}

func GetHost(cl ClientAPI, id string) (h Host, err error) {
	pathQuery := buildGetHost(id)
	reader, err := cl.SendRequest(pathQuery, http.MethodGet, nil)
	if err != nil {
		return h, fmt.Errorf("unable to send request: %w", err)
	}

	if err := jsonutil.DecodeJSON(reader, &h); err != nil {
		return h, fmt.Errorf("unable to get json: %w", err)
	}

	if ok := (h.Data.ID == id); !ok {
		return h, fmt.Errorf("received incorrect host: %v", h.Data.ID)
	}

	return h, nil
}

func (h *Hosts) NextPage() NextPage {
	return NextPage(strings.Split((h.Links.Next), "=")[1])
}

type NextPage string

func (np NextPage) Apply(o *option.Options) {
	o.Pagination.StartingAfter = string(np)
}

func (h *Hosts) HasMore() bool {
	return h.Metadata.HasMore
}

func buildListHosts(opts *option.Options) string {
	var sb strings.Builder
	var params []string

	fmt.Fprint(&sb, "/hosts")

	if opts.Pagination.StartingAfter != "" {
		params = append(params, fmt.Sprintf("starting_after=%v", opts.Pagination.StartingAfter))
	}

	if opts.Pagination.EndingBefore != "" {
		params = append(params, fmt.Sprintf("ending_before=%v", opts.Pagination.EndingBefore))
	}

	if opts.Limit != 0 {
		params = append(params, fmt.Sprintf("limit=%v", opts.Limit))
	}

	if len(params) != 0 {
		fmt.Fprintf(&sb, "?%v", strings.Join(params, "&"))
	}

	return sb.String()
}

func buildGetHost(id string) string {
	return fmt.Sprintf("/hosts/%v", id)
}
