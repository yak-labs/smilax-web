all:
	rm -rf _lib_ _chirp
	cd ../goapi && make
	cd ../ && go build
	cp ../smilax-web _chirp
	mkdir _lib_
	set $$(LD_LIBRARY_PATH=/usr/local/lib ldd _chirp) ; \
	  for x ; do case $$x in /* ) cp -v "$$x" _lib_/ ;; esac ; done

docker-build:
	sudo docker build -t smilax5 .

docker-serve:
	sudo docker run  smilax5

docker-cli:
	sudo docker run -i -t  smilax5  /opt/chirp

clean:
	-rm -rf ./_* ./,*
