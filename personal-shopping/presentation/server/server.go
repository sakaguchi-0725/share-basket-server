package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	proto "share-basket/personal-shopping/presentation/proto/gen"
	"share-basket/personal-shopping/registry"
	"syscall"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	addr string
	*grpc.Server
}

func New(addr string) *Server {
	return &Server{
		addr:   addr,
		Server: grpc.NewServer(),
	}
}

func (s *Server) MapServices(services registry.Services) {
	proto.RegisterPersonalShoppingServiceServer(s.Server, services.PersonalShopping)
}

func (s *Server) Run() {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("gRPC server is running on", s.addr)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown timeout exceeded")
		s.Stop()
	case <-done:
		log.Println("gRPC server shutdown completed")
	}
}
