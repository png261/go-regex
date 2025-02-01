package main

import (
	"fmt"
)

type tokenType unit8

const (
	group           tokenType = iota
	bracket         tokenType = iota
	or              tokenType = iota
	repeat          tokenType = iota
	literal         tokenType = iota
	gruopUncaptured tokenType = iota
)

type token struct {
	tokenType tokenType
	value     interface{}
}

type parseContext struct {
	pos    int
	tokens []token
}

func parse(regedx string) *parseContext {
	ctx := &parseContext{
		pos:    0,
		tokens: []token{},
	}

	for ctx.pos < len(regex) {
		process(regex, ctx)
		ctx.pos++
	}

	return ctx
}

func process(regex string, ctx *parseContext) {
	ch := regex[ctx.pos]
	if ch == '(' {
		groupCtx := &parseContext{
			pos:    ctx.pos,
			tokens: []token{},
		}

		parseGroup(regex, groupCtx)
		ctx.tokens = append(ctx.tokens, token{
			tokenType: group,
			value:     groupCtx.tokens,
		})
	} else if ch == '[' {
		parseBracket(regex, ctx)
	} else if ch == '|' {
		parseOr(regex, ctx)
	} else if ch == '*' || ch == '?' || ch == '+' {
		parseRepeat(regex, ctx))
	} else if ch == '{' {
		parseRepeatSpecified(regex, ctx)
	} else {
		t := token {
			tokenType: literal, 
			value: ch,
		}

		ctx.tokens = append(ctx.tokens, t)
	}

}

func main() {
	fmt.Println("Go Regex")
}
