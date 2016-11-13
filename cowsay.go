package main

import (
	"fmt"
	"math/rand"
	"strings"
)

type Cow struct {
	Phrase   string
	Eyes     string
	Tongue   string
	Random   bool
	Type     string
	Thinking bool
}

func main() {
	say, err := Say(&Cow{
		Phrase: "Hello, World\nWWWWWWWWWWWWWWWWWWWmkcnrin",
		Eyes:   "oo",
		Tongue: "  ",
		Random: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(say)
}

func Say(cow *Cow) (string, error) {

	if cow.Random {
		cow.Type = PickCow()
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

	return fmt.Sprintf("%s%s", cow.balloon(40), mow), nil
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
	return strings.Join(mow, "\n"), nil
}

func PickCow() string {
	cows := AssetNames()
	return pickup(cows)
}

func pickup(cows []string) string {
	n := len(cows)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		cows[i], cows[j] = cows[j], cows[i]
	}
	return cows[rand.Intn(n)]
}
