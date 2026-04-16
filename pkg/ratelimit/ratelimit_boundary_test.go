// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package ratelimit

import (
	"testing"
	"time"
)

func TestNewRateLimiterRejectsInvalidArguments(t *testing.T) {
	if _, err := NewRateLimiter(time.Second, 0); err != ErrRateLimitSize {
		t.Fatalf("unexpected error for invalid channel size: got %v, want %v", err, ErrRateLimitSize)
	}
	if _, err := NewRateLimiter(0, 1); err != ErrRateLimitTime {
		t.Fatalf("unexpected error for invalid duration: got %v, want %v", err, ErrRateLimitTime)
	}
}

func TestRateLimiterGetReturnsRateLimitedAfterTokensExhausted(t *testing.T) {
	rl, err := NewRateLimiter(time.Hour, 1)
	if err != nil {
		t.Fatalf("NewRateLimiter returned error: %v", err)
	}
	defer rl.Stop()

	if err := rl.Get(); err != nil {
		t.Fatalf("first Get returned error: %v", err)
	}
	if err := rl.Get(); err != ErrRateLimited {
		t.Fatalf("unexpected error after tokens exhausted: got %v, want %v", err, ErrRateLimited)
	}
}