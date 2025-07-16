package gclplugin_test

import (
	"testing"

	"github.com/delarean/consterrorreturn"
	gclplugin "github.com/delarean/consterrorreturn/cmd/gclplugin"
	"github.com/golangci/plugin-module-register/register"
	"github.com/stretchr/testify/require"
)

func TestPlugin(t *testing.T) {
	plugin, err := gclplugin.New(nil)
	require.NoError(t, err)
	require.NotNil(t, plugin)

	analyzers, err := plugin.BuildAnalyzers()
	require.NoError(t, err)
	require.Len(t, analyzers, 1)
	require.Equal(t, consterrorreturn.Analyzer, analyzers[0])

	loadMode := plugin.GetLoadMode()
	require.Equal(t, register.LoadModeTypesInfo, loadMode)
} 