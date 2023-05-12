package fileManager_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gmalka/Client/pkg/fileManager"
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

			fm, err = fileManager.NewMusicFileManager("./someDir")
			Expect(err).To(Succeed())
			Expect(fm.Contains("music1.mp3")).To(BeTrue())
			Expect(fm.Contains("music2.mp3")).To(BeTrue())
			Expect(fm.Contains("music3.mp3")).To(BeFalse())
		})

		AfterEach(func() {
			os.RemoveAll("./someDir")
		})

		It("Add", func() {
			Expect(fm.Add("music.mp3", []byte{73, 68, 51})).To(Succeed())
			Expect(os.ReadFile("./someDir/music.mp3")).To(Equal([]byte{73, 68, 51}))
			Expect(fm.Add("music.mp3", []byte{73, 68, 51})).To(Not(BeNil())) //?! уточнить ?!
		})

		It("Get", func() {
			Expect(fm.Get("music1.mp3")).To(Equal([]byte{73, 68, 51, 1}))
			Expect(fm.Get("music2.mp3")).To(Equal([]byte{73, 68, 51, 2}))
			b, err := fm.Get("music3.mp3")
			Expect(b).To(BeNil())
			Expect(err).NotTo(BeNil())
		})

		It("GetAll", func() {
			Expect(fm.GetAll()).To(Equal([]string{"music1.mp3", "music2.mp3"}))
		})

		It("Delete", func() {
			Expect(fm.Delete("music1.mp3")).To(BeNil())
			Expect(fm.Delete("music1.mp3")).NotTo(BeNil())
			Expect(fm.GetAll()).To(Equal([]string{"music2.mp3"}))
			_, err := fm.Get("music1.mp3")
			Expect(err).NotTo(BeNil())
			_, err = fm.Get("music2.mp3")
			Expect(err).To(BeNil())
		})
	})

})
