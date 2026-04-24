// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFaviconServeHTTPUsesHandlerBody(t *testing.T) {
	icon := &favicon{body: []byte{0x01, 0x02, 0x03}}
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)

	icon.ServeHTTP(recorder, req)

	if got := recorder.Header().Get("Content-Type"); got != "image/x-icon" {
		t.Fatalf("unexpected content type: got %q", got)
	}
	if got := recorder.Body.Bytes(); string(got) != string(icon.body) {
		t.Fatalf("unexpected favicon body: got %v, want %v", got, icon.body)
	}
}

func TestFaviconServeHTTPHeadHasNoBody(t *testing.T) {
	icon := &favicon{body: []byte{0x01, 0x02, 0x03}}
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodHead, "/favicon.ico", nil)

	icon.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != "image/x-icon" {
		t.Fatalf("unexpected content type: got %q", got)
	}
	if recorder.Body.Len() != 0 {
		t.Fatalf("expected empty body for HEAD request, got %q", recorder.Body.String())
	}
}

func TestFaviconServeHTTPRejectsPost(t *testing.T) {
	icon := &favicon{body: []byte{0x01, 0x02, 0x03}}
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/favicon.ico", nil)

	icon.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("unexpected status code: got %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
	if got := recorder.Header().Get("Allow"); got != "GET, HEAD" {
		t.Fatalf("unexpected allow header: got %q", got)
	}
}