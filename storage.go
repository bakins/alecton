package alecton

import (
	"github.com/bakins/alecton/api"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
)

// StorageProvider provides underlying storage
type StorageProvider interface {
	CreateImage(context.Context, *api.Image) (*api.Image, error)
	ListImages(context.Context, *api.ListImagesRequest) (*api.ListImagesResponse, error)
	GetImage(ctx context.Context, r *api.GetImageRequest) (*api.Image, error)
	CreateApplication(context.Context, *api.Application) (*api.Application, error)
	UpdateApplication(context.Context, *api.Application) (*api.Application, error)
	ListApplications(context.Context, *api.ListApplicationsRequest) (*api.ListApplicationsResponse, error)
	GetApplication(context.Context, *api.GetApplicationRequest) (*api.Application, error)
}

var storageProviders = map[string]StorageProviderCreateFunc{}

// StorageProviderCreateFunc is used by config to create a provider
type StorageProviderCreateFunc func(map[string]interface{}) (StorageProvider, error)

// RegisterStorageProvider is used by storage providers to register themselve
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
