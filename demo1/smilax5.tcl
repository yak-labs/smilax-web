# smilax5.t -- the trusted (unsafe) part of Smilax5, that configs & launches the safe interps.

#### New Storage Hierarchy:
# s: site (old bundle)
# v: volume (old dir)
# p: page (old file)
# f: file.  Special files: "@wiki".
# r: revision; also v: varient.

set SiteNameRx     [/regexp/MustCompile {^[A-Z]_[a-z0-9_]*$}]
set VolNameRx      [/regexp/MustCompile {^[A-Z]+[a-z0-9_]*$}]
set PageNameRx     [/regexp/MustCompile {^[A-Z]+[a-z]+[A-Z][A-Za-z0-9_]*$}]
set FileNameRx     [/regexp/MustCompile {^[A-Za-z0-9_.@%~][-A-Za-z0-9_.@%~]*[.][-A-Za-z0-9_.@%~]+$}]

set SiteDirRx     [/regexp/MustCompile {^s[.]([-A-Za-z0-9_.]+)$}]
set VolDirRx      [/regexp/MustCompile {^v[.]([-A-Za-z0-9_.]+)$}]
set PageDirRx     [/regexp/MustCompile {^p[.]([-A-Za-z0-9_.]+)$}]
set FileDirRx     [/regexp/MustCompile {^f[.]([-A-Za-z0-9_.]+)$}]
set RevFileRx     [/regexp/MustCompile {^r[.]([-A-Za-z0-9_.]+)$}]
set VarientFileRx [/regexp/MustCompile {^v[.]([-A-Za-z0-9_.]+)$}]

set MarkForSubinterpRx [/regexp/MustCompile {^[@]([-A-Za-z0-9_]+)$}]
set BasicAuthUserPwSplitterRx [/regexp/MustCompile {^([^:]*)[:](.*)$}]
set BASE64 [/encoding/base64/StdEncoding *]

yproc @ListSites {} {
	check-site Root
	foreach f [/io/ioutil/ReadDir "data"] {
		set fname [$f Name]
		set m [$SiteDirRx FindStringSubmatch $fname]
		if {[notnull $m]} {
			yield [lindex $m 1]
		}
	}
}

yproc @ListVols { site } {
	check-site $site
	foreach f [/io/ioutil/ReadDir "data/s.$site"] {
		set fname [$f Name]
		set m [$VolDirRx FindStringSubmatch $fname]
		if {[notnull $m]} {
			yield [lindex $m 1]
		}
	}
}

yproc @ListPages { site vol } {
	check-site $site
	foreach f [/io/ioutil/ReadDir "data/s.$site/v.$vol"] {
		set fname [$f Name]
		set m [$PageDirRx FindStringSubmatch $fname]
		if {[notnull $m]} {
			yield [lindex $m 1]
		}
	}
}

yproc @ListFiles { site vol page } {
	check-site $site
	foreach f [/io/ioutil/ReadDir "data/s.$site/v.$vol/p.$page"] {
		set fname [$f Name]
		set m [$FileDirRx FindStringSubmatch $fname]
		if {[notnull $m]} {
			yield [lindex $m 1]
		}
	}
}

yproc @ListRevs { site vol page file } {
	check-site $site
	foreach f [/io/ioutil/ReadDir "data/s.$site/v.$vol/p.$page/f.$file"] {
		set fname [$f Name]
		set m [$RevFileRx FindStringSubmatch $fname]
		if {[notnull $m]} {
			yield [lindex $m 1]
		}
	}
}

proc @ReadFile { site vol page file } {
	check-site $site
  # TODO: Use [lsort -decreasing] so we can do this in less commands.
	set revs [lsort [concat [@ListRevs $site $vol $page $file]]]
	set rev [lindex $revs [expr [llength $revs] - 1]]

	return [/io/ioutil/ReadFile "data/s.$site/v.$vol/p.$page/f.$file/r.$rev"]
}

proc @WriteFile { site vol page file contents } {
  check-site $site
  set now [/time/Now]
  set nowUnix [$now Unix]

  # Need to use strconv, otherwise the int64 gets turned into a float and the
  # timestamp will get represented as scientific notation.
  set timestamp [/strconv/FormatInt $nowUnix 10]

	/os/MkdirAll "data/s.$site/v.$vol/p.$page/f.$file" 448
	/io/ioutil/WriteFile "data/s.$site/v.$vol/p.$page/f.$file/r.$timestamp" $contents 384

	# Save no records, but stupid side-effect is to reread all files.
	db-save-records "data" {}
}

proc @Route { path query } {
	/fmt/Fprintf [cred w] %s "This is the base Router.  Replace me."
	/fmt/Fprintf [cred w] {path: %s | query: %s} $path $query
}

