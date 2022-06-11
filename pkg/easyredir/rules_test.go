package easyredir

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	filename string
}

func (m *mockClient) sendRequest(baseURL, path, method string, body io.Reader) (io.ReadCloser, error) {
	fh, err := os.Open(m.filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}
	return fh, nil
}

func TestListRules(t *testing.T) {
	type Args struct {
		options []func(*RulesOptions)
	}

	type Fields struct {
		filename string
	}

	type Want struct {
		rules Rules
		err   string
	}

	tests := []struct {
		name   string
		args   Args
		fields Fields
		want   Want
	}{
		{
			name: "default",
			fields: Fields{
				filename: "testdata/rules_list_default.json",
			},
			want: Want{
				rules: Rules{
					Data: []RuleData{
						{
							ID:   ptr.String("abc-def"),
							Type: ptr.String("rule"),
							Attributes: RuleAttributes{
								ForwardParams: ptr.Bool(true),
								ForwardPath:   ptr.Bool(true),
								ResponseType:  ptr.String("moved_permanently"),
								SourceURLs: []*string{
									ptr.String("abc.com"),
									ptr.String("abc.com/123"),
								},
								TargetURL: ptr.String("otherdomain.com"),
							},
						},
					},
					Metadata: Metadata{
						HasMore: true,
					},
					Links: Links{
						Next: "/v1/rules?starting_after=abc-def",
						Prev: "/v1/rules?ending_before=abc-def",
					},
				},
			},
		}, {
			name: "with_source_filter",
			args: Args{
				options: []func(*RulesOptions){
					WithSourceFilter("https://www1.example.org"),
				},
			},
			fields: Fields{
				filename: "testdata/rules_list_minimal.json",
			},
			want: Want{
				rules: Rules{
					Data: []RuleData{
						{
							ID:   ptr.String("abc-def"),
							Type: ptr.String("rule"),
						},
					},
				},
			},
		}, {
			name: "with_target_filter",
			args: Args{
				options: []func(*RulesOptions){
					WithTargetFilter("https://www2.example.org"),
				},
			},
			fields: Fields{
				filename: "testdata/rules_list_minimal.json",
			},
			want: Want{
				rules: Rules{
					Data: []RuleData{
						{
							ID:   ptr.String("abc-def"),
							Type: ptr.String("rule"),
						},
					},
				},
			},
		}, {
			name: "with_both_source_target_filter",
			args: Args{
				options: []func(*RulesOptions){
					WithSourceFilter("https://www1.example.org"),
					WithTargetFilter("https://www2.example.org"),
				},
			},
			fields: Fields{
				filename: "testdata/rules_list_minimal.json",
			},
			want: Want{
				rules: Rules{
					Data: []RuleData{
						{
							ID:   ptr.String("abc-def"),
							Type: ptr.String("rule"),
						},
					},
				},
			},
		}, {
			name: "with_limit",
			args: Args{
				options: []func(*RulesOptions){
					WithLimit(1),
				},
			},
			fields: Fields{
				filename: "testdata/rules_list_minimal.json",
			},
			want: Want{
				rules: Rules{
					Data: []RuleData{
						{
							ID:   ptr.String("abc-def"),
							Type: ptr.String("rule"),
						},
					},
				},
			},
		}, {
			name: "error_invalid_json",
			fields: Fields{
				filename: "testdata/rules_list_invalid.json",
			},
			want: Want{
				rules: Rules{
					Data: []RuleData{},
				},
				err: "unable to get json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Easyredir{
				client: &mockClient{
					filename: tt.fields.filename,
				},
				config: &Config{},
			}

			got, err := e.ListRules(tt.args.options...)
			if tt.want.err != "" {
				assert.NotNil(t, err)
				td.CmpContains(t, err, tt.want.err)
				return
			}
			assert.Nil(t, err)
			td.Cmp(t, got, tt.want.rules)
		})
	}
}
