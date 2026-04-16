// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package parser

import "testing"

func TestParseCounters(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		wantPackets uint64
		wantBytes   uint64
		wantOK      bool
	}{
		{name: "valid counters", token: "[12:34]", wantPackets: 12, wantBytes: 34, wantOK: true},
		{name: "valid counters with spaces", token: "[ 12 : 34 ]", wantPackets: 12, wantBytes: 34, wantOK: true},
		{name: "empty token", token: "", wantOK: false},
		{name: "invalid packets", token: "[abc:34]", wantOK: false},
		{name: "invalid bytes", token: "[12:xyz]", wantOK: false},
		{name: "missing separator", token: "[1234]", wantOK: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			packets, bytes, ok := parseCounters(test.token)
			if ok != test.wantOK {
				t.Fatalf("unexpected ok value: got %v, want %v", ok, test.wantOK)
			}
			if packets != test.wantPackets || bytes != test.wantBytes {
				t.Fatalf("unexpected counters: got (%d, %d), want (%d, %d)", packets, bytes, test.wantPackets, test.wantBytes)
			}
		})
	}
}