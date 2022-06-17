package rule

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir/client"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/option"

	"github.com/gotidy/ptr"
	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/assert"
)

func TestListRules(t *testing.T) {
	type Fields struct {
		data string
	}

	type Want struct {
		rules Rules
		err   string
	}

	tests := []struct {
		name   string
		fields Fields
		want   Want
	}{
		{
			name: "default",
			fields: Fields{
				data: `
					{
					  "data": [
					    {
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
					    }
					  ],
					  "meta": {
					    "has_more": true
					  },
					  "links": {
					    "next": "/v1/rules?starting_after=abc-def",
					    "prev": "/v1/rules?ending_before=abc-def"
					  }
					}
				`,
			},
			want: Want{
				rules: Rules{
					Data: []Data{
						{
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
					},
					Metadata: option.Metadata{
						HasMore: true,
					},
					Links: option.Links{
						Next: "/v1/rules?starting_after=abc-def",
						Prev: "/v1/rules?ending_before=abc-def",
					},
				},
			},
		}, {
			name: "error_invalid_json",
			fields: Fields{
				data: "notjson",
			},
			want: Want{
				rules: Rules{
					Data: []Data{},
				},
				err: "unable to get json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/rules/", func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.fields.data))
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			cl := client.New(WithBaseURL(server.URL))

			got, err := ListRules(cl)
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

func TestListRulesPaginator(t *testing.T) {
	type MockData struct {
		status int
		body   string
	}

	type Fields struct {
		data []MockData
	}

	type Want struct {
		rules Rules
		err   string
	}

	tests := []struct {
		name   string
		fields Fields
		want   Want
	}{
		{
			name: "one",
			fields: Fields{
				data: []MockData{
					{
						status: http.StatusOK,
						body: `
							{
							  "data": [
							    {
							      "id": "abc-def",
							      "type": "rule"
							    }
							  ]
							}
						`,
					},
				},
			},
			want: Want{
				rules: Rules{
					Data: []Data{
						{
							ID:   "abc-def",
							Type: "rule",
						},
					},
					Metadata: option.Metadata{
						HasMore: false,
					},
				},
			},
		},
		{
			name: "many",
			fields: Fields{
				data: []MockData{
					{
						status: http.StatusOK,
						body: `
							{
							  "data": [
							    {
							      "id": "abc-def",
							      "type": "rule"
							    }
							  ],
							  "meta": {
								  "has_more": true
							  },
							  "links": {
								  "next": "/v1/rules?starting_after=abc-def"
							  }
							}
						`,
					},
					{
						status: http.StatusOK,
						body: `
							{
							  "data": [
							    {
							      "id": "bcd-efg",
							      "type": "rule"
							    }
							  ]
							}
						`,
					},
				},
			},
			want: Want{
				rules: Rules{
					Data: []Data{
						{
							ID:   "abc-def",
							Type: "rule",
						}, {
							ID:   "bcd-efg",
							Type: "rule",
						},
					},
				},
			},
		},
		{
			name: "none",
			want: Want{
				rules: Rules{
					Data: []Data{},
				},
			},
		},
		{
			name: "invalid_page",
			fields: Fields{
				data: []MockData{
					{
						status: http.StatusOK,
						body: `
							{
							  "data": [
							    {
							      "id": "abc-def",
							      "type": "rule"
							    }
							  ],
							  "meta": {
								  "has_more": true
							  },
							  "links": {
								  "next": "/v1/rules?starting_after=abc-def"
							  }
							}
						`,
					},
					{
						status: http.StatusOK,
						body:   `{ notjson }`,
					},
				},
			},
			want: Want{
				rules: Rules{
					Data: []Data{
						{
							ID:   "abc-def",
							Type: "rule",
						},
					},
				},
				err: "unable to get a rules page",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := func() func() int {
				i := 0
				return func() int {
					i++
					return i - 1
				}
			}()

			mux := http.NewServeMux()
			mux.HandleFunc("/rules/", func(w http.ResponseWriter, req *http.Request) {
				if len(tt.fields.data) < 1 {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("{}"))
					return
				}

				data := tt.fields.data[page()]
				w.WriteHeader(data.status)
				w.Write([]byte(data.body))
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			cl := client.New(WithBaseURL(server.URL))

			got, err := ListRulesPaginator(cl)
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

func TestBuildListRules(t *testing.T) {
	type Args struct {
		options *option.Options
	}

	type Want struct {
		pathQuery string
	}

	tests := []struct {
		name string
		args Args
		want Want
	}{
		{
			name: "no_options",
			args: Args{
				options: &option.Options{},
			},
			want: Want{
				pathQuery: "/rules",
			},
		}, {
			name: "starting_after",
			args: Args{
				options: &option.Options{
					Pagination: option.Pagination{
						StartingAfter: "96b30ce8-6331-4c18-ae49-4155c3a2136c",
					},
				},
			},
			want: Want{
				pathQuery: "/rules?starting_after=96b30ce8-6331-4c18-ae49-4155c3a2136c",
			},
		}, {
			name: "ending_before",
			args: Args{
				options: &option.Options{
					Pagination: option.Pagination{
						EndingBefore: "c6312a3c5514-94ea-81c4-1336-8ec03b69",
					},
				},
			},
			want: Want{
				pathQuery: "/rules?ending_before=c6312a3c5514-94ea-81c4-1336-8ec03b69",
			},
		}, {
			name: "source_filter",
			args: Args{
				options: &option.Options{
					SourceFilter: "http://www1.example.org",
				},
			},
			want: Want{
				pathQuery: "/rules?sq=http://www1.example.org",
			},
		}, {
			name: "target_filter",
			args: Args{
				options: &option.Options{
					TargetFilter: "http://www2.example.org",
				},
			},
			want: Want{
				pathQuery: "/rules?tq=http://www2.example.org",
			},
		}, {
			name: "source_target_filter",
			args: Args{
				options: &option.Options{
					SourceFilter: "http://www1.example.org",
					TargetFilter: "http://www2.example.org",
				},
			},
			want: Want{
				pathQuery: "/rules?sq=http://www1.example.org&tq=http://www2.example.org",
			},
		}, {
			name: "limit",
			args: Args{
				options: &option.Options{
					Limit: 100,
				},
			},
			want: Want{
				pathQuery: "/rules?limit=100",
			},
		}, {
			name: "all",
			args: Args{
				options: &option.Options{
					SourceFilter: "http://www1.example.org",
					TargetFilter: "http://www2.example.org",
					Limit:        100,
					Pagination: option.Pagination{
						StartingAfter: "96b30ce8-6331-4c18-ae49-4155c3a2136c",
					},
				},
			},
			want: Want{
				pathQuery: "/rules?starting_after=96b30ce8-6331-4c18-ae49-4155c3a2136c&sq=http://www1.example.org&tq=http://www2.example.org&limit=100",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildListRules(tt.args.options)
			td.Cmp(t, got, tt.want.pathQuery)
		})
	}
}
