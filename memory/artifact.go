package memory

import (
	"github.com/bakins/alecton"
	"github.com/bakins/alecton/api"
	context "golang.org/x/net/context"
)

func artifactKey(name, version string) string {
	return name + "/" + version
}

func (m *Memory) GetArtifact(ctx context.Context, r *api.GetArtifactRequest) (*api.Artifact, error) {
	m.Lock()
	defer m.Unlock()
	key := artifactKey(r.Name, r.Version)
	a, ok := m.artifacts[key]
	if !ok || a == nil {
		return nil, alecton.NewNotFoundError("artifact", key)
	}
	return a, nil
}

func (m *Memory) ListArtifacts(ctx context.Context, r *api.ListArtifactsRequest) (*api.ArtifactList, error) {
	m.Lock()
	defer m.Unlock()

	list := &api.ArtifactList{}
	for _, v := range m.artifacts {
		if r.Name == "" || r.Name == v.Name {
			list.Items = append(list.Items, v)
		}
	}
	return list, nil
}

func (m *Memory) CreateArtifact(ctx context.Context, r *api.CreateArtifactRequest) (*api.Artifact, error) {
	m.Lock()
	defer m.Unlock()

	key := artifactKey(r.Artifact.Name, r.Artifact.Version)

	if !r.Overwrite {
		_, ok := m.artifacts[key]
		if ok {
			return nil, alecton.NewAlreadyExistsError("artifact", key)
		}
	}

	m.artifacts[key] = r.Artifact

	return r.Artifact, nil
}
