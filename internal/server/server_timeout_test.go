// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package server

import "testing"

func TestSetupHTTPServerSetsTimeouts(t *testing.T) {
	server := NewServer("iptables-exporter", "1.0.0")
	server.CommonConfig.MetricsPath = "/metrics"

	if err := server.setupHttpServer(); err != nil {
		t.Fatalf("setupHttpServer returned error: %v", err)
	}
	defer func() {
		if server.scrapeCancel != nil {
			server.scrapeCancel()
		}
	}()

	if server.server == nil {
		t.Fatal("expected http server to be initialized")
	}
	if server.server.ReadHeaderTimeout != defaultReadHeaderTimeout {
		t.Fatalf("unexpected ReadHeaderTimeout: got %s, want %s", server.server.ReadHeaderTimeout, defaultReadHeaderTimeout)
	}
	if server.server.IdleTimeout != defaultIdleTimeout {
		t.Fatalf("unexpected IdleTimeout: got %s, want %s", server.server.IdleTimeout, defaultIdleTimeout)
	}
}
