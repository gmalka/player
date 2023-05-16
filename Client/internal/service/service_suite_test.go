package service_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/gcapizzi/moka"
)

type songmanager interface {
	Get(name string) ([]byte, error)
	Add(name string) error
	Next() ([]byte, error)
	Pre() ([]byte, error)
	GetCurrent() string
	GetPlayList() []string
	Delete(id int) error
	DeleteLocal(name string) error
	SaveLocal(name string) error
	GetAllLocal() ([]string, error)
	GetAllRemote() ([]string, error)
}

type songmanagerDouble struct {
	Double
}

func (m songmanagerDouble) Get(name string) ([]byte, error) {
	returnValues, _ := m.Call("Get", name)
	returnedRollFirst, _ := returnValues[0].([]byte)
	returnedRollSecond, _ := returnValues[1].(error)
	return returnedRollFirst, returnedRollSecond
}

func (m songmanagerDouble) Add(name string) error {
	returnValues, _ := m.Call("Add", name)
	returnedRolls, _ := returnValues[0].(error)
	return returnedRolls
}

func (m songmanagerDouble) Next() ([]byte, error) {
	returnValues, _ := m.Call("Next")
	returnedRollFirst, _ := returnValues[0].([]byte)
	returnedRollSecond, _ := returnValues[1].(error)
	return returnedRollFirst, returnedRollSecond
}

func (m songmanagerDouble) Pre() ([]byte, error) {
	returnValues, _ := m.Call("Pre")
	returnedRollFirst, _ := returnValues[0].([]byte)
	returnedRollSecond, _ := returnValues[1].(error)
	return returnedRollFirst, returnedRollSecond
}

func (m songmanagerDouble) GetCurrent() string {
	returnValues, _ := m.Call("GetCurrent")
	returnedRollFirst, _ := returnValues[0].(string)
	return returnedRollFirst
}

func (m songmanagerDouble) GetPlayList() []string {
	returnValues, _ := m.Call("GetPlayList")
	returnedRollFirst, _ := returnValues[0].([]string)
	return returnedRollFirst
}

func (m songmanagerDouble) Delete(id int) error {
	returnValues, _ := m.Call("Delete", id)
	returnedRollFirst, _ := returnValues[0].(error)
	return returnedRollFirst
}

func (m songmanagerDouble) DeleteLocal(name string) error {
	returnValues, _ := m.Call("DeleteLocal", name)
	returnedRollFirst, _ := returnValues[0].(error)
	return returnedRollFirst
}

func (m songmanagerDouble) SaveLocal(name string) error {
	returnValues, _ := m.Call("SaveLocal", name)
	returnedRollFirst, _ := returnValues[0].(error)
	return returnedRollFirst
}

func (m songmanagerDouble) GetAllLocal() ([]string, error) {
	returnValues, _ := m.Call("GetAllLocal")
	returnedRollFirst, _ := returnValues[0].([]string)
	returnedRollSecond, _ := returnValues[1].(error)
	return returnedRollFirst, returnedRollSecond
}

func (m songmanagerDouble) GetAllRemote() ([]string, error) {
	returnValues, _ := m.Call("GetAllRemote")
	returnedRollFirst, _ := returnValues[0].([]string)
	returnedRollSecond, _ := returnValues[1].(error)
	return returnedRollFirst, returnedRollSecond
}

type player interface {
	Play()
	Pause()
	Stop()
	SetVolume(v int) error
 	Load(data []byte) error
	IsPlaying() bool
}

type playerDouble struct {
	Double
}

func (m playerDouble) Play() {
	m.Call("Play")
}

func (m playerDouble) Pause() {
	m.Call("Play")
}

func (m playerDouble) Stop() {
	m.Call("Play")
}

func (m playerDouble) SetVolume(v int) error {
	returnValues, _ := m.Call("SetVolume", v)
	returnedRollFirst, _ := returnValues[0].(error)
	return returnedRollFirst
}

func (m playerDouble) Load(data []byte) error {
	returnValues, _ := m.Call("Load", data)
	returnedRollFirst, _ := returnValues[0].(error)
	return returnedRollFirst
}

func (m playerDouble) IsPlaying() bool {
	returnValues, _ := m.Call("IsPlaying")
	returnedRollFirst, _ := returnValues[0].(bool)
	return returnedRollFirst
}

func TestService(t *testing.T) {
	RegisterDoublesFailHandler(func(message string, callerSkip ...int) {
		t.Fatal(message)
	})
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

func NewPlayerDouble() playerDouble {
	return playerDouble{Double: NewStrictDouble()}
}

func NewSongmanagerDouble() songmanagerDouble {
	return songmanagerDouble{Double: NewStrictDouble()}
}