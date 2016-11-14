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
	Rainbow     bool
	BallonWidth int
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

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
		mow = makeRainbow(mow)
	} else if cow.Aurora {
		mow = makeAurora(mow)
	}

	return mow, nil
}

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

	var mow []string
	for _, line := range strings.Split(newsrc, "\n") {
		if strings.Contains(line, "$the_cow = <<EOC") || strings.HasPrefix(line, "##") {
			continue
		}

		if strings.HasPrefix(line, "EOC") {
			break
		}

		mow = append(mow, line)
	}
	return cow.balloon() + strings.Join(mow, "\n"), nil
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
