package main

import (
	"fmt"
	"strings"
)

type StateFn func(*Lexer) StateFn

func StateAction(l *Lexer) StateFn {
	r := l.Next()
	switch r {
	case '\t', '\n', '\r', ' ':
		l.Ignore()
		return StateAction
	case '[':
		return StateLeftArray
	case ']':
		return StateRightArray
	case '"':
		return StateString
	case ',':
		return StateComma
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return StateNumber
	case '{':
		return StateLeftDict
	case '}':
		return StateRightDict
	case ':':
		return StateColon
	case EOF:
		return nil
	default:
		fmt.Println("88888888888888")
		fmt.Println(r)
		return l.Errorf("wrong element")
	}
	return nil
}

func StateColon(l *Lexer) StateFn {
	l.Emit(ItemColon)
	l.Start = l.Pos
	return StateAction
}

func StateLeftDict(l *Lexer) StateFn {
	l.Emit(ItemDictLeft)
	l.LeftDictParen += 1
	l.Start = l.Pos
	return StateAction
}

func StateRightDict(l *Lexer) StateFn {
	if l.LeftDictParen == 0 {
		return l.Errorf("no matching left dict paren")
	}
	if l.LeftDictParen > 1 && l.Peek() == EOF {
		return l.Errorf("no matching right dict paren")
	}
	l.Emit(ItemDictRight)
	l.Start = l.Pos
	l.LeftDictParen -= 1
	return StateAction
}

func StateLeftArray(l *Lexer) StateFn {
	l.Emit(ItemArrayLeft)
	l.LeftArrayParen += 1
	l.Start = l.Pos
	return StateAction
}

func StateRightArray(l *Lexer) StateFn {
	if l.LeftArrayParen == 0 {
		return l.Errorf("no matching left array paren")
	}
	if l.LeftArrayParen > 1 && l.Peek() == EOF {
		return l.Errorf("no matching right array paren")
	}
	l.Emit(ItemArrayRight)
	l.Start = l.Pos
	l.LeftArrayParen -= 1
	return StateAction
}

func StateString(l *Lexer) StateFn {
	for {
		switch l.Next() {
		case '"':
			l.Emit(ItemString)
			l.Start = l.Pos
			return StateAction
		case '\\':
			// 判断下一个是否是"，处理转移字符
			r := l.Peek()
			if r == '"' {
				r = l.Next()
				if r == EOF {
					return l.Errorf("unterminated string")
				}
			}
		}
	}
}

func StateComma(l *Lexer) StateFn {
	l.Emit(ItemComma)
	l.Start = l.Pos
	return StateAction
}

func StateNumber(l *Lexer) StateFn {
	l.Back()
	firstNumber := true
	isFloat := false
	for {
		r := l.Next()
		switch r {
		case '-':
			if !isNumber(l.Peek()) {
				return l.Errorf("minus must before number")
			}
		case '.':
			if !isNumber(l.Peek()) {
				return l.Errorf("decimal point must before number")
			}
			isFloat = true
		case '0':
			if firstNumber {
				if isNumber(l.Peek()) {
					return l.Errorf("zero can no be followed by number")
				}
				firstNumber = false
			}
			fallthrough
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			firstNumber = false
			if isElemEnd(l.Peek()) {
				if isFloat {
					l.Emit(ItemFloat)
				} else {
					l.Emit(ItemInteger)
				}
				l.Start = l.Pos
				return StateAction
			}
		default:
			return l.Errorf("wrong number in %d", l.Pos)
		}
	}
}

func isNumber(n rune) bool {
	return strings.ContainsRune("0123456789", n)
}

func isElemEnd(n rune) bool {
	return strings.ContainsRune(",]}\n\r\t ", n) || n == EOF
}
