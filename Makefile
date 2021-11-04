PROJECT = github.com/Code-Hex/Neo-cowsay

.PHONY: build
build: build/cowsay build/cowthink

.PHONY: build/cowsay
build/cowsay:
	CGO_ENABLED=0 go build -o bin/cowsay -ldflags "-w -s" \
		$(PROJECT)/cmd/cowsay

.PHONY: build/cowthink
build/cowthink:
	CGO_ENABLED=0 go build -o bin/cowthink -ldflags "-w -s" \
		$(PROJECT)/cmd/cowthink

.PHONY: lint
lint:
	golint ./...

.PHONY: test
test:
	go test ./...

.PHONY: man
man:
	asciidoctor --doctype manpage --backend manpage doc/neo-cowsay.1.txt.tpl -o doc/neo-cowsay.1

.PHONY: man/preview
man/preview:
	cat doc/neo-cowsay.1 | groff -man -Tascii | less
