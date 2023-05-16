package grpc_test

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"time"

	. "github.com/gcapizzi/moka"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/gmalka/MyServer/build/proto"
	myGrpc "github.com/gmalka/MyServer/internal/transport/grpc"
)

var (
	rfm proto.MusicPlayerServiceClient
	lis *bufconn.Listener
	fm Mp3GetFileManagerDouble
	srv *grpc.Server
	closer func()
	conn *grpc.ClientConn
)

var _ = BeforeSuite(func() {
	lis = bufconn.Listen(1024 * 1024)
	srv = grpc.NewServer()
	fm = NewMp3GetFileManagerDouble()
	songService := myGrpc.MusicPlayerService{Manager: fm}
	proto.RegisterMusicPlayerServiceServer(srv, &songService)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	DeferCleanup(cancel)
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	DeferCleanup(conn.Close)
	if err != nil {
		log.Fatalf("grpc.DialContext %v", err)
	}
	DeferCleanup(srv.Stop)
	rfm = proto.NewMusicPlayerServiceClient(conn)
	if err != nil {
		log.Fatalf("grpc.Create Client %v", err)
	}
})	

var _ = Describe("Grpc", func() {
	Context("Grpc Testing: ", func() {
		It("Get Song", func () {
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music1.mp3").AndReturn([]byte{1, 2, 3}, nil))
			p, err := rfm.LoadSong(context.Background(), &proto.SongRequest{Name: "music1.mp3"})
			Expect(err).Should(Succeed())
			
			b := make([]byte, 0)
			for {
				resp, err := p.Recv()
				if err == io.EOF {
					break
				}
				buf := resp.GetSong()
				b = append(b, buf...)
			}
			Expect(b).To(Equal([]byte{1, 2, 3}))

			AllowDouble(fm).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn(nil, errors.New("Some error")))
			p, err = rfm.LoadSong(context.Background(), &proto.SongRequest{Name: "music2.mp3"})
			Expect(err).Should(Succeed())
			resp, err := p.Recv()
			Expect(err).ShouldNot(Succeed())
			Expect(resp).To(BeNil())
		})

		It("Get All Songs", func () {
			AllowDouble(fm).To(ReceiveCallTo("GetAll").With().AndReturn([]string{"1", "2", "3"}))
			p, err := rfm.GetSongs(context.Background(), &proto.None{})
			Expect(err).Should(Succeed())
			
			b := make([]string, 0)
			for {
				resp, err := p.Recv()
				if err == io.EOF {
					break
				}
				buf := resp.GetName()
				b = append(b, buf)
			}
			Expect(b).To(Equal([]string{"1", "2", "3"}))
		})
	})
})

var _ = AfterSuite(func() {
	lis.Close()
	srv.Stop()
})

