package fileManager

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"sync"
)

const (
	DefaultPath string = "/Users/gmalka/Player/Client/music"
)

var CantFindFile error = errors.New("Error: File does not exists")

type Mp3FileManager interface {
	Add(name string, input []byte) error
	Get(name string) ([]byte, error)
	GetAll() []string
	Delete(name string) error
}

type myMp3FileManager struct {
	path  string
	files map[string]interface{}
	mutex *sync.Mutex
	doLog bool
}

func checkForMp3(input []byte) error {
	mp3Signature := []byte{73, 68, 51}

	buf := input[:3]

	for i, b := range buf {
		if b != mp3Signature[i] {
			return errors.New("Error: File hasn't type mp3")
		}
	}

	return nil
}

func NewMusicFileManager(path string, doLog bool) (Mp3FileManager, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{}, len(files))
	buf := make([]byte, 3)
	for _, f := range files {
		if !f.IsDir() {
			file, err := os.Open(fmt.Sprintf("%s/%s", path, f.Name()))
			if err != nil {
				return nil, err
			}
			_, err = file.Read(buf)
			if err != nil && err != io.EOF {
				return nil, err
			}
			err = checkForMp3(buf)
			if err == nil {
				m[f.Name()] = nil
			}
		}
	}

	return myMp3FileManager{path: path, files: m, mutex: &sync.Mutex{}, doLog: doLog}, nil
}

func (m myMp3FileManager) Add(name string, input []byte) error {
	err := checkForMp3(input[:3])
	if err != nil {
		if m.doLog {
			log.Println("Not mp3 file")
		}
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.files[name]; ok {
		if m.doLog {
			log.Printf("Cant save file %s, File is already exists\n", name)
		}
		return errors.New("Error: File is already exists")
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
	if m.doLog {
		log.Printf("Save song %s\n", name)
	}
	return nil
}

func (m myMp3FileManager) Get(name string) ([]byte, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.files[name]; !ok {
		if m.doLog {
			log.Printf("Cant get file %s, File is not exists\n", name)
		}
		return nil, CantFindFile
	}
	path := fmt.Sprintf("%s/%s", m.path, name)
	data, err := os.ReadFile(path)
	if m.doLog {
		log.Printf("Returning song named %s\n", name)
	}
	return data, err
}

func (m myMp3FileManager) GetAll() []string {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	result := make([]string, len(m.files))
	i := 0
	for s := range m.files {
		result[i] = s
		i++
	}
	sort.Strings(result)
	if m.doLog {
		log.Printf("Returning all songs\n")
	}
	return result
}

func (m myMp3FileManager) Delete(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.files[name]; !ok {
		if m.doLog {
			log.Printf("Cant delete file %s, File doesn't exists\n", name)
		}
		return errors.New("Error: File does not exists")
	}

	path := fmt.Sprintf("%s/%s", m.path, name)
	err := os.Remove(path)
	delete(m.files, name)
	if m.doLog {
		log.Printf("Deleted song %s\n", name)
	}
	return err
}
