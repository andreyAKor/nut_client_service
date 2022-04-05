package server

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestClose(t *testing.T) {
	t.Run("server not init", func(t *testing.T) {
		srv, err := New(nil, nil, nil, "", 0, 0)
		require.NoError(t, err)

		err = srv.Close()
		require.Equal(t, err, errors.Cause(ErrServerNotInit))
	})
}
