PROJS = r2d2
GO = GOPATH=$(shell pwd):$(shell go env GOROOT)/bin go

all: $(PROJS)

depends:
	$(GO) get code.google.com/p/gcfg
	$(GO) get code.google.com/p/goauth2/oauth
	$(GO) get github.com/google/go-github/github
	$(GO) get github.com/thoj/go-ircevent

r2d2:
	$(GO) install r2d2

clean:
	rm -f bin/r2d2
