package grpc_test

import (
	"testing"

	. "github.com/gcapizzi/moka"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Mp3GetFileManager interface {
	Get(name string) ([]byte, error)
	GetAll() []string
}

type Mp3GetFileManagerDouble struct {
	Double
}

func (m Mp3GetFileManagerDouble) Get(name string) ([]byte, error) {
	returnValues, _ := m.Call("Get", name)
	returnedRollFirst, _ := returnValues[0].([]byte)
	returnedRollSecond, _ := returnValues[1].(error)
	return returnedRollFirst, returnedRollSecond
}

func (m Mp3GetFileManagerDouble) GetAll() []string {
	returnValues, _ := m.Call("GetAll")
	returnedRollFirst, _ := returnValues[0].([]string)
	return returnedRollFirst
}

func NewMp3GetFileManagerDouble() Mp3GetFileManagerDouble {
	return Mp3GetFileManagerDouble{Double: NewStrictDouble()}
}

/*type grpcServer struct {
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
}*/

func TestGrpc(t *testing.T) {
	RegisterDoublesFailHandler(func(message string, callerSkip ...int) {
		t.Fatal(message)
	})
	RegisterFailHandler(Fail)
	RunSpecs(t, "Grpc Suite")
}
