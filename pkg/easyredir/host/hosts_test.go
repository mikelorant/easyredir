package host

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/mikelorant/easyredir-cli/pkg/easyredir"

	"github.com/MakeNowJust/heredoc"
	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	data string
}

func (m *mockClient) SendRequest(baseURL, path, method string, body io.Reader) (io.ReadCloser, error) {
	r := strings.NewReader(m.data)
	rc := io.NopCloser(r)
	return rc, nil
}

func TestListHosts(t *testing.T) {
	type Args struct {
		options []func(*HostsOptions)
	}
	type Fields struct {
		data string
	}
	type Want struct {
		hosts Hosts
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
					Data: []HostData{
						{
							ID:   "abc-def",
							Type: "host",
							Attributes: HostAttributes{
								Name:              "easyredir.com",
								DNSStatus:         "active",
								CertificateStatus: "active",
							},
							Links: HostLinks{
								Self: "/v1/hosts/abc-def",
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
		},
		{
			name: "minimal",
			fields: Fields{
				data: `
					{
					  "data": [
					    {
					      "id": "abc-def",
					      "type": "host"
					    }
					  ]
					}
				`,
			},
			want: Want{
				hosts: Hosts{
					Data: []HostData{
						{
							ID:   "abc-def",
							Type: "host",
						},
					},
				},
			},
		},
		{
			name: "with_limit",
			args: Args{
				options: []func(*HostsOptions){
					WithHostsLimit(1),
				},
			},
			fields: Fields{
				data: `
					{
					  "data": [
					    {
					      "id": "abc-def",
					      "type": "host"
					    }
					  ]
					}
				`,
			},
			want: Want{
				hosts: Hosts{
					Data: []HostData{
						{
							ID:   "abc-def",
							Type: "host",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &easyredir.Easyredir{
				Client: &mockClient{
					data: tt.fields.data,
				},
				Config: &easyredir.Config{},
			}

			got, err := ListHosts(e, tt.args.options...)
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

func TestBuildListHosts(t *testing.T) {
	type Args struct {
		options *HostsOptions
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
				options: &HostsOptions{},
			},
			want: Want{
				pathQuery: "/hosts",
			},
		}, {
			name: "starting_after",
			args: Args{
				options: &HostsOptions{
					pagination: Pagination{
						startingAfter: "96b30ce8-6331-4c18-ae49-4155c3a2136c",
					},
				},
			},
			want: Want{
				pathQuery: "/hosts?starting_after=96b30ce8-6331-4c18-ae49-4155c3a2136c",
			},
		}, {
			name: "ending_before",
			args: Args{
				options: &HostsOptions{
					pagination: Pagination{
						endingBefore: "c6312a3c5514-94ea-81c4-1336-8ec03b69",
					},
				},
			},
			want: Want{
				pathQuery: "/hosts?ending_before=c6312a3c5514-94ea-81c4-1336-8ec03b69",
			},
		}, {
			name: "limit",
			args: Args{
				options: &HostsOptions{
					limit: 100,
				},
			},
			want: Want{
				pathQuery: "/hosts?limit=100",
			},
		}, {
			name: "all",
			args: Args{
				options: &HostsOptions{
					limit: 100,
					pagination: Pagination{
						startingAfter: "96b30ce8-6331-4c18-ae49-4155c3a2136c",
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

type mockPaginatorClient struct {
	idx  int
	data string
}

func (m *mockPaginatorClient) SendRequest(baseURL, path, method string, body io.Reader) (io.ReadCloser, error) {
	data := strings.NewReader(m.data)
	docs := make(map[int]interface{})
	dec := json.NewDecoder(data)

	i := 0
	for {
		var doc map[string]interface{}
		err := dec.Decode(&doc)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("unable to decode json page %v: %w", i, err)
		}
		docs[i] = doc
		i++
	}

	b, err := json.Marshal(docs[m.idx])
	if err != nil {
		return nil, fmt.Errorf("unable to encode json page %v: %w", i, err)
	}

	r := strings.NewReader(string(b))
	rc := io.NopCloser(r)

	m.idx++

	return rc, nil
}

func TestListHostsPaginator(t *testing.T) {
	type Args struct {
		options []func(*HostsOptions)
	}

	type Fields struct {
		data string
	}

	type Want struct {
		hosts Hosts
		err   string
	}

	tests := []struct {
		name   string
		args   Args
		fields Fields
		want   Want
	}{
		{
			name: "one",
			fields: Fields{
				data: `
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
			want: Want{
				hosts: Hosts{
					Data: []HostData{
						{
							ID:   "abc-def",
							Type: "rule",
						},
					},
					Metadata: Metadata{
						HasMore: false,
					},
				},
			},
		},
		{
			name: "many",
			fields: Fields{
				data: `
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
			want: Want{
				hosts: Hosts{
					Data: []HostData{
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
					Data: []HostData{},
				},
			},
		},
		{
			name: "invalid_page",
			fields: Fields{
				data: `
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
					{ notjson }
				`,
			},
			want: Want{
				hosts: Hosts{
					Data: []HostData{
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
			e := &easyredir.Easyredir{
				Client: &mockPaginatorClient{
					data: tt.fields.data,
				},
				Config: &easyredir.Config{},
			}

			got, err := ListHostsPaginator(e, tt.args.options...)
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

func TestHostsDataStringer(t *testing.T) {
	tests := []struct {
		name string
		give HostData
		want string
	}{
		{
			name: "minimal",
			give: HostData{
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
			give: HostData{},
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
				Data: []HostData{
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

func TestGetHosts(t *testing.T) {
	type Args struct {
		id string
	}

	type Fields struct {
		data string
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
					Data: HostDataExtended{
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
				data: `
					{
						"data": {
							"id": "def-456",
							"type": "host"
						}
					}
				`,
			},
			want: Want{
				err: "received incorrect host",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &easyredir.Easyredir{
				Client: &mockClient{
					data: tt.fields.data,
				},
				Config: &easyredir.Config{},
			}

			got, err := GetHost(e, tt.args.id)
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
