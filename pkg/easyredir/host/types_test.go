package host

import (
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/maxatome/go-testdeep/td"
)

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
				attributes:
				  name: ""
				  dns_status: ""
				  certificate_status: ""
				links:
				  self: ""
			`),
		},
		{
			name: "empty",
			give: Data{},
			want: heredoc.Doc(`
				id: ""
				type: ""
				attributes:
				  name: ""
				  dns_status: ""
				  certificate_status: ""
				links:
				  self: ""
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
			name: "minimal",
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
				attributes:
				  name: ""
				  dns_status: ""
				  certificate_status: ""
				links:
				  self: ""

				Total: 1
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
