// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

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

	if !strings.EqualFold(u.Scheme, "http") && !strings.EqualFold(u.Scheme, "https") {
		return url.InvalidHostError("scheme must be http or https")
	}

	return nil
}
