package host

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir/client"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/option"

	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/assert"
)

type WithBaseURL string

func (u WithBaseURL) Apply(o *option.Options) {
	o.BaseURL = string(u)
}

func TestListHosts(t *testing.T) {
	type Fields struct {
		data string
	}
	type Want struct {
		hosts Hosts
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
					      "type": "host",
					      "attributes": {
					        "name": "easyredir.com",
					        "dns_status": "active",
					        "certificate_status": "active"
					      },
					      "links": {
					        "self": "/v1/hosts/abc-def"
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
				hosts: Hosts{
					Data: []Data{
						{
							ID:   "abc-def",
							Type: "host",
							Attributes: Attributes{
								Name:              "easyredir.com",
								DNSStatus:         "active",
								CertificateStatus: "active",
							},
							Links: Links{
								Self: "/v1/hosts/abc-def",
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
				hosts: Hosts{
					Data: []Data{},
				},
				err: "unable to get json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/hosts/", func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.fields.data))
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			cl := client.New(WithBaseURL(server.URL))

			got, err := ListHosts(cl)
			if tt.want.err != "" {
				assert.NotNil(t, err)
				td.CmpContains(t, err, tt.want.err)
				return
			}
			assert.Nil(t, err)
			td.Cmp(t, got, tt.want.hosts)
		})
	}
}

func TestListHostsPaginator(t *testing.T) {
	type MockData struct {
		status int
		body   string
	}

	type Fields struct {
		data []MockData
	}

	type Want struct {
		hosts Hosts
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
				hosts: Hosts{
					Data: []Data{
						{
							ID:   "abc-def",
							Type: "rule",
						},
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
							      "type": "host"
							    }
							  ],
							  "meta": {
								  "has_more": true
							  },
							  "links": {
								  "next": "/v1/hosts?starting_after=abc-def"
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
							      "type": "host"
							    }
							  ]
							}
						`,
					},
				},
			},
			want: Want{
				hosts: Hosts{
					Data: []Data{
						{
							ID:   "abc-def",
							Type: "host",
						}, {
							ID:   "bcd-efg",
							Type: "host",
						},
					},
				},
			},
		},
		{
			name: "none",
			want: Want{
				hosts: Hosts{
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
							      "type": "host"
							    }
							  ],
							  "meta": {
								  "has_more": true
							  },
							  "links": {
								  "next": "/v1/hosts?starting_after=abc-def"
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
				hosts: Hosts{
					Data: []Data{
						{
							ID:   "abc-def",
							Type: "host",
						},
					},
				},
				err: "unable to get a hosts page",
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
			mux.HandleFunc("/hosts/", func(w http.ResponseWriter, req *http.Request) {
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

			got, err := ListHostsPaginator(cl)
			if tt.want.err != "" {
				assert.NotNil(t, err)
				td.CmpContains(t, err, tt.want.err)
				return
			}
			assert.Nil(t, err)
			td.Cmp(t, got, tt.want.hosts)
		})
	}
}

func TestGetHosts(t *testing.T) {
	type Args struct {
		id string
	}

	type Fields struct {
		status int
		data   string
	}

	type Want struct {
		host Host
		err  string
	}

	tests := []struct {
		name   string
		args   Args
		fields Fields
		want   Want
	}{
		{
			name: "valid",
			args: Args{
				id: "abc-123",
			},
			fields: Fields{
				status: http.StatusOK,
				data: `
					{
						"data": {
							"id": "abc-123",
							"type": "host"
						}
					}
				`,
			},
			want: Want{
				host: Host{
					Data: DataExtended{
						ID:   "abc-123",
						Type: "host",
					},
				},
			},
		},
		{
			name: "invalid",
			args: Args{
				id: "abc-123",
			},
			fields: Fields{
				status: http.StatusNotFound,
				data: `
					{
					  "type": "record_not_found_error",
					  "message": "Record not found"
					}
				`,
			},
			want: Want{
				err: "record_not_found_error: Record not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/hosts/", func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(tt.fields.status)
				w.Write([]byte(tt.fields.data))
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			cl := client.New(WithBaseURL(server.URL))

			got, err := GetHost(cl, tt.args.id)
			if tt.want.err != "" {
				assert.NotNil(t, err)
				td.CmpContains(t, err, tt.want.err)
				return
			}
			assert.Nil(t, err)
			td.Cmp(t, got.Data.ID, tt.args.id)
			td.Cmp(t, got, tt.want.host)
		})
	}
}

func TestBuildListHosts(t *testing.T) {
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
				pathQuery: "/hosts",
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
				pathQuery: "/hosts?starting_after=96b30ce8-6331-4c18-ae49-4155c3a2136c",
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
				pathQuery: "/hosts?ending_before=c6312a3c5514-94ea-81c4-1336-8ec03b69",
			},
		}, {
			name: "limit",
			args: Args{
				options: &option.Options{
					Limit: 100,
				},
			},
			want: Want{
				pathQuery: "/hosts?limit=100",
			},
		}, {
			name: "all",
			args: Args{
				options: &option.Options{
					Limit: 100,
					Pagination: option.Pagination{
						StartingAfter: "96b30ce8-6331-4c18-ae49-4155c3a2136c",
					},
				},
			},
			want: Want{
				pathQuery: "/hosts?starting_after=96b30ce8-6331-4c18-ae49-4155c3a2136c&limit=100",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildListHosts(tt.args.options)
			td.Cmp(t, got, tt.want.pathQuery)
		})
	}
}