proc @RxCompile { pattern } {
	/regexp/MustCompile $pattern
}

proc @FindStringSubmatch { rx str } {
	$rx FindStringSubmatch $str
}

proc check-site site {
	set s [cred site]
	# If site empty, then no enforcement yet.
	if [null $s] return
	# If correct site, OK.
	if [eq $s $site] return
	# Must be super-user to access a different site.
	@auth-require-level 90
}

proc @EntityGet {site table id field tag} {
	check-site $site
	entity-get $site $table $id $field $tag
}

proc @EntityPut {site table id field tag values} {
	check-site $site
	entity-put $site $table $id $field $tag $values
}

proc @EntityLike {site table field tag value} {
	check-site $site
	entity-like $site $table $field $tag $value
}

proc @EntityTriples {site table id field tag value} {
	check-site $site
	entity-triples $site $table $id $field $tag $value
}

proc @ShowValue v { 
    /fmt/Sprintf %v $v
}

proc @Puts { str} {
    /fmt/Fprintf [cred w] %s $str
}

proc @ModeHtml {} {
    [[cred w] Header] Set "Content-Type" "text/html"
}

#proc @TemporaryRedirect url {
#	set rh [/net/http/RedirectHandler $url 307]
#	$rh ServeHTTP [cred w] [cred r]
#}
proc @TemporaryRedirect url {
	# set rh [/net/http/RedirectHandler $url 307]
	# $rh ServeHTTP [cred w] [cred r]
	throw 307 $url
}

proc @auth-require-level {level} {
	if {[cred level] < $level} {
		@RequestBasicAuth
	}
}

proc @RequestBasicAuth {} {
	set h [[cred w] Header]
	$h Set "WWW-Authenticate" "Basic realm=\"[cred site]\""
	[cred w] WriteHeader 401
}

proc @user {} {
	cred user
}
proc @host {} {
	cred host
}
proc @level {} {
	cred level
}
proc @site {} {
	cred site
}
proc @r {} {
	cred r
}
proc @w {} {
	cred w
}

# Dir name is "data"
db-rebuild "data"

######  DEFINE @-procs ABOVE.

set Zygote [interp]

foreach cmd [info commands] {
  set m [$MarkForSubinterpRx FindStringSubmatch $cmd]
  if [notnull $m] {
    $Zygote Alias - [lindex $m 1] $cmd
  }
}

$Zygote Alias - DB "set DB"

# -- Load our mixins into our sub-interpreter
set mixins [@ListPages Root Mixin]
foreach m $mixins {
	$Zygote Eval [list mixin $m [@ReadFile Root Mixin $m @wiki]]
}

proc gold-level {user pw} {
	foreach r [db-select-like [cred site] pw Sys PassWord "$user:$pw" *] {
		return [$r . Values @ 1]
	}
	return 0
}

proc lookup-site {} {
	foreach r [db-select-like Root serve Sys ServeSite [cred host] *] {
		return [$r . Values @ 1]
	}
	error "Unknown Site for HOST=[cred host]"
}

# NOW HANDLE REQUESTS
proc ZygoteHandler {w r} { # TODO: get rid of the args.i w & r.
	cred-put site [lookup-site]

	set headers [[cred r] . Header]
	set authorization [$headers Get Authorization]
	if [notnull $authorization] {
		set obfuscated [lindex $authorization 1]
		set decoded [$BASE64 DecodeString $obfuscated]
		set m [$BasicAuthUserPwSplitterRx FindStringSubmatch $decoded]
		if [notnull $m] {
			set _,USER,PASSWORD $m
		}
		set level [gold-level $USER $PASSWORD]
		if {$level <= 0} {
			set level [gold-level * $PASSWORD]
		}
	} else {
		set level [gold-level * *]
	}

	  cred-put level $level

	  set clone [$Zygote Clone]
	  $clone CopyCredFrom -

	set e [catch {
	  if {$level <= 0} {
	  }
	  $clone Eval [ list Route [$r . URL . Path] [[$r . URL] Query] ]
	} what]

    case $e in {
	  0 {
		# OK
	  }
	  307 {
		set url [lindex [split $what "\n"] 0]
		set rh [/net/http/RedirectHandler $url 307]
		$rh ServeHTTP [cred w] [cred r]
	  }
	  default {
		# TODO: something better.
		[cred w] Write [[ht cat $what] Html]
	  }
	}
}

/net/http/HandleFunc / [ http-handler-lambda {w r} {ZygoteHandler $w $r} ]
/net/http/ListenAndServe 127.0.0.1:8080 ""
