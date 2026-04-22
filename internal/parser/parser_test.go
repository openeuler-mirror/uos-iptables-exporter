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

func TestParseRuleWithSourceAndDestination(t *testing.T) {
	line := "[8:256] -A INPUT -s 192.168.1.10/32 -d 10.0.0.0/24 -j ACCEPT"

	p := &parser{
		result:           make(Tables),
		currentTableName: "filter",
		currentTable: &Table{
			Name:   "filter",
			Chains: make(map[string]Chain),
		},
	}
	p.currentTable.Chains["INPUT"] = Chain{Name: "INPUT", Rules: []Rule{}}

	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) < 3 {
		t.Fatalf("unexpected fields length: %d", len(fields))
	}

	p.handleRuleWithCounter(strings.Join(fields[1:], " "), 8, 256)

	rules := p.currentTable.Chains["INPUT"].Rules
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}

	r := rules[0]
	if r.Source != "192.168.1.10/32" {
		t.Fatalf("expected Source=192.168.1.10/32, got %q", r.Source)
	}
	if r.Destination != "10.0.0.0/24" {
		t.Fatalf("expected Destination=10.0.0.0/24, got %q", r.Destination)
	}
}

func TestParseRuleWithLongAddressFlags(t *testing.T) {
	line := "[3:128] -A OUTPUT --source 2001:db8::1 --destination 2001:db8::2 -j ACCEPT"

	p := &parser{
		result:           make(Tables),
		currentTableName: "filter",
		currentTable: &Table{
			Name:   "filter",
			Chains: make(map[string]Chain),
		},
	}
	p.currentTable.Chains["OUTPUT"] = Chain{Name: "OUTPUT", Rules: []Rule{}}

	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) < 3 {
		t.Fatalf("unexpected fields length: %d", len(fields))
	}

	p.handleRuleWithCounter(strings.Join(fields[1:], " "), 3, 128)

	rules := p.currentTable.Chains["OUTPUT"].Rules
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}

	r := rules[0]
	if r.Source != "2001:db8::1" {
		t.Fatalf("expected Source=2001:db8::1, got %q", r.Source)
	}
	if r.Destination != "2001:db8::2" {
		t.Fatalf("expected Destination=2001:db8::2, got %q", r.Destination)
	}
}

func TestHandleRuleWithCounterRejectsUnknownChain(t *testing.T) {
	p := &parser{
		result:           make(Tables),
		currentTableName: "filter",
		currentTable: &Table{
			Name:   "filter",
			Chains: make(map[string]Chain),
		},
	}

	p.handleRuleWithCounter("-A UNKNOWN -j ACCEPT", 1, 64)

	if len(p.errors) != 1 {
		t.Fatalf("expected 1 parse error, got %d", len(p.errors))
	}
	if got := p.errors[0].Error(); !strings.Contains(got, "unknown chain") {
		t.Fatalf("expected unknown chain error, got %q", got)
	}
	if len(p.currentTable.Chains) != 0 {
		t.Fatalf("expected no chains to be created, got %d", len(p.currentTable.Chains))
	}
}
