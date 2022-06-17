package rule

import (
	"fmt"
	"io"
	"strings"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir/host"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/option"
	"github.com/mikelorant/easyredir-cli/pkg/structutil"
)

type ClientAPI interface {
	SendRequest(path, method string, body io.Reader) (io.ReadCloser, error)
}

type Rules struct {
	Data     []Data          `json:"data"`
	Metadata option.Metadata `json:"meta"`
	Links    option.Links    `json:"links"`
}

type Rule struct {
	Data          Data
	Relationships Relationships `json:"relationships,omitempty"` // API docs are incorrect
	Included	  []host.Data	`json:"included,omitempty"`
}

type Data struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Attributes    Attributes    `json:"attributes,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"` // API docs are incorrect
}

type Attributes struct {
	ForwardParams *bool    `json:"forward_params,omitempty"`
	ForwardPath   *bool    `json:"forward_path,omitempty"`
	ResponseType  *ResponseType  `json:"response_type,omitempty"`
	SourceURLs    []string `json:"source_urls,omitempty"`
	TargetURL     *string  `json:"target_url,omitempty"`
}

type ResponseType string

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

const (
	ResponseMovedPermanently ResponseType = "moved_permanently"
	ResponseFound            ResponseType = "found"
)

func (r Rule) String() string {
	str, _ := structutil.Sprint(r)

	return str
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

func ref[T any](x T) *T {
    return &x
}
