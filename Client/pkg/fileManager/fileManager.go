package fileManager

import (
	"errors"
	"fmt"
	"os"
)

const (
	DefaultPath string = "/Users/gmalka/Player/Client/music"
)

var CantFindFile error = errors.New("File does not exists")

type Mp3FileManager interface {
	Add(name string, input []byte) error
	Get(name string) ([]byte, error)
	GetAll() []string
	Delete(name string) error
}

type myMp3FileManager struct {
	path	string
	files	map[string]interface{}
}

func NewMusicFileManager(path string) (Mp3FileManager, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{}, len(files))
	buf := make([]byte, 3)
	mp3Signature := []byte{73, 68, 51}
	for _, f := range files {
		if !f.IsDir() {
			file, err := os.Open(fmt.Sprintf("%s/%s", path, f.Name()))
			if err != nil {
				return nil, err
			}
			_, err = file.Read(buf)
			if err != nil {
				return nil, err
			}
			for i, b := range buf {
				if b !=  mp3Signature[i] {
					return nil, err
				}
			}
			m[f.Name()] = nil
		}
	}

	return myMp3FileManager{path: path, files: m}, nil
}

func (m myMp3FileManager) Add(name string, input []byte) error {
	if _, ok := m.files[name]; ok {
		return errors.New("File is already exists")
	}
	path := fmt.Sprintf("%s/%s", m.path, name)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = file.Write(input)
	if err != nil {
		return err
	}
	m.files[name] = nil
	return nil
}

func (m myMp3FileManager) Get(name string) ([]byte, error) {
	if _, ok := m.files[name]; !ok {
		return nil, CantFindFile
	}
	path := fmt.Sprintf("%s/%s", m.path, name)
	data, err := os.ReadFile(path)
	return data, err
}

func (m myMp3FileManager) GetAll() []string {
	result := make([]string, len(m.files))
	i := 0
	for s := range m.files {
		result[i] = s
		i++
	}
	return result
}

func (m myMp3FileManager) Delete(name string) error {
	if _, ok := m.files[name]; !ok {
		return errors.New("File does not exists")
	}

	path := fmt.Sprintf("%s/%s", m.path, name)
	err := os.Remove(path)
	return err
}

/*func (m myMp3FileManager) Save(name string, data []byte) error {
	if _, ok := m.files[name]; ok {
		return errors.New(fmt.Sprintf("File, named %s already exists", name))
	}
	path := fmt.Sprintf("%s/%s", m.path, name)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}*/