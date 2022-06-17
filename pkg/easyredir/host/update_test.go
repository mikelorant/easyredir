package host

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/maxatome/go-testdeep/td"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/client"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHost(t *testing.T) {
	type Args struct {
		id 			string
		attributes	Attributes
	}
	type Fields struct {
		status	int
		body	string
	}
	type Want struct {
		host	Host
		err		string
	}

	tests := []struct{
		name	string
		args	Args
		fields	Fields
		want	Want
	}{
		{
			name: "success",
			args: Args{
				id: "b8a2287c-5580-41e8-8b8c-438231dd6875",
				attributes: Attributes{
					MatchOptions: MatchOptions{
						CaseInsensitive: ptr.Bool(true),
						SlashInsensitive: ptr.Bool(true),
					},
					NotFoundAction: NotFoundAction{
						ForwardParams: ptr.Bool(true),
						ForwardPath: ptr.Bool(true),
						Custom404Body: ptr.String("<html><body>My Custom 404 content.</body></html>"),
						ResponseCode: ref(ResponseCodeFound),
						ResponseURL: ptr.String("https://www.example.com"),
					},
					Security: Security{
						HTTPSUpgrade: ptr.Bool(true),
						PreventForeignEmbedding: ptr.Bool(true),
						HSTSIncludeSubDomains: ptr.Bool(true),
						HSTSMaxAge: ptr.Int(31536000),
						HSTSPreload: ptr.Bool(true),
					},
				},
			},
			fields: Fields{
				status: http.StatusOK,
				body: `
					{
					  "data": {
					    "id": "abc-def",
					    "type": "host",
					    "attributes": {
					      "name": "easyredir.com",
					      "dns_status": "active",
					      "dns_tested_at": "2020-11-24T22:33:35Z",
					      "certificate_status": "active",
					      "acme_enabled": true,
					      "match_options": {
					        "case_insensitive": true,
					        "slash_insensitive": true
					      },
					      "not_found_action": {
					        "forward_params": true,
					        "forward_path": true,
					        "custom_404_body_present": true,
					        "response_code": 302,
					        "response_url": "https://www.example.com"
					      },
					      "security": {
					        "https_upgrade": true,
					        "prevent_foreign_embedding": true,
					        "hsts_include_sub_domains": true,
					        "hsts_max_age": 31536000,
					        "hsts_preload": true
					      },
					      "required_dns_entries": {
					        "recommended": {
					          "type": "A",
					          "values": [
					            "34.213.106.51",
					            "54.68.182.72"
					          ]
					        },
					        "alternatives": [
					          {
					            "type": "A",
					            "values": [
					              "34.213.106.51",
					              "54.68.182.72"
					            ]
					          }
					        ]
					      },
					      "detected_dns_entries": [
					        {
					          "type": "A",
					          "values": [
					            "34.213.106.51",
					            "54.68.182.72"
					          ]
					        }
					      ]
					    },
					    "links": {
					      "self": "/v1/hosts/abc-def"
					    }
					  }
					}
				`,
			},
			want: Want{
				host: Host{
					Data: Data{
						ID: "abc-def",
						Type: "host",
						Attributes: Attributes{
							Name: "easyredir.com",
							DNSStatus: "active",
							DNSTestedAt: "2020-11-24T22:33:35Z",
							CertificateStatus: "active",
							ACMEEnabled: ptr.Bool(true),
							MatchOptions: MatchOptions{
								CaseInsensitive: ptr.Bool(true),
								SlashInsensitive: ptr.Bool(true),
							},
							NotFoundAction: NotFoundAction{
								ForwardParams: ptr.Bool(true),
								ForwardPath: ptr.Bool(true),
								Custom404BodyPresent: ptr.Bool(true),
								ResponseCode: ref(ResponseCodeFound),
								ResponseURL: ptr.String("https://www.example.com"),
							},
							Security: Security{
								HTTPSUpgrade: ptr.Bool(true),
								PreventForeignEmbedding: ptr.Bool(true),
								HSTSIncludeSubDomains: ptr.Bool(true),
								HSTSMaxAge: ptr.Int(31536000),
								HSTSPreload: ptr.Bool(true),
							},
							RequiredDNSEntries: RequiredDNSEntries{
								Recommended: DNSValues{
									Type: "A",
									Values: []string{
										"34.213.106.51",
										"54.68.182.72",
									},
								},
								Alternatives: []DNSValues{
									{
										Type: "A",
										Values: []string{
											"34.213.106.51",
											"54.68.182.72",
										},
									},
								},
							},
							DetectedDNSEntries: []DNSValues{
								{
									Type: "A",
									Values: []string{
										"34.213.106.51",
										"54.68.182.72",
									},
								},
							},

						},
						Links: Links{
							Self: "/v1/hosts/abc-def",
						},
					},
				},
			},
		},
		{
			name: "failure",
			args: Args{
				id: "b8a2287c-5580-41e8-8b8c-438231dd6875",
				attributes: Attributes{
					MatchOptions: MatchOptions{
						CaseInsensitive: ptr.Bool(true),
						SlashInsensitive: ptr.Bool(true),
					},
					NotFoundAction: NotFoundAction{
						ForwardParams: ptr.Bool(true),
						ForwardPath: ptr.Bool(true),
						Custom404Body: ptr.String("<html><body>My Custom 404 content.</body></html>"),
						ResponseCode: ref(ResponseCodeFound),
						ResponseURL: ptr.String("https://www.example.com"),
					},
					Security: Security{
						HTTPSUpgrade: ptr.Bool(true),
						PreventForeignEmbedding: ptr.Bool(true),
						HSTSIncludeSubDomains: ptr.Bool(true),
						HSTSMaxAge: ptr.Int(31536000),
						HSTSPreload: ptr.Bool(true),
					},
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
					      "resource": "host",
					      "param": "id",
					      "code": "required",
					      "message": ""
					    }
					  ]
					}
				`,
			},
			want: Want{
				host: Host{},
				err: "invalid_request_error: Invalid Request",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/hosts/", func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(tt.fields.status)
				w.Write([]byte(tt.fields.body))
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			cl := client.New(WithBaseURL(server.URL))

			got, err := UpdateHost(cl, tt.args.id, tt.args.attributes)
			if tt.want.err != "" {
				assert.NotNil(t, err)
				td.CmpContains(t, err, tt.want.err)
				return
			}
			assert.Nil(t, err)
			td.Cmp(t, got, tt.want.host)
		})
	}
}
