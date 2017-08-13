package memory

import (
	"github.com/bakins/darrell"
	"github.com/bakins/darrell/api"
	context "golang.org/x/net/context"
)

func (m *Memory) GetArtifact(ctx context.Context, r *api.GetArtifactRequest) (*api.Artifact, error) {
	m.Lock()
	defer m.Unlock()
	a, ok := m.artifacts[r.Name]
	if !ok || a == nil {
		return nil, darrell.NewNotFoundError("artifact", r.Name)
	}
	return a, nil
}

func (m *Memory) ListArtifacts(ctx context.Context, r *api.ListArtifactsRequest) (*api.ArtifactList, error) {
	m.Lock()
	defer m.Unlock()

	list := &api.ArtifactList{}
	for _, v := range m.artifacts {
		list.Items = append(list.Items, v)
	}
	return list, nil
}

func (m *Memory) GetArtifactBuild(ctx context.Context, r *api.GetArtifactBuildRequest) (*api.ArtifactBuild, error) {
	m.Lock()
	defer m.Unlock()
	b, ok := m.artifactBuilds[r.Name]
	if !ok || b == nil {
		return nil, darrell.NewNotFoundError("artifactBuild", r.Name)
	}
	return b, nil
}

func (m *Memory) ListArtifactBuilds(ctx context.Context, r *api.ListArtifactBuildsRequest) (*api.ArtifactBuildList, error) {
	m.Lock()
	defer m.Unlock()

	list := &api.ArtifactBuildList{}
	for _, v := range m.artifactBuilds {
		list.Items = append(list.Items, v)
	}

	return list, nil
}

func (m *Memory) CreateArtifact(ctx context.Context, r *api.CreateArtifactRequest) (*api.Artifact, error) {
	m.Lock()
	defer m.Unlock()

	m.artifacts[r.Artifact.Name] = r.Artifact

	return r.Artifact, nil
}
