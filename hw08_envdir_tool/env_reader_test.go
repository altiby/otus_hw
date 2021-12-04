package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	ev, err := ReadDir("./testdata/env")
	t.Run("read dir", func(t *testing.T) {
		require.NoError(t, err)
		require.Len(t, ev, 5)
	})

	t.Run("bar case", func(t *testing.T) {
		require.Equal(t, "bar", ev["BAR"].Value)
		require.Equal(t, false, ev["BAR"].NeedRemove)
	})

	t.Run("empty case", func(t *testing.T) {
		require.Equal(t, "", ev["EMPTY"].Value)
		require.Equal(t, true, ev["EMPTY"].NeedRemove)
	})

	t.Run("foo case", func(t *testing.T) {
		require.Equal(t, "   foo\nwith new line", ev["FOO"].Value)
		require.Equal(t, false, ev["FOO"].NeedRemove)
	})

	t.Run("hello case", func(t *testing.T) {
		require.Equal(t, "\"hello\"", ev["HELLO"].Value)
		require.Equal(t, false, ev["HELLO"].NeedRemove)
	})

	t.Run("unset case", func(t *testing.T) {
		require.Equal(t, "", ev["UNSET"].Value)
		require.Equal(t, true, ev["UNSET"].NeedRemove)
	})
}
