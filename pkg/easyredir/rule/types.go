package rule

import (
	"fmt"
	"io"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mikelorant/easyredir/pkg/easyredir/host"
	"github.com/mikelorant/easyredir/pkg/easyredir/option"
	"github.com/mikelorant/easyredir/pkg/structutil"
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
	Included      []host.Data   `json:"included,omitempty"`
}

type Data struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Attributes    Attributes    `json:"attributes,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"` // API docs are incorrect
}

type Attributes struct {
	ForwardParams *bool         `json:"forward_params,omitempty"`
	ForwardPath   *bool         `json:"forward_path,omitempty"`
	ResponseType  *ResponseType `json:"response_type,omitempty"`
	SourceURLs    []string      `json:"source_urls,omitempty"`
	TargetURL     *string       `json:"target_url,omitempty"`
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

const (
	ErrNoRules = "No rules."
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
	switch len(r.Data) {
	case 0:
		return ErrNoRules
	case 1:
		return fmt.Sprint(r.Data[0])
	default:
		t := table.NewWriter()
		t.SetStyle(table.StyleColoredBright)
		t.Style().Options.DrawBorder = false
		t.Style().Color = table.ColorOptions{}
		t.Style().Box.PaddingLeft = ""
		t.Style().Box.PaddingRight = "\t"
		t.Style().Color.Header = text.Colors{text.Bold}
		t.AppendHeader(table.Row{"ID", "SOURCE URLS", "TARGET URL"})
		for _, d := range r.Data {
			t.AppendRow(table.Row{
				d.ID,
				strings.Join(d.Attributes.SourceURLs, "\n"),
				*d.Attributes.TargetURL,
			})
		}
		return t.Render()
	}
}

func ref[T any](x T) *T {
	return &x
}
