// Copyright 2011 AUTHORS. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"strings"
	"time"
	"strconv"
	tdb "ypb/github.com/ypb/gotdb"
)

const (
	R = "q"

	// drfe has/had: atomic, event, head, info, meta, tail...
	HERE  = "Here"
	EVER  = "Ever"
	NOW   = "Now"
	LATER = "Later"
	META  = "Meta"

	EXT = ".tdb"

	atom = "atom"
)

// this should really go into tdb, experimenting...
var db map[string]*tdb.DB //+root path of course...

var BNAME string // basename $0
var argc int

func exit(rety tdb.Error) {
	println(rety.String())
	os.Exit(rety.Errno())
}

func done() {
	println("# closing DB")
	for k, v := range db {
		if ret := v.Close(); ret != nil {
			println(k, ret.String())
		}
	}
}

func init() {
	db = make(map[string]*tdb.DB)
	db[HERE+EXT] = nil
	db[EVER+EXT] = nil
	db[NOW+EXT] = nil
	db[LATER+EXT] = nil
	db[META+EXT] = nil

	for k, _ := range db {
		if tdb, err := tdb.New(k); err != nil {
			done()
			exit(err)
		} else {
			db[k] = &tdb
		}
	}

	argc = len(os.Args)
	println("# len(os.Args) =", argc)

	if argc >= 1 {
		s := strings.SplitN(os.Args[0], "/", -1)
		BNAME = s[len(s)-1]
	} else {
		BNAME = R
	}
}

func main() {
	defer done()

	r := P(false)           // this is silly TOFIX...
	q := P(BNAME).Heldby(r) // and so is this: Hold'ing on should auto-Heldby
	// q := P(R)
	q.Hold(P(HERE).Heldby(q))
	q.Hold(P(EVER).Heldby(q))
	q.Hold(P(NOW).Heldby(q))
	q.Hold(P(LATER).Heldby(q))
	q.Hold(P(META).Heldby(q))

	// db := Opendb()
	var err tdb.Error

	time := time.UTC()
	// stamp := time.Format("2006 Jan 02 | 05 04 15  MST (-0700)")
	stamp := strconv.Itob64(time.Seconds(), 36)

	// TODO tdb.Exists()
	if err = db[META+EXT].Store(NOW, stamp, tdb.MODIFY); err != nil {
		// presumably it yet exists not...
		if err = db[META+EXT].Store(NOW, stamp, tdb.INSERT); err != nil {
			exit(err)
		}
	}

	var now string
	if now, err = db[META+EXT].Fetch(NOW); err != nil {
		exit(err)
	}
	println("# now:", now)
	q.Heldby(P(NOW)).Hold(P(now))

	if argc <= 1 {
		q.Print()
	} else {
		if cmd := q.Heldby(P(os.Args[1])); cmd != nil {
			cmd.Print()
		} else {
			println(os.Args[1], "uknown command")
			os.Exit(-1)
		}
	}
	return
}

// Local Variables:
// mode: Go
// End:
