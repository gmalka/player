package fileManager

import (
	"errors"
	"fmt"
	"os"
)

const (
	DefaultPath string = "/Users/gmalka/Player/music"
)

type myMusicFileManager struct {
	path	string
	files	map[string]interface{}
}

func NewMusicFileManager(path string) (MusicFileManager, error) {
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

	return myMusicFileManager{path: path, files: m}, nil
}

type MusicFileManager interface {
	Add(name string, input []byte) error
	Open(name string) ([]byte, error)
	GetAll() []string
	Delete(name string) error
}

func (m myMusicFileManager) Add(name string, input []byte) error {
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

func (m myMusicFileManager) Open(name string) ([]byte, error) {
	if _, ok := m.files[name]; !ok {
		return nil, errors.New("File does not exists")
	}
	path := fmt.Sprintf("%s/%s", m.path, name)
	data, err := os.ReadFile(path)
	return data, err
}

func (m myMusicFileManager) GetAll() []string {
	result := make([]string, len(m.files))
	i := 0
	for s, _ := range m.files {
		result[i] = s
		i++
	}
	return result
}

func (m myMusicFileManager) Delete(name string) error {
	if _, ok := m.files[name]; !ok {
		return errors.New("File does not exists")
	}

	path := fmt.Sprintf("%s/%s", m.path, name)
	err := os.Remove(path)
	return err
}