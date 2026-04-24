// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileRotatorCloseIsIdempotent(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "exporter.log")
	rotator := NewFileRotator(filePath, 1024, time.Hour)

	if _, err := rotator.Write([]byte("first line")); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if err := rotator.Close(); err != nil {
		t.Fatalf("first Close returned error: %v", err)
	}
	if err := rotator.Close(); err != nil {
		t.Fatalf("second Close returned error: %v", err)
	}
}

func TestFileRotatorWriteReopensAfterClose(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "exporter.log")
	rotator := NewFileRotator(filePath, 1024, time.Hour)

	if _, err := rotator.Write([]byte("first line\n")); err != nil {
		t.Fatalf("first Write returned error: %v", err)
	}
	if err := rotator.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	if _, err := rotator.Write([]byte("second line\n")); err != nil {
		t.Fatalf("second Write returned error: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(content) != "first line\nsecond line\n" {
		t.Fatalf("unexpected file content: got %q", string(content))
	}
}