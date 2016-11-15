PWD          = $(shell pwd)
REPO         = $(shell basename $(PWD))
COWSAYPATH   = ./cmd/cowsay/
COWTHINKPATH = ./cmd/cowthink/
GOVERSION    = $(shell go version)
GOOS         = $(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH       = $(word 2,$(subst /, ,$(lastword $(GOVERSION))))
XC_OS        = "darwin windows linux"
XC_ARCH      = "i386 amd64"
VERSION      = $(patsubst "%",%,$(lastword $(shell grep 'version =' cmd/cowsay/main.go)))
RELEASE      = ./releases/$(VERSION)

GITHUB_USERNAME = "Code-Hex"

rm-build:
	@rm -rf build

rm-releases:
	@rm -rf releases

rm-all: rm-build rm-releases

release: all
	@mkdir -p $(RELEASE)
	@for DIR in $(shell ls ./build/$(VERSION)/) ; do \
		echo "Processing in build/$(VERSION)/$$DIR"; \
		cd $(PWD); \
		cp README.md ./build/$(VERSION)/$$DIR; \
		cp LICENSE ./build/$(VERSION)/$$DIR; \
		tar -czf ./$(RELEASE)/cowsay_$(VERSION)_$$DIR.tar.gz -C ./build/$(VERSION) $$DIR; \
	done

prepare-github: github-token
	@echo "'github-token' file is required"

release-upload: prepare-github release
	@echo "Uploading..."
	@ghr -u $(GITHUB_USERNAME) -t $(shell cat github-token) --draft --replace $(VERSION) $(RELEASE)

lint:
	go get github.com/golang/lint/golint
	@golint

gen:
	@go generate
	
all: gen
	@gox -os=$(XC_OS) -arch=$(XC_ARCH) -output="build/$(VERSION)/{{.OS}}_{{.Arch}}/{{.Dir}}" $(COWSAYPATH)
	@gox -os=$(XC_OS) -arch=$(XC_ARCH) -output="build/$(VERSION)/{{.OS}}_{{.Arch}}/{{.Dir}}" $(COWTHINKPATH)

