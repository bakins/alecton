package alecton

import (
	"github.com/bakins/alecton/api"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
)

// StorageProvider provides underlying storage
type StorageProvider interface {
	GetArtifact(context.Context, *api.GetArtifactRequest) (*api.Artifact, error)
	ListArtifacts(context.Context, *api.ListArtifactsRequest) (*api.ArtifactList, error)
	CreateArtifact(context.Context, *api.CreateArtifactRequest) (*api.Artifact, error)
	GetApplication(context.Context, *api.GetApplicationRequest) (*api.Application, error)
	ListApplications(context.Context, *api.ListApplicationsRequest) (*api.ApplicationList, error)
	CreateApplication(context.Context, *api.CreateApplicationRequest) (*api.Application, error)
	GetDeployment(context.Context, *api.GetDeploymentRequest) (*api.Deployment, error)
	ListDeployments(context.Context, *api.ListDeploymentsRequest) (*api.DeploymentList, error)
	CreateDeployment(ctx context.Context, r *api.CreateDeploymentRequest) (*api.Deployment, error)
	// UpdateDeployment is used by the server to store status about the deployment
	// the full deployment is passed in.
	UpdateDeployment(ctx context.Context, r *api.Deployment) (*api.Deployment, error)
}

// XXX: ^^^ is a large interface, but all these functions really
// go together

var storageProviders = map[string]StorageProviderCreateFunc{}

type StorageProviderCreateFunc func(map[string]interface{}) (StorageProvider, error)

func RegisterStorageProvider(name string, p StorageProviderCreateFunc) {
	storageProviders[name] = p
}

func getStorageProvider(config map[string]interface{}) (StorageProvider, error) {
	name, err := getProviderName(config)
	if err != nil {
		return nil, errors.Wrap(err, "unable to determine storage provider")
	}
	p, ok := storageProviders[name]
	if !ok || p == nil {
		return nil, errors.Wrapf(err, "no storage provider found for \"%s\"", name)
	}

	s, err := p(config)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create storage provider for \"%s\"", name)
	}
	return s, nil
}
