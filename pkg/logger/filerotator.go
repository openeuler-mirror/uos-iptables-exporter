// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT

package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileRotator struct {
	filePath    string
	maxSize     int64
	maxAge      time.Duration
	currentFile *os.File
	currentSize int64
	mu          sync.Mutex
}

func NewFileRotator(filePath string, maxSize int64, maxAge time.Duration) *FileRotator {
	return &FileRotator{
		filePath: filePath,
		maxSize:  maxSize,
		maxAge:   maxAge,
	}
}

func (fr *FileRotator) Write(p []byte) (n int, err error) {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	if fr.currentFile == nil {
		if err := fr.openNewFile(); err != nil {
			return 0, err
		}
	}

	if fr.currentSize >= fr.maxSize {
		if err := fr.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = fr.currentFile.Write(p)
	fr.currentSize += int64(n)
	return n, err
}

func (fr *FileRotator) openNewFile() error {
	dir := filepath.Dir(fr.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(fr.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	fr.currentFile = file
	info, err := file.Stat()
	if err != nil {
		return err
	}
	fr.currentSize = info.Size()
	return nil
}

func (fr *FileRotator) rotate() error {
	if fr.currentFile != nil {
		fr.currentFile.Close()
	}

	timestamp := time.Now().Format("20060102-150405")
	archivePath := fmt.Sprintf("%s.%s.gz", fr.filePath, timestamp)

	if err := fr.compressFile(fr.filePath, archivePath); err != nil {
		return err
	}

	os.Remove(fr.filePath)
	return fr.openNewFile()
}

func (fr *FileRotator) compressFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzipWriter := gzip.NewWriter(dstFile)
	defer gzipWriter.Close()

	_, err = io.Copy(gzipWriter, srcFile)
	return err
}

func (fr *FileRotator) Close() error {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	if fr.currentFile != nil {
		err := fr.currentFile.Close()
		fr.currentFile = nil
		return err
	}
	return nil
}
