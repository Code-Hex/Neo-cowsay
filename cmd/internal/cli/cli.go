package cli

import (
	"bufio"
	cryptorand "crypto/rand"
	"errors"
	"fmt"
	"github.com/Code-Hex/go-wordwrap"
	"github.com/ktr0731/go-fuzzyfinder"
	"io"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Code-Hex/Neo-cowsay/cmd/v2/internal/super"
	cowsay "github.com/Code-Hex/Neo-cowsay/v2"
	"github.com/Code-Hex/Neo-cowsay/v2/decoration"
	"github.com/jessevdk/go-flags"
	"github.com/mattn/go-colorable"
)

func init() {
	// safely set the seed globally so we generate random ids. Tries to use a
	// crypto seed before falling back to time.
	var seed int64
	cryptoseed, err := cryptorand.Int(cryptorand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		// This should not happen, but worst-case fallback to time-based seed.
		seed = time.Now().UnixNano()
	} else {
		seed = cryptoseed.Int64()
	}
	rand.Seed(seed)
}

// options struct for parse command line arguments
type options struct {
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

// CLI prepare for running command-line.
type CLI struct {
	Version  string
	Thinking bool
	stderr   io.Writer
	stdout   io.Writer
	stdin    io.Reader
}

func (c *CLI) program() string {
	if c.Thinking {
		return "cowthink"
	}
	return "cowsay"
}

// Run runs command-line.
func (c *CLI) Run(argv []string) int {
	if c.stderr == nil {
		c.stderr = os.Stderr
	}
	if c.stdout == nil {
		c.stdout = colorable.NewColorableStdout()
	}
	if c.stdin == nil {
		c.stdin = os.Stdin
	}
	if err := c.mow(argv); err != nil {
		fmt.Fprintf(c.stderr, "%s: %s\n", c.program(), err.Error())
		return 1
	}
	return 0
}

// mow will parsing for cowsay command line arguments and invoke cowsay.
func (c *CLI) mow(argv []string) error {
	var opts options
	args, err := c.parseOptions(&opts, argv)
	if err != nil {
		return err
	}

	if opts.List {
		cowPaths, err := cowsay.Cows()
		if err != nil {
			return err
		}
		for _, cowPath := range cowPaths {
			if cowPath.LocationType == cowsay.InBinary {
				fmt.Fprintf(c.stdout, "Cow files in binary:\n")
			} else {
				fmt.Fprintf(c.stdout, "Cow files in %s:\n", cowPath.Name)
			}
			fmt.Fprintln(c.stdout, wordwrap.WrapString(strings.Join(cowPath.CowFiles, " "), 80))
			fmt.Fprintln(c.stdout)
		}
		return nil
	}

	if err := c.mowmow(&opts, args); err != nil {
		return err
	}

	return nil
}

func (c *CLI) parseOptions(opts *options, argv []string) ([]string, error) {
	p := flags.NewParser(opts, flags.None)
	args, err := p.ParseArgs(argv)
	if err != nil {
		return nil, err
	}

	if opts.Help {
		c.stdout.Write(c.usage())
		os.Exit(0)
	}

	return args, nil
}

func (c *CLI) usage() []byte {
	year := strconv.Itoa(time.Now().Year())
	return []byte(c.program() + ` version ` + c.Version + `, (c) ` + year + ` codehex
Usage: ` + c.program() + ` [-bdgpstwy] [-h] [-e eyes] [-f cowfile] [--random]
          [-l] [-n] [-T tongue] [-W wrapcolumn]
          [--bold] [--rainbow] [--aurora] [--super] [message]

Original Author: (c) 1999 Tony Monroe
`)
}

func (c *CLI) generateOptions(opts *options) []cowsay.Option {
	o := make([]cowsay.Option, 0, 8)
	if opts.File == "-" {
		cows := cowList()
		idx, _ := fuzzyfinder.Find(cows, func(i int) string {
			return cows[i]
		})
		opts.File = cows[idx]
	}
	o = append(o, cowsay.Type(opts.File))
	if c.Thinking {
		o = append(o,
			cowsay.Thinking(),
			cowsay.Thoughts('o'),
		)
	}
	if opts.Random {
		o = append(o, cowsay.Random())
	}
	if opts.Eyes != "" {
		o = append(o, cowsay.Eyes(opts.Eyes))
	}
	if opts.Tongue != "" {
		o = append(o, cowsay.Tongue(opts.Tongue))
	}
	if opts.Width > 0 {
		o = append(o, cowsay.BallonWidth(uint(opts.Width)))
	}
	if opts.NewLine {
		o = append(o, cowsay.DisableWordWrap())
	}
	return selectFace(opts, o)
}

func cowList() []string {
	cows, err := cowsay.Cows()
	if err != nil {
		return cowsay.CowsInBinary()
	}
	list := make([]string, 0)
	for _, cow := range cows {
		list = append(list, cow.CowFiles...)
	}
	return list
}

func (c *CLI) phrase(opts *options, args []string) string {
	if len(args) > 0 {
		return strings.Join(args, " ")
	}
	lines := make([]string, 0, 40)
	scanner := bufio.NewScanner(c.stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return strings.Join(lines, "\n")
}

func (c *CLI) mowmow(opts *options, args []string) error {
	phrase := c.phrase(opts, args)
	o := c.generateOptions(opts)
	if opts.Super {
		return super.RunSuperCow(phrase, opts.Bold, o...)
	}

	say, err := cowsay.Say(phrase, o...)
	if err != nil {
		var notfound *cowsay.NotFound
		if errors.As(err, &notfound) {
			return fmt.Errorf("could not find %s cowfile", notfound.Cowfile)
		}
		return err
	}

	options := make([]decoration.Option, 0)

	if opts.Bold {
		options = append(options, decoration.WithBold())
	}
	if opts.Rainbow {
		options = append(options, decoration.WithRainbow())
	}
	if opts.Aurora {
		options = append(options, decoration.WithAurora(rand.Intn(256)))
	}

	w := decoration.NewWriter(c.stdout, options...)
	fmt.Fprintln(w, say)

	return nil
}

func selectFace(opts *options, o []cowsay.Option) []cowsay.Option {
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
