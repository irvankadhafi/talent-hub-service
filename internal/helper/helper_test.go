package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelper_FormatEmail(t *testing.T) {
	var (
		a = "john.doe@mail.com"
		b = "JOHN.DOE@mail.com"
		c = " john.DOE@mail.com "
	)

	t.Run("success", func(t *testing.T) {
		require.Equal(t, a, FormatEmail(a))
		require.Equal(t, a, FormatEmail(b))
		require.Equal(t, a, FormatEmail(c))
	})
}
