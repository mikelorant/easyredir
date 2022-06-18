package rule

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/gotidy/ptr"
	"github.com/leaanthony/go-ansi-parser"
	"github.com/maxatome/go-testdeep/td"
	"github.com/mikelorant/easyredir/pkg/easyredir/option"
)

type WithBaseURL string

func (u WithBaseURL) Apply(o *option.Options) {
	o.BaseURL = string(u)
}

func TestRulesDataStringer(t *testing.T) {
	tests := []struct {
		name string
		give Data
		want string
	}{
		{
			name: "minimal",
			give: Data{
				ID:   "abc-def",
				Type: "rule",
			},
			want: heredoc.Doc(`
				id: abc-def
				type: rule
			`),
		},
		{
			name: "typical",
			give: Data{
				ID:   "abc-def",
				Type: "rule",
				Attributes: Attributes{
					ForwardParams: ptr.Bool(true),
					ForwardPath:   ptr.Bool(true),
					ResponseType:  ref(ResponseMovedPermanently),
					SourceURLs: []string{
						"http://www1.example.org",
						"http://www2.example.org",
					},
					TargetURL: ptr.String("http://www3.example.org"),
				},
			},
			want: heredoc.Doc(`
				id: abc-def
				type: rule
				attributes:
				  forward_params: true
				  forward_path: true
				  response_type: moved_permanently
				  source_urls:
				  - http://www1.example.org
				  - http://www2.example.org
				  target_url: http://www3.example.org
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

func TestRulesStringer(t *testing.T) {
	tests := []struct {
		name string
		give Rules
		want string
	}{
		{
			name: "one",
			give: Rules{
				Data: []Data{
					{
						ID:   "abc-def",
						Type: "rule",
					},
				},
			},
			want: heredoc.Doc(`
				id: abc-def
				type: rule
			`),
		},
		{
			name: "multiple",
			give: Rules{
				Data: []Data{
					{
						ID:   "abc-def",
						Type: "rule",
						Attributes: Attributes{
							SourceURLs: []string{
								"source.example.com",
							},
							TargetURL: ref("target.example.com"),
						},
					},
					{
						ID:   "def-abc",
						Type: "rule",
						Attributes: Attributes{
							SourceURLs: []string{
								"source2.example.com",
							},
							TargetURL: ref("target2.example.com"),
						},
					},
				},
			},
			want: heredoc.Doc(`
				ID		SOURCE URLS			TARGET URL
				abc-def	source.example.com	target.example.com
				def-abc	source2.example.com	target2.example.com
			`),
		},
		{
			name: "none",
			give: Rules{},
			want: "No rules.",
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
