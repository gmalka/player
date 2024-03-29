package songsManager

import (
	"errors"
	"fmt"
)

type Mp3FileManager interface {
	Add(name string, input []byte) error
	Get(name string) ([]byte, error)
	GetAll() ([]string, error)
	Delete(name string) error
	Contains(str string) bool
}

type RemoteFileUploadService interface {
	Get(name string) ([]byte, error)
	GetAll() ([]string, error)
}

type SongsManager interface {
	Get(name string) ([]byte, error)
	Add(name string) error
	Next() ([]byte, error)
	Pre() ([]byte, error)
	GetPlayList() []string
	GetCurrent() string
	Delete(id int) error
	DeleteLocal(name string) error
	SaveLocal(name string) error
	GetAllLocal() ([]string, error)
	GetAllRemote() ([]string, error)
}

type mySongsManager struct {
	id                int
	list, first, last *song
	fileManager       Mp3FileManager
	rFileManager      RemoteFileUploadService
}

type song struct {
	id        int
	name      string
	local     bool
	next, pre *song
	data      []byte
}

func NewSongManager(f Mp3FileManager, r RemoteFileUploadService) SongsManager {
	return &mySongsManager{fileManager: f, rFileManager: r, id: 1}
}

func (sm *mySongsManager) Get(name string) ([]byte, error) {
	if name == "" {
		if sm.list == nil {
			return nil, errors.New("No songs in playlist")
		}
		if sm.list.local == true {
			return sm.fileManager.Get(sm.list.name)
		}
		return sm.list.data, nil
	}
	result, err := sm.fileManager.Get(name)
	if err == nil {
		return result, nil
	}
	result, err = sm.rFileManager.Get(name)
	if err == nil {
		return result, nil
	}
	return nil, err
}

func (sm *mySongsManager) Add(name string) error {
	_, err := sm.fileManager.Get(name)
	if err == nil {
		if sm.list == nil {
			sm.list = &song{name: name, local: true, id: sm.id}
			sm.id++
			sm.first = sm.list
			sm.last = sm.list
		} else {
			sm.last.next = &song{name: name, local: true, pre: sm.last, id: sm.id}
			sm.id++
			sm.last = sm.last.next
		}
		return nil
	}
	res, err := sm.rFileManager.Get(name)
	if err == nil {
		if sm.list == nil {
			sm.list = &song{name: name, local: false, data: res, id: sm.id}
			sm.id++
			sm.first = sm.list
			sm.last = sm.list
		} else {
			sm.last.next = &song{name: name, local: false, pre: sm.last, data: res, id: sm.id}
			sm.id++
			sm.last = sm.last.next
		}
		return err
	}

	return errors.New("Can't find the file")
}

func (sm *mySongsManager) Next() ([]byte, error) {
	if sm.list == nil {
		return nil, errors.New("Empty list")
	}
	result := sm.list.next
	if result == nil {
		return nil, errors.New("Already last song in list")
	}
	sm.list = sm.list.next
	if result.local == false {
		return result.data, nil
	} else {
		return sm.Get(result.name)
	}
}

func (sm *mySongsManager) Pre() ([]byte, error) {
	if sm.list == nil {
		return nil, errors.New("Empty list")
	}
	result := sm.list.pre
	if result == nil {
		return nil, errors.New("Already last song in list")
	}
	sm.list = sm.list.pre
	if result.local == false {
		return result.data, nil
	} else {
		return sm.Get(result.name)
	}
}

func (sm *mySongsManager) GetPlayList() []string {
	s := make([]string, 0, 10)
	cur := sm.first
	for ; cur != nil; cur = cur.next {
		s = append(s, cur.name)
	}
	if len(s) == 0 {
		return nil
	}
	return s
}

func (sm *mySongsManager) Delete(id int) error {
	if sm.list == nil {
		return errors.New("No songs in list")
	}
	if id <= 0 {
		return errors.New("Incorrect id")
	}

	cur := sm.first
	for i := 1; cur != nil && i < id; cur, i = cur.next, i+1 {
	}
	if cur == nil {
		return errors.New(fmt.Sprintf("Cant find song with id: %d", id))
	}
	sm.deleteFromList(cur)

	return nil
}

func (sm *mySongsManager) DeleteLocal(name string) error {
	err := sm.fileManager.Delete(name)
	if err != nil {
		return err
	}
	if sm.list == nil {
		return nil
	}
	cur := sm.first
	for ; cur != nil; cur = cur.next {
		if cur.name == name && cur.local == true {
			sm.deleteFromList(cur)
		}
	}
	return nil
}

func (sm *mySongsManager) deleteFromList(cur *song) {
	if sm.list.id == cur.id {
		if cur.next != nil {
			sm.list = sm.list.next
		} else {
			sm.list = sm.list.pre
		}
	}

	if cur.pre == nil {
		if cur.next == nil {
			sm.last, sm.first, sm.list = nil, nil, nil
			return
		} else {
			cur.next.pre = nil
		}
		sm.first = cur.next
		cur.pre = nil
		cur.next = nil
	} else {
		if cur.next == nil {
			sm.last = cur.pre
			sm.last.next = nil
		} else {
			cur.next.pre = cur.pre
		}
		cur.pre.next = cur.next
		cur.next = nil
		cur.pre = nil
	}
}

func (sm *mySongsManager) SaveLocal(name string) error {
	if name == "" {
		name = sm.list.name
	}
	cur := sm.first
	for ; cur != nil; cur = cur.next {
		if cur.name == name && cur.local == true {
			return errors.New(fmt.Sprintf("File already exists: %s", name))
		}
	}
	data, err := sm.rFileManager.Get(name)
	if err != nil {
		return err
	}

	err = sm.fileManager.Add(name, data)
	return err
}

func (sm *mySongsManager) GetAllLocal() ([]string, error) {
	return sm.fileManager.GetAll()
}

func (sm *mySongsManager) GetAllRemote() ([]string, error) {
	s, err := sm.rFileManager.GetAll()
	if err != nil {
		return nil, err
	}
	j := 0
	for i := 0; i < len(s); i++ {
		if !sm.fileManager.Contains(s[i]) {
			s[j] = s[i]
			j++
		}
	}
	s = s[:j]
	return s, nil
}

func (sm *mySongsManager) GetCurrent() string {
	if sm.list == nil {
		return "nothing"
	}
	return sm.list.name
}
