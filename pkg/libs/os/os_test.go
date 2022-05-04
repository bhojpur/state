package os_test

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	bos "github.com/bhojpur/state/pkg/libs/os"
)

func TestCopyFile(t *testing.T) {
	tmpfile, err := os.CreateTemp(t.TempDir(), "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	content := []byte("hello world")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}

	copyfile := fmt.Sprintf("%s.copy", tmpfile.Name())
	if err := bos.CopyFile(tmpfile.Name(), copyfile); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(copyfile); os.IsNotExist(err) {
		t.Fatal("copy should exist")
	}
	data, err := os.ReadFile(copyfile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, content) {
		t.Fatalf("copy file content differs: expected %v, got %v", content, data)
	}
	os.Remove(copyfile)
}

func TestEnsureDir(t *testing.T) {
	tmp := t.TempDir()

	// Should be possible to create a new directory.
	err := bos.EnsureDir(filepath.Join(tmp, "dir"), 0755)
	require.NoError(t, err)
	require.DirExists(t, filepath.Join(tmp, "dir"))

	// Should succeed on existing directory.
	err = bos.EnsureDir(filepath.Join(tmp, "dir"), 0755)
	require.NoError(t, err)

	// Should fail on file.
	err = os.WriteFile(filepath.Join(tmp, "file"), []byte{}, 0644)
	require.NoError(t, err)
	err = bos.EnsureDir(filepath.Join(tmp, "file"), 0755)
	require.Error(t, err)

	// Should allow symlink to dir.
	err = os.Symlink(filepath.Join(tmp, "dir"), filepath.Join(tmp, "linkdir"))
	require.NoError(t, err)
	err = bos.EnsureDir(filepath.Join(tmp, "linkdir"), 0755)
	require.NoError(t, err)

	// Should error on symlink to file.
	err = os.Symlink(filepath.Join(tmp, "file"), filepath.Join(tmp, "linkfile"))
	require.NoError(t, err)
	err = bos.EnsureDir(filepath.Join(tmp, "linkfile"), 0755)
	require.Error(t, err)
}

// Ensure that using CopyFile does not truncate the destination file before
// the origin is positively a non-directory and that it is ready for copying.
func TestTrickedTruncation(t *testing.T) {
	tmpDir := t.TempDir()

	originalWALPath := filepath.Join(tmpDir, "wal")
	originalWALContent := []byte("I AM BECOME DEATH, DESTROYER OF ALL WORLDS!")
	if err := os.WriteFile(originalWALPath, originalWALContent, 0755); err != nil {
		t.Fatal(err)
	}

	// 1. Sanity check.
	readWAL, err := os.ReadFile(originalWALPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(readWAL, originalWALContent) {
		t.Fatalf("Cannot proceed as the content does not match\nGot:  %q\nWant: %q", readWAL, originalWALContent)
	}

	// 2. Now cause the truncation of the original file.
	// It is absolutely legal to invoke os.Open on a directory.
	if err := bos.CopyFile(tmpDir, originalWALPath); err == nil {
		t.Fatal("Expected an error")
	}

	// 3. Check the WAL's content
	reReadWAL, err := os.ReadFile(originalWALPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(reReadWAL, originalWALContent) {
		t.Fatalf("Oops, the WAL's content was changed :(\nGot:  %q\nWant: %q", reReadWAL, originalWALContent)
	}
}
