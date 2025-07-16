package consterrorreturn_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/delarean/consterrorreturn"
)

func TestLinter(t *testing.T) {
	testdata := analysistest.TestData()

	results := analysistest.Run(t, testdata, consterrorreturn.Analyzer, "a")
	require.Len(t, results, 1)
	require.NoError(t, results[0].Err)
} 