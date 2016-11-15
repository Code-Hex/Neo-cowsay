package cowsay

import (
	"math/rand"
	"sort"
	"strings"
	"time"
)

type Cow struct {
	Phrase      string
	Eyes        string
	Tongue      string
	Type        string
	Random      bool
	Aurora      bool
	Thinking    bool
	Bold        bool
	Rainbow     bool
	BallonWidth int
}

//go:generate go-bindata -pkg cowsay cows
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Say to return cowsay string.
func Say(cow *Cow) (string, error) {

	if cow.Random {
		cow.Type = pickCow()
	}

	if cow.Type == "" {
		cow.Type = "cows/default.cow"
	}

	if !strings.HasSuffix(cow.Type, ".cow") {
		cow.Type += ".cow"
	}

	if !strings.HasPrefix(cow.Type, "cows/") {
		cow.Type = "cows/" + cow.Type
	}

	mow, err := cow.getCow()
	if err != nil {
		return "", err
	}

	if cow.Rainbow {
		mow = cow.makeRainbow(cow.balloon() + mow)
	} else if cow.Aurora {
		mow = cow.makeAurora(rand.Intn(256), cow.balloon()+mow)
	}

	return mow, nil
}

// Cows to get list of cows
func Cows() []string {
	cows := make([]string, 0, len(AssetNames()))
	for _, key := range AssetNames() {
		cows = append(cows, strings.TrimSuffix(strings.TrimPrefix(key, "cows/"), ".cow"))
	}

	sort.Strings(cows)
	return cows
}

func (cow *Cow) getCow() (string, error) {
	src, err := Asset(cow.Type)
	if err != nil {
		return "", err
	}

	var thoughts string
	if cow.Thinking {
		thoughts = "o"
	} else {
		thoughts = "\\"
	}

	if len(cow.Eyes) > 2 {
		cow.Eyes = cow.Eyes[0:2]
	} else if cow.Eyes == "" {
		cow.Eyes = "oo"
	}

	if len(cow.Tongue) > 2 {
		cow.Tongue = cow.Tongue[0:2]
	} else if cow.Tongue == "" {
		cow.Tongue = "  "
	}

	r := strings.NewReplacer(
		"\\\\", "\\",
		"\\@", "@",
		"\\$", "$",
		"$eyes", cow.Eyes,
		"${eyes}", cow.Eyes,
		"$tongue", cow.Tongue,
		"${tongue}", cow.Tongue,
		"$thoughts", thoughts,
		"${thoughts}", thoughts,
	)
	newsrc := r.Replace(string(src))
	separate := strings.Split(newsrc, "\n")
	mow := make([]string, 0, len(separate))
	for _, line := range separate {
		if strings.Contains(line, "$the_cow = <<EOC") || strings.HasPrefix(line, "##") {
			continue
		}

		if strings.HasPrefix(line, "EOC") {
			break
		}

		mow = append(mow, line)
	}
	return strings.Join(mow, "\n"), nil
}

func pickCow() string {
	cows := AssetNames()
	n := len(cows)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		cows[i], cows[j] = cows[j], cows[i]
	}
	return cows[rand.Intn(n)]
}
