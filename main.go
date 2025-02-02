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
	groupUncaptured tokenType = iota
)

type token struct {
	tokenType tokenType
	value     interface{}
}

type parseContext struct {
	pos    int
	tokens []token
}

func parse(regex string) *parseContext {
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

func parseGroup(regex string, ctx * parseContext) {
	ctx.pos += 1
	for regex[ctx.pos] != ')' {
		process(regex, ctx)
		ctx.pos += 1
	}
}

func parseBracket(regex string, ctx * parseContext ){
	ctx.pos ++
	var literals [] string
	for regex[ctx.pos] != ']' {
		ch := regex[ctx.pos]

		if ch == '-' {
			next := regex[ctx.pos + 1]
			prev := literals[len(literals) - 1][0]
			literals[len(literals) - 1] = fmt.Sprintf("%c%c", prev, next)
			ctx.pos++
		} else {
			literals = append(literals, fmt.Sprintf("%c",ch))
		}
		 ctx.pos++
	}

	literalsSet := map[uint8]bool {}


	for _, l := range literals {
		for i := l[0]; i <= l[len(l) - 1]; i++ {
			literalsSet[i] = true
		}
	}

	ctx.tokens = append(ctx.tokens, token{
		tokenType: bracket, 
		value: literalsSet,
	})
}  

func main() {
	fmt.Println("Go Regex")
}
