package cowsay

import (
	"math/rand"
	"sort"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Say to return cowsay string.
func Say(options ...Option) (string, error) {
	cow, err := New(options...)
	if err != nil {
		return "", err
	}
	mow, err := cow.GetCow()
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

// Cows to get list of cows
func Cows() []string {
	assets := AssetNames()
	cows := make([]string, 0, len(assets))
	for _, key := range assets {
		cows = append(cows, strings.TrimSuffix(strings.TrimPrefix(key, "cows/"), ".cow"))
	}

	sort.Strings(cows)
	return cows
}

// GetCow to get cow's ascii art
func (cow *Cow) GetCow() (string, error) {
	src, err := Asset(cow.typ)
	if err != nil {
		return "", err
	}

	if len(cow.eyes) > 2 {
		cow.eyes = cow.eyes[0:2]
	}

	if len(cow.tongue) > 2 {
		cow.tongue = cow.tongue[0:2]
	}

	r := strings.NewReplacer(
		"\\\\", "\\",
		"\\@", "@",
		"\\$", "$",
		"$eyes", cow.eyes,
		"${eyes}", cow.eyes,
		"$tongue", cow.tongue,
		"${tongue}", cow.tongue,
		"$thoughts", string(cow.thoughts),
		"${thoughts}", string(cow.thoughts),
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
