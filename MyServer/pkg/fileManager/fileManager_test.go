package fileManager_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gmalka/MyServer/pkg/fileManager"
)

var _ = Describe("FileManager", func() {
	var fm fileManager.Mp3FileManager

	Context("File Manager Testing: ", func() {
		BeforeEach(func() {
			var err error
			os.Mkdir("someDir", 0777)
			f, _ := os.Create("someDir/music1.mp3")
			f.Write([]byte{73, 68, 51, 1})
			f.Close()
			f, _ = os.Create("someDir/music2.mp3")
			f.Write([]byte{73, 68, 51, 2})
			f.Close()
			f, _ = os.Create("someDir/music3.mp3")
			f.Write([]byte{73, 68, 52})
			f.Close()

			fm, err = fileManager.NewMusicFileManager("./someDir", false)
			Expect(err).Should(Succeed())
		})

		It("Add", func() {
			Expect(fm.Add("music3.mp3", []byte{73, 68, 51, 3})).Should(Succeed())
			Expect(fm.Get("music3.mp3")).To(Equal([]byte{73, 68, 51, 3}))

			Expect(fm.Add("music4.mp3", []byte{73, 68, 52, 3})).ShouldNot(Succeed())
			b, err := fm.Get("music4.mp3")
			Expect(b).To(BeNil())
			Expect(err).ShouldNot(Succeed())
		})

		It("Get", func() {
			Expect(fm.Get("music2.mp3")).To(Equal([]byte{73, 68, 51, 2}))

			Expect(fm.Add("music3.mp3", []byte{73, 68, 51, 3})).Should(Succeed())
			Expect(fm.Get("music3.mp3")).To(Equal([]byte{73, 68, 51, 3}))

			Expect(fm.Add("music4.mp3", []byte{73, 68, 52, 3})).ShouldNot(Succeed())
			b, err := fm.Get("music4.mp3")
			Expect(b).To(BeNil())
			Expect(err).ShouldNot(Succeed())
		})

		It("GetAll", func() {
			Expect(fm.GetAll()).To(Equal([]string{"music1.mp3", "music2.mp3"}))

			Expect(fm.Add("music3.mp3", []byte{73, 68, 51, 3})).Should(Succeed())
			Expect(fm.GetAll()).To(Equal([]string{"music1.mp3", "music2.mp3", "music3.mp3"}))
		})

		It("Delete", func() {
			Expect(fm.GetAll()).To(Equal([]string{"music1.mp3", "music2.mp3"}))

			Expect(fm.Delete("music1.mp3")).Should(Succeed())
			Expect(fm.GetAll()).To(Equal([]string{"music2.mp3"}))

			Expect(fm.Delete("music2.mp3")).Should(Succeed())
			Expect(fm.GetAll()).To(Equal([]string{}))

			Expect(fm.Delete("music1.mp3")).ShouldNot(Succeed())
		})

		AfterEach(func() {
			os.RemoveAll("./someDir")
		})
	})
})
