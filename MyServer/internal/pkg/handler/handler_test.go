package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/gcapizzi/moka"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gmalka/MyServer/internal/pkg/handler"
)

var (
	sm Mp3FileManagerDouble
)

var _ = Describe("Handler", func() {
	var ts *httptest.Server
	Context("Song Manager testing: ", func() {
		BeforeEach(func() {
			sm = NewMp3FileManagerDouble()
			h := handler.NewHandler(sm, false)
			ts = httptest.NewServer(h.InitRouter())
		})

		It("Get song", func() {
			AllowDouble(sm).To(ReceiveCallTo("Get").With("music1.mp3").AndReturn([]byte{1, 2, 3}, nil))
			resp, err := http.Get(fmt.Sprintf("%s/music1.mp3", ts.URL))
			Expect(err).Should(Succeed())
			Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte{1, 2, 3}))
			Expect(resp.Status).To(Equal("200 OK"))

			AllowDouble(sm).To(ReceiveCallTo("Get").With("music2.mp3").AndReturn(nil, errors.New("Some error")))
			resp, err = http.Get(fmt.Sprintf("%s/music2.mp3", ts.URL))
			Expect(err).Should(Succeed())
			Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte{83, 111, 109, 101, 32, 101, 114, 114, 111, 114, 10}))
			Expect(resp.Status).To(Equal("400 Bad Request"))
		})

		It("Save song", func() {
			AllowDouble(sm).To(ReceiveCallTo("Add").With("music1.mp3", []byte{1, 2, 3}).AndReturn(nil))
			resp, err := http.Post(fmt.Sprintf("%s/music1.mp3", ts.URL), "audio/mpeg", bytes.NewReader([]byte{1, 2, 3}))
			Expect(err).Should(Succeed())
			Expect(resp.Status).To(Equal("200 OK"))
		})

		It("Get All songs", func() {
			AllowDouble(sm).To(ReceiveCallTo("GetAll").With().AndReturn([]string{"1", "2", "3"}))
			resp, err := http.Get(fmt.Sprintf("%s", ts.URL))
			b, err := ioutil.ReadAll(resp.Body)
			Expect(err).Should(Succeed())
			s := make([]string, 0)
			json.Unmarshal(b, &s)
			Expect(s).Should(Equal([]string{"1", "2", "3"}))
			Expect(resp.Status).To(Equal("200 OK"))
		})

		It("Get All songs", func() {
			AllowDouble(sm).To(ReceiveCallTo("Delete").With("music1.mp3").AndReturn(nil))
			client := &http.Client{}
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/music1.mp3", ts.URL), nil)
			resp, err := client.Do(req)
			Expect(err).Should(Succeed())
			Expect(resp.Status).To(Equal("200 OK"))

			AllowDouble(sm).To(ReceiveCallTo("Delete").With("music2.mp3").AndReturn(errors.New("Some error")))
			client = &http.Client{}
			req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/music2.mp3", ts.URL), nil)
			resp, err = client.Do(req)
			Expect(err).Should(Succeed())
			Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte{83, 111, 109, 101, 32, 101, 114, 114, 111, 114, 10}))
			Expect(resp.Status).To(Equal("400 Bad Request"))
		})

		AfterEach(func() {
			ts.Close()
		})
	})
})
