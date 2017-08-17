// Package memory implements a darrel store in memory
package memory

import (
	"sync"

	"github.com/bakins/alecton"
	"github.com/bakins/alecton/api"
	context "golang.org/x/net/context"
)

type Memory struct {
	sync.Mutex
	images       []*api.Image
	applications map[string]*api.Application
}

func New() *Memory {
	return &Memory{
		images:       make([]*api.Image, 0),
		applications: make(map[string]*api.Application),
	}
}

func provider(map[string]interface{}) (alecton.StorageProvider, error) {
	return New(), nil
}

func init() {
	alecton.RegisterStorageProvider("memory", provider)
}

func (m *Memory) ListImages(ctx context.Context, r *api.ListImagesRequest) (*api.ListImagesResponse, error) {
	m.Lock()
	defer m.Unlock()

	resp := &api.ListImagesResponse{}

	for _, i := range m.images {
		if (r.Name == "" || r.Name == i.Name) && (r.Version == "" || r.Version == i.Version) {
			resp.Images = append(resp.Images, i)
		}
	}
	return resp, nil
}

func (m *Memory) GetImage(ctx context.Context, r *api.GetImageRequest) (*api.Image, error) {
	m.Lock()
	defer m.Unlock()

	for _, i := range m.images {
		if (r.Name == i.Name) && (r.Version == i.Version) {
			return i, nil
		}
	}
	return nil, alecton.NewNotFoundError("image", r.Name+"/"+r.Version)
}

func (m *Memory) CreateImage(ctx context.Context, r *api.Image) (*api.Image, error) {
	m.Lock()
	defer m.Unlock()

	for _, i := range m.images {
		if (r.Name == i.Name) && (r.Version == i.Version) {
			return nil, alecton.NewAlreadyExistsError("image", r.Name+"/"+r.Version)
		}
	}

	// we should make a copy, but this is fine for the mock case
	m.images = append(m.images, r)
	return r, nil
}

func (m *Memory) GetApplication(ctx context.Context, r *api.GetApplicationRequest) (*api.Application, error) {
	m.Lock()
	defer m.Unlock()

	a, ok := m.applications[r.Name]
	if !ok || a == nil {
		return nil, alecton.NewNotFoundError("application", r.Name)
	}
	return a, nil
}

func (m *Memory) ListApplications(ctx context.Context, r *api.ListApplicationsRequest) (*api.ListApplicationsResponse, error) {
	m.Lock()
	defer m.Unlock()

	res := &api.ListApplicationsResponse{}
	for _, a := range m.applications {
		if r.Chart != "" && r.Chart != a.Chart {
			continue

		}

		if r.Image != "" && r.Image != a.Image {
			continue

		}

		res.Applications = append(res.Applications, a)
	}
	return res, nil
}

func (m *Memory) CreateApplication(ctx context.Context, r *api.Application) (*api.Application, error) {
	m.Lock()
	defer m.Unlock()

	_, ok := m.applications[r.Name]
	if ok {
		return nil, alecton.NewAlreadyExistsError("application", r.Name)
	}

	m.applications[r.Name] = r

	return r, nil
}

func (m *Memory) UpdateApplication(ctx context.Context, r *api.Application) (*api.Application, error) {
	m.Lock()
	defer m.Unlock()

	m.applications[r.Name] = r

	return r, nil
}
