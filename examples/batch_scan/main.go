package main

import (
	"fmt"

	"github.com/monkey92t/afrog/v2"
)

func main() {
	if err := afrog.NewScanner([]string{}, afrog.Scanner{
		TargetsFile: "./urls.txt",
	}); err != nil {
		fmt.Println(err.Error())
	}
}
