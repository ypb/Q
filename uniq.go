// Copyright 2011 AUTHORS. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"time"
	"strconv"
	"os"
	tdb "ypb/github.com/ypb/gotdb"
)

const (
	BASE = 36 // why, no upper case letters?
	CNT  = "now"
)

type uniq struct {
	utc *time.Time
	sec string
	cnt uint16 // 64kB of memory is enough for anybody!
}

func Uniq() (uniq, os.Error) {
	y := uniq{nil, "0", 0}
	if err := y.Tack(); err != nil {
		// return nil, err // wishy/washy thinking
		return y, err
	}
	return y, nil
}

// Tack starts a fresh "this second Epoch". It does not mean it lasts a
// second, it's mere resolution limitation. Given a "Tack" all following
// events receive monotonically increasing "Ticks". If Tack detects
// that last second is still up it does not reset the "Ticks" counter.
func (o *uniq) Tack() os.Error {
	var sec string
	var cnt uint16
	// stamp := time.Format("2006 Jan 02 | 05 04 15  MST (-0700)")
	var err tdb.Error
	meta := db[META+EXT] // GLOBAL ENV rox!
	store := true
	// ahh, those intermunged concerns...
	if meta == nil {
		return os.NewError("uniq.Tack(): nil META db")
	} else {
		o.utc = time.UTC()
		o.sec = strconv.Itob64(o.utc.Seconds(), BASE)
		var cnt_s string
		if sec, err = meta.Fetch(NOW); err == nil {
			// trin // fuck // order
			if o.sec == sec {
				if cnt_s, err = meta.Fetch(CNT); err == nil {
					cnt_ugh, _ := strconv.Btoui64(cnt_s, BASE)
					cnt = uint16(cnt_ugh)
					if cnt > o.cnt {
						o.cnt = cnt
						store = false
					}
				}
			} else {
				o.cnt = 0
			}
		}
		// TODO tdb.Exists()
		if store {
			if err = meta.Store(NOW, o.sec, tdb.MODIFY); err != nil {
				// presumably it yet exists not...
				if err = meta.Store(NOW, o.sec, tdb.INSERT); err != nil {
					return err
				}
			}
			// lazy logic.
			cnt_s = strconv.Itob64(int64(o.cnt), BASE)
			// mip map
			if err = meta.Store(CNT, cnt_s, tdb.MODIFY); err != nil {
				// presumably it yet exists not...
				if err = meta.Store(CNT, cnt_s, tdb.INSERT); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Tick ticks me off... we should store each one in db or will well clubber
// himself on Tack. But for now... we Tack once per forkproc and Tick in hopes
// ... oh, we can always Tack on exit...
//
func (o *uniq) Tick() string {
	cnt_s := strconv.Itob64(int64(o.cnt), BASE)
	o.cnt++
	return o.sec + "." + cnt_s
}

// Local Variables:
// mode: Go
// End:
