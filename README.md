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
[![GoDoc](https://godoc.org/github.com/Code-Hex/Neo-cowsay?status.svg)](https://godoc.org/github.com/Code-Hex/Neo-cowsay) [![Build Status](https://travis-ci.org/Code-Hex/Neo-cowsay.svg?branch=master)](https://travis-ci.org/Code-Hex/Neo-cowsay) [![Go Report Card](https://goreportcard.com/badge/github.com/Code-Hex/Neo-cowsay)](https://goreportcard.com/report/github.com/Code-Hex/Neo-cowsay)
# What's?
```
cowsay is a configurable talking cow, written in Perl.  It operates
much as the figlet program does, and it written in the same spirit
of silliness.
```  
by [Original](https://github.com/schacon/cowsay).  
Neo Cowsay written in Go. This cowsay extended the original and added fun more options. And it can be used as a library.

# Usage
## As command
```
cow{say,think} version 0.0.1, (c) 2016 CodeHex
Usage: cowsay [-bdgpstwy] [-h] [-e eyes] [-f cowfile] [--random]
          [-l] [-n] [-T tongue] [-W wrapcolumn]
          [--rainbow] [--aurora] [--super] [message]

Original Author: (c) 1999 Tony Monroe
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
	say, err := cowsay.Say(&cowsay.Cow{
		Phrase:      "Hello!!",
		Type:        "default",
		BallonWidth: 40,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(say)
}
```
 [Example](https://github.com/Code-Hex/Neo-cowsay/blob/master/cmd/cowsay/main.go#L75) or [GoDoc](https://godoc.org/github.com/Code-Hex/Neo-cowsay)
# New options
## Random
[![asciicast](https://asciinema.org/a/avq390avlf6ddb4jn7d0n0y37.png)](https://asciinema.org/a/avq390avlf6ddb4jn7d0n0y37)
## Rainbow and Aurora, Bold
[![asciicast](https://asciinema.org/a/d3k3a182rsndlgez5sdzhqprk.png)](https://asciinema.org/a/d3k3a182rsndlgez5sdzhqprk)
# And, Super Cows mode
asciinema is heavy...
[![asciicast](https://asciinema.org/a/crf5crjim1d2nw01ioigug0ks.png)](https://asciinema.org/a/crf5crjim1d2nw01ioigug0ks)
# Install
## library

    go get -u github.com/Code-Hex/Neo-cowsay
## cowsay

    go get -u github.com/Code-Hex/Neo-cowsay/cmd/cowsay
## cowthink

    go get -u github.com/Code-Hex/Neo-cowsay/cmd/cowthink

# License
[cowsay license](https://github.com/Code-Hex/Neo-cowsay/blob/master/LICENSE)  
(The Artistic License or The GNU General Public License)

# Author
Neo Cowsay: [codehex](https://twitter.com/CodeHex)  
Original: (c) 1999 Tony Monroe
