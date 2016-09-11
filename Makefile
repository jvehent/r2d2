GO          := GO15VENDOREXPERIMENT=1 go
GOGETTER    := GOPATH=$(shell pwd)/.tmpdeps go get -d

all: install

install:
		$(GO) install github.com/jvehent/r2d2

vendor:
	govend -u

.PHONY: vendor
