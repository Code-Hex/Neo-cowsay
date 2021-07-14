# Neo Cowsay
Fast, funny, everyone wanted? new cowsay!!
```
 ______________
< I'm Neo cows >
 --------------
       \   ^__^
        \  (oo)\_______
           (__)\       )\/\
               ||----w |
               ||     ||
```

for GitHub Actions users: [Code-Hex/neo-cowsay-action](https://github.com/marketplace/actions/neo-cowsay)

[![GoDoc](https://godoc.org/github.com/Code-Hex/Neo-cowsay?status.svg)](https://godoc.org/github.com/Code-Hex/Neo-cowsay) [![.github/workflows/main.yml](https://github.com/Code-Hex/Neo-cowsay/actions/workflows/main.yml/badge.svg)](https://github.com/Code-Hex/Neo-cowsay/actions/workflows/main.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/Code-Hex/Neo-cowsay)](https://goreportcard.com/report/github.com/Code-Hex/Neo-cowsay)
# What's?
```
cowsay is a configurable talking cow, written in Perl.  It operates
much as the figlet program does, and it written in the same spirit
of silliness.
```  
by [Original](https://github.com/schacon/cowsay).  

Neo Cowsay is written in Go. This cowsay is extended the original cowsay. added fun more options, and you can be used as a library.

# Usage
## As command
```
cow{say,think} version 1.0.0, (c) 2021 codehex
Usage: cowsay [-bdgpstwy] [-h] [-e eyes] [-f cowfile] [--random]
      [-l] [-n] [-T tongue] [-W wrapcolumn]
      [--rainbow] [--aurora] [--super] [message]

Original Author: (c) 1999 Tony Monroe
Repository: https://github.com/Code-Hex/Neo-cowsay
```
Normal
```
$ cowsay Hello
 _______
< Hello >
 -------
       \   ^__^
        \  (oo)\_______
           (__)\       )\/\
               ||----w |
               ||     ||
```
Borg mode
```
$ cowsay -b Hello
 _______
< Hello >
 -------
       \   ^__^
        \  (==)\_______
           (__)\       )\/\
               ||----w |
               ||     ||
```
## As library
```go
package main

import (
	"fmt"

	cowsay "github.com/Code-Hex/Neo-cowsay"
)

func main() {
	say, err := cowsay.Say(
		cowsay.Phrase("Hello"),
		cowsay.Type("default"),
		cowsay.BallonWidth(40),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(say)
}
```
[Example](https://github.com/Code-Hex/Neo-cowsay/blob/master/eg/main.go) or [GoDoc](https://godoc.org/github.com/Code-Hex/Neo-cowsay)
# New options
## Random
[![asciicast](https://asciinema.org/a/228210.svg)](https://asciinema.org/a/228210)
## Rainbow and Aurora, Bold
[![asciicast](https://asciinema.org/a/228213.svg)](https://asciinema.org/a/228213)
# And, Super Cows mode
asciinema is heavy...
[![asciicast](https://asciinema.org/a/228215.svg)](https://asciinema.org/a/228215)

# Install
## library

    $ go get -u github.com/Code-Hex/Neo-cowsay

## Go

### cowsay

    $ go get -u github.com/Code-Hex/Neo-cowsay/cmd/cowsay

### cowthink

    $ go get -u github.com/Code-Hex/Neo-cowsay/cmd/cowthink
    
## Mac and Linux users via Homebrew

    $ brew update
    $ brew install Code-Hex/tap/neo-cowsay

## Binary
You can download from [here](https://github.com/Code-Hex/Neo-cowsay/releases)

# License
[cowsay license](https://github.com/Code-Hex/Neo-cowsay/blob/master/LICENSE)  
(The Artistic License or The GNU General Public License)

# Author
Neo Cowsay: [codehex](https://twitter.com/CodeHex)  
Original: (c) 1999 Tony Monroe
