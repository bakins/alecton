package darrell

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/bakins/darrell/api"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/hkwi/h2c"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	ctx              context.Context
	cancel           context.CancelFunc
	darrellInterface DarrellInterface
	grpc             *grpc.Server
	server           *http.Server
}

func NewServer(ctx context.Context, d DarrellInterface) *Server {
	c, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:              c,
		cancel:           cancel,
		darrellInterface: d,
	}
}

// Run starts the server. This generally does not return.
func (s *Server) Run(address string) error {
	logger, err := NewDefaultLogger()
	if err != nil {
		return errors.Wrapf(err, "failed to create logger")
	}

	l, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on %s", address)
	}

	grpc_zap.ReplaceGrpcLogger(logger)
	grpc_prometheus.EnableHandlingTimeHistogram()

	s.grpc = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_validator.UnaryServerInterceptor(),
				grpc_prometheus.UnaryServerInterceptor,
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_recovery.UnaryServerInterceptor(),
			),
		),
	)

	api.RegisterApplicationServiceServer(s.grpc, s)
	// not exactly sure what this is used for, but examples
	// always do it:
	// https://godoc.org/google.golang.org/grpc/reflection
	reflection.Register(s.grpc)

	gwmux := runtime.NewServeMux()

	_, port, err := net.SplitHostPort(address)
	if err != nil {
		return errors.Wrapf(err, "invalid address %s", address)
	}

	if err := api.RegisterApplicationServiceHandlerFromEndpoint(s.ctx, gwmux, net.JoinHostPort("127.0.0.1", port), []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		return errors.Wrap(err, "failed to register grpc gateway")
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", gwmux)

	s.server = &http.Server{
		Handler: h2c.Server{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.ProtoMajor == 2 &&
					strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
					s.grpc.ServeHTTP(w, r)
				} else {
					mux.ServeHTTP(w, r)
				}
			}),
		},
	}

	if err := s.server.Serve(l); err != nil {
		if err != http.ErrServerClosed {
			return errors.Wrap(err, "failed to start http server")
		}
	}

	return nil
}

// NewDefaultLogger creates a new logger with our prefered options
func NewDefaultLogger() (*zap.Logger, error) {
	config := zap.Config{
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		Encoding:          "json",
		ErrorOutputPaths:  []string{"stdout"},
		Level:             zap.NewAtomicLevel(),
		OutputPaths:       []string{"stdout"},
	}
	l, err := config.Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create logger")
	}
	return l, nil
}
func (s *Server) GetApplication(context.Context, *api.GetApplicationRequest) (*api.Application, error) {
	panic("not implemented")
}

func (s *Server) ListApplications(context.Context, *api.ListApplicationRequest) (*api.ApplicationList, error) {
	panic("not implemented")
}

func (s *Server) GetArtifact(context.Context, *api.GetArtifactRequest) (*api.Artifact, error) {
	panic("not implemented")
}

func (s *Server) ListArtifacts(ctx context.Context, r *api.ListArtifactsRequest) (*api.ArtifactList, error) {
	fmt.Println("server", "ListArtifacts")
	return s.darrellInterface.ListArtifacts(ctx, r)
}

func (s *Server) GetArtifactBuild(context.Context, *api.GetArtifactBuildRequest) (*api.ArtifactBuild, error) {
	panic("not implemented")
}

func (s *Server) ListArtifactBuilds(context.Context, *api.ListArtifactBuildsRequest) (*api.ArtifactBuildList, error) {
	panic("not implemented")
}

func (s *Server) CreateArtifact(ctx context.Context, r *api.CreateArtifactRequest) (*api.Artifact, error) {
	// TODO: validate?
	return s.darrellInterface.CreateArtifact(ctx, r)
}

func (s *Server) GetDeployment(context.Context, *api.GetDeploymentRequest) (*api.Deployment, error) {
	panic("not implemented")
}

func (s *Server) ListDeployments(context.Context, *api.ListDeploymentsRequest) (*api.DeploymentList, error) {
	panic("not implemented")
}
