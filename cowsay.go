package main

import (
	"fmt"
	"math/rand"
	"strings"
)

type Cow struct {
	Phrase      string
	Eyes        string
	Tongue      string
	Type        string
	Random      bool
	Thinking    bool
	Rainbow     bool
	BallonWidth int
}

func main() {
	say, err := Say(&Cow{
		Phrase:      "オッピハートあああああああああああ",
		Eyes:        "oo",
		Tongue:      "  ",
		Random:      true,
		Rainbow:     true,
		BallonWidth: 40,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(say)
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
	}

	return mow, nil
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
