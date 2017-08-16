package memory

import (
	"github.com/bakins/alecton"
	"github.com/bakins/alecton/api"
	context "golang.org/x/net/context"
)

func (m *Memory) GetApplication(ctx context.Context, r *api.GetApplicationRequest) (*api.Application, error) {
	m.Lock()
	defer m.Unlock()
	a, ok := m.applications[r.Name]
	if !ok || a == nil {
		return nil, alecton.NewNotFoundError("application", r.Name)
	}
	return a, nil
}

func (m *Memory) ListApplications(context.Context, *api.ListApplicationsRequest) (*api.ApplicationList, error) {
	m.Lock()
	defer m.Unlock()

	list := &api.ApplicationList{}
	for _, v := range m.applications {
		list.Items = append(list.Items, v)
	}
	return list, nil
}

func (m *Memory) CreateApplication(ctx context.Context, r *api.CreateApplicationRequest) (*api.Application, error) {
	m.Lock()
	defer m.Unlock()

	if !r.Overwrite {
		_, ok := m.applications[r.Application.Name]
		if ok {
			return nil, alecton.NewAlreadyExistsError("application", r.Application.Name)
		}
	}

	m.applications[r.Application.Name] = r.Application

	return r.Application, nil
}
