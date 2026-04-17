// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthzHandlerAllowsGet(t *testing.T) {
	server := NewServer("iptables-exporter", "1.0.0")
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	server.healthzHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("unexpected content type: got %q", got)
	}
	if !strings.Contains(recorder.Body.String(), "iptables-exporter is running normally") {
		t.Fatalf("unexpected healthz body: %q", recorder.Body.String())
	}
}

func TestHealthzHandlerRejectsPost(t *testing.T) {
	server := NewServer("iptables-exporter", "1.0.0")
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)

	server.healthzHandler(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("unexpected status code: got %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
	if got := recorder.Header().Get("Allow"); got != "GET, HEAD" {
		t.Fatalf("unexpected allow header: got %q", got)
	}
	if !strings.Contains(recorder.Body.String(), http.StatusText(http.StatusMethodNotAllowed)) {
		t.Fatalf("unexpected response body: %q", recorder.Body.String())
	}
}
