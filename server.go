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
	address string
	StorageProvider
	chart  ChartProvider
	deploy DeployProvider
}

// NewServer creates a new server.
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

	if s.StorageProvider == nil {
		return nil, errors.New("storage provider is required")
	}

	if s.chart == nil {
		return nil, errors.New("chart provider is required")
	}

	if s.deploy == nil {
		return nil, errors.New("deploy provider is required")
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
		s.StorageProvider = p
		return nil
	}
}

// SetChartProvider sets the storage provide. There is no default
func SetChartProvider(p ChartProvider) ServerOptionFunc {
	return func(s *Server) error {
		s.chart = p
		return nil
	}
}

// SetDeployProvider sets the storage provide. There is no default
func SetDeployProvider(p DeployProvider) ServerOptionFunc {
	return func(s *Server) error {
		s.deploy = p
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

	api.RegisterDeployServiceServer(s.grpc, s)

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
	if err := api.RegisterDeployServiceHandlerFromEndpoint(s.ctx, gwmux, net.JoinHostPort("127.0.0.1", port), []grpc.DialOption{grpc.WithInsecure()}); err != nil {
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
