.PHONY: build
build: build/cowsay build/cowthink

.PHONY: build/cowsay
build/cowsay:
	CGO_ENABLED=0 cd cmd && go build -o ../bin/cowsay -ldflags "-w -s" ./cowsay

.PHONY: build/cowthink
build/cowthink:
	CGO_ENABLED=0 cd cmd && go build -o ../bin/cowthink -ldflags "-w -s" ./cowthink

.PHONY: lint
lint:
	golint ./...
	cd cmd && golint ./...

.PHONY: vet
vet:
	go vet ./...
	cd cmd && go vet ./...

.PHONY: test
test: test/pkg test/cli

.PHONY: test/pkg
test/pkg:
	go test ./...

.PHONY: test/cli
test/cli:
	cd cmd && go test ./...

.PHONY: man
man:
	asciidoctor --doctype manpage --backend manpage doc/cowsay.1.txt.tpl -o doc/cowsay.1

.PHONY: man/preview
man/preview:
	cat doc/cowsay.1 | groff -man -Tascii | less
