package memory

import (
	"github.com/bakins/alecton"
	"github.com/bakins/alecton/api"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
)

func (m *Memory) GetDeployment(ctx context.Context, r *api.GetDeploymentRequest) (*api.Deployment, error) {
	m.Lock()
	defer m.Unlock()
	a, ok := m.deployments[r.Id]
	if !ok || a == nil {
		return nil, alecton.NewNotFoundError("deployment", r.Id)
	}
	return a, nil
}

func (m *Memory) ListDeployments(ctx context.Context, r *api.ListDeploymentsRequest) (*api.DeploymentList, error) {
	m.Lock()
	defer m.Unlock()

	list := &api.DeploymentList{}
	for _, v := range m.deployments {
		if r.Application == "" || r.Application == v.Application {
			list.Items = append(list.Items, v)
		}
	}
	return list, nil
}

func (m *Memory) CreateDeployment(ctx context.Context, r *api.CreateDeploymentRequest) (*api.Deployment, error) {
	m.Lock()
	defer m.Unlock()

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create id")
	}

	d := &api.Deployment{
		Id:          id.String(),
		Application: r.Application,
		Target:      r.Target,
	}
	m.deployments[d.Id] = d

	return d, nil
}

func (m *Memory) UpdateDeployment(ctx context.Context, r *api.Deployment) (*api.Deployment, error) {
	m.Lock()
	defer m.Unlock()

	m.deployments[r.Id] = r

	return r, nil
}
