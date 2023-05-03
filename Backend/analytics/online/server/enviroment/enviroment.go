package environment

import (
	"github.com/jialunzhai/crimemap/analytics/online/server/interfaces"
)

type Env interface {
	GetConfig() *interfaces.Config
	SetConfig(*interfaces.Config)
	GetHTTPServer() interfaces.HTTPServer
	SetHTTPServer(interfaces.HTTPServer)
	GetGRPCServer() interfaces.GRPCServer
	SetGRPCServer(interfaces.GRPCServer)
	GetGRPCWebServer() interfaces.GRPCWebServer
	SetGRPCWebServer(interfaces.GRPCWebServer)
	GetCrimeMapService() interfaces.CrimeMapService
	SetCrimeMapService(interfaces.CrimeMapService)
	GetDatabaseClient() interfaces.DatabaseClient
	SetDatabaseClient(interfaces.DatabaseClient)
}
