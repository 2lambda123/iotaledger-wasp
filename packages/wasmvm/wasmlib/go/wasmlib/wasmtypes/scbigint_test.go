package wasmtypes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DivMod(t *testing.T) {
	dividend, divisor, quotient, remainder := "0", "1", "0", "0"
	quo, rem := BigIntFromString(dividend).DivMod(BigIntFromString(divisor))
	require.EqualValues(t, quotient, quo.String())
	require.EqualValues(t, remainder, rem.String())
	dividend, divisor, quotient, remainder = "1", "1", "1", "0"
	quo, rem = BigIntFromString(dividend).DivMod(BigIntFromString(divisor))
	require.EqualValues(t, quotient, quo.String())
	require.EqualValues(t, remainder, rem.String())
	dividend, divisor, quotient, remainder = "123456789012345678901234567", "63531", "1943252727209483227105", "26812"
	quo, rem = BigIntFromString(dividend).DivMod(BigIntFromString(divisor))
	require.EqualValues(t, quotient, quo.String())
	require.EqualValues(t, remainder, rem.String())
}

func BenchmarkDivMod(b *testing.B) {
	dividend := BigIntFromString("123456789012345678901234567")
	divisor := BigIntFromString("63531")
	for i := 0; i < b.N; i++ {
		dividend.DivMod(divisor)
	}
}

func Test_divModEstimate(t *testing.T) {
	dividend := "1481481481481481481474074074074074074074"
	divisor := "44444444444444444444"
	quotient := "33333333333333333333"
	remainder := "22222222222222222222"
	quo, rem := BigIntFromString(dividend).divModEstimate(BigIntFromString(divisor))
	require.EqualValues(t, quotient, quo.String())
	require.EqualValues(t, remainder, rem.String())
}
