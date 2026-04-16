// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package parser

import "testing"

func TestRuleParserCapturesProtocol(t *testing.T) {
	var rp ruleParser
	for _, token := range []string{"-A", "INPUT", "-p", "tcp"} {
		rp.handleToken(token)
	}
	rp.flush()

	if rp.chain != "INPUT" {
		t.Fatalf("unexpected chain: got %q", rp.chain)
	}
	if rp.protocol != "tcp" {
		t.Fatalf("unexpected protocol: got %q", rp.protocol)
	}
}