# Chirp script to unbundle scenarios into _data directory.

set TEXT_rx [/regexp/MustCompile {^/TEXT (.*)$}]
# set all [/io/ioutil/ReadFile /dev/stdin]
set all [/io/ioutil/ReadAll [/os/Stdin]]

proc Create path {
	set mode777 511
    if {$path eq "/table_log.txt"} {
        /os/MkdirAll "_data" $mode777
        return [/os/Create "_data$path"]
    }
    set s,v,p,f [dropnull [split $path /]]
    set dir "_data/s.$s/v.$v/p.$p/f.$f"
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
