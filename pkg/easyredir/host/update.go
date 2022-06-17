package host

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/mikelorant/easyredir/pkg/easyredir/option"
	"github.com/mikelorant/easyredir/pkg/jsonutil"
)

func UpdateHost(cl ClientAPI, id string, attr Attributes, opts ...option.Option) (h Host, err error) {
	h = Host{}

	o := &option.Options{}
	for _, opt := range opts {
		opt.Apply(o)
	}

	var b bytes.Buffer
	if err := jsonutil.EncodeJSON(&attr, &b); err != nil {
		return h, fmt.Errorf("unable to encode to json: %w", err)
	}

	pathQuery := buildUpdateHost(id)
	reader, err := cl.SendRequest(pathQuery, http.MethodPatch, &b)
	if err != nil {
		return h, fmt.Errorf("unable to send request: %w", err)
	}

	if err := jsonutil.DecodeJSON(reader, &h); err != nil {
		return h, fmt.Errorf("unable to get json: %w", err)
	}

	return h, nil
}

func buildUpdateHost(id string) string {
	return fmt.Sprintf("/hosts/%v", id)
}
