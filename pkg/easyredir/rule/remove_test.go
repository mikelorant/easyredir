package rule

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/mikelorant/easyredir-cli/pkg/easyredir/client"
	"github.com/stretchr/testify/assert"
)

func TestRemoveRule(t *testing.T) {
	type Args struct {
		id string
	}

	type Fields struct {
		status int
		body   string
	}

	type Want struct {
		result bool
		err    string
	}

	tests := []struct {
		name   string
		args   Args
		fields Fields
		want   Want
	}{
		{
			name: "success",
			args: Args{
				id: "5d29f819-302f-40c0-8b5a-6d670267615b",
			},
			fields: Fields{
				status: http.StatusNoContent,
			},
			want: Want{
				result: true,
			},
		},
		{
			name: "failure",
			args: Args{
				id: "5d29f819-302f-40c0-8b5a-6d670267615b",
			},
			fields: Fields{
				status: http.StatusNotFound,
				body: `
					{
						"type": "record_not_found_error",
						"message": "Record not found"
					}
				`,
			},
			want: Want{
				result: true,
				err:    "record_not_found_error: Record not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/rules/", func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(tt.fields.status)
				w.Write([]byte(tt.fields.body))
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			cl := client.New(WithBaseURL(server.URL))

			got, err := RemoveRule(cl, tt.args.id)
			if tt.want.err != "" {
				assert.NotNil(t, err)
				td.CmpContains(t, err, tt.want.err)
				return
			}
			assert.Nil(t, err)
			td.Cmp(t, got, tt.want.result)
		})
	}
}
