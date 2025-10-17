package utils

import (
	"net/url"
	"strings"
)

func ValidateURI(uri string) error {
	if uri == "" {
		return nil
	}

	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	if u.Scheme == "" {
		return url.InvalidHostError("missing scheme")
	}

	if !strings.HasPrefix(u.Scheme, "http") {
		return url.InvalidHostError("scheme must be http or https")
	}

	return nil
}
