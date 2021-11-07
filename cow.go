package cowsay

import (
	"fmt"
	"math/rand"
	"strings"
)

// Cow struct!!
type Cow struct {
	eyes            string
	tongue          string
	typ             string
	thoughts        rune
	thinking        bool
	bold            bool
	isAurora        bool
	isRainbow       bool
	ballonWidth     int
	disableWordWrap bool

	buf strings.Builder
}

// New returns pointer of Cow struct that made by options
func New(options ...Option) (*Cow, error) {
	cow := &Cow{
		eyes:        "oo",
		tongue:      "  ",
		thoughts:    '\\',
		typ:         "cows/default.cow",
		ballonWidth: 40,
	}
	for _, o := range options {
		if err := o(cow); err != nil {
			return nil, err
		}
	}
	return cow, nil
}

// Say returns string that said by cow
func (cow *Cow) Say(phrase string) (string, error) {
	mow, err := cow.GetCow()
	if err != nil {
		return "", err
	}

	said := cow.Balloon(phrase) + mow

	if cow.isRainbow {
		return cow.Rainbow(said), nil
	}
	if cow.isAurora {
		return cow.Aurora(rand.Intn(256), said), nil
	}
	return said, nil
}

// Clone returns a copy of cow.
//
// If any options are specified, they will be reflected.
func (cow *Cow) Clone(options ...Option) (*Cow, error) {
	ret := new(Cow)
	*ret = *cow
	ret.buf.Reset()
	for _, o := range options {
		if err := o(ret); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

// Option defined for Options
type Option func(*Cow) error

// Eyes specifies eyes
// The specified string will always be adjusted to be equal to two characters.
func Eyes(s string) Option {
	return func(c *Cow) error {
		c.eyes = adjustTo2Chars(s)
		return nil
	}
}

// Tongue specifies tongue
// The specified string will always be adjusted to be less than or equal to two characters.
func Tongue(s string) Option {
	return func(c *Cow) error {
		c.tongue = adjustTo2Chars(s)
		return nil
	}
}

func adjustTo2Chars(s string) string {
	if len(s) >= 2 {
		return s[:2]
	}
	if len(s) == 1 {
		return s + " "
	}
	return "  "
}

func containCows(t string) bool {
	for _, cow := range AssetNames() {
		if t == cow {
			return true
		}
	}
	return false
}

// NotFound is indicated not found the cowfile.
type NotFound struct {
	Cowfile string
}

var _ error = (*NotFound)(nil)

func (n *NotFound) Error() string {
	return fmt.Sprintf("not found %q cowfile", n.Cowfile)
}

// Type specify name of the cowfile
func Type(s string) Option {
	if s == "" {
		s = "cows/default.cow"
	}
	if !strings.HasSuffix(s, ".cow") {
		s += ".cow"
	}
	if !strings.HasPrefix(s, "cows/") {
		s = "cows/" + s
	}
	return func(c *Cow) error {
		if containCows(s) {
			c.typ = s
			return nil
		}
		s = strings.TrimPrefix(s, "cows/")
		s = strings.TrimSuffix(s, ".cow")
		return &NotFound{Cowfile: s}
	}
}

// Thinking enables thinking mode
func Thinking() Option {
	return func(c *Cow) error {
		c.thinking = true
		return nil
	}
}

// Thoughts Thoughts allows you to specify
// the rune that will be drawn between
// the speech bubbles and the cow
func Thoughts(thoughts rune) Option {
	return func(c *Cow) error {
		c.thoughts = thoughts
		return nil
	}
}

// Random specifies something .cow from cows directory
func Random() Option {
	pick := pickCow()
	return func(c *Cow) error {
		c.typ = pick
		return nil
	}
}

func pickCow() string {
	cows := AssetNames()
	n := len(cows)
	rand.Shuffle(n, func(i, j int) {
		cows[i], cows[j] = cows[j], cows[i]
	})
	return cows[rand.Intn(n)]
}

// Bold enables bold mode
func Bold() Option {
	return func(c *Cow) error {
		c.bold = true
		return nil
	}
}

// Aurora enables aurora mode
func Aurora() Option {
	return func(c *Cow) error {
		c.isAurora = true
		return nil
	}
}

// Rainbow enables raibow mode
func Rainbow() Option {
	return func(c *Cow) error {
		c.isRainbow = true
		return nil
	}
}

// BallonWidth specifies ballon size
func BallonWidth(size uint) Option {
	return func(c *Cow) error {
		c.ballonWidth = int(size)
		return nil
	}
}

// DisableWordWrap disables word wrap.
// Ignoring width of the ballon.
func DisableWordWrap() Option {
	return func(c *Cow) error {
		c.disableWordWrap = true
		return nil
	}
}
