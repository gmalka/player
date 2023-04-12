package songsManager

import (
	"errors"
	"fmt"

	"google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"google.golang.org/grpc/credentials/local"
)

type Mp3FileManager interface {
	Add(name string, input []byte) error
	Get(name string) ([]byte, error)
	GetAll() []string
	Delete(name string) error
}

type RemoteFileUploadService interface {
	Get(name string) ([]byte, error)
	GetAll() ([]string, error)
}

type SongsManager interface {
	Get(name string) ([]byte, error)
	Add(name string) ([]byte, error)
	GetNext() ([]byte, error)
	GetPre() ([]byte, error)
	GetPlayList() []string
	Delete(name string) error
	DeleteLocal(name string) error
	SaveLocal(name string) error
	GetAllLocal() []string
	GetAllRemote() ([]string, error)
}

type mySongsManager struct {
	list, first, last	*song
	fileManager			Mp3FileManager
	rFileManager		RemoteFileUploadService
}

type song struct {
	name		string
	local		bool
	next, pre	*song
	data		[]byte
}

func NewSongManager(f Mp3FileManager, r RemoteFileUploadService) SongsManager {
	return mySongsManager{fileManager: f, rFileManager: r}
}

func (sm mySongsManager) Get(name string) ([]byte, error) {
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

func (sm mySongsManager) Next() ([]byte, error) {
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

func (sm mySongsManager) Pre() ([]byte, error) {
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

func (sm mySongsManager) Delete(name string) error {
	if sm.list == nil {
		return errors.New("No songs in list")
	}

	if name == "" {
		if sm.list.pre == nil {
			if sm.last == sm.list {
				sm.last = nil
			}
			cur := sm.list
			sm.list = sm.list.next
			sm.first = sm.list
			cur.pre = nil
			cur.next = nil
		} else {
			cur := sm.list
			if sm.last == cur {
				sm.last = cur.pre
			}
			sm.list.pre = sm.list.next
			cur.next = nil
			cur.pre = nil
		}
	} else {
		cur := sm.list
		for ; cur != nil && cur.name != name; cur = cur.next {
		}
		if cur == nil {
			return errors.New(fmt.Sprintf("Cant find song: %s", name))
		}

		if sm.list.pre == nil {
			if sm.last == cur {
				sm.last = nil
			}
			sm.first = cur.next
			cur.pre = nil
			cur.next = nil
		} else {
			if sm.last == cur {
				sm.last = cur.pre
			}
			cur.pre = cur.next
			cur.next = nil
			cur.pre = nil
		}
	}

	return nil
}

func (sm mySongsManager) Add(name string) error {
	_, err := sm.fileManager.Get(name)
	if err == nil {
		if sm.list == nil {
			sm.list = &song{name:name, local: true}
			sm.first = sm.list
			sm.last = sm.list
		} else {
			sm.last.next = &song{name:name, local: true, pre: sm.last}
			sm.last = sm.last.next
		}
		return nil
	}
	res, err := sm.rFileManager.Get(name)
	if err == nil {
		if sm.list == nil {
			sm.list = &song{name:name, local: false, data: res}
			sm.first = sm.list
			sm.last = sm.list
		} else {
			sm.last.next = &song{name:name, local: false, pre: sm.last, data: res}
			sm.last = sm.last.next
		}
		return err
	}

	return errors.New("Can't find the file")
}

func (sm mySongsManager) DeleteLocal(name string) error {
	return sm.fileManager.Delete(name)
}

func (sm mySongsManager) SaveLocal(name string) error {
	cur := sm.first
	for ; cur != nil && cur.name != name; cur = cur.next {
	}
	if cur != nil {
		if cur.local == true {
			return errors.New(fmt.Sprintf("File already exists: %s", name))
		} else {
			err := sm.fileManager.Add(name, cur.data)
			if err != nil {
				return err
			}
			cur.data = nil
			cur.local = true
			return errors.New(fmt.Sprintf("File with same name already exists: %s", name))
		}
	}
	data, err :=  sm.rFileManager.Get(name)
	if err != nil {
		return err
	}

	err = sm.fileManager.Add(name, data)
	return err
}

func (sm mySongsManager) GetAllLocal() []string {
	return sm.fileManager.GetAll()
}

func (sm mySongsManager) GetAllRemote() ([]string, error) {
	return sm.rFileManager.GetAll()
}