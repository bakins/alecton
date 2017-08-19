package alecton

import (
	"github.com/bakins/alecton/api"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/strvals"
)

// DeployProvider provides functions for interacing with helm or similar.
type DeployProvider interface {
	InstallRelease(context.Context, *services.InstallReleaseRequest) (*release.Release, error)
	UpdateRelease(context.Context, *services.UpdateReleaseRequest) (*release.Release, error)
	RollbackRelease(context.Context, *services.RollbackReleaseRequest) (*release.Release, error)
	ReleaseHistory(context.Context, *services.GetHistoryRequest) ([]*release.Release, error)
}

var deployProviders = map[string]DeployProviderCreateFunc{}

// DeployProviderCreateFunc is used by config to create a provider
type DeployProviderCreateFunc func(map[string]interface{}) (DeployProvider, error)

// RegisterDeployProvider is used by deploy providers to register themselve
func RegisterDeployProvider(name string, p DeployProviderCreateFunc) {
	deployProviders[name] = p
}

func getDeployProvider(config map[string]interface{}) (DeployProvider, error) {
	name, err := getProviderName(config)
	if err != nil {
		return nil, errors.Wrap(err, "unable to determine deploy provider")
	}
	p, ok := deployProviders[name]
	if !ok || p == nil {
		return nil, errors.Errorf("no deploy provider found for \"%s\"", name)
	}

	s, err := p(config)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create deploy provider for \"%s\"", name)
	}
	return s, nil
}

func getReleaseName(namespace, app string) string {
	return namespace + "-" + app
}

func (s *Server) DeployApplication(ctx context.Context, r *api.DeployRequest) (*api.DeployResponse, error) {
	a, err := s.GetApplication(ctx, &api.GetApplicationRequest{Name: r.Application})
	if err != nil {
		// logging? or just let normal log handle it?
		return nil, err
	}

	var target *api.Target
	for _, t := range a.Targets {
		if t.Name == r.Target {
			target = t
			break
		}
	}

	if target == nil {
		return nil, NewNotFoundError("target", r.Target)
	}

	image, err := s.GetImage(ctx, &api.GetImageRequest{Name: a.Image, Version: r.Version})
	if err != nil {
		return nil, err
	}

	c, err := s.chart.GetChart(ctx, a.Chart)
	if err != nil {
		return nil, err
	}

	// this is a bit cheesy, but is an easy way to do nested maps without
	// needing to do that in protos
	values := map[string]interface{}{}
	for _, source := range []map[string]string{a.Values, target.Values} {
		for k, v := range source {
			value := k + "=" + v
			if err := strvals.ParseInto(value, values); err != nil {
				return nil, NewInvalidArgumentError("values", value)
			}
		}
	}

	// our deployment specific values. these have highest precendence
	// and cannot be overwritten
	values["Deploy"] = map[string]interface{}{
		"Image":   image.Image,
		"Version": image.Version,
		"App":     a.Name,
		"Target":  target.Name,
	}

	data, err := yaml.Marshal(values)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal values")
	}

	hisReq := &services.GetHistoryRequest{
		Name: getReleaseName(target.Namespace, a.Name),
		Max:  1,
	}

	isInstall := false
	_, err = s.deploy.ReleaseHistory(ctx, hisReq)

	switch {
	case err == nil:
		isInstall = false
	case IsNotFoundError(err):
		isInstall = true
	default:
		return nil, err
	}

	if isInstall {
		install := &services.InstallReleaseRequest{
			Chart:     c,
			Values:    &chart.Config{Raw: string(data)},
			Name:      getReleaseName(target.Namespace, a.Name),
			Namespace: target.Namespace,
			Timeout:   5, //todo confugurable?
			Wait:      false,
		}

		rel, err := s.deploy.InstallRelease(ctx, install)
		if err != nil {
			return nil, err
		}
		return &api.DeployResponse{Release: rel}, nil
	}

	update := &services.UpdateReleaseRequest{
		Chart:       c,
		Values:      &chart.Config{Raw: string(data)},
		Name:        getReleaseName(target.Namespace, a.Name),
		Timeout:     5, //todo confugurable?
		ResetValues: true,
		Wait:        false,
	}

	rel, err := s.deploy.UpdateRelease(ctx, update)
	if err != nil {
		return nil, err
	}
	return &api.DeployResponse{Release: rel}, nil
}

func (s *Server) RollbackApplication(ctx context.Context, r *api.RollbackRequest) (*api.RollbackResponse, error) {
	a, err := s.GetApplication(ctx, &api.GetApplicationRequest{Name: r.Application})
	if err != nil {
		return nil, err
	}

	var target *api.Target
	for _, t := range a.Targets {
		if t.Name == r.Target {
			target = t
			break
		}
	}
	if target == nil {
		return nil, NewNotFoundError("target", r.Target)
	}

	req := &services.RollbackReleaseRequest{
		Name:    getReleaseName(target.Namespace, r.Application),
		Version: r.Version,
	}

	res, err := s.deploy.RollbackRelease(ctx, req)
	if err != nil {
		return nil, err
	}
	return &api.RollbackResponse{Release: res}, nil
}

func (s *Server) ListReleases(ctx context.Context, r *api.ListReleasesRequest) (*api.ListReleasesResponse, error) {
	a, err := s.GetApplication(ctx, &api.GetApplicationRequest{Name: r.Application})
	if err != nil {
		return nil, err
	}

	var target *api.Target
	for _, t := range a.Targets {
		if t.Name == r.Target {
			target = t
			break
		}
	}
	if target == nil {
		return nil, NewNotFoundError("target", r.Target)
	}

	req := &services.GetHistoryRequest{
		Name: getReleaseName(target.Namespace, a.Name),
		Max:  256,
	}

	res, err := s.deploy.ReleaseHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	return &api.ListReleasesResponse{Releases: res}, nil
}
