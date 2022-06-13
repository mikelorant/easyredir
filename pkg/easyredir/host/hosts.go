package host

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir"
	"github.com/mikelorant/easyredir-cli/pkg/structutil"
)

type Hosts struct {
	Data     []Data   `json:"data"`
	Metadata Metadata `json:"meta"`
	Links    Links    `json:"links"`
}

type Data struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
	Links      HostLinks  `json:"links"`
}

type Attributes struct {
	Name              string            `json:"name"`
	DNSStatus         DNSStatus         `json:"dns_status"`
	CertificateStatus CertificateStatus `json:"certificate_status"`
}

type HostLinks struct {
	Self string `json:"self"`
}

type Options struct {
	limit      int
	pagination Pagination
}

type Host struct {
	Data DataExtended
}

type DataExtended struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes AttributesExtended `json:"attributes"`
	Links      HostLinks          `json:"links"`
}

type AttributesExtended struct {
	Name               string             `json:"name"`
	DNSStatus          DNSStatus          `json:"dns_status"`
	DNSTestedAt        string             `json:"dns_tested_at"` // TODO: time.Time
	CertificateStatus  string             `json:"certificate_status"`
	ACMEEnabled        bool               `json:"acme_enabled"`
	MatchOptions       MatchOptions       `json:"match_options"`
	NotFoundAction     NotFoundActions    `json:"not_found_action"`
	Security           Security           `json:"security"`
	RequiredDNSEntries RequiredDNSEntries `json:"required_dns_entries"`
	DetectedDNSEntries []DNSValues        `json:"detected_dns_entries"`
}

type MatchOptions struct {
	CaseInsensitive  *bool `json:"case_insensitive"`
	SlashInsensitive *bool `json:"slash_insensitive"`
}

type NotFoundActions struct {
	ForwardParams        *bool         `json:"forward_params"`
	ForwardPath          *bool         `json:"forward_path"`
	Custom404Body        *string       `json:"custom_404_body"`
	Custom404BodyPresent bool          `json:"custom_404_body_present"`
	ResponseCode         *ResponseCode `json:"response_code"`
	ResponseURL          *string       `json:"response_url"`
}

type Security struct {
	HTTPSUpgrade            *bool `json:"https_upgrade"`
	PreventForeignEmbedding *bool `json:"prevent_foreign_embedding"`
	HSTSIncludeSubDomains   *bool `json:"hsts_include_sub_domains"`
	HSTSMaxAge              *int  `json:"hsts_max_age"`
	HSTSPreload             *bool `json:"hsts_preload"`
}

type RequiredDNSEntries struct {
	Recommended  DNSValues   `json:"recommended"`
	Alternatives []DNSValues `json:"alternatives"`
}

type DNSValues struct {
	Type   string          `json:"type"`
	Values []DNSValuesType `json:"values"`
}

type DNSValuesType string

const (
	DNSARecord     DNSValuesType = "A"
	DNSCNAMERecord DNSValuesType = "CNAME"
)

type ResponseCode int

const (
	ResponseCodeMovedPermanently ResponseCode = 301
	ResponseCodeFound            ResponseCode = 302
	ResponseCodeNotFound         ResponseCode = 401
)

type DNSStatus string

const (
	DNSStatusActive  DNSStatus = "active"
	DNSStatusInvalid DNSStatus = "invalid"
)

type CertificateStatus string

const (
	CertificateStatusActive                     CertificateStatus = "active"
	CertificateStatusProcessing                 CertificateStatus = "processing"
	CertificateStatusInvalidDNS                 CertificateStatus = "invalid_dns"
	CertificateStatusAutoSSLNotSupported        CertificateStatus = "auto_ssl_not_supported"
	CertificateStatusHostnameContainsUnderscore CertificateStatus = "hostname_contains_underscore"
	CertificateStatusInvalidCAARecord           CertificateStatus = "invalid_caa_record"
	CertificateStatusAAAARecordPresent          CertificateStatus = "aaaa_record_present"
)

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

func (h Data) String() string {
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

func (h DataExtended) String() string {
	str, _ := structutil.Sprint(h)

	return str
}

func (h Host) String() string {
	return fmt.Sprint(h.Data)
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
