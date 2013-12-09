run1:
	rm -rf _data
	cd ../goapi && make
	cd ../ && go install
	../../../../../bin/smilax-web -d=c ./unbundle.tcl < ./run1.txt 2>&1 | tee _unbundle_log
	# smilax/smilax -d=relapz ./smilax5.tcl 2>&1 | tee _smilax5_log
	../../../../../bin/smilax-web -d=a ./smilax5.tcl 2>&1 | tee _smilax5_log
cont:
	cd ../goapi && make
	cd ../ && go install
	../../../../../bin/smilax-web -d=c ./unbundle.tcl < ./run1.txt 2>&1 | tee _unbundle_log
	# smilax/smilax -d=relapz ./smilax5.tcl 2>&1 | tee _smilax5_log
	../../../../../bin/smilax-web -d=a ./smilax5.tcl 2>&1 | tee _smilax5_log
clean:
	rm -rf ./_data
	rm -f ./_* ./,*
