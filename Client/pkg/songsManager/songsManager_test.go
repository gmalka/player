package songsManager_test

import (
	"errors"

	. "github.com/gcapizzi/moka"
	"github.com/gmalka/Client/pkg/songsManager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SongsManager", func() {
	var sm SongsManager
	var fm Mp3FileManagerDouble
	var rmf RemoteFileUploadServiceDouble
	Context("Song Manager testing: ", func() {
		BeforeEach(func() {
			fm = NewMp3FileManagerDouble()
			rmf = NewRemoteFileUploadServiceDouble()
			sm = songsManager.NewSongManager(fm, rmf)
		})

		It("Get", func () {
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music1.mp3").AndReturn([]byte{1}, nil))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn([]byte{2}, nil))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))

			Expect(sm.Get("music1.mp3")).To(Equal([]byte{1}))
			Expect(sm.Get("music2.mp3")).To(Equal([]byte{2}))

			b, err := sm.Get("music3.mp3")
			Expect(b).To(BeNil())
			Expect(err).NotTo(BeNil())
		})

		It("Add", func () {
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music1.mp3").AndReturn([]byte{1}, nil))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn([]byte{2}, nil))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))

			Expect(sm.Add("music1.mp3")).To(BeNil())
			Expect(sm.GetPlayList()).To(Equal([]string{"music1.mp3"}))
			Expect(sm.Add("music2.mp3")).To(BeNil())
			Expect(sm.GetPlayList()).To(Equal([]string{"music1.mp3", "music2.mp3"}))
			Expect(sm.Add("music3.mp3")).NotTo(BeNil())
			Expect(sm.GetPlayList()).To(Equal([]string{"music1.mp3", "music2.mp3"}))
		})

		It("Next", func () {
			b, err := sm.Next()
			Expect(b).To(BeNil())
			Expect(err).NotTo(BeNil())
			
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music1.mp3").AndReturn([]byte{1}, nil))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn([]byte{2}, nil))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))

			sm.Add("music1.mp3")
			b, err = sm.Next()
			Expect(b).To(BeNil())
			Expect(err).NotTo(BeNil())
			sm.Add("music2.mp3")

			Expect(sm.Next()).To(Equal([]byte{2}))

			b, err = sm.Next()
			Expect(b).To(BeNil())
			Expect(err).NotTo(BeNil())
		})

		It("Pre", func () {
			b, err := sm.Pre()
			Expect(b).To(BeNil())
			Expect(err).NotTo(BeNil())
			
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music1.mp3").AndReturn([]byte{1}, nil))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn([]byte{2}, nil))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))

			sm.Add("music1.mp3")
			b, err = sm.Pre()
			Expect(b).To(BeNil())
			Expect(err).NotTo(BeNil())

			sm.Add("music2.mp3")
			sm.Next()
			
			Expect(sm.Pre()).To(Equal([]byte{1}))

			b, err = sm.Pre()
			Expect(b).To(BeNil())
			Expect(err).ShouldNot(Succeed())
		})

		It("GetPlayList", func () {
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music1.mp3").AndReturn([]byte{1}, nil))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn([]byte{2}, nil))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))

			Expect(sm.GetPlayList()).To(BeNil())

			sm.Add("music1.mp3")
			Expect(sm.GetPlayList()).To(Equal([]string{"music1.mp3"}))
			sm.Add("music2.mp3")
			Expect(sm.GetPlayList()).To(Equal([]string{"music1.mp3", "music2.mp3"}))
		})

		It("Delete", func () {
			Expect(sm.Delete(1)).ShouldNot(Succeed())
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music1.mp3").AndReturn([]byte{1}, nil))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn([]byte{2}, nil))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))

			sm.Add("music1.mp3")
			sm.Add("music2.mp3")
			Expect(sm.GetPlayList()).To(Equal([]string{"music1.mp3", "music2.mp3"}))
			Expect(sm.Delete(2)).Should(Succeed())
			Expect(sm.GetPlayList()).To(Equal([]string{"music1.mp3"}))
			Expect(sm.Delete(1)).Should(Succeed())
			Expect(sm.GetPlayList()).To(BeNil())
			Expect(sm.Delete(1)).ShouldNot(Succeed())
			Expect(sm.Delete(0)).ShouldNot(Succeed())
		})

		It("DeleteLocal", func () {
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music1.mp3").AndReturn([]byte{1}, nil))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn([]byte{2}, nil))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn([]byte{3}, nil))

			AllowDouble(fm).To(ReceiveCallTo("Delete").With("music1.mp3").AndReturn(nil))
			AllowDouble(fm).To(ReceiveCallTo("Delete").With("music3.mp3").AndReturn(nil))
			AllowDouble(fm).To(ReceiveCallTo("Delete").With("music.mp3").AndReturn(errors.New("Some error")))

			sm.Add("music1.mp3")
			sm.Add("music2.mp3")
			sm.Add("music3.mp3")
			Expect(sm.DeleteLocal("music1.mp3")).Should(Succeed())
			Expect(sm.GetPlayList()).To(Equal([]string{"music2.mp3", "music3.mp3"}))
			sm.Next()
			Expect(sm.DeleteLocal("music3.mp3")).Should(Succeed())
			Expect(sm.GetPlayList()).To(Equal([]string{"music2.mp3"}))
			Expect(sm.DeleteLocal("music.mp3")).ShouldNot(Succeed())
			Expect(sm.GetPlayList()).To(Equal([]string{"music2.mp3"}))
		})

		It("SaveLocal", func () {
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music1.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music1.mp3").AndReturn([]byte{1}, nil))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn([]byte{2}, nil))

			sm.Add("music1.mp3")
			sm.Add("music2.mp3")
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music1.mp3").AndReturn([]byte{1}, nil))
			AllowDouble(fm).To(ReceiveCallTo("Add").With("music1.mp3", []byte{1}).AndReturn(nil))
			Expect(sm.SaveLocal("music1.mp3")).Should(Succeed())
			Expect(sm.SaveLocal("music2.mp3")).ShouldNot(Succeed())
		})

		It("GetAllLocal", func () {
			AllowDouble(fm).To(ReceiveCallTo("GetAll").With().AndReturn([]string{"1", "2"}, nil))
			Expect(sm.GetAllLocal()).Should(Equal([]string{"1", "2"}))
		})

		It("GetAllRemote", func () {
			AllowDouble(rmf).To(ReceiveCallTo("GetAll").With().AndReturn([]string{"1", "2", "3"}, nil))
			AllowDouble(fm).To(ReceiveCallTo("Contains").With("1").AndReturn(true))
			AllowDouble(fm).To(ReceiveCallTo("Contains").With("2").AndReturn(false))
			AllowDouble(fm).To(ReceiveCallTo("Contains").With("3").AndReturn(false))
			Expect(sm.GetAllRemote()).Should(Equal([]string{"2", "3"}))
		})

		It("GetCurrent", func () {
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music1.mp3").AndReturn([]byte{1}, nil))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(fm).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn([]byte{2}, nil))
			AllowDouble(rmf).To(ReceiveCallTo("Get").With("music3.mp3").AndReturn(nil, errors.New("Some error")))

			Expect(sm.GetCurrent()).Should(Equal("nothing"))
			sm.Add("music1.mp3")
			Expect(sm.GetCurrent()).Should(Equal("music1.mp3"))
			sm.Add("music2.mp3")
			Expect(sm.GetCurrent()).Should(Equal("music1.mp3"))
			sm.Next()
			Expect(sm.GetCurrent()).Should(Equal("music2.mp3"))
		})
	})
})
