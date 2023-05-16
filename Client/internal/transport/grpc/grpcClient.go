package grpc

import (
	"context"
	"errors"
	"io"

	"google.golang.org/grpc"

	"github.com/gmalka/Player/build/proto"
)

var CantFindFile error = errors.New("File does not exists")

type RemoteFileUploadService interface {
	Get(name string) ([]byte, error)
	GetAll() ([]string, error)
}

type myGrpcClient struct {
	client proto.MusicPlayerServiceClient
}

/*func NewGrpcClient(ip, port string) (RemoteFileUploadService, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	path := fmt.Sprintf("%s:%s", ip, port)
	conn, err := grpc.Dial(path, opts...)
	if err != nil {
		return nil, err
	}

	client := proto.NewMusicPlayerServiceClient(conn)
	return myGrpcClient{client: client}, nil
}*/

func NewGrpcClient(cc grpc.ClientConnInterface) (RemoteFileUploadService, error) {
	client := proto.NewMusicPlayerServiceClient(cc)
	return myGrpcClient{client: client}, nil
}


func (g myGrpcClient) Get(name string) ([]byte, error) {
	res, err := g.client.LoadSong(context.Background(), &proto.SongRequest{Name: name})
	if err != nil {
		return nil, CantFindFile
	}

	result := make([]byte, 0, 1000)
	for {
		resp, err := res.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		result = append(result, resp.GetSong()...)
	}

	return result, nil
}

func (g myGrpcClient) GetAll() ([]string, error) {
	res, err := g.client.GetSongs(context.Background(), &proto.None{})
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, 10)
	for {
		resp, err := res.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		result = append(result, resp.GetName())
	}

	return result, nil
}