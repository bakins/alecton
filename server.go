package alecton

import (
	"net"
	"net/http"
	"strings"

	"github.com/bakins/alecton/api"
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

// ServerOptionFunc is a function passed to new for setting options on a new server.
type ServerOptionFunc func(*Server) error

// Server is a deployment server
type Server struct {
	ctx     context.Context
	cancel  context.CancelFunc
	grpc    *grpc.Server
	server  *http.Server
	storage StorageProvider
	address string
}

func NewServer(options ...ServerOptionFunc) (*Server, error) {
	s := &Server{
		address: "127.0.0.1:8080",
	}

	for _, f := range options {
		if err := f(s); err != nil {
			return nil, errors.Wrap(err, "failed to set options")
		}
	}

	if s.ctx == nil || s.cancel == nil {
		s.ctx, s.cancel = context.WithCancel(context.Background())
	}

	return s, nil
}

// SetAddress sets the listening address.
func SetAddress(addr string) ServerOptionFunc {
	return func(s *Server) error {
		s.address = addr
		return nil
	}
}

// SetContext sets the context. if not set, a default Background context
// is used
func SetContext(ctx context.Context) ServerOptionFunc {
	return func(s *Server) error {
		s.ctx, s.cancel = context.WithCancel(ctx)
		return nil
	}
}

// SetStorageProvider sets the storage provide. There is no default
func SetStorageProvider(p StorageProvider) ServerOptionFunc {
	return func(s *Server) error {
		s.storage = p
		return nil
	}
}

// Run starts the server. This generally does not return.
func (s *Server) Run() error {
	logger, err := NewDefaultLogger()
	if err != nil {
		return errors.Wrapf(err, "failed to create logger")
	}

	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on %s", s.address)
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

	_, port, err := net.SplitHostPort(s.address)
	if err != nil {
		return errors.Wrapf(err, "invalid address %s", s.address)
	}

	// TODO: need to determine if we can actually connect to localhost
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

func (s *Server) GetApplication(ctx context.Context, r *api.GetApplicationRequest) (*api.Application, error) {
	return s.storage.GetApplication(ctx, r)
}

func (s *Server) ListApplications(ctx context.Context, r *api.ListApplicationsRequest) (*api.ApplicationList, error) {
	return s.storage.ListApplications(ctx, r)
}

func (s *Server) CreateApplication(ctx context.Context, r *api.CreateApplicationRequest) (*api.Application, error) {
	return s.storage.CreateApplication(ctx, r)
}

func (s *Server) GetArtifact(ctx context.Context, r *api.GetArtifactRequest) (*api.Artifact, error) {
	return s.storage.GetArtifact(ctx, r)
}

func (s *Server) ListArtifacts(ctx context.Context, r *api.ListArtifactsRequest) (*api.ArtifactList, error) {
	return s.storage.ListArtifacts(ctx, r)
}

func (s *Server) CreateArtifact(ctx context.Context, r *api.CreateArtifactRequest) (*api.Artifact, error) {
	return s.storage.CreateArtifact(ctx, r)
}

func (s *Server) GetDeployment(ctx context.Context, r *api.GetDeploymentRequest) (*api.Deployment, error) {
	return s.storage.GetDeployment(ctx, r)
}

func (s *Server) ListDeployments(ctx context.Context, r *api.ListDeploymentsRequest) (*api.DeploymentList, error) {
	return s.storage.ListDeployments(ctx, r)
}

func (s *Server) CreateDeployment(ctx context.Context, r *api.CreateDeploymentRequest) (*api.Deployment, error) {
	// TODO: after creating deployment, then actually try to
	// deploy it if activate is true
	// need storage to be able to set status
	return s.storage.CreateDeployment(ctx, r)
}
