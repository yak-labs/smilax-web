# Chirp script to unbundle scenarios into _data directory.

set TEXT_rx [/regexp/MustCompile {^/TEXT (.*)$}]
set all [/io/ioutil/ReadFile /dev/stdin]

proc Create path {
	set s,v,p,f [dropnull [split $path /]]
	set dir "_data/s.$s/v.$v/p.$p/f.$f"
	set mode777 511
	/os/MkdirAll $dir $mode777
	return [/os/Create $dir/r.0]
}

set f ""
foreach line [split $all \n] {
	set m [$TEXT_rx FindStringSubmatch $line]
	if [notnull $m] {
		# is a magic /TEXT line
		set path [lindex $m 1]
		if [notnull $f] {
			$f Close
		}
		set f [Create $path]
	} else {
		# not a magic /TEXT line
		if [notnull $f] {
			$f WriteString "$line\n"
		}
	}
}

if [notnull $f] {
	$f Close
}
