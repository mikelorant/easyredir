package rule

import (
	"fmt"
	"strings"

	"github.com/mikelorant/easyredir-cli/pkg/structutil"
)

type Rules struct {
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
	ResponseMovedPermanently ResponseType = "moved permanently"
	ResponseFound            ResponseType = "found"
)

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
