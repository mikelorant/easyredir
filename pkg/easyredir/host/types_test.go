package host

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/leaanthony/go-ansi-parser"
	"github.com/maxatome/go-testdeep/td"
	"github.com/mikelorant/easyredir/pkg/easyredir/option"
)

type WithBaseURL string

func (u WithBaseURL) Apply(o *option.Options) {
	o.BaseURL = string(u)
}

func TestHostsDataStringer(t *testing.T) {
	tests := []struct {
		name string
		give Data
		want string
	}{
		{
			name: "minimal",
			give: Data{
				ID:   "abc-def",
				Type: "host",
			},
			want: heredoc.Doc(`
				id: abc-def
				type: host
			`),
		},
		{
			name: "empty",
			give: Data{},
			want: heredoc.Doc(`
				id: ""
				type: ""
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give
			td.CmpString(t, got, strings.ReplaceAll(tt.want, "\t", "    "))
		})
	}
}

func TestHostsStringer(t *testing.T) {
	tests := []struct {
		name string
		give Hosts
		want string
	}{
		{
			name: "one",
			give: Hosts{
				Data: []Data{
					{
						ID:   "abc-def",
						Type: "host",
					},
				},
			},
			want: heredoc.Doc(`
				id: abc-def
				type: host
			`),
		},
		{
			name: "multiple",
			give: Hosts{
				Data: []Data{
					{
						ID:   "abc-def",
						Type: "host",
						Attributes: Attributes{
							DNSStatus:         DNSStatusActive,
							CertificateStatus: CertificateStatusActive,
						},
					},
					{
						ID:   "def-abc",
						Type: "host",
						Attributes: Attributes{
							DNSStatus:         DNSStatusInvalid,
							CertificateStatus: CertificateStatusInvalidDNS,
						},
					},
				},
			},
			want: heredoc.Doc(`
				ID     	DNS STATUS	CERTIFICATE STATUS
				abc-def	active    	active
				def-abc	invalid   	invalid_dns
			`),
		},
		{
			name: "none",
			give: Hosts{
				Data: []Data{},
			},
			want: "No hosts.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fmt.Sprint(tt.give)
			got, _ = ansi.Cleanse(got)
			re := regexp.MustCompile(`[\t\n ]`)
			td.CmpString(t, re.ReplaceAllString(got, ""), re.ReplaceAllString(tt.want, ""))
		})
	}
}
