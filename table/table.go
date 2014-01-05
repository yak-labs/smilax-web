package db

import . "github.com/yak-labs/chirp-lang"

// import "github.com/yak-labs/smilax-web/table/level"
import "github.com/yak-labs/smilax-web/table/mappy"

import (
	. "fmt"
	"strings"
)

/*
	table get Site Table Row -> []Value
	table set Site Table Row []Value
	table match Site Table RowPattern ValuePattern -> []{row value}
*/

func cmdTableGet(fr *Frame, argv []T) T {
	site, table, row := Arg3(argv)
	return MkList(TableGet(site.String(), table.String(), row.String()))
}

func cmdTableSet(fr *Frame, argv []T) T {
	site, table, row, values := Arg4(argv)
	TableSet(site.String(), table.String(), row.String(), values.String())
	return Empty
}

func cmdTableMatch(fr *Frame, argv []T) T {
	site, table, rowPat, valuePat := Arg4(argv)
	return MkList(TableMatch(site.String(), table.String(), rowPat.String(), valuePat.String()))
}

// For now, one global instance of Level, since it's in-memory.
var Lev = mappy.New("SINGLETON")

func TableGet(site, table, row string) []T {
	key := Sprintf("/%s/%s/%s", site, table, row)
	vals := Lev.Get(key)
	return ParseList(vals)
}

func TableSet(site, table, row, values string) {
	key := Sprintf("/%s/%s/%s", site, table, row)
	Lev.Set(key, values)
}

func TableMatch(site, table, rowPat, valuePat string) []T {
	threeSlashLen := len(site) + len(table) + 3
	keyPattern := Sprintf("/%s/%s/%s", site, table, rowPat)
	prefix := keyPattern
	star := strings.IndexAny(prefix, "*[")
	if star >= 0 {
		prefix = prefix[:star] // Shorten prefix, stopping before '*'.
	}

	zz := make([]T, 0)
	iter := Lev.Find(prefix)
	for iter.Next() {
		key := string(iter.Key())
		value := string(iter.Value())

		if !strings.HasPrefix(key, prefix) {
			break // Gone too far.
		}

		if StringMatch(keyPattern, key) {
			z := make([]T, 0, 0)
			vv := ParseList(value)
			for _, v := range vv {
				if StringMatch(valuePat, v.String()) {
					z = append(z, v)
				}
			}
			if len(z) > 0 {
				zz = append(zz, MkList([]T{MkString(key[threeSlashLen:]), MkList(z)}))
			}
		}
	}

	return zz
}

var tableEnsemble = []EnsembleItem{
	EnsembleItem{Name: "get", Cmd: cmdTableGet},
	EnsembleItem{Name: "set", Cmd: cmdTableSet},
	EnsembleItem{Name: "match", Cmd: cmdTableMatch},
}

func init() {
	if Unsafes == nil {
		Unsafes = make(map[string]Command, 333)
	}

	Unsafes["table"] = MkEnsemble(tableEnsemble)
}
