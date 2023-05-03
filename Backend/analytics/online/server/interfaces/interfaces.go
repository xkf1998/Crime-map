package interfaces

import (
	"context"

	cmspb "github.com/jialunzhai/crimemap/analytics/online/proto/crimemap_service"
	"google.golang.org/grpc"
)

type Config struct {
	Database struct {
		Address string
		Table   string
	}
	GRPC struct {
		Address string
	}
	GRPCWeb struct {
		Address string
	}
	HTTP struct {
		Address string
		Bundle  string
	}
}

type Crime struct {
	Longitude   float64
	Latitude    float64
	Time        int64
	Description string
}

type HTTPServer interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type GRPCServer interface {
	GetServer() *grpc.Server
	Run() error
	Shutdown()
}

type GRPCWebServer interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type CrimeMapService interface {
	GetCrimes(ctx context.Context, req *cmspb.GetCrimesRequest) (*cmspb.GetCrimesResponse, error)
}

type DatabaseClient interface {
	Conn(ctx context.Context) error
	GetCrimes(ctx context.Context, minLongitude, maxLongitude, minLaitude, maxLaitude float64, minTime, maxTime int64) ([]*Crime, error)
	Close() error
}
