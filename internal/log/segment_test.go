package log

import (
	"io"
	"os"
	"testing"

	api "github.com/justagabriel/proglog/api/v1"
	"github.com/justagabriel/proglog/internal"
	"github.com/stretchr/testify/require"
)

func TestSegment(t *testing.T) {
	dir := internal.GetTempDir(t, "segment-test")
	defer os.RemoveAll(dir)

	want := &api.Record{
		Value: []byte("hello world"),
	}

	configuredOffset := uint64(16)

	c := Config{}
	c.Segment.InitialOffset = 0
	c.Segment.MaxStoreBytes = 1024
	c.Segment.MaxIndexBytes = entWidth * 3

	s, err := newSegment(dir, configuredOffset, c)
	require.NoError(t, err)
	require.Equal(t, configuredOffset, s.nextOffset, s.nextOffset)
	require.False(t, s.IsMaxed())

	for i := uint64(0); i < 3; i++ {
		off, err := s.Append(want)
		require.NoError(t, err)
		require.Equal(t, configuredOffset+i, off)

		got, err := s.Read(off)
		require.NoError(t, err)
		require.Equal(t, want.Value, got.Value)
	}

	_, err = s.Append(want)
	require.Equal(t, err, io.EOF)

	// maxed index
	require.True(t, s.IsMaxed())
	c.Segment.MaxStoreBytes = uint64(len(want.Value) * 3)
	c.Segment.MaxIndexBytes = 1024

	s, err = newSegment(dir, configuredOffset, c)
	require.NoError(t, err)

	// maxed store
	require.True(t, s.IsMaxed())
	err = s.Remove()
	require.NoError(t, err)
	s, err = newSegment(dir, configuredOffset, c)
	require.NoError(t, err)
	require.False(t, s.IsMaxed())
}
