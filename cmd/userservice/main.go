package main

import (
	"context"
	"io"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
	jsonlog "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/logger"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/server"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"userservice/api/authenticationservice"
	"userservice/api/authorizationservice"
	"userservice/api/userservice"
	migrationsembedder "userservice/data/mysql"
	"userservice/pkg/userservice/infrastructure"
	"userservice/pkg/userservice/infrastructure/transport"
)

var appID = "UNKNOWN"

func main() {
	logger, err := initLogger()
	if err != nil {
		stdlog.Fatal("failed to initialize logger")
	}

	config, err := parseEnv()
	if err != nil {
		logger.FatalError(err)
	}

	err = runService(config, logger)
	if err == server.ErrStopped {
		logger.Info("service is successfully stopped")
	} else if err != nil {
		logger.FatalError(err)
	}
}

func runService(config *config, logger log.MainLogger) error {
	dsn := mysql.DSN{
		User:     config.DatabaseUser,
		Password: config.DatabasePassword,
		Host:     config.DatabaseHost,
		Database: config.DatabaseName,
	}
	connector := mysql.NewConnector()
	err := connector.MigrateUp(dsn, migrationsembedder.MigrationsEmbedder)
	if err != nil {
		logger.Error(err, "failed to migrate")
	}
	err = connector.Open(dsn, config.MaxDatabaseConnections)
	if err != nil {
		return err
	}
	defer connector.Close()

	stopChan := make(chan struct{})
	listenForKillSignal(stopChan)

	container := infrastructure.NewDependencyContainer(connector.TransactionalClient(), config)

	userServiceServer := transport.NewUserServiceServer(container)
	authServiceServer := transport.NewAuthServer(container)
	serverHub := server.NewHub(stopChan)

	baseServer := grpc.NewServer(grpc.UnaryInterceptor(makeGRPCUnaryInterceptor(logger)))
	userservice.RegisterUserServiceServer(baseServer, userServiceServer)
	authenticationservice.RegisterAuthenticationServiceServer(baseServer, authServiceServer)
	authorizationservice.RegisterAuthorizationServiceServer(baseServer, authServiceServer)

	serverHub.AddServer(server.NewGrpcServer(
		baseServer,
		server.GrpcServerConfig{ServeAddress: config.ServeGRPCAddress},
		logger),
	)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var httpServer *http.Server

	serverHub.AddServer(&server.FuncServer{
		ServeImpl: func() error {
			grpcGatewayMux := runtime.NewServeMux()
			opts := []grpc.DialOption{grpc.WithInsecure()}
			err := userservice.RegisterUserServiceHandlerFromEndpoint(ctx, grpcGatewayMux, config.ServeGRPCAddress, opts)
			if err != nil {
				return err
			}

			router := mux.NewRouter()
			router.PathPrefix("/api/").Handler(grpcGatewayMux)

			router.HandleFunc("/resilience/ready", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
			}).Methods(http.MethodGet)

			httpServer = &http.Server{
				Handler:      transport.NewLoggingMiddleware(router, logger),
				Addr:         config.ServeRESTAddress,
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			logger.Info("REST server started")
			return httpServer.ListenAndServe()
		},
		StopImpl: func() error {
			cancel()
			return httpServer.Shutdown(context.Background())
		},
	})

	return serverHub.Run()
}

func initLogger() (log.MainLogger, error) {
	return jsonlog.NewLogger(&jsonlog.Config{AppName: appID}), nil
}

func listenForKillSignal(stopChan chan<- struct{}) {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		<-ch
		stopChan <- struct{}{}
	}()
}

func makeGRPCUnaryInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	loggerInterceptor := transport.NewLoggerServerInterceptor(logger)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = loggerInterceptor(ctx, req, info, handler)
		return resp, err
	}
}
