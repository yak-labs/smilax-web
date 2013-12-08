run1:
	rm -rf _data
	cd ../goapi && make
	cd ../smilax && go build
	../smilax/smilax -d=c ./unbundle.tcl < ./run1.txt 2>&1 | tee _unbundle_log
	# smilax/smilax -d=relapz ./smilax5.tcl 2>&1 | tee _smilax5_log
	../smilax/smilax -d=a ./smilax5.tcl 2>&1 | tee _smilax5_log
clean:
	rm -rf ./_data
	rm -f ./_* ./,*
