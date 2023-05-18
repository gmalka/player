package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gmalka/MyServer/build/proto"
	"github.com/gmalka/MyServer/internal/pkg/handler"
	mygrpc "github.com/gmalka/MyServer/internal/transport/grpc"
	"github.com/gmalka/MyServer/pkg/fileManager"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func Start() {
	list, err := net.Listen("tcp", fmt.Sprintf("%s:%s", viper.GetString("ip"), viper.GetString("grpc_port")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	manager, err := fileManager.NewMusicFileManager("./music", true)
	if err != nil {
		log.Fatal(err)
	}
	MusicPlayerService := &mygrpc.MusicPlayerService{Manager: manager, DoLog: true}

	proto.RegisterMusicPlayerServiceServer(grpcServer, MusicPlayerService)

	handlers := handler.NewHandler(manager, true)

	serv := new(Server)

	go serv.Run(viper.GetString("http_port"), handlers.InitRouter())
	log.Println("Waiting for connect...")
	grpcServer.Serve(list)
}

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
