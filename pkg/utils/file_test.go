// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "exists.txt")
	missingFile := filepath.Join(tempDir, "missing.txt")

	if err := os.WriteFile(existingFile, []byte("ok"), 0o644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if !FileExists(existingFile) {
		t.Fatalf("expected existing file to be reported as present")
	}

	if FileExists(missingFile) {
		t.Fatalf("expected missing file to be reported as absent")
	}
}