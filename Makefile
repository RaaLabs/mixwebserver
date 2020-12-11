NAME := $(shell basename `git rev-parse --show-toplevel`)
VERSION := $(shell git tag --sort=-version:refname | head -n 1)
BUILDDIR := $(NAME)_$(VERSION)

define CONFIGURE
Package: $(NAME)
Version: $(VERSION)
Architecture: amd64
Maintainer: Bjørn Tore Svinningen postmannen@gmail.com
Description: mixwebserver
  Small web server that lets you set up a simple web server to serve a file directory,
  and it will also register with letsencypt for ssl.
endef
export CONFIGURE


deb: main.go
	mkdir -p $(BUILDDIR)/DEBIAN
	mkdir -p $(BUILDDIR)/usr/local/bin
	go build -o $(BUILDDIR)/usr/local/bin
	echo "$$CONFIGURE" > $(BUILDDIR)/DEBIAN/control
	dpkg-deb --build $(BUILDDIR)