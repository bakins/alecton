package alecton

import (
	"github.com/bakins/alecton/api"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
)

// DeployProvider handles actual deployment.
type DeployProvider interface {
	EnsureNamespace(ctx context.Context, cluster string, namespace string) error
	Deploy(ctx context.Context, cluster string, namespace string, resources []*api.Resource) error
}

var deployProviders = map[string]DeployProviderCreateFunc{}

type DeployProviderCreateFunc func(map[string]interface{}) (DeployProvider, error)

func RegisterdeployProvider(name string, p DeployProviderCreateFunc) {
	deployProviders[name] = p
}

func getdeployProvider(config map[string]interface{}) (DeployProvider, error) {
	name, err := getProviderName(config)
	if err != nil {
		return nil, errors.Wrap(err, "unable to determine deploy provider")
	}
	p, ok := deployProviders[name]
	if !ok || p == nil {
		return nil, errors.Wrapf(err, "no deploy provider found for \"%s\"", name)
	}

	s, err := p(config)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create deploy provider for \"%s\"", name)
	}
	return s, nil
}
