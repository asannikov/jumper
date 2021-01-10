#
# Makefile
#
VERSION = snapshot
GHRFLAGS =
.PHONY: build release

default: build

build:
	goxc -t -wd=. -d=pkg -pv=$(VERSION)

release:
	ghr  -u asannikov  $(GHRFLAGS) v$(VERSION) pkg/$(VERSION)
