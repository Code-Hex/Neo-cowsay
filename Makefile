PROJECT = github.com/Code-Hex/Neo-cowsay

build: build/cowsay build/cowthink

build/cowsay:
	CGO_ENABLED=0 go build -o bin/cowsay -ldflags "-w -s" \
		$(PROJECT)/cmd/cowsay

build/cowthink:
	CGO_ENABLED=0 go build -o bin/cowthink -ldflags "-w -s" \
		$(PROJECT)/cmd/cowthink

lint:
	golint ./...
