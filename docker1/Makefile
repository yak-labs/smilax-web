all:
	rm -rf _opt_
	cd ../goapi && make
	cd ../ && go build
	mkdir _opt_
	cp ../smilax-web _opt_/
	set $$(LD_LIBRARY_PATH=/usr/local/lib ldd _opt_/smilax-web) ; \
	  for x ; do case $$x in /* ) cp -v "$$x" _opt_/ ;; esac ; done

clean:
	rm -f ./_* ./,*
