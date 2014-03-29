package level

import "code.google.com/p/leveldb-go/leveldb"
import "code.google.com/p/leveldb-go/leveldb/db"
import "code.google.com/p/leveldb-go/leveldb/memfs"

import (
	"path/filepath"
)

type Level struct {
	Db *leveldb.DB
	Fs db.FileSystem
}

func New(dbName string) *Level {
	fs := memfs.New()
	opts := &db.Options{
		FileSystem: fs,
	}
	fp := filepath.Join("level", dbName)
	d, err := leveldb.Open(fp, opts)
	if err != nil {
		panic(err)
	}
	return &Level{Db: d, Fs: fs}
}

// Get returns value at key, or "" if record is not found.
func (lev *Level) Get(k string) string {
	v, err := lev.Db.Get([]byte(k), nil)
	if err == db.ErrNotFound {
		return ""
	}
	if err != nil {
		panic(err)
	}
	return string(v)
}

// Set sets a record with key & value.  If value is empty, the record is deleted.
func (lev *Level) Set(k, v string) {
	if v == "" {
		lev.Db.Delete([]byte(k), nil)
	} else {
		lev.Db.Set([]byte(k), []byte(v), nil)
	}
}

func (lev *Level) Find(k string) db.Iterator {
	return lev.Db.Find([]byte(k), nil)
}

func (lev *Level) Close() {
	if err := lev.Db.Close(); err != nil {
		panic(err)
	}
}
