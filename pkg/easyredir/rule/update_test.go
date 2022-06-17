package rule

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/maxatome/go-testdeep/td"
	"github.com/mikelorant/easyredir/pkg/easyredir/client"
	"github.com/stretchr/testify/assert"
)

func TestUpdateRule(t *testing.T) {
	type Args struct {
		id         string
		attributes Attributes
	}
	type Fields struct {
		status int
		body   string
	}
	type Want struct {
		rule Rule
		err  string
	}

	tests := []struct {
		name   string
		args   Args
		fields Fields
		want   Want
	}{
		{
			name: "success",
			args: Args{
				id: "5d29f819-302f-40c0-8b5a-6d670267615b",
				attributes: Attributes{
					ForwardParams: ptr.Bool(true),
					ForwardPath:   ptr.Bool(true),
					ResponseType:  ref(ResponseMovedPermanently),
					SourceURLs: []string{
						"abc.com",
						"abc.com/123",
					},
					TargetURL: ptr.String("otherdomain.com"),
				},
			},
			fields: Fields{
				status: http.StatusOK,
				body: `
					{
					  "data": {
					    "id": "abc-def",
					    "type": "rule",
					    "attributes": {
					      "forward_params": true,
					      "forward_path": true,
					      "response_type": "moved_permanently",
					      "source_urls": [
					        "abc.com",
					        "abc.com/123"
					      ],
					      "target_url": "otherdomain.com"
					    }
					  },
					  "relationships": {
					    "source_hosts": {
					      "data": [
					        {
					          "id": "abc-123",
					          "type": "host"
					        }
					      ],
					      "links": {
					        "related": "/api/v1/rules/e819dbbb-892c-487f-a2e0-7b60b0362b6c/hosts"
					      }
					    }
					  }
					}
				`,
			},
			want: Want{
				rule: Rule{
					Data: Data{
						ID:   "abc-def",
						Type: "rule",
						Attributes: Attributes{
							ForwardParams: ptr.Bool(true),
							ForwardPath:   ptr.Bool(true),
							ResponseType:  ref(ResponseMovedPermanently),
							SourceURLs: []string{
								"abc.com",
								"abc.com/123",
							},
							TargetURL: ptr.String("otherdomain.com"),
						},
					},
					Relationships: Relationships{
						SourceHosts: SourceHosts{
							Data: []SourceHostData{
								{
									ID:   "abc-123",
									Type: "host",
								},
							},
							Links: SourceHostsLinks{
								Related: "/api/v1/rules/e819dbbb-892c-487f-a2e0-7b60b0362b6c/hosts",
							},
						},
					},
				},
			},
		},
		{
			name: "failure",
			args: Args{
				id: "5d29f819-302f-40c0-8b5a-6d670267615b",
				attributes: Attributes{
					ForwardParams: ptr.Bool(true),
					ForwardPath:   ptr.Bool(true),
					ResponseType:  ref(ResponseMovedPermanently),
					SourceURLs:    []string{},
					TargetURL:     ptr.String("otherdomain.com"),
				},
			},
			fields: Fields{
				status: http.StatusUnprocessableEntity,
				body: `
					{
					  "type": "invalid_request_error",
					  "message": "Invalid Request",
					  "errors": [
					    {
					      "resource": "rule",
					      "param": "source_urls",
					      "code": "required",
					      "message": "You need to include at least one URL for us to redirect"
					    }
					  ]
					}
				`,
			},
			want: Want{
				rule: Rule{},
				err:  "invalid_request_error: Invalid Request",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/rules/", func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(tt.fields.status)
				w.Write([]byte(tt.fields.body))
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			cl := client.New(WithBaseURL(server.URL))

			got, err := UpdateRule(cl, tt.args.id, tt.args.attributes)
			if tt.want.err != "" {
				assert.NotNil(t, err)
				td.CmpContains(t, err, tt.want.err)
				return
			}
			assert.Nil(t, err)
			td.Cmp(t, got, tt.want.rule)
		})
	}
}
