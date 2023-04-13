package grpc

import (
	"github.com/gmalka/Server/build/proto"
)

type GrpcServer interface {
	LoadSong(req proto.SongRequest, stream proto.MusicPlayerService_LoadSongServer) error
	GetSongs(req proto.None, stream proto.MusicPlayerService_GetSongsServer) error
	mustEmbedUnimplementedMusicPlayerServiceServer()
}

type MusicPlayerService struct {
	proto.UnimplementedMusicPlayerServiceServer
}

func (g MusicPlayerService) LoadSong(req proto.SongRequest, stream proto.MusicPlayerService_LoadSongServer) error {

}

func (g MusicPlayerService) GetSongs(req proto.None, stream proto.MusicPlayerService_GetSongsServer) error {

}