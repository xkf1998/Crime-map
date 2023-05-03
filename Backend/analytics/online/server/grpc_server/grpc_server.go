package grpc_server

import (
	"errors"
	"fmt"
	"log"
	"net"

	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	env      env_interface.Env
	listener net.Listener
	server   *grpc.Server
}

func Register(env env_interface.Env) error {
	config := env.GetConfig()
	if config == nil || config.GRPC.Address == "" {
		return errors.New("gRPC server not configured")
	}
	s, err := NewGRPCServer(env, config.GRPC.Address)
	if err != nil {
		return err
	}
	env.SetGRPCServer(s)
	return nil
}

func NewGRPCServer(env env_interface.Env, address string) (*GRPCServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer()
	return &GRPCServer{
		env:      env,
		listener: listener,
		server:   server,
	}, nil
}

func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

func (s *GRPCServer) Run() error {
	if s.env.GetCrimeMapService() == nil {
		fmt.Errorf("CrimeMapServer.Register should be called before GRPCServer.Run")
	}
	log.Printf("Start to serve gRPPC requests from address: %v\n", s.listener.Addr())
	return s.server.Serve(s.listener)
}

func (s *GRPCServer) Shutdown() {
	s.server.Stop()
}
