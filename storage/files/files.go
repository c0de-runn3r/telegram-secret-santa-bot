package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"main/storage"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"
)

const defaultPerm = 0774

var ErrNoSavedPages = errors.New("no saved pages")

type Storage struct {
	basePath string
}

func New(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

func (s *Storage) Save(page *storage.Page) (err error) {
	filePath := path.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return fmt.Errorf("can't save page %w", err)
	}
	fName, err := fileName(page)
	if err != nil {
		return err
	}
	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	gob.NewEncoder(file).Encode(page)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	path := path.Join(s.basePath, userName)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]
	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return fmt.Errorf("can't remove file %w", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("can't remove file %w", err)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, fmt.Errorf("can't check if file exists %w", err)
	}
	path := path.Join(s.basePath, p.UserName, fileName)

	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("can't check if file exists %w", err)
	}
	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't decode page %w", err)
	}
	var p storage.Page
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("can't decode page %w", err)
	}
	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()

}
