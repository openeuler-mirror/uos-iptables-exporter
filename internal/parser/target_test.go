// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package parser

import "testing"

func TestRuleParserCapturesJumpTarget(t *testing.T) {
	var rp ruleParser
	for _, token := range []string{"-A", "INPUT", "-j", "ACCEPT"} {
		rp.handleToken(token)
	}
	rp.flush()

	if rp.chain != "INPUT" {
		t.Fatalf("unexpected chain: got %q", rp.chain)
	}
	if rp.target != "ACCEPT" {
		t.Fatalf("unexpected target: got %q", rp.target)
	}
}