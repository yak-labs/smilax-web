all: stuff demo1
stuff:
	cd goapi && make
	cd smilax && go build
demo1:
	cd demo1; make
clean:
	rm -f smilax/smilax
	cd demo1; make clean
