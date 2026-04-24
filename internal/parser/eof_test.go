// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package parser

import (
	"strings"
	"testing"
)

func TestParseIptablesSaveHandlesFinalLineWithoutNewline(t *testing.T) {
	input := strings.NewReader("*filter\n:INPUT ACCEPT [0:0]\nCOMMIT")

	tables, err := ParseIptablesSave(input)
	if err != nil {
		t.Fatalf("ParseIptablesSave returned error: %v", err)
	}

	table, ok := tables["filter"]
	if !ok {
		t.Fatalf("expected filter table to be parsed")
	}
	if _, ok := table.Chains["INPUT"]; !ok {
		t.Fatalf("expected INPUT chain to be parsed from final line without newline")
	}
}

func TestParseIptablesSaveRejectsMissingCommit(t *testing.T) {
	input := strings.NewReader("*filter\n:INPUT ACCEPT [0:0]\n")

	_, err := ParseIptablesSave(input)
	if err == nil {
		t.Fatal("expected parse error for missing COMMIT, got nil")
	}
	if !strings.Contains(err.Error(), "missing COMMIT") {
		t.Fatalf("expected missing COMMIT error, got %q", err.Error())
	}
}