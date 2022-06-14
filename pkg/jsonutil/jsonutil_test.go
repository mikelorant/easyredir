package jsonutil

import (
	"io"
	"strings"
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/assert"
)

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
			err := DecodeJSON(io.NopCloser(tt.args.src), &got)
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
