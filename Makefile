all: _stuff_ _demo1_
_stuff_:
	cd goapi && make
	go install
_demo1_:
	cd demo1; make
clean:
	rm -f smilax-web
	cd demo1; make clean
