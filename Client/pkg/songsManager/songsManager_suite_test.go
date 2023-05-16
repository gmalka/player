package songsManager_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/gcapizzi/moka"
)


type RemoteFileUploadService interface {
	Get(name string) ([]byte, error)
	GetAll() ([]string, error)
}

type Mp3FileManager interface {
	Add(name string, input []byte) error
	Get(name string) ([]byte, error)
	GetAll() ([]string, error)
	Delete(name string) error
	Contains(str string) bool
}

type Mp3FileManagerDouble struct {
	Double
}

func (m Mp3FileManagerDouble) Add(name string, input []byte) error {
	returnValues, _ := m.Call("Add", name, input)
	returnedRolls, _ := returnValues[0].(error)
	return returnedRolls
}

func (m Mp3FileManagerDouble) Get(name string) ([]byte, error) {
	returnValues, _ := m.Call("Get", name)
	returnedRollFirst, _ := returnValues[0].([]byte)
	returnedRollSecond, _ := returnValues[1].(error)
	return returnedRollFirst, returnedRollSecond
}

func (m Mp3FileManagerDouble) GetAll() ([]string, error) {
	returnValues, _ := m.Call("GetAll")
	returnedRollFirst, _ := returnValues[0].([]string)
	returnedRollSecond, _ := returnValues[1].(error)
	return returnedRollFirst, returnedRollSecond
}

func (m Mp3FileManagerDouble) Delete(name string) error {
	returnValues, _ := m.Call("Delete", name)
	returnedRolls, _ := returnValues[0].(error)
	return returnedRolls
}

func (m Mp3FileManagerDouble) Contains(str string) bool {
	returnValues, _ := m.Call("Contains", str)
	returnedRolls, _ := returnValues[0].(bool)
	return returnedRolls
}

type RemoteFileUploadServiceDouble struct {
	Double
}

func (r RemoteFileUploadServiceDouble) Get(name string) ([]byte, error) {
	returnValues, _ := r.Call("Get", name)
	returnedRollFirst, _ := returnValues[0].([]byte)
	returnedRollSecond, _ := returnValues[1].(error)
	return returnedRollFirst, returnedRollSecond
}

func (r RemoteFileUploadServiceDouble) GetAll() ([]string, error) {
	returnValues, _ := r.Call("GetAll")
	returnedRollFirst, _ := returnValues[0].([]string)
	returnedRollSecond, _ := returnValues[1].(error)
	return returnedRollFirst, returnedRollSecond
}

func NewRemoteFileUploadServiceDouble() RemoteFileUploadServiceDouble {
	return RemoteFileUploadServiceDouble{Double: NewStrictDouble()}
}

func NewMp3FileManagerDouble() Mp3FileManagerDouble {
	return Mp3FileManagerDouble{Double: NewStrictDouble()}
}

func TestSongsManager(t *testing.T) {
	RegisterDoublesFailHandler(func(message string, callerSkip ...int) {
		t.Fatal(message)
	})
	RegisterFailHandler(Fail)
	RunSpecs(t, "SongsManager Suite")
}
