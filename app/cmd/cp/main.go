package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"cp/lib/cfg"
	grpc_ptid "cp/lib/grpc/ptid"
	grpc_validate "cp/lib/grpc/validate"
	pb_lib "cp/lib/proto"
	"cp/lib/psql"
	"github.com/caarlos0/env/v6"

	pb "cp/helloworld/helloworld"
	"google.golang.org/grpc"
)

// server is used to implement hello-world.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

func main() {

	config := cfg.Config{
		Service: pb_lib.Services_backend.String(),
	}

	if err := env.Parse(&config); err != nil {
		panic(err.Error())
	}

	_, cancel := context.WithCancel(context.Background())

	log := config.GetSrvLogger()

	if len(config.GetListenerAddress()) == 0 {
		log.Fatal().Msgf("Empty listener settings")
	}

	pool, err := psql.InitPgxPool(&config)
	if err != nil {
		log.Fatal().Msgf("Unable set database connection: %s", err.Error())
		return
	}
	// -- Настройка gRPC конфигураций сервиса ------------------------------------------------------------------------------
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpc_ptid.Interceptor,     // Инициация перехватчика uuid для трассировки
		grpc_validate.Interceptor, // Инициация перехватчика валидации запросов
	))
	// -- Запуск сервиса backend -------------------------------------------------------------------------------------------
	listener, err := net.Listen("tcp", config.GetListenerAddress())
	if err != nil {
		log.Fatal().Msgf("Unable to raise Listener (%s), error: %s", config.GetListenerAddress(), err.Error())
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		s := <-sigCh
		log.Info().Msgf("Stopping server ... (Signal: %s)", s.String())

		//grpc.GracefulStop()
		//strNode.Close()
		//orcNode.Close()
		pool.Close()
		cancel()
	}()

	pb.RegisterGreeterServer(grpcServer, &server{})
	log.Info().Msgf("Server is starting %s ... ", config.GetListenerAddress())

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Msgf("Unable to serve server, error: %s", err.Error())
	}

	wg.Wait()
	log.Info().Msg("Stopped server")
	os.Exit(1)
}
