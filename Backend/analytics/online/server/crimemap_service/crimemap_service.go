package crimemap_server

import (
	"context"
	"fmt"

	cmspb "github.com/jialunzhai/crimemap/analytics/online/proto/crimemap_service"
	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
)

type CrimeMapService struct {
	env env_interface.Env
}

func Register(env env_interface.Env) error {
	grpcServer := env.GetGRPCServer()
	if grpcServer == nil {
		return fmt.Errorf("GRPCServer.Register should be called before CrimeMapService.Register")
	}
	s, err := NewCrimeMapService(env)
	if err != nil {
		return err
	}
	cmspb.RegisterCrimeMapServer(grpcServer.GetServer(), s)
	env.SetCrimeMapService(s)
	return nil
}

func NewCrimeMapService(env env_interface.Env) (*CrimeMapService, error) {
	return &CrimeMapService{
		env: env,
	}, nil
}

func (s *CrimeMapService) GetCrimes(ctx context.Context, req *cmspb.GetCrimesRequest) (*cmspb.GetCrimesResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("crimemap service warnning: empty gRPC request for GetCrimes")
	}
	if req.LongitudeMin > req.LongitudeMax || req.LatitudeMin > req.LatitudeMax || req.TimeMin > req.TimeMax {
		return nil, fmt.Errorf("crimemap service warnning: bad gRPC arguments for GetCrimes")
	}

	rsp := cmspb.GetCrimesResponse{
		Crimes: make([]*cmspb.Crime, 0),
	}
	crimes, err := s.env.GetDatabaseClient().GetCrimes(ctx, req.LongitudeMin, req.LongitudeMax, req.LatitudeMin, req.LatitudeMax, req.TimeMin, req.TimeMax)
	if err != nil {
		return nil, err
	}
	for _, crime := range crimes {
		if crime == nil {
			return nil, fmt.Errorf("crimemap service internal error: database client returned result contains nil pointer")
		}
		rsp.Crimes = append(rsp.Crimes, &cmspb.Crime{
			Time:        crime.Time,
			Longitude:   crime.Longitude,
			Latitude:    crime.Latitude,
			Description: crime.Description,
		})
	}

	// rsp.Crimes = append(rsp.Crimes, &cmspb.Crime{
	// 	Time:        10002,
	// 	Longitude:   7.809007,
	// 	Latitude:    51.678418,
	// 	Description: "test1",
	// })
	// rsp.Crimes = append(rsp.Crimes, &cmspb.Crime{
	// 	Time:        10002,
	// 	Longitude:   7.810007,
	// 	Latitude:    51.679418,
	// 	Description: "test2",
	// })
	// rsp.Crimes = append(rsp.Crimes, &cmspb.Crime{
	// 	Time:        10002,
	// 	Longitude:   7.811007,
	// 	Latitude:    51.680418,
	// 	Description: "test3",
	// })
	// rsp.Crimes = append(rsp.Crimes, &cmspb.Crime{
	// 	Time:        10002,
	// 	Longitude:   7.812007,
	// 	Latitude:    51.681418,
	// 	Description: "test4",
	// })
	// rsp.Crimes = append(rsp.Crimes, &cmspb.Crime{
	// 	Time:        10002,
	// 	Longitude:   7.812007,
	// 	Latitude:    51.681418,
	// 	Description: "test5",
	// })
	// rsp.Crimes = append(rsp.Crimes, &cmspb.Crime{
	// 	Time:        10002,
	// 	Longitude:   7.812007,
	// 	Latitude:    51.681418,
	// 	Description: "test5",
	// })
	// rsp.Crimes = append(rsp.Crimes, &cmspb.Crime{
	// 	Time:        10002,
	// 	Longitude:   7.812007,
	// 	Latitude:    51.681418,
	// 	Description: "test5",
	// })
	// rsp.Crimes = append(rsp.Crimes, &cmspb.Crime{
	// 	Time:        10002,
	// 	Longitude:   7.812007,
	// 	Latitude:    51.681418,
	// 	Description: "test5",
	// })
	// rsp.Crimes = append(rsp.Crimes, &cmspb.Crime{
	// 	Time:        10002,
	// 	Longitude:   7.812007,
	// 	Latitude:    51.681418,
	// 	Description: "test5",
	// })
	// rsp.Crimes = append(rsp.Crimes, &cmspb.Crime{
	// 	Time:        10002,
	// 	Longitude:   7.812007,
	// 	Latitude:    51.681418,
	// 	Description: "test5",
	// })
	// rsp.Crimes = append(rsp.Crimes, &cmspb.Crime{
	// 	Time:        10002,
	// 	Longitude:   7.812007,
	// 	Latitude:    51.681418,
	// 	Description: "test5",
	// })
	// rsp.Crimes = append(rsp.Crimes, &cmspb.Crime{
	// 	Time:        10002,
	// 	Longitude:   7.812007,
	// 	Latitude:    51.681418,
	// 	Description: "test5",
	// })
	return &rsp, nil
}
