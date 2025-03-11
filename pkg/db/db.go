package db

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/brunoga/deep"
	"github.com/studio-b12/elk"
	"github.com/vmihailenco/msgpack"
)

type LocalDb[T any] struct {
	mtx  sync.RWMutex
	data T
	dir  string
}

func OpenLocalDb[T any](dir string) (*LocalDb[T], error) {
	baseDir := filepath.Dir(dir)
	stat, err := os.Stat(baseDir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(baseDir, 0771); err != nil {
			return nil, elk.Wrap(ErrDirectory, err, "failed creating directory to database file")
		}
	} else if err != nil {
		return nil, elk.Wrap(ErrDirectory, err, "failed to stat database directory")
	}
	if stat != nil && !stat.IsDir() {
		return nil, elk.NewError(ErrDirectory, "directory to database file is not a directory")
	}

	t := &LocalDb[T]{
		dir: dir,
	}

	f, err := os.Open(dir)
	if os.IsNotExist(err) {
		return t, nil
	}
	if err != nil {
		return nil, elk.Wrap(ErrFile, err, "failed to open database file")
	}
	defer f.Close()

	err = msgpack.NewDecoder(f).Decode(&t.data)
	if err != nil {
		return nil, elk.Wrap(ErrDecode, err, "failed to decode database file")
	}

	return t, nil
}

func (t *LocalDb[T]) Store(data T) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.data = data

	f, err := os.Create(t.dir)
	if err != nil {
		return elk.Wrap(ErrFile, err, "failed to open database file for write")
	}
	defer f.Close()

	err = msgpack.NewEncoder(f).Encode(t.data)
	if err != nil {
		return elk.Wrap(ErrEncode, err, "failed to encode data to file")
	}

	return nil
}

func (t *LocalDb[T]) Load() (v T, err error) {
	t.mtx.RLock()
	defer t.mtx.RUnlock()

	v, err = deep.Copy(t.data)
	if err != nil {
		return v, elk.Wrap(ErrDeepCopy, err, "failed to deep copy internal data")
	}

	return v, nil
}
