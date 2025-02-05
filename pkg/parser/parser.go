package parser

import (
	"fmt"
	"strconv"
	"strings"
)

func Parse(regex string) *ParseContext {
	ctx := &ParseContext{
		Pos:    0,
		Tokens: []Token{},
	}

	for ctx.Pos < len(regex) {
		process(regex, ctx)
		ctx.Pos++
	}

	return ctx
}

func process(regex string, ctx *ParseContext) {
	ch := regex[ctx.Pos]
	if ch == '(' {
		groupCtx := &ParseContext{
			Pos:    ctx.Pos,
			Tokens: []Token{},
		}

		parseGroup(regex, groupCtx)
		ctx.Tokens = append(ctx.Tokens, Token{
			TokenType: Group,
			Value:     groupCtx.Tokens,
		})
	} else if ch == '[' {
		parseBracket(regex, ctx)
	} else if ch == '|' {
		parseOr(regex, ctx)
	} else if ch == '*' || ch == '?' || ch == '+' {
		parseRepeat(regex, ctx)
	} else if ch == '{' {
		parseRepeatSpecified(regex, ctx)
	} else {
		t := Token{
			TokenType: Literal,
			Value:     ch,
		}

		ctx.Tokens = append(ctx.Tokens, t)
	}

}

func parseGroup(regex string, ctx *ParseContext) {
	ctx.Pos += 1
	for regex[ctx.Pos] != ')' {
		process(regex, ctx)
		ctx.Pos += 1
	}
}

func parseBracket(regex string, ctx *ParseContext) {
	ctx.Pos++
	var literals []string
	for regex[ctx.Pos] != ']' {
		ch := regex[ctx.Pos]

		if ch == '-' {
			next := regex[ctx.Pos+1]
			prev := literals[len(literals)-1][0]
			literals[len(literals)-1] = fmt.Sprintf("%c%c", prev, next)
			ctx.Pos++
		} else {
			literals = append(literals, fmt.Sprintf("%c", ch))
		}
		ctx.Pos++
	}

	literalsSet := map[uint8]bool{}

	for _, l := range literals {
		for i := l[0]; i <= l[len(l)-1]; i++ {
			literalsSet[i] = true
		}
	}

	ctx.Tokens = append(ctx.Tokens, Token{
		TokenType: Bracket,
		Value:     literalsSet,
	})
}

func parseOr(regex string, ctx *ParseContext) {
	rhsContext := &ParseContext{
		Pos:    ctx.Pos,
		Tokens: []Token{},
	}

	rhsContext.Pos += 1
	for rhsContext.Pos < len(regex) && regex[rhsContext.Pos] != ')' {
		process(regex, rhsContext)
		rhsContext.Pos += 1
	}

	left := Token{
		TokenType: GroupUncaptured,
		Value:     ctx.Tokens,
	}

	right := Token{
		TokenType: GroupUncaptured,
		Value:     rhsContext.Tokens,
	}

	ctx.Pos = rhsContext.Pos
	ctx.Tokens = []Token{{
		TokenType: Or,
		Value:     []Token{left, right},
	}}
}

const RepeatInfinity = -1

func parseRepeat(regex string, ctx *ParseContext) {
	ch := regex[ctx.Pos]
	var min, max int
	if ch == '*' {
		min = 0
		max = RepeatInfinity
	} else if ch == '?' {
		min = 0
		max = 1
	} else {
		min = 1
		max = RepeatInfinity
	}

	lastToken := ctx.Tokens[len(ctx.Tokens)-1]
	ctx.Tokens[len(ctx.Tokens)-1] = Token{
		TokenType: Repeat,
		Value: RepeatPayload{
			Min:   min,
			Max:   max,
			Token: lastToken,
		},
	}
}

type RepeatPayload struct {
	Min   int
	Max   int
	Token Token
}

func parseRepeatSpecified(regex string, ctx *ParseContext) {
	start := ctx.Pos + 1

	for regex[ctx.Pos] != '}' {
		ctx.Pos++
	}

	boundariesStr := regex[start:ctx.Pos]
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
			max = RepeatInfinity
		} else if value, err := strconv.Atoi(pieces[1]); err != nil {
			panic(err.Error())
		} else {
			max = value
		}
	} else {
		panic(fmt.Sprintf("There must be either 1 or 2 values specified for the quantifier: provided '%s'", boundariesStr))
	}

	lastToken := ctx.Tokens[len(ctx.Tokens)-1]
	ctx.Tokens[len(ctx.Tokens)-1] = Token{
		TokenType: Repeat,
		Value: RepeatPayload{
			Min:   min,
			Max:   max,
			Token: lastToken,
		},
	}
}
