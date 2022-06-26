package host

import (
	"fmt"
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mikelorant/easyredir/pkg/easyredir/option"
	"github.com/mikelorant/easyredir/pkg/structutil"
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
	CertificateStatus  CertificateStatus  `json:"certificate_status,omitempty"`
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
	HTTPSUpgrade            *bool `json:"https_upgrade,omitempty"`
	PreventForeignEmbedding *bool `json:"prevent_foreign_embedding,omitempty"`
	HSTSIncludeSubDomains   *bool `json:"hsts_include_sub_domains,omitempty"`
	HSTSMaxAge              *int  `json:"hsts_max_age,omitempty"`
	HSTSPreload             *bool `json:"hsts_preload,omitempty"`
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

const (
	ErrNoHosts = "No hosts."
)

func (h Data) String() string {
	str, _ := structutil.Sprint(h)

	return str
}

func (h Hosts) String() string {
	switch len(h.Data) {
	case 0:
		return ErrNoHosts
	case 1:
		return fmt.Sprint(h.Data[0])
	default:
		t := table.NewWriter()
		t.SetStyle(table.StyleColoredBright)
		t.Style().Options.DrawBorder = false
		t.Style().Color = table.ColorOptions{}
		t.Style().Box.PaddingLeft = ""
		t.Style().Box.PaddingRight = "\t"
		t.Style().Color.Header = text.Colors{text.Bold}
		t.AppendHeader(table.Row{"ID", "DNS STATUS", "CERTIFICATE STATUS"})
		for _, d := range h.Data {
			t.AppendRow(table.Row{
				d.ID,
				d.Attributes.DNSStatus,
				d.Attributes.CertificateStatus,
			})
		}
		return t.Render()
	}
}

func (h Host) String() string {
	return fmt.Sprint(h.Data)
}

func ref[T any](x T) *T {
	return &x
}
