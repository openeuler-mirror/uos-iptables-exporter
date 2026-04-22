// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestFailWithError(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	request := NewRequest(recorder, req)
	request.Error = errors.New("rate limited")

	request.Fail(http.StatusTooManyRequests)

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("unexpected status code: got %d, want %d", recorder.Code, http.StatusTooManyRequests)
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Fatalf("unexpected content type: got %q", got)
	}
	if got := recorder.Body.String(); got != "rate limited" {
		t.Fatalf("unexpected response body: got %q", got)
	}
}

func TestRequestFailWithoutErrorUsesStatusText(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	request := NewRequest(recorder, req)

	request.Fail(http.StatusInternalServerError)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: got %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	if got := recorder.Body.String(); got != http.StatusText(http.StatusInternalServerError) {
		t.Fatalf("unexpected response body: got %q", got)
	}
}

func TestRequestFailHeadHasNoBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodHead, "/metrics", nil)
	request := NewRequest(recorder, req)
	request.Error = errors.New("rate limited")

	request.Fail(http.StatusTooManyRequests)

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("unexpected status code: got %d, want %d", recorder.Code, http.StatusTooManyRequests)
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Fatalf("unexpected content type: got %q", got)
	}
	if recorder.Body.Len() != 0 {
		t.Fatalf("expected empty body for HEAD request, got %q", recorder.Body.String())
	}
}