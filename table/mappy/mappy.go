package mappy

import (
	"sort"
)

type Mappy struct {
	m map[string]string
}

func New(dbName string) *Mappy {
	return &Mappy{m: make(map[string]string)}
}

// Get returns value at key, or "" if record is not found.
func (mp *Mappy) Get(k string) string {
	return mp.m[k]
}

// Set sets a record with key & value.  If value is empty, the record is deleted.
func (mp *Mappy) Set(k, v string) {
	mp.m[k] = v // TODO: delete when v empty.
}

func (mp *Mappy) Close() {
}

// pair is a single key/value tuple.
type pair struct {
	k, v string
}

// pairs is a slice of pairs that sorts by key.
type pairs []pair

func (pp pairs) Len() int {
	return len(pp)
}

func (pp pairs) Less(i, j int) bool {
	return pp[i].k < pp[j].k
}

func (pp pairs) Swap(i, j int) {
	pp[i], pp[j] = pp[j], pp[i]
}

// iterator presents db.Iterator interface.
type iterator struct {
	pp   pairs
	next int // Next position; subtract 1 for current, so it initializes to 0, for the "valid Zero Struct" principle.
}

func (it *iterator) Next() bool {
	it.next++
	return it.next-1 < len(it.pp)
}
func (it *iterator) Key() []byte {
	return []byte(it.pp[it.next-1].k)
}
func (it *iterator) Value() []byte {
	return []byte(it.pp[it.next-1].v)
}
func (it *iterator) Close() error {
	it.pp = nil
	return nil
}

func (mp *Mappy) Find(start string) /*db.*/ Iterator {
	// Copy entire map into z, starting at start.
	pp := make(pairs, 0, 8)
	for k, v := range mp.m {
		println("for", k, v)
		if k >= start {
			println("keep", k, v)
			pp = append(pp, pair{k: k, v: v})
		}
	}
	sort.Sort(pp)

	println("return", len(pp))
	return &iterator{pp: pp}
}

// Copied from code.google.com/p/leveldb-go/leveldb/db/db.go
// so we do not have to import a large unused project.
type Iterator interface {
	// Next moves the iterator to the next key/value pair.
	// It returns whether the iterator is exhausted.
	Next() bool

	// Key returns the key of the current key/value pair, or nil if done.
	// The caller should not modify the contents of the returned slice, and
	// its contents may change on the next call to Next.
	Key() []byte

	// Value returns the value of the current key/value pair, or nil if done.
	// The caller should not modify the contents of the returned slice, and
	// its contents may change on the next call to Next.
	Value() []byte

	// Close closes the iterator and returns any accumulated error. Exhausting
	// all the key/value pairs in a table is not considered to be an error.
	// It is valid to call Close multiple times. Other methods should not be
	// called after the iterator has been closed.
	Close() error
}
