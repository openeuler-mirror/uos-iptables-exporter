// main.go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	Name    = "iptables_exporter"
	Version = "1.0.0"
)

func main() {
	// 为了兼容 kingpin 和标准库 flag，先让 flag 跳过已经被 kingpin 解析的参数。
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	flag.CommandLine.Parse([]string{})

	err := Run(Name, Version)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}



