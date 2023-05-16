package service_test

import (
	"errors"

	. "github.com/gcapizzi/moka"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gmalka/Client/internal/service"
)

var _ = Describe("Service", func() {
	var (
		c service.Controller
		p playerDouble
		sm songmanagerDouble
	)

	BeforeEach(func() {
		p = NewPlayerDouble()
		sm = NewSongmanagerDouble()
		c = service.NewController(p, sm)
	})

	It("SetVolume", func () {
		AllowDouble(p).To(ReceiveCallTo("SetVolume").With(10).AndReturn(nil))
		Expect(c.SetVolume("10")).Should(Succeed())

		Expect(c.SetVolume("")).ShouldNot(Succeed())
		Expect(c.SetVolume("r1")).ShouldNot(Succeed())
		AllowDouble(p).To(ReceiveCallTo("SetVolume").With(101).AndReturn(errors.New("Some error")))
		Expect(c.SetVolume("101")).ShouldNot(Succeed())
		AllowDouble(p).To(ReceiveCallTo("SetVolume").With(-1).AndReturn(errors.New("Some error")))
		Expect(c.SetVolume("-1")).ShouldNot(Succeed())
	})

	It("GetCurrent", func () {
		AllowDouble(sm).To(ReceiveCallTo("GetCurrent").With().AndReturn("Some string"))
		Expect(c.GetCurrent()).To(Equal("Some string"))
	})

	It("AddSong", func () {
		AllowDouble(sm).To(ReceiveCallTo("Add").With("name").AndReturn(nil))
		Expect(c.AddSong("name")).Should(Succeed())

		AllowDouble(sm).To(ReceiveCallTo("Add").With("another name").AndReturn(errors.New("Some error")))
		Expect(c.AddSong("another name")).ShouldNot(Succeed())
	})

	It("GetPlayList", func () {
		AllowDouble(sm).To(ReceiveCallTo("GetPlayList").With().AndReturn([]string{"m1", "m2"}, nil))
		Expect(c.GetPlayList()).To(Equal([]string{"m1", "m2"}))
	})

	It("GetPlayList With error", func () {
		AllowDouble(sm).To(ReceiveCallTo("GetPlayList").With().AndReturn(nil, errors.New("some error")))
		s, err := c.GetPlayList()
		Expect(s).To(BeNil())
		Expect(err).ShouldNot(BeNil())
	})

	It("DeleteLocal", func () {
		AllowDouble(sm).To(ReceiveCallTo("DeleteLocal").With("name").AndReturn(nil))
		Expect(c.DeleteLocal("name")).Should(Succeed())

		AllowDouble(sm).To(ReceiveCallTo("DeleteLocal").With("another name").AndReturn(errors.New("some error")))
		Expect(c.DeleteLocal("another name")).ShouldNot(Succeed())
	})

	It("DeleteSong", func () {
		AllowDouble(sm).To(ReceiveCallTo("Delete").With(0).AndReturn(nil))
		Expect(c.DeleteSong(0)).Should(Succeed())

		AllowDouble(sm).To(ReceiveCallTo("Delete").With(1).AndReturn(errors.New("some error")))
		Expect(c.DeleteSong(1)).ShouldNot(Succeed())
	})

	It("GetAllSongs", func () {
		AllowDouble(sm).To(ReceiveCallTo("GetAllLocal").With().AndReturn([]string{"1", "2"}, nil))
		AllowDouble(sm).To(ReceiveCallTo("GetAllRemote").With().AndReturn([]string{"3"}, nil))

		Expect(c.GetAllSongs("")).To(Equal([][]string{{"1", "2"}, {"3"}}))
		Expect(c.GetAllSongs("local")).To(Equal([][]string{{"1", "2"}, nil}))
		Expect(c.GetAllSongs("remote")).To(Equal([][]string{nil, {"3"}}))

		s, err := c.GetAllSongs("something else")
		Expect(s).To(BeNil())
		Expect(err).ShouldNot(BeNil())
	})

	It("SaveSong", func () {
		AllowDouble(sm).To(ReceiveCallTo("SaveLocal").With("").AndReturn(nil))
		AllowDouble(sm).To(ReceiveCallTo("SaveLocal").With("something").AndReturn(nil))
		AllowDouble(sm).To(ReceiveCallTo("SaveLocal").With("1").AndReturn(errors.New("some error")))

		Expect(c.SaveSong("")).Should(Succeed())
		Expect(c.SaveSong("something")).Should(Succeed())
		Expect(c.SaveSong("1")).ShouldNot(Succeed())
	})

	It("PreSong", func () {
		AllowDouble(sm).To(ReceiveCallTo("Pre").With().AndReturn([]byte{1}, nil))
		AllowDouble(p).To(ReceiveCallTo("Load").With([]byte{1}).AndReturn(nil))

		Expect(c.PreSong()).Should(Succeed())
	})

	It("PreSong with error in Pre", func () {
		AllowDouble(sm).To(ReceiveCallTo("Pre").With().AndReturn(nil, errors.New("some error")))
		AllowDouble(p).To(ReceiveCallTo("Load").With([]byte{1}).AndReturn(nil))

		Expect(c.PreSong()).ShouldNot(Succeed())
	})

	It("PreSong with error in Load", func () {
		AllowDouble(sm).To(ReceiveCallTo("Pre").With().AndReturn([]byte{1}, nil))
		AllowDouble(p).To(ReceiveCallTo("Load").With([]byte{1}).AndReturn(errors.New("some error")))

		Expect(c.PreSong()).ShouldNot(Succeed())
	})

	It("NextSong", func () {
		AllowDouble(sm).To(ReceiveCallTo("Next").With().AndReturn([]byte{1}, nil))
		AllowDouble(p).To(ReceiveCallTo("Load").With([]byte{1}).AndReturn(nil))

		Expect(c.NextSong()).Should(Succeed())
	})

	It("NextSong with error in Next", func () {
		AllowDouble(sm).To(ReceiveCallTo("Next").With().AndReturn(nil, errors.New("some error")))
		AllowDouble(p).To(ReceiveCallTo("Load").With([]byte{1}).AndReturn(nil))

		Expect(c.NextSong()).ShouldNot(Succeed())
	})

	It("NextSong with error in Load", func () {
		AllowDouble(sm).To(ReceiveCallTo("Next").With().AndReturn([]byte{1}, nil))
		AllowDouble(p).To(ReceiveCallTo("Load").With([]byte{1}).AndReturn(errors.New("some error")))

		Expect(c.NextSong()).ShouldNot(Succeed())
	})

	Context("PlaySong", func () {
		It ("While playing", func () {
			AllowDouble(p).To(ReceiveCallTo("IsPlaying").With().AndReturn(true))
			AllowDouble(p).To(ReceiveCallTo("Play").With().AndReturn())

			Expect(c.PlaySong("some name")).Should(Succeed())
		})

		It ("Not playing", func () {
			AllowDouble(p).To(ReceiveCallTo("IsPlaying").With().AndReturn(false))
			AllowDouble(sm).To(ReceiveCallTo("Get").With("some name").AndReturn([]byte{1, 2}, nil))
			AllowDouble(p).To(ReceiveCallTo("Load").With([]byte{1, 2}).AndReturn(nil))
			AllowDouble(p).To(ReceiveCallTo("Play").With().AndReturn())

			Expect(c.PlaySong("some name")).Should(Succeed())
		})

		It ("Not playing, error while Get", func () {
			AllowDouble(p).To(ReceiveCallTo("IsPlaying").With().AndReturn(false))
			AllowDouble(sm).To(ReceiveCallTo("Get").With("some name").AndReturn(nil, errors.New("some error")))

			Expect(c.PlaySong("some name")).ShouldNot(Succeed())
		})

		It ("Not playing, error while Load", func () {
			AllowDouble(p).To(ReceiveCallTo("IsPlaying").With().AndReturn(false))
			AllowDouble(sm).To(ReceiveCallTo("Get").With("some name").AndReturn([]byte{1, 2}, nil))
			AllowDouble(p).To(ReceiveCallTo("Load").With([]byte{1, 2}).AndReturn(errors.New("some error")))

			Expect(c.PlaySong("some name")).ShouldNot(Succeed())
		})
	})
})
