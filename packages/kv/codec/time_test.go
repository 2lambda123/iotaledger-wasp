package codec

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestZeroTimeEncoding(t *testing.T) {
	z := time.Time{}
	require.True(t, z.IsZero())
	bin0 := Time.Encode(z)
	zback, err := Time.Decode(bin0)
	require.NoError(t, err)
	require.True(t, zback.IsZero())
	require.True(t, zback.Equal(z))
	require.True(t, zback.IsZero())
}
