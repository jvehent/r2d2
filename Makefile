GO          := GO15VENDOREXPERIMENT=1 go
GOGETTER    := GOPATH=$(shell pwd)/.tmpdeps go get -d

all: install

install:
		$(GO) install github.com/jvehent/r2d2

go_vendor_dependencies::
	$(GOGETTER) gopkg.in/gcfg.v1
	$(GOGETTER) golang.org/x/oauth2
	$(GOGETTER) github.com/google/go-github/github
	$(GOGETTER) github.com/thoj/go-ircevent
	#$(GOGETTER) github.com/oschwald/geoip2-golang
	echo 'removing .git from vendored pkg and moving them to vendor'
	find .tmpdeps/src -name ".git" ! -name ".gitignore" -exec rm -rf {} \; || exit 0
	[ -d vendor ] && git rm -rf vendor/ || exit 0
	mkdir vendor/ || exit 0
	cp -ar .tmpdeps/src/* vendor/
	git add vendor/
	rm -rf .tmpdeps
