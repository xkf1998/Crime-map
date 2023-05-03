package http_server

import (
	"context"
	"errors"
	"log"
	"net/http"

	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
)

type HTTPServer struct {
	env    env_interface.Env
	server *http.Server
}

func Register(env env_interface.Env) error {
	config := env.GetConfig()
	if config == nil || config.HTTP.Address == "" {
		return errors.New("HTTP server not configured")
	}
	s, err := NewHTTPServer(env, config.HTTP.Address, config.HTTP.Bundle)
	if err != nil {
		return err
	}
	env.SetHTTPServer(s)
	return nil
}

func NewHTTPServer(env env_interface.Env, address string, bundle string) (*HTTPServer, error) {
	server := &http.Server{
		Addr:    address,
		Handler: http.FileServer(http.Dir(bundle)),
	}
	return &HTTPServer{
		env:    env,
		server: server,
	}, nil
}

func (s *HTTPServer) Run() error {
	log.Printf("Start to serve HTTP requests from address: %v\n", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
