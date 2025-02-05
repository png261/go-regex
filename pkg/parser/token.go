package parser

type TokenType uint8

const (
	Group           TokenType = iota
	Bracket         TokenType = iota
	Or              TokenType = iota
	Repeat          TokenType = iota
	Literal         TokenType = iota
	GroupUncaptured TokenType = iota
)

type Token struct {
	TokenType TokenType
	Value     interface{}
}
