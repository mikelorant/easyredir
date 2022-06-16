package easyredir

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/assert"
)

func TestSendRequest(t *testing.T) {
	type Args struct {
		path   string
		method string
		body   string
	}

	type Fields struct {
		status int
		header map[string]string
		body   string
	}

	type Want struct {
		body string
		err  string
	}

	tests := []struct {
		name   string
		args   Args
		fields Fields
		want   Want
	}{
		{
			name: "ok",
			fields: Fields{
				body: "payload",
			},
			want: Want{
				body: "payload",
			},
		}, {
			name: "ok_empty",
			fields: Fields{
				body: "",
			},
			want: Want{
				body: "",
			},
		}, {
			name: "custom_error",
			fields: Fields{
				status: http.StatusBadRequest,
				body: `
					{
					  "type": "invalid_request_error",
					  "message": "Invalid Request",
					  "errors": [
					    {
					      "resource": "rule",
					      "param": "forward_params",
					      "code": "invalid_option",
					      "message": "Must be true or false"
					    }
					  ]
					}
				`,
			},
			want: Want{
				body: "",
				err:  "invalid_request_error: Invalid Request",
			},
		}, {
			name: "generic_error",
			fields: Fields{
				status: http.StatusInternalServerError,
				body:   "",
			},
			want: Want{
				body: "",
				err:  "received status code: 500",
			},
		}, {
			name: "method_invalid",
			args: Args{
				method: "invalid method",
			},
			fields: Fields{
				body: "",
			},
			want: Want{
				body: "",
				err:  `net/http: invalid method "invalid method"`,
			},
		}, {
			name: "rate_limited",
			fields: Fields{
				status: http.StatusTooManyRequests,
				header: map[string]string{
					"X-Ratelimit-Limit":     "1",
					"X-Ratelimit-Remaining": "2",
					"X-Ratelimit-Reset":     "3",
				},
				body: "",
			},
			want: Want{
				body: "",
				err:  "rate limited with limit: 1, remaining: 2, reset: 3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/"
			if tt.args.path != "" {
				path = tt.args.path
			}

			method := http.MethodGet
			if tt.args.method != "" {
				method = tt.args.method
			}

			status := http.StatusOK
			if tt.fields.status != 0 {
				status = tt.fields.status
			}

			mux := http.NewServeMux()
			mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
				for k, v := range tt.fields.header {
					w.Header().Set(k, v)
				}
				w.WriteHeader(status)
				w.Write([]byte(tt.fields.body))
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			cl := NewClient(WithBaseURL(server.URL))
			r, err := cl.SendRequest(path, method, strings.NewReader(tt.args.body))
			if tt.want.err != "" {
				assert.NotNil(t, err)
				td.CmpContains(t, err, tt.want.err)
				return
			}
			assert.Nil(t, err)

			got, err := io.ReadAll(r)
			assert.Nil(t, err)
			td.Cmp(t, string(got), tt.want.body)
		})
	}
}
