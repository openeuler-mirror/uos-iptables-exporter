// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package utils

import (
	"os"
	"syscall"
	"testing"
)

func TestWatchSignalsCallsCallbackAndCleanup(t *testing.T) {
	sigChan := make(chan os.Signal, 2)
	callbackCalls := 0
	cleanupCalls := 0

	sigChan <- syscall.SIGTERM
	sigChan <- syscall.SIGINT

	watchSignals(sigChan, func() {
		callbackCalls++
	}, func() {
		cleanupCalls++
	})

	if callbackCalls != 1 {
		t.Fatalf("expected callback to be called once, got %d", callbackCalls)
	}
	if cleanupCalls != 1 {
		t.Fatalf("expected cleanup to be called once, got %d", cleanupCalls)
	}
	if len(sigChan) != 1 {
		t.Fatalf("expected one unread signal to remain, got %d", len(sigChan))
	}
}

func TestWatchSignalsReturnsOnClosedChannel(t *testing.T) {
	sigChan := make(chan os.Signal)
	close(sigChan)
	callbackCalls := 0

	watchSignals(sigChan, func() {
		callbackCalls++
	}, nil)

	if callbackCalls != 0 {
		t.Fatalf("expected callback not to be called, got %d", callbackCalls)
	}
}