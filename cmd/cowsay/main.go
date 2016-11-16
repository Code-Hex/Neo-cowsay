package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	cowsay "github.com/Code-Hex/Neo-cowsay"
	wordwrap "github.com/Code-Hex/go-wordwrap"
	flags "github.com/jessevdk/go-flags"
	colorable "github.com/mattn/go-colorable"
)

// Options struct for parse command line arguments
type Options struct {
	Help     bool   `short:"h"`
	Eyes     string `short:"e"`
	Tongue   string `short:"T"`
	Width    int    `short:"W"`
	Borg     bool   `short:"b"`
	Dead     bool   `short:"d"`
	Greedy   bool   `short:"g"`
	Paranoia bool   `short:"p"`
	Stoned   bool   `short:"s"`
	Tired    bool   `short:"t"`
	Wired    bool   `short:"w"`
	Youthful bool   `short:"y"`
	List     bool   `short:"l"`
	NewLine  bool   `short:"n"`
	File     string `short:"f"`
	Bold     bool   `long:"bold"`
	Super    bool   `long:"super"`
	Random   bool   `long:"random"`
	Rainbow  bool   `long:"rainbow"`
	Aurora   bool   `long:"aurora"`
}

const (
	version = "0.0.3"
)

func main() {
	os.Exit(run())
}

func run() int {
	if err := mow(); err != nil {
		fmt.Fprintf(os.Stderr, "Error:\n  %v\n", err)
		return 1
	}
	return 0
}

// mow will parsing for cowsay command line arguments and invoke cowsay.
func mow() error {
	var opts Options
	args, err := parseOptions(&opts, os.Args[1:])
	if err != nil {
		return err
	}

	if opts.List {
		fmt.Println(wordwrap.WrapString(strings.Join(cowsay.Cows(), " "), 80))
		return nil
	}

	if err := mowmow(&opts, args); err != nil {
		return err
	}

	return nil
}

func mowmow(opts *Options, args []string) error {
	cow := &cowsay.Cow{
		Type:      opts.File,
		Bold:      opts.Bold,
		IsRandom:  opts.Random,
		IsRainbow: opts.Rainbow,
		IsAurora:  opts.Aurora,
	}

	selectFace(opts, cow)

	if opts.Width <= 0 {
		cow.BallonWidth = 40
	}

	if len(args) > 0 {
		cow.Phrase = strings.Join(args, " ")
	} else {
		lines := make([]string, 0, 40)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if opts.NewLine {
			cow.Phrase = strings.Join(lines, "\n")
		} else {
			cow.Phrase = strings.Join(lines, " ")
		}
	}

	if opts.Super {
		return runSuperCow(cow)
	}

	say, err := cowsay.Say(cow)
	if err != nil {
		return err
	}

	fmt.Fprintln(colorable.NewColorableStdout(), say)

	return nil
}

func selectFace(opts *Options, cow *cowsay.Cow) {
	switch {
	case opts.Borg:
		cow.Eyes = "=="
		cow.Tongue = "  "
	case opts.Dead:
		cow.Eyes = "xx"
		cow.Tongue = "U "
	case opts.Greedy:
		cow.Eyes = "$$"
		cow.Tongue = "  "
	case opts.Paranoia:
		cow.Eyes = "@@"
		cow.Tongue = "  "
	case opts.Stoned:
		cow.Eyes = "**"
		cow.Tongue = "U "
	case opts.Tired:
		cow.Eyes = "--"
		cow.Tongue = "  "
	case opts.Wired:
		cow.Eyes = "OO"
		cow.Tongue = "  "
	case opts.Youthful:
		cow.Eyes = ".."
		cow.Tongue = "  "
	}
}

func parseOptions(opts *Options, argv []string) ([]string, error) {
	p := flags.NewParser(opts, flags.PrintErrors)
	args, err := p.ParseArgs(argv)
	if err != nil {
		return nil, err
	}

	if opts.Help {
		os.Stdout.Write(opts.usage())
		os.Exit(0)
	}

	return args, nil
}

func (opts Options) usage() []byte {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, `cow{say,think} version `+version+`, (c) 2016 CodeHex
Usage: cowsay [-bdgpstwy] [-h] [-e eyes] [-f cowfile] [--random]
          [-l] [-n] [-T tongue] [-W wrapcolumn]
          [--rainbow] [--aurora] [--super] [message]

Original Author: (c) 1999 Tony Monroe
`)
	return buf.Bytes()
}
