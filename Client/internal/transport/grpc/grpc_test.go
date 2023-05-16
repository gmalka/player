package grpc_test

import (
	"context"
	"log"
	"net"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/test/bufconn"

	myGrpc "github.com/gmalka/Client/internal/transport/grpc"
	"github.com/gmalka/Player/build/proto"
	"google.golang.org/grpc"
)

var (
	rfm myGrpc.RemoteFileUploadService
	lis *bufconn.Listener
	srv *grpc.Server
	cancel context.CancelFunc
	conn *grpc.ClientConn
)

var _ = BeforeSuite(func() {
	lis = bufconn.Listen(1024 * 1024)
	srv = grpc.NewServer()
	songService := grpcServer{}
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
	rfm, err = myGrpc.NewGrpcClient(conn)
	if err != nil {
		log.Fatalf("grpc.Create Client %v", err)
	}
})	

var _ = Describe("Grpc", func() {

	Context("GRPC test: ", func() {
		It("Get Song", func () {
			Expect(rfm.Get("1")).To(Equal([]byte{1}))
			Expect(rfm.Get("2")).To(Equal([]byte{2}))
			b, err := rfm.Get("3")
			Expect(b).To(BeNil())
			Expect(err).ShouldNot(Succeed())
		})

		It("Get All Songs", func () {
			Expect(rfm.GetAll()).To(Equal([]string{"1", "2", "3"}))
			Expect(rfm.GetAll()).To(Equal([]string{"this"}))
			b, err := rfm.GetAll()
			Expect(b).To(BeNil())
			Expect(err).ShouldNot(Succeed())
		})
	})
})

var _ = AfterSuite(func() {
	lis.Close()
	srv.Stop()
})
