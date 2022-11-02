package api

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetFormulaApi(t *testing.T) {
	user, err := randomUser()
	require.NoError(t, err)

	require.NotEmpty(t, user)

	// TODO kerok - implement tests for creating, updating and deleting Formulas!
}
