package host

import (
	"fmt"
	"strings"

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

func (h DataExtended) String() string {
	str, _ := structutil.Sprint(h)

	return str
}

func (h Host) String() string {
	return fmt.Sprint(h.Data)
}
