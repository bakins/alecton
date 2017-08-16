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
	//ChartsConfig map[string]interface{} `json:"charts"`
}

// ChartProvider provides the chart
//type ChartProvider interface {
//	Load(string) *chart.Chart
//}

// ProviderConfigDecode will decode a generic config into a
// specific config
func ProviderConfigDecode(in map[string]interface{}, rawVal interface{}) error {
	return mapstructure.Decode(in, rawVal)
}

// ServerFromConfigFile creates a server from a config file
func ServerFromConfigFile(filename string) (*Server, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", filename)
	}

	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, errors.Wrapf(err, "failed to parse file %s", filename)
	}

	return ServerFromConfig(&c)
}

// ServerFromConfig creates a Server from a Config
func ServerFromConfig(c *Config) (*Server, error) {
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
