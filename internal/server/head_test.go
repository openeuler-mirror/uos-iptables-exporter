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