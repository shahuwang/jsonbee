package main

import (
	"fmt"
	"unicode/utf8"
)

type Item struct {
	Start int64
	Type  ItemType
	Val   []byte
}

func (this Item) String() string {
	return string(this.Val)
}

type Lexer struct {
	Pos            int64
	Width          int64
	Start          int64
	LeftArrayParen int
	LeftDictParen  int
	Length         int64
	Input          []byte
	Items          chan Item
}

func NewLexer(data []byte) *Lexer {
	l := new(Lexer)
	l.Input = data
	l.Items = make(chan Item)
	l.Length = int64(len(l.Input))
	return l
}

const EOF = -1

func (l *Lexer) Next() rune {
	if l.Pos == l.Length {
		return EOF
	}
	r, n := utf8.DecodeRune(l.Input[l.Pos:])
	l.Width = int64(n)
	l.Pos = l.Pos + l.Width
	return r
}

func (l *Lexer) Back() {
	l.Pos = l.Pos - l.Width
}

func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Back()
	return r
}

func (l *Lexer) Ignore() {
	l.Start = l.Pos
}

func (l *Lexer) Emit(Type ItemType) {
	l.Items <- Item{Start: l.Pos, Type: Type, Val: l.Input[l.Start:l.Pos]}
	l.Start = l.Pos
}

func (l *Lexer) Errorf(format string, args ...interface{}) StateFn {
	l.Items <- Item{Start: l.Pos, Type: ItemError, Val: []byte(fmt.Sprintf(format, args...))}
	return nil
}

func (l *Lexer) Run() {
	fn := StateAction(l)
	for {
		if fn == nil {
			break
		}
		fn = fn(l)
	}
}

func main() {

}
