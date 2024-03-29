package grpc

import (
	"log"

	"github.com/gmalka/MyServer/build/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	frame int = 4000000
)

type GrpcServer interface {
	LoadSong(req *proto.SongRequest, stream proto.MusicPlayerService_LoadSongServer) error
	GetSongs(req *proto.None, stream proto.MusicPlayerService_GetSongsServer) error
	mustEmbedUnimplementedMusicPlayerServiceServer()
}

type Mp3GetFileManager interface {
	Get(name string) ([]byte, error)
	GetAll() []string
}

type MusicPlayerService struct {
	DoLog bool
	Manager Mp3GetFileManager
	proto.UnimplementedMusicPlayerServiceServer
}

func (g MusicPlayerService) LoadSong(req *proto.SongRequest, stream proto.MusicPlayerService_LoadSongServer) error {
	if g.DoLog {
		log.Println("New GRPC connection")
	}
	data, err := g.Manager.Get(req.Name)
	if err != nil {
		return status.Error(codes.Canceled, err.Error())
	}

	for i := frame; i < len(data); i += frame {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Stream has ended")
		default:
			err := stream.SendMsg(&proto.LoadSongResponse{
				Song: data[:i],
			})
			data = data[i:]
			if err != nil {
				return status.Error(codes.Canceled, err.Error())
			}
		}
	}
	select {
	case <-stream.Context().Done():
		return status.Error(codes.Canceled, "Stream has ended")
	default:
		err := stream.SendMsg(&proto.LoadSongResponse{
			Song: data,
		})
		if err != nil {
			return status.Error(codes.Canceled, err.Error())
		}
	}
	return nil
}

func (g MusicPlayerService) GetSongs(req *proto.None, stream proto.MusicPlayerService_GetSongsServer) error {
	if g.DoLog {
		log.Println("New GRPC connection")
	}
	data := g.Manager.GetAll()

	for _, s := range data {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Stream has ended")
		default:
			err := stream.SendMsg(&proto.SongRequest{
				Name: s,
			})
			if err != nil {
				return status.Error(codes.Canceled, err.Error())
			}
		}
	}
	return nil
}