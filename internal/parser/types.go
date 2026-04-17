// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package parser

type Tables map[string]*Table

type Table struct {
	Name   string
	Chains map[string]Chain
}

type Chain struct {
	Name    string
	Policy  string
	Packets uint64
	Bytes   uint64
	Rules   []Rule
}

type Rule struct {
	RuleSpec    string // 规则的完整文本（去掉前缀计数器），例如 "-A INPUT -i ens3 -j ACCEPT"
	Packets     uint64
	Bytes       uint64
	InIface     string
	OutIface    string
	Source      string
	Destination string
	Target      string
	Protocol    string
	Index       int
}

type parser struct {
	result           Tables
	currentTableName string
	currentTable     *Table
	line             int
	errors           []error
}

type ruleParser struct {
	packets       uint64
	bytes         uint64
	countersOk    bool
	current       string
	currentValues []string
	chain         string
	flags         []string
	inIface       string
	outIface      string
	source        string
	destination   string
	target        string
	protocol      string
}

type ParseError struct {
	Message    string
	LineNumber int
	LineText   string
	Table      string
	Chain      string
}
