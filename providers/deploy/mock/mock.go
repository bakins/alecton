// Package mock implements a mock deployer
package mock

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bakins/alecton"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/timeconv"
)

type Mock struct {
	sync.Mutex
	releases map[string]*release.Release
}

func New() *Mock {
	return &Mock{
		releases: make(map[string]*release.Release),
	}
}

func renderChart(c *chart.Chart, config *chart.Config, name string, namespace string) (string, error) {
	e := engine.New()
	options := chartutil.ReleaseOptions{
		Name:      name,
		Time:      timeconv.Now(),
		Namespace: namespace,
	}
	vals, err := chartutil.ToRenderValues(c, config, options)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert values")
	}
	out, err := e.Render(c, vals)
	if err != nil {
		return "", errors.Wrap(err, "failed to render chart")
	}

	var buff bytes.Buffer
	writer := bufio.NewWriter(&buff)
	for k, v := range out {
		if filepath.Ext(k) != ".yaml" {
			continue
		}
		if strings.HasPrefix(filepath.Base(k), "_") {
			continue
		}
		fmt.Fprintf(writer, "---\n# Source: %s\n", k)
		fmt.Fprintln(writer, v)
	}
	writer.Flush()

	return buff.String(), nil
}

func (m *Mock) InstallRelease(ctx context.Context, req *services.InstallReleaseRequest) (*release.Release, error) {
	m.Lock()
	defer m.Unlock()

	manifest, err := renderChart(req.Chart, req.Values, req.Name, req.Namespace)
	if err != nil {
		return nil, err
	}
	rel := &release.Release{
		Name:      req.Name,
		Namespace: req.Namespace,
		Chart:     req.Chart,
		Config:    req.Values,
		Version:   1,
		Manifest:  manifest,
	}

	return rel, nil
}

func (m *Mock) UpdateRelease(ctx context.Context, req *services.UpdateReleaseRequest) (*release.Release, error) {
	m.Lock()
	defer m.Unlock()

	rel, ok := m.releases[req.Name]
	if !ok {
		return nil, alecton.NewNotFoundError("release", req.Name)
	}

	manifest, err := renderChart(req.Chart, req.Values, req.Name, rel.Namespace)
	if err != nil {
		return nil, err
	}

	rel.Chart = req.Chart
	rel.Config = req.Values
	rel.Manifest = manifest

	return rel, nil
}

func (m *Mock) RollbackRelease(ctx context.Context, req *services.RollbackReleaseRequest) (*release.Release, error) {
	m.Lock()
	defer m.Unlock()

	rel, ok := m.releases[req.Name]
	if !ok {
		return nil, alecton.NewNotFoundError("release", req.Name)
	}

	rel.Version = req.Version

	return rel, nil
}

func (m *Mock) ReleaseHistory(ctx context.Context, req *services.GetHistoryRequest) ([]*release.Release, error) {
	m.Lock()
	defer m.Unlock()

	rel, ok := m.releases[req.Name]
	if !ok {
		return nil, alecton.NewNotFoundError("release", req.Name)
	}

	return []*release.Release{rel}, nil
}

func provider(map[string]interface{}) (alecton.DeployProvider, error) {
	return New(), nil
}

func init() {
	alecton.RegisterDeployProvider("mock", provider)
}
