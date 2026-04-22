// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func HandleSignals(callback func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go watchSignals(sigChan, callback, func() {
		signal.Stop(sigChan)
		close(sigChan)
	})
}

func watchSignals(sigChan <-chan os.Signal, callback func(), cleanup func()) {
	sig, ok := <-sigChan
	if !ok {
		return
	}

	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		if cleanup != nil {
			cleanup()
		}
		callback()
	}
}
