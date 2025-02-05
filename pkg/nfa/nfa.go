package nfa

import (
	"github.com/png261/go-regex/pkg/parser"
)

type state struct {
	start       bool
	terminal    bool
	transitions map[uint8][]*state
}

const epsilonChar uint8 = 0

func ToNfa(ctx *parser.ParseContext) *state {
	startState, endState := tokenToNfa(&ctx.Tokens[0])

	for i := 1; i < len(ctx.Tokens); i++ {
		startNext, endNext := tokenToNfa(&ctx.Tokens[i])
		endState.transitions[epsilonChar] = append(
			endState.transitions[epsilonChar],
			startNext,
		)
		endState = endNext
	}

	start := &state{
		transitions: map[uint8][]*state{
			epsilonChar: {startState},
		},
		start: true,
	}

	end := &state{
		transitions: map[uint8][]*state{},
		terminal:    true,
	}

	endState.transitions[epsilonChar] = append(
		endState.transitions[epsilonChar],
		end,
	)

	return start
}

func tokenToNfa(t *parser.Token) (*state, *state) {
	start := &state{
		transitions: map[uint8][]*state{},
	}
	end := &state{
		transitions: map[uint8][]*state{},
	}

	switch t.TokenType {
	case parser.Literal:
		ch := t.Value.(uint8)
		start.transitions[ch] = []*state{end}
	case parser.Or:
		values := t.Value.([]parser.Token)
		left := values[0]
		right := values[1]

		s1, e1 := tokenToNfa(&left)
		s2, e2 := tokenToNfa(&right)

		start.transitions[epsilonChar] = []*state{s1, s2}
		e1.transitions[epsilonChar] = []*state{end}
		e2.transitions[epsilonChar] = []*state{end}
	case parser.Bracket:
		literals := t.Value.(map[uint8]bool)

		for l := range literals {
			start.transitions[l] = []*state{end}
		}
	case parser.Group, parser.GroupUncaptured:
		tokens := t.Value.([]parser.Token)
		start, end = tokenToNfa(&tokens[0])
		for i := 1; i < len(tokens); i++ {
			ts, te := tokenToNfa(&tokens[i])
			end.transitions[epsilonChar] = append(
				end.transitions[epsilonChar],
				ts,
			)
			end = te
		}
	case parser.Repeat:
		p := t.Value.(parser.RepeatPayload)

		if p.Min == 0 { // <1>
			start.transitions[epsilonChar] = []*state{end}
		}

		var copyCount int // <2>

		if p.Max == parser.RepeatInfinity {
			if p.Min == 0 {
				copyCount = 1
			} else {
				copyCount = p.Min
			}
		} else {
			copyCount = p.Max
		}

		from, to := tokenToNfa(&p.Token)
		start.transitions[epsilonChar] = append(
			start.transitions[epsilonChar],
			from,
		)

		for i := 2; i <= copyCount; i++ {
			s, e := tokenToNfa(&p.Token)

			to.transitions[epsilonChar] = append(
				to.transitions[epsilonChar],
				s,
			)

			from = s
			to = e

			if i > p.Min {
				s.transitions[epsilonChar] = append(
					s.transitions[epsilonChar],
					end,
				)
			}
		}

		to.transitions[epsilonChar] = append(
			to.transitions[epsilonChar],
			end,
		)

		if p.Max == parser.RepeatInfinity {
			end.transitions[epsilonChar] = append(
				end.transitions[epsilonChar],
				from,
			)
		}
	default:
		panic("unknown type of token")
	}
	return start, end
}

func (s *state) Check(input string, pos int) bool {
	ch := getChar(input, pos)

	if ch == endOfText && s.terminal {
		return true
	}

	if states := s.transitions[ch]; len(states) > 0 {
		nextState := states[0]
		if nextState.Check(input, pos+1) {
			return true
		}
	}

	for _, state := range s.transitions[epsilonChar] {
		if state.Check(input, pos) {
			return true
		}

		if ch == startOfText && state.Check(input, pos+1) {
			return true
		}
	}

	return false
}

const (
	startOfText uint8 = 1
	endOfText   uint8 = 2
)

func getChar(input string, pos int) uint8 {
	if pos >= len(input) {
		return endOfText
	}

	if pos < 0 {
		return startOfText
	}
	return input[pos]
}
