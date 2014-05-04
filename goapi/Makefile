wrap.go : ../../chirp-lang/goapi/goapi_wrap.tcl Makefile
	tclsh ../../chirp-lang/goapi/goapi_wrap.tcl < ../../chirp-lang/goapi/api/go1.txt > wrap.go reflect bytes strings strconv fmt bufio os net/http regexp  encoding/base64 io/ioutil math math/big time
	go build
clean:
	-rm wrap.go
