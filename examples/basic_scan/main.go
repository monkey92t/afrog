package main

import (
	"fmt"

	"github.com/monkey92t/afrog/v2"
)

func main() {
	if err := afrog.NewScanner([]string{"http://example.com"}, afrog.Scanner{}); err != nil {
		fmt.Println(err.Error())
	}
}
