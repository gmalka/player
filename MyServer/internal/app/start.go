package app

import (
	"log"
	"net"

	"github.com/gmalka/MyServer/build/proto"
	"github.com/gmalka/MyServer/pkg/fileManager"
	mygrpc "github.com/gmalka/MyServer/internal/transport/grpc"
	"google.golang.org/grpc"
)

func Start() {
	list, err := net.Listen("tcp", "localhost:9879")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	manager, err := fileManager.NewMusicFileManager("/Users/gmalka/Player/MyServer/music")
	if err != nil {
		log.Fatal(err)
	}
	MusicPlayerService := &mygrpc.MusicPlayerService{Manager: manager}

	proto.RegisterMusicPlayerServiceServer(grpcServer, MusicPlayerService)

	grpcServer.Serve(list)
}