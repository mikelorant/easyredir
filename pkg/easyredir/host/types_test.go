package host

import (
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/maxatome/go-testdeep/td"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/option"
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
