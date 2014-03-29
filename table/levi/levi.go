package levi

import "github.com/jmhodges/levigo"

const CACHE_SIZE_MB = 10

type Level struct {
	Db *levigo.DB
	ROpts  *levigo.ReadOptions
	WOpts  *levigo.WriteOptions
}

func New(dbName string) *Level {
	opts := levigo.NewOptions()
	opts.SetCreateIfMissing(true)
	opts.SetCache(levigo.NewLRUCache(CACHE_SIZE_MB<<20))
	opts.SetFilterPolicy(levigo.NewBloomFilter(10))

	db, err := levigo.Open(dbName, opts)
	if err != nil {
		panic(err)
	}
	return &Level{
		Db: db,
		ROpts : levigo.NewReadOptions(),
		WOpts : levigo.NewWriteOptions(),
	}
}

// Get returns value at key, or "" if record is not found.
func (lev *Level) Get(k string) string {
	data, err := lev.Db.Get(lev.ROpts, []byte(k))
/*
	if err == lev.Db.ErrNotFound {
		return ""
	}
*/
	if err != nil {
		panic(err)
	}
	return string(data)
}

// Set sets a record with key & value.  If value is empty, the record is deleted.
func (lev *Level) Set(k, v string) {
	var err error
	if v == "" {
		err = lev.Db.Delete(lev.WOpts, []byte(k))
	} else {
		err = lev.Db.Put(lev.WOpts, []byte(k), []byte(v))
	}
	if err != nil {
		panic(err)
	}
}

func (lev *Level) Find(k string) *levigo.Iterator {
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)

	it := lev.Db.NewIterator(ro)
	defer it.Close()
	it.Seek([]byte(k))
	return it
}

func (lev *Level) Close() {
    lev.ROpts.Close()
    lev.WOpts.Close()
	lev.Db.Close()
}
