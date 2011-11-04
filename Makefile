# Copyright 2011 AUTHORS. All rights reserved.
# Use of this source code is governed by a GPL-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=Q
GOFILES=\
	main.go\
	hold.go\
	uniq.go\

CLEANFILES += *.tdb

include $(GOROOT)/src/Make.cmd

fmt: $(GOFILES)
	gofmt -d $^

# Local Variables:
# mode: Makefile
# End:
