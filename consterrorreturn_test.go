package consterrorreturn_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/delarean/consterrorreturn"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLinter(t *testing.T) {
	testdata := analysistest.TestData()

	results := analysistest.Run(t, testdata, consterrorreturn.Analyzer, "a")
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
}

func TestIsErrorType(t *testing.T) {
	require.False(t, consterrorreturn.IsErrorType(nil))
}