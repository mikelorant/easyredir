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
			cl := New(&Config{
				Key:    "",
				Secret: "",
			})

			got := cl.Ping()
			td.Cmp(t, got, tt.want)
		})
	}
}

func TestDecodeJSON(t *testing.T) {
	type Data struct {
		Key string `json:"key"`
	}

	type Args struct {
		src io.Reader
		dst Data
	}

	type Want struct {
		dst Data
		err string
	}

	tests := []struct {
		name string
		args Args
		want Want
	}{
		{
			name: "exactfields",
			args: Args{
				src: strings.NewReader(`{ "Key": "Value" }`),
				dst: Data{},
			},
			want: Want{
				dst: Data{
					Key: "Value",
				},
			},
		},
		{
			name: "extrafields",
			args: Args{
				src: strings.NewReader(`{ "Key": "Value", "Key2": "Value2" }`),
				dst: Data{},
			},
			want: Want{
				dst: Data{
					Key: "Value",
				},
			},
		},
		{
			name: "nofields",
			args: Args{
				src: strings.NewReader(`{}`),
				dst: Data{},
			},
			want: Want{
				dst: Data{},
			},
		},
		{
			name: "notjson",
			args: Args{
				src: strings.NewReader(`not json`),
				dst: Data{},
			},
			want: Want{
				dst: Data{},
				err: "unable to json decode",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Data{}
			err := decodeJSON(io.NopCloser(tt.args.src), &got)
			if tt.want.err != "" {
				assert.NotNil(t, err)
				td.CmpContains(t, err, tt.want.err)
				return
			}
			assert.Nil(t, err)
			td.Cmp(t, got, tt.want.dst)
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
		},
		{
			name: "ok_empty",
			fields: Fields{
				body: []byte{},
			},
			want: Want{
				body: "",
			},
		},
		{
			name: "status_code_error",
			fields: Fields{
				statusCode: http.StatusInternalServerError,
				body:       []byte("payload"),
			},
			want: Want{
				body: "",
				err:  "received status code: 500",
			},
		},
		{
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

			e := New(&Config{})
			r, err := e.client.sendRequest(server.URL, path, method, body)
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