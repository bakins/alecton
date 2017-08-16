// Package chartdir loads charts from a directory
package chartdir

import (
	"os"
	"path/filepath"

	"github.com/bakins/alecton"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type chartDirConfig struct {
	Dir string
}

// ChartDir is a directory of charts.
type ChartDir struct {
	dir string
}

// New creates a new ChartDir
func New(dir string) (*ChartDir, error) {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine absolute path of %s", dir)
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to stat %s", absPath)
	}

	if !info.IsDir() {
		return nil, errors.Wrapf(err, "not a directory %s", absPath)
	}

	return &ChartDir{dir: absPath}, nil
}

// GetChart gets a single Chart from the directory.
func (d *ChartDir) GetChart(ctx context.Context, name string) (*chart.Chart, error) {
	return chartutil.LoadDir(filepath.Join(d.dir, name))
}

func provider(config map[string]interface{}) (alecton.ChartProvider, error) {
	var d chartDirConfig
	if err := alecton.ProviderConfigDecode(config, &d); err != nil {
		return nil, errors.Wrap(err, "failed to get chartdir config")
	}
	if d.Dir == "" {
		return nil, errors.New("dir cannot be empty")
	}
	return New(d.Dir)
}

func init() {
	alecton.RegisterChartProvider("dir", provider)
}
