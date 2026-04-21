// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      logrus.Level
		wantValid bool
	}{
		{name: "debug", input: "debug", want: logrus.DebugLevel, wantValid: true},
		{name: "info", input: "info", want: logrus.InfoLevel, wantValid: true},
		{name: "warn", input: "warn", want: logrus.WarnLevel, wantValid: true},
		{name: "warning alias", input: "warning", want: logrus.WarnLevel, wantValid: true},
		{name: "error", input: "error", want: logrus.ErrorLevel, wantValid: true},
		{name: "unknown", input: "trace", want: logrus.WarnLevel, wantValid: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := parseLevel(test.input)
			if ok != test.wantValid {
				t.Fatalf("parseLevel(%q) valid = %v, want %v", test.input, ok, test.wantValid)
			}
			if got != test.want {
				t.Fatalf("parseLevel(%q) = %v, want %v", test.input, got, test.want)
			}
		})
	}
}