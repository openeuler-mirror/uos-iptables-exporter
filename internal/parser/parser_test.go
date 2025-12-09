// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package parser

import (
	"strings"
	"testing"
)

func TestParseRuleWithInterfacesAndCounters(t *testing.T) {
	// 模拟 iptables-save -c 输出中的一条规则行
	line := "[10:1024] -A INPUT -i ens3 -j ACCEPT"

	// 直接走 ParseIptablesSave 的内部逻辑：先去掉计数器再交给 handleRuleWithCounter
	p := &parser{
		result:           make(Tables),
		currentTableName: "filter",
		currentTable: &Table{
			Name:   "filter",
			Chains: make(map[string]Chain),
		},
	}
	p.currentTable.Chains["INPUT"] = Chain{Name: "INPUT", Rules: []Rule{}}

	// 仿照 handleLine 中的逻辑
	line = strings.TrimSpace(line)
	if line == "" {
		t.Fatal("test line is empty")
	}

	// 去掉前缀计数器部分，构造 ruleLine
	fields := strings.Fields(line)
	if len(fields) < 3 {
		t.Fatalf("unexpected fields length: %d", len(fields))
	}
	// fields[0] 是计数器，[1] 应该是 -A
	ruleLine := strings.Join(fields[1:], " ")

	p.handleRuleWithCounter(ruleLine, 10, 1024)

	tables := p.currentTable.Chains
	chain, ok := tables["INPUT"]
	if !ok {
		t.Fatalf("chain INPUT not found")
	}
	if len(chain.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(chain.Rules))
	}

	r := chain.Rules[0]
	if r.InIface != "ens3" {
		t.Errorf("expected InIface=ens3, got %q", r.InIface)
	}
	if r.OutIface != "" {
		t.Errorf("expected OutIface empty, got %q", r.OutIface)
	}
	if r.Packets != 10 || r.Bytes != 1024 {
		t.Errorf("unexpected counters: packets=%d bytes=%d", r.Packets, r.Bytes)
	}
}

func TestParseRuleWithOutInterface(t *testing.T) {
	line := "[5:512] -A OUTPUT -o eth1 -j ACCEPT"

	p := &parser{
		result:           make(Tables),
		currentTableName: "filter",
		currentTable: &Table{
			Name:   "filter",
			Chains: make(map[string]Chain),
		},
	}
	p.currentTable.Chains["OUTPUT"] = Chain{Name: "OUTPUT", Rules: []Rule{}}

	line = strings.TrimSpace(line)
	fields := strings.Fields(line)
	if len(fields) < 3 {
		t.Fatalf("unexpected fields length: %d", len(fields))
	}
	ruleLine := strings.Join(fields[1:], " ")

	p.handleRuleWithCounter(ruleLine, 5, 512)

	chain, ok := p.currentTable.Chains["OUTPUT"]
	if !ok {
		t.Fatalf("chain OUTPUT not found")
	}
	if len(chain.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(chain.Rules))
	}

	r := chain.Rules[0]
	if r.InIface != "" {
		t.Errorf("expected InIface empty, got %q", r.InIface)
	}
	if r.OutIface != "eth1" {
		t.Errorf("expected OutIface=eth1, got %q", r.OutIface)
	}
	if r.Packets != 5 || r.Bytes != 512 {
		t.Errorf("unexpected counters: packets=%d bytes=%d", r.Packets, r.Bytes)
	}
}
