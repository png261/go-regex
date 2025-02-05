package main

import (
	"flag"
	"fmt"
	"github.com/png261/go-regex/pkg/nfa"
	"github.com/png261/go-regex/pkg/parser"
	"log"
)

func main() {
	var regex string
	var str string
	flag.StringVar(&regex, "r", "", "The regular expression to test")
	flag.StringVar(&str, "s", "", "The string to test against the regular expression")
	flag.Parse()

	if regex == "" || str == "" {
		log.Fatal("You must provide both regex and string arguments.")
	}

	ctx := parser.Parse(regex)
	nfa := nfa.ToNfa(ctx)
	fmt.Println(nfa.Check(str, -1))
}
