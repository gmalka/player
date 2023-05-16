package grpc_test

import (
	"testing"

	"github.com/gmalka/Player/build/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type grpcServer struct {
	i	int
	proto.UnimplementedMusicPlayerServiceServer
}

func (g *grpcServer) LoadSong(req *proto.SongRequest, stream proto.MusicPlayerService_LoadSongServer) error {

	select {
	case <-stream.Context().Done():
		return status.Error(codes.Canceled, "Stream has ended")
	default:
		switch req.Name {
		case "1":
			stream.Send(&proto.LoadSongResponse{
				Song: []byte{1},
			})
		case "2":
			stream.Send(&proto.LoadSongResponse{
				Song: []byte{2},
			})
		case "3":
			return status.Error(codes.Canceled, "Some error")
		}
	}
	return nil
}

func (g *grpcServer) GetSongs(req *proto.None, stream proto.MusicPlayerService_GetSongsServer) error {
	g.i++
	switch g.i {
	case 1:
		for _, v := range []string{"1", "2", "3"} {
			stream.Send(&proto.SongRequest{
				Name: v,
			})
		}
	case 2:
		for _, v := range []string{"this"} {
			stream.Send(&proto.SongRequest{
				Name: v,
			})
		}
	case 3:
		return status.Error(codes.Canceled, "Some error")
	}
	select {
	case <-stream.Context().Done():
		return status.Error(codes.Canceled, "Stream has ended")
	default:
		
	}
	return nil
}

func TestGrpc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Grpc Suite")
}
