package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func cmpFiles(file1, file2 string) bool {
	f1, _ := os.Open(file1)
	defer f1.Close()
	f2, _ := os.Open(file2)
	defer f2.Close()

	data1, _ := io.ReadAll(f1)
	data2, _ := io.ReadAll(f2)
	if len(data1) != len(data2) {
		return false
	}

	for i := 0; i < len(data1); i++ {
		if data1[i] != data2[i] {
			return false
		}
	}

	return true
}

func TestCopy(t *testing.T) {
	outfile, err := os.CreateTemp("", "out_test")
	outfile.Close()
	defer os.Remove(outfile.Name())

	t.Run("offset 0 limit 0", func(t *testing.T) {
		require.NoError(t, err)
		require.NoError(t, Copy("./testdata/input.txt", outfile.Name(), 0, 0))
		require.FileExists(t, outfile.Name())
		require.True(t, cmpFiles("./testdata/out_offset0_limit0.txt", outfile.Name()))
	})

	t.Run("offset 0 limit 10", func(t *testing.T) {
		require.NoError(t, err)
		require.NoError(t, Copy("./testdata/input.txt", outfile.Name(), 0, 10))
		require.FileExists(t, outfile.Name())
		require.True(t, cmpFiles("./testdata/out_offset0_limit10.txt", outfile.Name()))
	})

	t.Run("offset 0 limit 1000", func(t *testing.T) {
		require.NoError(t, err)
		require.NoError(t, Copy("./testdata/input.txt", outfile.Name(), 0, 1000))
		require.FileExists(t, outfile.Name())
		require.True(t, cmpFiles("./testdata/out_offset0_limit1000.txt", outfile.Name()))
	})

	t.Run("offset 0 limit 10000", func(t *testing.T) {
		require.NoError(t, err)
		require.NoError(t, Copy("./testdata/input.txt", outfile.Name(), 0, 10000))
		require.FileExists(t, outfile.Name())
		require.True(t, cmpFiles("./testdata/out_offset0_limit10000.txt", outfile.Name()))
	})

	t.Run("offset 100 limit 1000", func(t *testing.T) {
		require.NoError(t, err)
		require.NoError(t, Copy("./testdata/input.txt", outfile.Name(), 100, 1000))
		require.FileExists(t, outfile.Name())
		require.True(t, cmpFiles("./testdata/out_offset100_limit1000.txt", outfile.Name()))
	})

	t.Run("offset 6000 limit 1000", func(t *testing.T) {
		require.NoError(t, err)
		require.NoError(t, Copy("./testdata/input.txt", outfile.Name(), 6000, 1000))
		require.FileExists(t, outfile.Name())
		require.True(t, cmpFiles("./testdata/out_offset6000_limit1000.txt", outfile.Name()))
	})
}
