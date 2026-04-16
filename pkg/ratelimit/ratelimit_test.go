// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package ratelimit

import (
	"testing"
	"time"
)

func TestRateLimiterStopIsIdempotent(t *testing.T) {
	rl, err := NewRateLimiter(10*time.Millisecond, 1)
	if err != nil {
		t.Fatalf("NewRateLimiter returned error: %v", err)
	}

	rl.Stop()
	rl.Stop()
}

func TestRateLimiterGetAfterStopReturnsRateLimited(t *testing.T) {
	rl, err := NewRateLimiter(10*time.Millisecond, 1)
	if err != nil {
		t.Fatalf("NewRateLimiter returned error: %v", err)
	}

	rl.Stop()

	if err := rl.Get(); err != ErrRateLimited {
		t.Fatalf("unexpected error after Stop: got %v, want %v", err, ErrRateLimited)
	}
}