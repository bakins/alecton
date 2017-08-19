package alecton

import (
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

// ChartProvider provides functions for interacing with helm or similar.
type ChartProvider interface {
	GetChart(context.Context, string) (*chart.Chart, error)
	// XXX: need list charts?
}

var chartProviders = map[string]ChartProviderCreateFunc{}

// ChartProviderCreateFunc is used by config to create a provider
type ChartProviderCreateFunc func(map[string]interface{}) (ChartProvider, error)

// RegisterChartProvider is used by chart providers to register themselve
func RegisterChartProvider(name string, p ChartProviderCreateFunc) {
	chartProviders[name] = p
}

func getChartProvider(config map[string]interface{}) (ChartProvider, error) {
	name, err := getProviderName(config)
	if err != nil {
		return nil, errors.Wrap(err, "unable to determine chart provider")
	}
	p, ok := chartProviders[name]
	if !ok || p == nil {
		return nil, errors.Errorf("no chart provider found for \"%s\"", name)
	}
	s, err := p(config)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create chart provider for \"%s\"", name)
	}
	return s, nil
}
