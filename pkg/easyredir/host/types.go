package host

import (
	"fmt"
	"io"
	"strings"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir/option"
	"github.com/mikelorant/easyredir-cli/pkg/structutil"
)

type ClientAPI interface {
	SendRequest(path, method string, body io.Reader) (io.ReadCloser, error)
}

type Hosts struct {
	Data     []Data          `json:"data"`
	Metadata option.Metadata `json:"meta,omitempty"`
	Links    option.Links    `json:"links,omitempty"`
}

type Host struct {
	Data Data
}

type Data struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes,omitempty"`
	Links      Links      `json:"links,omitempty"`
}

type Attributes struct {
	Name               string             `json:"name,omitempty"`
	DNSStatus          DNSStatus          `json:"dns_status,omitempty"`
	DNSTestedAt        string             `json:"dns_tested_at,omitempty"` // TODO: time.Time
	CertificateStatus  string             `json:"certificate_status,omitempty"`
	ACMEEnabled        *bool              `json:"acme_enabled,omitempty"`
	MatchOptions       MatchOptions       `json:"match_options,omitempty"`
	NotFoundAction     NotFoundAction     `json:"not_found_action,omitempty"`
	Security           Security           `json:"security,omitempty"`
	RequiredDNSEntries RequiredDNSEntries `json:"required_dns_entries,omitempty"`
	DetectedDNSEntries []DNSValues        `json:"detected_dns_entries,omitempty"`
}

type Links struct {
	Self string `json:"self,omitempty"`
}

type MatchOptions struct {
	CaseInsensitive  *bool `json:"case_insensitive,omitempty"`
	SlashInsensitive *bool `json:"slash_insensitive,omitempty"`
}

type NotFoundAction struct {
	ForwardParams        *bool         `json:"forward_params,omitempty"`
	ForwardPath          *bool         `json:"forward_path,omitempty"`
	Custom404Body        *string       `json:"my_custom_404_body,omitempty"`
	Custom404BodyPresent *bool         `json:"custom_404_body_present,omitempty"` // TODO: Marked as string in example
	ResponseCode         *ResponseCode `json:"response_code,omitempty"`
	ResponseURL          *string       `json:"response_url,omitempty"`
}

type Security struct {
	HTTPSUpgrade            *bool `json:"https_upgrade"`
	PreventForeignEmbedding *bool `json:"prevent_foreign_embedding"`
	HSTSIncludeSubDomains   *bool `json:"hsts_include_sub_domains"`
	HSTSMaxAge              *int  `json:"hsts_max_age"`
	HSTSPreload             *bool `json:"hsts_preload"`
}

type RequiredDNSEntries struct {
	Recommended  DNSValues   `json:"recommended,omitempty"`
	Alternatives []DNSValues `json:"alternatives,omitempty"`
}

type DNSValues struct {
	Type   DNSValuesType `json:"type,omitempty"`
	Values []string      `json:"values,omitempty"`
}

type DNSValuesType string

type ResponseCode int

type DNSStatus string

type CertificateStatus string

const (
	DNSARecord     DNSValuesType = "A"
	DNSCNAMERecord DNSValuesType = "CNAME"
)

const (
	ResponseCodeMovedPermanently ResponseCode = 301
	ResponseCodeFound            ResponseCode = 302
	ResponseCodeNotFound         ResponseCode = 401
)

const (
	DNSStatusActive  DNSStatus = "active"
	DNSStatusInvalid DNSStatus = "invalid"
)

const (
	CertificateStatusActive                     CertificateStatus = "active"
	CertificateStatusProcessing                 CertificateStatus = "processing"
	CertificateStatusInvalidDNS                 CertificateStatus = "invalid_dns"
	CertificateStatusAutoSSLNotSupported        CertificateStatus = "auto_ssl_not_supported"
	CertificateStatusHostnameContainsUnderscore CertificateStatus = "hostname_contains_underscore"
	CertificateStatusInvalidCAARecord           CertificateStatus = "invalid_caa_record"
	CertificateStatusAAAARecordPresent          CertificateStatus = "aaaa_record_present"
)

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

func (h Host) String() string {
	return fmt.Sprint(h.Data)
}

func ref[T any](x T) *T {
	return &x
}
