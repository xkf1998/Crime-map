package grpc_web_server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
	"google.golang.org/grpc"
)

type HTTPServer struct {
	env    env_interface.Env
	server *http.Server
}

func Register(env env_interface.Env) error {
	config := env.GetConfig()
	if config == nil || config.GRPCWeb.Address == "" {
		return errors.New("gRPC web server not configured")
	}
	grpcServer := env.GetGRPCServer()
	if grpcServer == nil {
		return errors.New("gRPC server should be registered before gRPC web server")
	}
	s, err := NewGRPCWebServer(env, config.GRPCWeb.Address, grpcServer.GetServer())
	if err != nil {
		return err
	}
	env.SetGRPCWebServer(s)
	return nil
}

func NewGRPCWebServer(env env_interface.Env, address string, grpcServer *grpc.Server) (*HTTPServer, error) {
	handler := grpcweb.WrapServer(
		grpcServer,
		// Enable CORS
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)
	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}
	return &HTTPServer{
		env:    env,
		server: server,
	}, nil
}

func (s *HTTPServer) Run() error {
	log.Printf("Start to serve gRPC web requests from address: %v\n", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
