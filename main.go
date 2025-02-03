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

func parseOr(regex string, ctx* parseContext) {
	rhsContext := &parseContext {
		pos: ctx.pos,
		tokens: [] token{},
	}

	rhsContext.pos += 1
	for rhsContext.pos <len(regex) && regex[rhsContext.pos] != ')' {
		process(regex,rhsContext)
		rhsContext.pos += 1
	}


	left := token {
		tokenType: groupUncaptured,
		value: ctx.tokens,
	}

	right := token {
		tokenType: groupUncaptured,
		value: rhsContext.tokens,
	}

	ctx.pos = rhsContext.pos
	ctx.tokens = []token{{
		tokenType: or,
		value: [] token{left, right},

	}}
}

const repeatInfinity = -1

func parseRepeat(regex string, ctx *parseContext) {
	ch := regex[ctx.pos]
	var min, max int
	if ch == '*' {
		min = 0
		max = repeatInfinity
	} else if ch == '?' {
		min = 0
		max = 1
	} else {
		min = 1
		max = repeatInfinity
	}


	lastToken := ctx.tokens[len(ctx.tokens) - 1]
	ctx.tokens[len(ctx.tokens) - 1] = token {
		tokenType: repeat,
		value: repeatPayload{ 
			min: min,
			max: max,
			token: lastToken,
		},
	}
}

func parseRepeatSpecified(regex string, ctx * parseContext) {
	start := ctx.pos + 1

	for regex[ctx.pos] != '}' {
		ctx.pos ++
	}

	boundariesStr := regex[start: ctx.pos]
	pieces := strings.Split(boundariesStr, ",")

	var min, max int
	if len(pieces) == 1 {
		if value, err := strconv.Atoi(pieces[0]); err != nil {
			panic(err.Error())
		} else {
			min = value
			max = value
		}
	} else if len(pieces) == 2 {
		if value, err := strconv.Atoi(pieces[0]); err != nil {
			panic(err.Error())
		} else {
			min = value
		}

		if pieces[1] == "" {
			max = repeatInfinity
		} else if value, err := strconv.Atoi(pieces[1]); err != nil {
			panic(err.Error())
		} else {
			max = value
		}
	} else {
		panic(fmt.Sprintf("There must be either 1 or 2 values specified for the quantifier: provided '%s'", boundariesStr))
	}

	lastToken := ctx.tokens[len(ctx.tokens) -1]
	ctx.tokens[len(ctx.tokens) -1] = token {
		tokenType: repeat, 
		value: repeatPayload {
			min: min,
			max: max,
			token: lastToken,
		},
	}
}

func main() {
	fmt.Println("Go Regex")
}
