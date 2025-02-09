package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var testFiles = map[string]string{
	"full":      "./testdata/input.txt",
	"0_0":       "./testdata/out_offset0_limit0.txt",
	"0_10":      "./testdata/out_offset0_limit10.txt",
	"0_1000":    "./testdata/out_offset0_limit1000.txt",
	"100_1000":  "./testdata/out_offset100_limit1000.txt",
	"6000_1000": "./testdata/out_offset6000_limit1000.txt",
}

func TestCopy(t *testing.T) {
	t.Run("Unsupported file", func(t *testing.T) {
		err := Copy("/dev/urandom", "/tmp/out.txt", 0, 0)
		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("Offset exceeds file size", func(t *testing.T) {
		err := Copy(testFiles["full"], "/tmp/out.txt", 10000, 0)
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})

	t.Run("offset:0 limit:0", func(t *testing.T) {
		runTest(t, 0, 0)
	})

	t.Run("offset:0 limit:10", func(t *testing.T) {
		runTest(t, 0, 10)
	})

	t.Run("offset:0 limit:1000", func(t *testing.T) {
		runTest(t, 0, 1000)
	})

	t.Run("offset:100 limit:1000", func(t *testing.T) {
		runTest(t, 100, 1000)
	})

	t.Run("offset:6000 limit:1000", func(t *testing.T) {
		runTest(t, 6000, 1000)
	})
}

func runTest(t *testing.T, offset, limit int64) {
	t.Helper()
	file, err := os.CreateTemp("/tmp", "out.txt")
	require.NoError(t, err)

	defer func() {
		err := file.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	srcFileKey := fmt.Sprintf("%d_%d", offset, limit)

	err = Copy(testFiles["full"], file.Name(), offset, limit)
	require.NoError(t, err)

	srcInfo, err := os.Stat(testFiles[srcFileKey])
	if err != nil {
		t.Fatal(err)
	}

	destInfo, err := os.Stat(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, srcInfo.Size(), destInfo.Size())
}
