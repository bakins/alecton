package alecton

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

//

// Config is a server configuration
type Config struct {
	// listen address
	Address       string
	StorageConfig map[string]interface{} `json:"storage"`
	ChartConfig   map[string]interface{} `json:"chart"`
	DeployConfig  map[string]interface{} `json:"deploy"`
}

// ProviderConfigDecode will decode a generic config into a
// specific config
func ProviderConfigDecode(in map[string]interface{}, rawVal interface{}) error {
	return mapstructure.Decode(in, rawVal)
}

// NewServerFromConfigFile creates a server from a config file
func NewServerFromConfigFile(filename string) (*Server, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", filename)
	}

	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, errors.Wrapf(err, "failed to parse file %s", filename)
	}

	return NewServerFromConfig(&c)
}

// NewServerFromConfig creates a Server from a Config
func NewServerFromConfig(c *Config) (*Server, error) {
	opts := []ServerOptionFunc{}

	if c.Address != "" {
		opts = append(opts, SetAddress(c.Address))
	}

	if c.StorageConfig != nil {
		s, err := getStorageProvider(c.StorageConfig)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get storage provider")
		}
		opts = append(opts, SetStorageProvider(s))
	}

	if c.ChartConfig != nil {
		prov, err := getChartProvider(c.ChartConfig)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get chart provider")
		}
		opts = append(opts, SetChartProvider(prov))
	}

	if c.DeployConfig != nil {
		prov, err := getDeployProvider(c.DeployConfig)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get deploy provider")
		}
		opts = append(opts, SetDeployProvider(prov))
	}

	return NewServer(opts...)
}

func getProviderName(in map[string]interface{}) (string, error) {
	raw := in["provider"]
	val, ok := raw.(string)
	if !ok {
		return "", errors.New("provider is not a string")
	}
	if val == "" {
		return "", errors.New("no provider given")
	}
	return val, nil
}
