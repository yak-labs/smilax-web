run1:
	cd ../goapi && make
	cd ../ && go install
	# smilax/smilax -d=relapz ./smilax5.tcl 2>&1 | tee _smilax5_log
	LD_LIBRARY_PATH=$(HOME)/leveldb-1.15.0 SMILAX_DATA_DIR=data ../../../../../bin/smilax-web -d=a ./smilax5.tcl 2>&1 | tee _smilax5_log

clean:
	rm -f ./_* ./,*
