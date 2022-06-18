package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/maxatome/go-testdeep/td"
	"github.com/mikelorant/easyredir/pkg/easyredir/option"
	"github.com/stretchr/testify/assert"
)

type WithAPIKey string

func (k WithAPIKey) Apply(o *option.Options) {
	o.APIKey = string(k)
}

type WithAPISecret string

func (s WithAPISecret) Apply(o *option.Options) {
	o.APISecret = string(s)
}

type WithBaseURL string

func (u WithBaseURL) Apply(o *option.Options) {
	o.BaseURL = string(u)
}

func TestSendRequest(t *testing.T) {
	type Args struct {
		path   string
		method string
		body   string
	}

	type Fields struct {
		apiKey    string
		apiSecret string
		status    int
		header    map[string]string
		body      string
	}

	type Want struct {
		authorization string
		body          string
		err           string
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
		},
		{
			name: "ok_empty",
			fields: Fields{
				body: "",
			},
			want: Want{
				body: "",
			},
		},
		{
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
				err: heredoc.Doc(`
					invalid_request_error: Invalid Request
					errors:
					- resource: rule
					  param: forward_params
					  code: invalid_option
					  message: Must be true or false
				`),
			},
		},
		{
			name: "generic_error",
			fields: Fields{
				status: http.StatusInternalServerError,
				body:   "",
			},
			want: Want{
				body: "",
				err:  "unknown error: status code: 500",
			},
		},
		{
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
		},
		{
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
		{
			name: "api_key_secret",
			fields: Fields{
				apiKey:    "test",
				apiSecret: "test",
				body:      "",
			},
			want: Want{
				authorization: "Basic dGVzdDp0ZXN0",
				body:          "",
			},
		},
		{
			name: "post",
			args: Args{
				method: http.MethodPost,
				body:   "payload",
			},
			fields: Fields{
				status: http.StatusOK,
				body:   "ok",
			},
			want: Want{
				body: "ok",
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

			authorization := "Basic Og=="
			if tt.want.authorization != "" {
				authorization = tt.want.authorization
			}

			mux := http.NewServeMux()
			mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
				if req.Method == http.MethodPatch || req.Method == http.MethodPost || req.Method == http.MethodPut {
					_, ok := req.Header["Idempotency-Key"]
					td.CmpTrue(t, ok)
				}

				td.Cmp(t, req.Header.Get("Authorization"), authorization)
				td.Cmp(t, req.Header.Get("Content-Type"), ResourceType)
				td.Cmp(t, req.Header.Get("Accept"), ResourceType)

				for k, v := range tt.fields.header {
					w.Header().Set(k, v)
				}
				w.WriteHeader(status)
				w.Write([]byte(tt.fields.body))
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			cl := New(
				WithAPIKey(tt.fields.apiKey),
				WithAPISecret(tt.fields.apiSecret),
				WithBaseURL(server.URL),
			)

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
