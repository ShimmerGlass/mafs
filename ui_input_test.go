package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUIInputEmpty(t *testing.T) {
	i := &uiInput{}
	i.Add('a')
	require.Equal(t, "a", i.Text())
	require.Equal(t, 1, i.Cursor)

	i = &uiInput{}
	i.Del()
	require.Equal(t, "", i.Text())
	require.Equal(t, 0, i.Cursor)

	i = &uiInput{}
	i.DelFwd()
	require.Equal(t, "", i.Text())
	require.Equal(t, 0, i.Cursor)

	i = &uiInput{}
	i.MoveCursor(-1)
	require.Equal(t, "", i.Text())
	require.Equal(t, 0, i.Cursor)

	i = &uiInput{}
	i.MoveCursor(1)
	require.Equal(t, "", i.Text())
	require.Equal(t, 0, i.Cursor)
}

func TestUIInputOne(t *testing.T) {
	i := &uiInput{}
	i.Add('a')
	i.Add('b')
	require.Equal(t, "ab", i.Text())
	require.Equal(t, 2, i.Cursor)

	i = &uiInput{}
	i.Add('a')
	i.Del()
	require.Equal(t, "", i.Text())
	require.Equal(t, 0, i.Cursor)

	i = &uiInput{}
	i.Add('a')
	i.DelFwd()
	require.Equal(t, "a", i.Text())
	require.Equal(t, 1, i.Cursor)

	i = &uiInput{}
	i.Add('a')
	i.MoveCursor(-1)
	i.DelFwd()
	require.Equal(t, "", i.Text())
	require.Equal(t, 0, i.Cursor)

	i = &uiInput{}
	i.Add('a')
	i.MoveCursor(-1)
	i.Add('b')
	require.Equal(t, "ba", i.Text())
	require.Equal(t, 1, i.Cursor)
}
