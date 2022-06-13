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
	Name              string                `json:"name"`
	DNSStatus         HostDNSStatus         `json:"dns_status"`
	CertificateStatus HostCertificateStatus `json:"certificate_status"`
}

type HostLinks struct {
	Self string `json:"self"`
}

type HostsOptions struct {
	limit      int
	pagination Pagination
}

type Host struct {
	Data HostDataExtended
}

type HostDataExtended struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes HostAttributesExtended `json:"attributes"`
	Links      HostLinks              `json:"links"`
}

type HostAttributesExtended struct {
	Name               string                 `json:"name"`
	DNSStatus          HostDNSStatus          `json:"dns_status"`
	DNSTestedAt        string                 `json:"dns_tested_at"` // TODO: time.Time
	CertificateStatus  string                 `json:"certificate_status"`
	ACMEEnabled        bool                   `json:"acme_enabled"`
	MatchOptions       HostMatchOptions       `json:"match_options"`
	NotFoundAction     HostNotFoundActions    `json:"not_found_action"`
	Security           HostSecurity           `json:"security"`
	RequiredDNSEntries HostRequiredDNSEntries `json:"required_dns_entries"`
	DetectedDNSEntries []DNSValues            `json:"detected_dns_entries"`
}

type HostMatchOptions struct {
	CaseInsensitive  *bool `json:"case_insensitive"`
	SlashInsensitive *bool `json:"slash_insensitive"`
}

type HostNotFoundActions struct {
	ForwardParams        *bool             `json:"forward_params"`
	ForwardPath          *bool             `json:"forward_path"`
	Custom404Body        *string           `json:"custom_404_body"`
	Custom404BodyPresent bool              `json:"custom_404_body_present"`
	ResponseCode         *HostResponseCode `json:"response_code"`
	ResponseURL          *string           `json:"response_url"`
}

type HostSecurity struct {
	HTTPSUpgrade            *bool `json:"https_upgrade"`
	PreventForeignEmbedding *bool `json:"prevent_foreign_embedding"`
	HSTSIncludeSubDomains   *bool `json:"hsts_include_sub_domains"`
	HSTSMaxAge              *int  `json:"hsts_max_age"`
	HSTSPreload             *bool `json:"hsts_preload"`
}

type HostRequiredDNSEntries struct {
	Recommended  DNSValues   `json:"recommended"`
	Alternatives []DNSValues `json:"alternatives"`
}

type DNSValues struct {
	Type   string          `json:"type"`
	Values []DNSValuesType `json:"values"`
}

type DNSValuesType string

const (
	HostDNSARecord     DNSValuesType = "A"
	HostDNSCNAMERecord DNSValuesType = "CNAME"
)

type HostResponseCode int

const (
	HostResponseMovedPermanently HostResponseCode = 301
	HostResponseFound            HostResponseCode = 302
	HostResponseNotFound         HostResponseCode = 401
)

type HostDNSStatus string

const (
	HostDNSStatusActive  HostDNSStatus = "active"
	HostDNSStatusInvalid HostDNSStatus = "invalid"
)

type HostCertificateStatus string

const (
	HostCertificateStatusActive                     HostCertificateStatus = "active"
	HostCertificateStatusProcessing                 HostCertificateStatus = "processing"
	HostCertificateStatusInvalidDNS                 HostCertificateStatus = "invalid_dns"
	HostCertificateStatusAutoSSLNotSupported        HostCertificateStatus = "auto_ssl_not_supported"
	HostCertificateStatusHostnameContainsUnderscore HostCertificateStatus = "hostname_contains_underscore"
	HostCertificateStatusInvalidCAARecord           HostCertificateStatus = "invalid_caa_record"
	HostCertificateStatusAAAARecordPresent          HostCertificateStatus = "aaaa_record_present"
)

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

func (h HostDataExtended) String() string {
	str, _ := structutil.Sprint(h)

	return str
}

func (h Host) String() string {
	return fmt.Sprint(h.Data)
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

func (e *Easyredir) GetHost(id string) (h Host, err error) {
	pathQuery := buildGetHost(id)
	reader, err := e.client.sendRequest(e.config.baseURL, pathQuery, http.MethodGet, nil)
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
