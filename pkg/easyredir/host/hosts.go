package host

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir/option"
)

type ClientAPI interface {
	SendRequest(baseURL, path, method string, body io.Reader) (io.ReadCloser, error)
}

type ConfigAPI interface {
	BaseURL() string
}

func ListHostsPaginator(cl ClientAPI, cfg ConfigAPI, opts ...func(*option.Options)) (h Hosts, err error) {
	h = Hosts{
		Data: []Data{},
	}

	hosts := Hosts{}
	for {
		optsWithPage := opts
		if hosts.HasMore() {
			optsWithPage = append(optsWithPage, hosts.NextPage())
		}

		hosts, err = ListHosts(cl, cfg, optsWithPage...)
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

func (h *Hosts) NextPage() func(o *option.Options) {
	return func(o *option.Options) {
		o.Pagination.StartingAfter = strings.Split(h.Links.Next, "=")[1]
	}
}

func (h *Hosts) HasMore() bool {
	return h.Metadata.HasMore
}

func ListHosts(cl ClientAPI, cfg ConfigAPI, opts ...func(*option.Options)) (h Hosts, err error) {
	options := &option.Options{}
	for _, o := range opts {
		o(options)
	}

	pathQuery := buildListHosts(options)
	reader, err := cl.SendRequest(cfg.BaseURL(), pathQuery, http.MethodGet, nil)
	if err != nil {
		return h, fmt.Errorf("unable to send request: %w", err)
	}

	if err := decodeJSON(reader, &h); err != nil {
		return h, fmt.Errorf("unable to get json: %w", err)
	}

	return h, nil
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

func GetHost(cl ClientAPI, cfg ConfigAPI, id string) (h Host, err error) {
	pathQuery := buildGetHost(id)
	reader, err := cl.SendRequest(cfg.BaseURL(), pathQuery, http.MethodGet, nil)
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
