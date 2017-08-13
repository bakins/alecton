package memory

import (
	"github.com/bakins/darrell/api"
	context "golang.org/x/net/context"
)

func (m *Memory) GetApplication(context.Context, *api.GetApplicationRequest) (*api.Artifact, error) {
	panic("not implemented")
}

func (m *Memory) ListApplications(context.Context, *api.ListApplicationRequest) (*api.ApplicationList, error) {
	panic("not implemented")
}
