package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	cowsay "github.com/Code-Hex/Neo-cowsay"
	"github.com/Code-Hex/Neo-cowsay/internal/super"
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
	version = "0.0.5"
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

func generateOptions(opts *Options, phrase string) []cowsay.Option {
	o := make([]cowsay.Option, 0, 8)
	o = append(o,
		cowsay.Phrase(phrase),
		cowsay.Type(opts.File),
		cowsay.Thinking(),
	)
	if opts.Bold {
		o = append(o, cowsay.Bold())
	}
	if opts.Random {
		o = append(o, cowsay.Random())
	}
	if opts.Rainbow {
		o = append(o, cowsay.Rainbow())
	}
	if opts.Aurora {
		o = append(o, cowsay.Aurora())
	}
	return selectFace(opts, o)
}

func phrase(opts *Options, args []string) string {
	if len(args) > 0 {
		return strings.Join(args, " ")
	}
	lines := make([]string, 0, 40)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if opts.NewLine {
		return strings.Join(lines, "\n")
	}
	return strings.Join(lines, " ")
}

func mowmow(opts *Options, args []string) error {
	phrase := phrase(opts, args)
	o := generateOptions(opts, phrase)
	cow, err := cowsay.NewCow(o...)
	if err != nil {
		return err
	}
	if opts.Super {
		return super.RunSuperCow(cow)
	}

	say, err := cow.Say()
	if err != nil {
		return err
	}
	fmt.Fprintln(colorable.NewColorableStdout(), say)

	return nil
}

func selectFace(opts *Options, o []cowsay.Option) []cowsay.Option {
	switch {
	case opts.Borg:
		o = append(o,
			cowsay.Eyes("=="),
			cowsay.Tongue("  "),
		)
	case opts.Dead:
		o = append(o,
			cowsay.Eyes("xx"),
			cowsay.Tongue("U "),
		)
	case opts.Greedy:
		o = append(o,
			cowsay.Eyes("$$"),
			cowsay.Tongue("  "),
		)
	case opts.Paranoia:
		o = append(o,
			cowsay.Eyes("@@"),
			cowsay.Tongue("  "),
		)
	case opts.Stoned:
		o = append(o,
			cowsay.Eyes("**"),
			cowsay.Tongue("U "),
		)
	case opts.Tired:
		o = append(o,
			cowsay.Eyes("--"),
			cowsay.Tongue("  "),
		)
	case opts.Wired:
		o = append(o,
			cowsay.Eyes("OO"),
			cowsay.Tongue("  "),
		)
	case opts.Youthful:
		o = append(o,
			cowsay.Eyes(".."),
			cowsay.Tongue("  "),
		)
	}
	return o
}

func parseOptions(opts *Options, argv []string) ([]string, error) {
	p := flags.NewParser(opts, flags.None)
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
	return []byte(`cow{say,think} version ` + version + `, (c) 2016 CodeHex
Usage: cowsay [-bdgpstwy] [-h] [-e eyes] [-f cowfile] [--random]
          [-l] [-n] [-T tongue] [-W wrapcolumn]
          [--rainbow] [--aurora] [--super] [message]

Original Author: (c) 1999 Tony Monroe
`)
}
