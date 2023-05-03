package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cms "github.com/jialunzhai/crimemap/analytics/online/server/crimemap_service"
	env_interface "github.com/jialunzhai/crimemap/analytics/online/server/enviroment"
	"github.com/jialunzhai/crimemap/analytics/online/server/grpc_server"
	"github.com/jialunzhai/crimemap/analytics/online/server/grpc_web_server"
	"github.com/jialunzhai/crimemap/analytics/online/server/hbase_client"
	"github.com/jialunzhai/crimemap/analytics/online/server/http_server"
	"github.com/jialunzhai/crimemap/analytics/online/server/interfaces"
	real_env "github.com/jialunzhai/crimemap/analytics/online/server/real_enviroment"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	env := real_env.NewRealEnv()

	if len(os.Args) != 2 {
		log.Fatalf("Usage: cmd ${PATH_TO_CONFIG}")
	}
	if err := loadConfig(env, os.Args[1]); err != nil {
		log.Fatalf("Load config failed with error: `%v`\n", err)
	}
	config := env.GetConfig()
	if config == nil {
		log.Fatalf("Load config failed with error: `empty config file`\n")
	}

	if err := hbase_client.Register(env); err != nil {
		log.Fatalf("HBaseClient.Register failed with error: `%v`\n", err)
	}
	if config.GRPC.Address != "" {
		if err := grpc_server.Register(env); err != nil {
			log.Fatalf("GRPCServer.Register failed with error: `%v`\n", err)
		}
	}
	if err := cms.Register(env); err != nil {
		log.Fatalf("CrimeMapServer.Register failed with error: `%v`\n", err)
	}
	if config.GRPCWeb.Address != "" {
		if err := grpc_web_server.Register(env); err != nil {
			log.Fatalf("GRPCWeb.Register failed with error: `%v`\n", err)
		}
	}
	if config.HTTP.Address != "" {
		if err := http_server.Register(env); err != nil {
			log.Fatalf("HTTPServer.Register failed with error: `%v`\n", err)
		}
	}

	g, ctx := errgroup.WithContext(ctx)

	if env.GetDatabaseClient() != nil {
		log.Println("Try to connect to database with timeout...")
		ctxWithTimeout, timeout := context.WithTimeout(ctx, time.Duration(time.Second*4))
		defer timeout()
		err := env.GetDatabaseClient().Conn(ctxWithTimeout)
		if err != nil {
			log.Fatalf("DatabaseClient.Conn failed with error: `%v`\n", err)
		}
		log.Println("Connected to database.")
	}

	if env.GetGRPCServer() != nil {
		g.Go(func() error {
			err := env.GetGRPCServer().Run()
			if err != nil {
				log.Printf("GRPCServer shutdowned with error: `%v`\n", err)
			}
			log.Printf("GRPCServer gracefully shutdowned\n")
			return err
		})
	}
	if env.GetGRPCWebServer() != nil {
		g.Go(func() error {
			err := env.GetGRPCWebServer().Run()
			if err != http.ErrServerClosed {
				log.Printf("GRPCWebServer shutdowned with error: `%v`\n", err)
				return err
			}
			log.Printf("GRPCWebServer gracefully shutdowned\n")
			return err
		})
	}
	if env.GetHTTPServer() != nil {
		g.Go(func() error {
			err := env.GetHTTPServer().Run()
			if err != http.ErrServerClosed {
				log.Printf("HTTPServer shutdowned with error: `%v`\n", err)
				return err
			}
			log.Printf("HTTPServer gracefully shutdowned\n")
			return err
		})
	}

	// wait for signals
	select {
	case sig := <-sigs:
		// received signal, cancel context in the reverse order
		log.Printf("Received signal: `%v`\n", sig)
		if env.GetHTTPServer() != nil {
			env.GetHTTPServer().Shutdown(ctx)
		}
		if env.GetGRPCWebServer() != nil {
			env.GetGRPCWebServer().Shutdown(ctx)
		}
		if env.GetGRPCServer() != nil {
			env.GetGRPCServer().Shutdown()
		}
		cancel()
		break
	case <-ctx.Done():
		// context cancelled, all goroutines have returned
		break
	}

	if env.GetDatabaseClient() != nil {
		if err := env.GetDatabaseClient().Close(); err != nil {
			log.Fatalf("DatabaseClient closed with error: `%v`\n", err)
		}
		log.Printf("DatabaseClient gracefully closed\n")
	}

	// wait for all go-routines in errgroup to return
	if err := g.Wait(); err != nil {
		log.Printf("main exited with error: `%v`\n", err)
	}
}

func loadConfig(env env_interface.Env, configFile string) error {
	rawConfig, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("read configure from file `%v` failed with error: `%v`\n", configFile, err)
	}

	config := interfaces.Config{}
	if err := yaml.Unmarshal(rawConfig, &config); err != nil {
		return err
	}

	env.SetConfig(&config)
	log.Printf("Loaded config from `%v`.", configFile)
	return nil
}
