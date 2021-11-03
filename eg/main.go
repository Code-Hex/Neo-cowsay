package main

import (
	"fmt"

	cowsay "github.com/Code-Hex/Neo-cowsay"
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
		cowsay.BallonWidth(40),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(say)
}

func complex() {
	cow, err := cowsay.New(
		cowsay.BallonWidth(40),
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
