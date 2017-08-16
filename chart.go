package alecton

import (
	"path/filepath"

	"github.com/bakins/alecton/api"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/releaseutil"
	"k8s.io/helm/pkg/timeconv"
)

// ChartProvider provides access to charts
type ChartProvider interface {
	GetChart(context.Context, string) (*chart.Chart, error)
}

var chartProviders = map[string]ChartProviderCreateFunc{}

type ChartProviderCreateFunc func(map[string]interface{}) (ChartProvider, error)

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
		return nil, errors.Wrapf(err, "no chart provider found for \"%s\"", name)
	}

	c, err := p(config)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create chart provider for \"%s\"", name)
	}
	return c, nil
}

// liberally borrows from Helm
// and https://github.com/technosophos/helm-template/blob/master/main.go
// merge multiple values. Later aruments win
func vals(values ...map[string]interface{}) ([]byte, error) {
	base := map[string]interface{}{}

	for _, v := range values {
		base = mergeValues(base, v)
	}

	return yaml.Marshal(base)
}

// Copied from Helm.
// and https://github.com/technosophos/helm-template/blob/master/main.go

func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = nextMap
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}

func mergeVals(values ...map[string]string) ([]byte, error) {
	base := map[string]string{}

	for _, vals := range values {
		for k, v := range vals {
			base[k] = v
		}
	}

	return yaml.Marshal(base)
}

func renderChart(c *chart.Chart, a *api.Application, t *api.Target) ([]*api.Resource, error) {
	options := chartutil.ReleaseOptions{
		Name:      a.Name + "-" + t.Name,
		Time:      timeconv.Now(),
		Namespace: t.Namespace,
	}

	overrides := map[string]string{
		"Cluster":     t.Cluster,
		"Target":      t.Name,
		"Application": a.Name,
	}

	vv, err := mergeVals(a.Defaults, t.Values, overrides)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to merge values: %s %s", a.Name, t.Name)
	}

	config := &chart.Config{
		Raw:    string(vv),
		Values: map[string]*chart.Value{},
	}

	renderer := engine.New()

	vals, err := chartutil.ToRenderValues(c, config, options)
	if err != nil {
		return nil, errors.Wrapf(err, "ToRenderValues failed: %s %s", a.Name, t.Name)
	}

	out, err := renderer.Render(c, vals)
	if err != nil {
		return nil, errors.Wrapf(err, "template render failed: %s %s", a.Name, t.Name)
	}

	output := make([]*api.Resource, 0, len(out))

	for name, data := range out {
		if filepath.Ext(name) != ".yaml" {
			continue
		}
		manifests := releaseutil.SplitManifests(data)
		for _, v := range manifests {
			k, err := parseAndValidatek8sResource(v)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to validate kubernetes resource in %s", name)
			}

			r := &api.Resource{
				Name:      k.Metadata.Name,
				Namespace: k.Metadata.Name,
				Kind:      k.Kind,
				Data:      v,
			}
			output = append(output, r)
		}
	}
	return output, nil
}

// just enough kubernetes to validate simple objects
type k8sResource struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Annotations map[string]string `json:"annotations"`
		Labels      map[string]string `json:"labels"`
		Name        string            `json:"name"`
		Namespace   string            `json:"namespace"`
	} `json:"metadata"`
}

func parseAndValidatek8sResource(data string) (*k8sResource, error) {
	var k k8sResource

	if err := yaml.Unmarshal([]byte(data), &k); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal resource")
	}

	return &k, nil
}
