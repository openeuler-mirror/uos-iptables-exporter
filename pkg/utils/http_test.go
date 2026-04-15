// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package utils

import "testing"

func TestValidateURI(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
	}{
		{name: "empty uri", uri: "", wantErr: false},
		{name: "http uri", uri: "http://127.0.0.1:9077/metrics", wantErr: false},
		{name: "https uri", uri: "https://example.com/metrics", wantErr: false},
		{name: "missing scheme", uri: "127.0.0.1:9077/metrics", wantErr: true},
		{name: "unsupported scheme", uri: "ftp://example.com/metrics", wantErr: true},
		{name: "http prefix but invalid scheme", uri: "httpx://example.com/metrics", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateURI(test.uri)
			if (err != nil) != test.wantErr {
				t.Fatalf("ValidateURI(%q) error = %v, wantErr %v", test.uri, err, test.wantErr)
			}
		})
	}
}