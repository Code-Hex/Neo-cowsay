package main

import (
	"fmt"

	cowsay "github.com/Code-Hex/Neo-cowsay/v2"
)

func main() {
	if false {
		simple()
	} else {
		complex()
	}
}

func simple() {
	say, err := cowsay.Say(
		"Hello",
		cowsay.Type("default"),
		cowsay.BalloonWidth(40),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(say)
}

func complex() {
	cow, err := cowsay.New(
		cowsay.BalloonWidth(40),
		//cowsay.Thinking(),
		cowsay.Random(),
	)
	if err != nil {
		panic(err)
	}
	say, err := cow.Say("Hello")
	if err != nil {
		panic(err)
	}
	fmt.Println(say)
}
