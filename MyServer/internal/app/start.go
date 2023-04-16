package app

import (
	"fmt"
	"log"
	"net"

	"github.com/gmalka/MyServer/build/proto"
	mygrpc "github.com/gmalka/MyServer/internal/transport/grpc"
	"github.com/gmalka/MyServer/pkg/fileManager"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func Start() {

	list, err := net.Listen("tcp", fmt.Sprintf("%s:%s", viper.GetString("ip"), viper.GetString("port")))
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