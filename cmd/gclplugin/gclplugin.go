package gclplugin

import (
	"github.com/delarean/consterrorreturn"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("consterrorreturn", New)
}

func New(settings any) (register.LinterPlugin, error) {
	return &ConstErrorReturnPlugin{}, nil
}

type ConstErrorReturnPlugin struct{}

func (p *ConstErrorReturnPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{consterrorreturn.Analyzer}, nil
}

func (p *ConstErrorReturnPlugin) GetLoadMode() string { return register.LoadModeTypesInfo } 