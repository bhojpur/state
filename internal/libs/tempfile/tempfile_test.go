package tempfile

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

// Need access to internal variables, so can't use _test package

import (
	"bytes"
	"fmt"
	mrand "math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	librand "github.com/bhojpur/state/pkg/libs/rand"
)

func TestWriteFileAtomic(t *testing.T) {
	var (
		data             = []byte(librand.Str(mrand.Intn(2048)))
		old              = librand.Bytes(mrand.Intn(2048))
		perm os.FileMode = 0600
	)

	f, err := os.CreateTemp(t.TempDir(), "write-atomic-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	if err = os.WriteFile(f.Name(), old, 0600); err != nil {
		t.Fatal(err)
	}

	if err = WriteFileAtomic(f.Name(), data, perm); err != nil {
		t.Fatal(err)
	}

	rData, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, rData) {
		t.Fatalf("data mismatch: %v != %v", data, rData)
	}

	stat, err := os.Stat(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	if have, want := stat.Mode().Perm(), perm; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

// This tests atomic write file when there is a single duplicate file.
// Expected behavior is for a new file to be created, and the original write file to be unaltered.
func TestWriteFileAtomicDuplicateFile(t *testing.T) {
	var (
		defaultSeed    uint64 = 1
		testString            = "This is a glorious test string"
		expectedString        = "Did the test file's string appear here?"

		fileToWrite = "/tmp/TestWriteFileAtomicDuplicateFile-test.txt"
	)
	// Create a file at the seed, and reset the seed.
	atomicWriteFileRand = defaultSeed
	firstFileRand := randWriteFileSuffix()
	atomicWriteFileRand = defaultSeed
	fname := "/tmp/" + atomicWriteFilePrefix + firstFileRand
	f, err := os.OpenFile(fname, atomicWriteFileFlag, 0777)
	defer os.Remove(fname)
	// Defer here, in case there is a panic in WriteFileAtomic.
	defer os.Remove(fileToWrite)

	require.NoError(t, err)
	_, err = f.WriteString(testString)
	require.NoError(t, err)
	err = WriteFileAtomic(fileToWrite, []byte(expectedString), 0777)
	require.NoError(t, err)
	// Check that the first atomic file was untouched
	firstAtomicFileBytes, err := os.ReadFile(fname)
	require.NoError(t, err, "Error reading first atomic file")
	require.Equal(t, []byte(testString), firstAtomicFileBytes, "First atomic file was overwritten")
	// Check that the resultant file is correct
	resultantFileBytes, err := os.ReadFile(fileToWrite)
	require.NoError(t, err, "Error reading resultant file")
	require.Equal(t, []byte(expectedString), resultantFileBytes, "Written file had incorrect bytes")

	// Check that the intermediate write file was deleted
	// Get the second write files' randomness
	atomicWriteFileRand = defaultSeed
	_ = randWriteFileSuffix()
	secondFileRand := randWriteFileSuffix()
	_, err = os.Stat("/tmp/" + atomicWriteFilePrefix + secondFileRand)
	require.True(t, os.IsNotExist(err), "Intermittent atomic write file not deleted")
}

// This tests atomic write file when there are many duplicate files.
// Expected behavior is for a new file to be created under a completely new seed,
// and the original write files to be unaltered.
func TestWriteFileAtomicManyDuplicates(t *testing.T) {
	var (
		defaultSeed    uint64 = 2
		testString            = "This is a glorious test string, from file %d"
		expectedString        = "Did any of the test file's string appear here?"

		fileToWrite = "/tmp/TestWriteFileAtomicDuplicateFile-test.txt"
	)
	// Initialize all of the atomic write files
	atomicWriteFileRand = defaultSeed
	for i := 0; i < atomicWriteFileMaxNumConflicts+2; i++ {
		fileRand := randWriteFileSuffix()
		fname := "/tmp/" + atomicWriteFilePrefix + fileRand
		f, err := os.OpenFile(fname, atomicWriteFileFlag, 0777)
		require.NoError(t, err)
		_, err = f.WriteString(fmt.Sprintf(testString, i))
		require.NoError(t, err)
		defer os.Remove(fname)
	}

	atomicWriteFileRand = defaultSeed
	// Defer here, in case there is a panic in WriteFileAtomic.
	defer os.Remove(fileToWrite)

	err := WriteFileAtomic(fileToWrite, []byte(expectedString), 0777)
	require.NoError(t, err)
	// Check that all intermittent atomic file were untouched
	atomicWriteFileRand = defaultSeed
	for i := 0; i < atomicWriteFileMaxNumConflicts+2; i++ {
		fileRand := randWriteFileSuffix()
		fname := "/tmp/" + atomicWriteFilePrefix + fileRand
		firstAtomicFileBytes, err := os.ReadFile(fname)
		require.NoError(t, err, "Error reading first atomic file")
		require.Equal(t, []byte(fmt.Sprintf(testString, i)), firstAtomicFileBytes,
			"atomic write file %d was overwritten", i)
	}

	// Check that the resultant file is correct
	resultantFileBytes, err := os.ReadFile(fileToWrite)
	require.NoError(t, err, "Error reading resultant file")
	require.Equal(t, []byte(expectedString), resultantFileBytes, "Written file had incorrect bytes")
}
