package easyredir

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/client"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/config"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "ping",
			want: "pong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New("", "")

			got := e.Ping()
			td.Cmp(t, got, tt.want)
		})
	}
}

func TestSendRequest(t *testing.T) {
	type Fields struct {
		method     string
		statusCode int
		body       []byte
	}

	type Want struct {
		body string
		err  string
	}

	tests := []struct {
		name   string
		fields Fields
		want   Want
	}{
		{
			name: "ok",
			fields: Fields{
				body: []byte("payload"),
			},
			want: Want{
				body: "payload",
			},
		}, {
			name: "ok_empty",
			fields: Fields{
				body: []byte{},
			},
			want: Want{
				body: "",
			},
		}, {
			name: "custom_error",
			fields: Fields{
				statusCode: http.StatusBadRequest,
				body: []byte(`
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
				`),
			},
			want: Want{
				body: "",
				err:  "invalid_request_error: Invalid Request",
			},
		}, {
			name: "status_code_error",
			fields: Fields{
				statusCode: http.StatusInternalServerError,
				body:       []byte("payload"),
			},
			want: Want{
				body: "",
				err:  "received status code: 500",
			},
		}, {
			name: "method_invalid",
			fields: Fields{
				method: "invalid method",
				body:   []byte("payload"),
			},
			want: Want{
				body: "",
				err:  `net/http: invalid method "invalid method"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/"

			method := http.MethodGet
			if tt.fields.method != "" {
				method = tt.fields.method
			}

			statusCode := http.StatusOK
			if tt.fields.statusCode != 0 {
				statusCode = tt.fields.statusCode
			}

			body := strings.NewReader("")

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(statusCode)
				w.Write(tt.fields.body)
			}))
			defer server.Close()

			cfg := config.New("", "")
			cl := client.New(cfg)

			r, err := cl.SendRequest(server.URL, path, method, body)
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
