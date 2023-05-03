package real_environment

import (
	"github.com/jialunzhai/crimemap/analytics/online/server/interfaces"
)

type RealEnv struct {
	config          *interfaces.Config
	httpServer      interfaces.HTTPServer
	grpcServer      interfaces.GRPCServer
	grpcWebServer   interfaces.GRPCWebServer
	crimemapService interfaces.CrimeMapService
	databaseClient  interfaces.DatabaseClient
}

func NewRealEnv() *RealEnv {
	return &RealEnv{}
}

func (r *RealEnv) GetConfig() *interfaces.Config {
	return r.config
}
func (r *RealEnv) SetConfig(c *interfaces.Config) {
	r.config = c
}

func (r *RealEnv) GetHTTPServer() interfaces.HTTPServer {
	return r.httpServer
}

func (r *RealEnv) SetHTTPServer(s interfaces.HTTPServer) {
	r.httpServer = s
}

func (r *RealEnv) GetGRPCServer() interfaces.GRPCServer {
	return r.grpcServer
}

func (r *RealEnv) SetGRPCServer(s interfaces.GRPCServer) {
	r.grpcServer = s
}

func (r *RealEnv) GetGRPCWebServer() interfaces.GRPCWebServer {
	return r.grpcWebServer
}

func (r *RealEnv) SetGRPCWebServer(s interfaces.GRPCWebServer) {
	r.grpcWebServer = s
}

func (r *RealEnv) GetCrimeMapService() interfaces.CrimeMapService {
	return r.crimemapService
}

func (r *RealEnv) SetCrimeMapService(s interfaces.CrimeMapService) {
	r.crimemapService = s
}

func (r *RealEnv) GetDatabaseClient() interfaces.DatabaseClient {
	return r.databaseClient
}

func (r *RealEnv) SetDatabaseClient(s interfaces.DatabaseClient) {
	r.databaseClient = s
}
