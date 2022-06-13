package easyredir

import (
	"testing"

	// "github.com/gotidy/ptr"
	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/assert"
)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Easyredir{
				client: &mockClient{
					data: tt.fields.data,
				},
				config: &Config{},
			}

			got, err := e.ListHosts(tt.args.options...)
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
			e := &Easyredir{
				client: &mockPaginatorClient{
					data: tt.fields.data,
				},
				config: &Config{},
			}

			got, err := e.ListHostsPaginator(tt.args.options...)
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
