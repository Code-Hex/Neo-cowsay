package cowsay

import (
	"math/rand"
	"strings"

	"github.com/pkg/errors"
)

// Cow struct!!
type Cow struct {
	phrase      string
	eyes        string
	tongue      string
	typ         string
	thinking    bool
	bold        bool
	isAurora    bool
	isRainbow   bool
	ballonWidth int
}

// NewCow returns pointer of Cow struct that made by options
func NewCow(options ...Option) (*Cow, error) {
	cow := &Cow{
		phrase:      "",
		eyes:        "oo",
		tongue:      "  ",
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
func (cow *Cow) Say() (string, error) {
	mow, err := cow.GetCow(0)
	if err != nil {
		return "", err
	}

	said := cow.Balloon() + mow

	if cow.isRainbow {
		return cow.Rainbow(said), nil
	}
	if cow.isAurora {
		return cow.Aurora(rand.Intn(256), said), nil
	}
	return said, nil
}

// Option defined for Options
type Option func(*Cow) error

// Phrase specifies you want to say
func Phrase(s string) Option {
	return func(c *Cow) error {
		c.phrase = s
		return nil
	}
}

// Eyes specifies eyes
// You must specify two length string
func Eyes(s string) Option {
	return func(c *Cow) error {
		if l := len(s); l != 2 {
			return errors.New("You should pass 2 length string because cow has only two eyes")
		}
		c.eyes = s
		return nil
	}
}

// Tongue specifies tongue
// You must specify two length string
func Tongue(s string) Option {
	return func(c *Cow) error {
		if l := len(s); l != 2 {
			return errors.New("You should pass 2 length string because cow has only two space on mouth")
		}
		c.tongue = s
		return nil
	}
}

func containCows(t string) bool {
	for _, cow := range AssetNames() {
		if t == cow {
			return true
		}
	}
	return false
}

// Type specify cow type that is file name of .cow
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
		return errors.Errorf("Could not find %s", s)
	}
}

// Thinking enables thinking mode
func Thinking() Option {
	return func(c *Cow) error {
		c.thinking = true
		return nil
	}
}

// Random specifies something .cow from cows directory
func Random() Option {
	return func(c *Cow) error {
		c.typ = pickCow()
		return nil
	}
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
