package darrell

import (
	"fmt"

	"github.com/bakins/darrell/api"
	context "golang.org/x/net/context"
)

type ArtifactGetter interface {
	GetArtifact(context.Context, *api.GetArtifactRequest) (*api.Artifact, error)
}

type ArtifactLister interface {
	ListArtifacts(context.Context, *api.ListArtifactsRequest) (*api.ArtifactList, error)
}

type ArtifactBuildGetter interface {
	GetArtifactBuild(context.Context, *api.GetArtifactBuildRequest) (*api.ArtifactBuild, error)
}

type ArtifactBuildLister interface {
	ListArtifactBuilds(context.Context, *api.ListArtifactBuildsRequest) (*api.ArtifactBuildList, error)
}

type ArtifactCreator interface {
	CreateArtifact(context.Context, *api.CreateArtifactRequest) (*api.Artifact, error)
}

type ArtifactInterface interface {
	ArtifactGetter
	ArtifactLister
	ArtifactBuildGetter
	ArtifactBuildLister
	ArtifactCreator
}

type ApplicationGetter interface {
	GetApplication(context.Context, *api.GetApplicationRequest) (*api.Artifact, error)
}

type ApplicationLister interface {
	ListApplications(context.Context, *api.ListApplicationRequest) (*api.ApplicationList, error)
}

type ApplicationInterface interface {
	ApplicationGetter
	ApplicationLister
}

type DeploymentGetter interface {
	GetDeployment(context.Context, *api.GetDeploymentRequest) (*api.Deployment, error)
}

type DeploymentLister interface {
	ListDeployments(context.Context, *api.ListDeploymentsRequest) (*api.DeploymentList, error)
}

type DeploymentInterface interface {
	DeploymentGetter
	DeploymentLister
}

type DarrellInterface interface {
	ArtifactInterface
	//	ApplicationInterface
	//	DeploymentInterface
}

type NotFoundError struct {
	name string
	kind string
}

func NewNotFoundError(kind, name string) error {
	return &NotFoundError{name: name, kind: kind}
}

func (err *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s/%s", err.kind, err.name)
}

func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
