package rule

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/mikelorant/easyredir/pkg/easyredir/option"
	"github.com/mikelorant/easyredir/pkg/jsonutil"
)

func CreateRule(cl ClientAPI, attr Attributes, opts ...option.Option) (r Rule, err error) {
	r = Rule{}

	o := &option.Options{}
	for _, opt := range opts {
		opt.Apply(o)
	}

	var b bytes.Buffer
	if err := jsonutil.EncodeJSON(&attr, &b); err != nil {
		return r, fmt.Errorf("unable to encode to json: %w", err)
	}

	pathQuery := buildCreateRule(o)
	reader, err := cl.SendRequest(pathQuery, http.MethodPost, &b)
	if err != nil {
		return r, fmt.Errorf("unable to send request: %w", err)
	}

	if err := jsonutil.DecodeJSON(reader, &r); err != nil {
		return r, fmt.Errorf("unable to get json: %w", err)
	}

	return r, nil
}

func buildCreateRule(opts *option.Options) string {
	var sb strings.Builder
	var params []string

	fmt.Fprint(&sb, "/rules")

	if opts.Include != "" {
		params = append(params, fmt.Sprintf("include[]=%v", opts.Include))
	}

	if len(params) != 0 {
		fmt.Fprintf(&sb, "?%v", strings.Join(params, "&"))
	}

	return sb.String()
}
