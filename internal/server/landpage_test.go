// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewLandingPageUsesCustomCSS(t *testing.T) {
	handler, err := NewLandingPage(LandingPageConfig{
		Name:    "iptables-exporter",
		Version: "1.0.0",
		CSS:     "body { color: red; }",
		Links: []LandingPageLinks{{
			Text:    "Metrics",
			Address: "/metrics",
		}},
	})
	if err != nil {
		t.Fatalf("NewLandingPage returned error: %v", err)
	}

	body := string(handler.landingPage)
	if !strings.Contains(body, "body { color: red; }") {
		t.Fatalf("landing page did not embed custom CSS")
	}
	if !strings.Contains(body, "iptables-exporter") {
		t.Fatalf("landing page did not render service name")
	}
}

func TestLandingPageHandlerReturns404ForNonRoot(t *testing.T) {
	handler, err := NewLandingPage(LandingPageConfig{Name: "iptables-exporter"})
	if err != nil {
		t.Fatalf("NewLandingPage returned error: %v", err)
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/not-found", nil)

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("unexpected status code: got %d, want %d", recorder.Code, http.StatusNotFound)
	}
}