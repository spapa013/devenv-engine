package templates

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunPostRender_NoOp(t *testing.T) {
	err := RunPostRender(t.TempDir(), NewPostRenderOptions())
	require.NoError(t, err)
}
