package host

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikelorant/easyredir/pkg/easyredir/client"

	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/assert"
)

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
					Data: Data{
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
