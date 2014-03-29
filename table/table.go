package db

import (
	"io/ioutil"
	"os"
	. "fmt"
	"strings"
	. "github.com/yak-labs/chirp-lang"
	"github.com/yak-labs/smilax-web/table/levi"
)

/*
	table get Site Table Row -> []Value
	table set Site Table Row []Value
	table match Site Table RowPattern ValuePattern -> []{row value}
*/

var Lev = levi.New("leveldb.dat")

var data_dir = os.Getenv("SMILAX_DATA_DIR")
var log_file = Sprintf("%s/table_log.txt", data_dir)

func init() {
	if len(data_dir) == 0 {
		data_dir = "."
	}
}

func cmdTableLoad(fr *Frame, argv []T) T {
	Arg0(argv)
	TableLoad()
	return Empty
}

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


func TableLoad() {
	// Start with a fresh levelDB database.
	Lev.Close()
	err := os.RemoveAll("leveldb.data")
	if err != nil {
		panic(err)
	}
	Lev = levi.New("leveldb.dat")
	text, err := ioutil.ReadFile(log_file)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(text), "\n")
	for _, line := range lines {
		line = strings.Trim(line, " \t")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		words := ParseList(line)
		leviSet(words[0].String(), words[1].String(), words[2].String(), words[3].String())
	}
}

func TableGet(site, table, row string) []T {
	key := Sprintf("/%s/%s/%s", site, table, row)
	vals := Lev.Get(key)
	return ParseList(vals)
}

func leviSet(site, table, row, values string) {
	key := Sprintf("/%s/%s/%s", site, table, row)
	Lev.Set(key, values)
}

func TableSet(site, table, row, values string) {
	leviSet(site, table, row, values)

	line := MkList(
		[]T {
			MkString(site),
			MkString(table),
			MkString(row),
			MkString(values),
		}).String()
	
	line = Sprintf("%s\n", line)
	
	f, err := os.OpenFile(log_file, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(line); err != nil {
		panic(err)
	}
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
	it := Lev.Find(prefix)
	// for it.Next() // IF mappy
	for _ = it; it.Valid() ; it.Next() {
		key := string(it.Key())
		value := string(it.Value())

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
	if err := it.GetError(); err != nil {
		panic(err)
	}

	return zz
}

var tableEnsemble = []EnsembleItem{
	EnsembleItem{Name: "load", Cmd: cmdTableLoad},
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
