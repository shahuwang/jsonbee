package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Parser struct {
	Input  []byte
	Lex    *Lexer
	result interface{}
	err    error
}

func NewParser(data []byte) *Parser {
	lex := NewLexer(data)
	p := new(Parser)
	p.Input = data
	p.Lex = lex
	return p
}

func (p *Parser) Parse() (interface{}, error) {
	go p.parse()
	p.Lex.Run()
	return p.result, p.err
}

func (p *Parser) parse() {
	flag := false
	for item := range p.Lex.Items {
		if flag {
			// 说明解析完字典或者数组后，还有其他元素，语法不对
			p.err = errors.New("Error: wrong end")
			return
		}
		if item.Type == ItemArrayLeft {
			p.result, p.err = p.parseArray()
		}
		if item.Type == ItemDictLeft {
			p.result, p.err = p.parseDict()
		}
	}
}

func (p *Parser) parseArray() ([]interface{}, error) {
	array := make([]interface{}, 0)
	needComma := false
	initStart := true
	for item := range p.Lex.Items {
		if initStart && item.Type == ItemArrayRight {
			return array, nil
		}
		initStart = false
		if needComma && item.Type == ItemArrayRight {
			return array, nil
		}
		if needComma && item.Type != ItemComma {
			return array, errors.New(
				fmt.Sprintf("Error: need comma between element, position: %d", item.Start))
		}
		if needComma && item.Type == ItemComma {
			needComma = false
			continue
		}
		switch item.Type {
		case ItemError:
			return array, errors.New(
				fmt.Sprintf("Error: %s, position: %d", item.String(), item.Start))
		case ItemString:
			array = append(array, string(item.Val))
			needComma = true
		case ItemFloat:
			f, err := strconv.ParseFloat(string(item.Val), 64)
			if err != nil {
				return array, errors.New(
					fmt.Sprintf("Error: %s, position: %d", err.Error(), item.Start))
			}
			array = append(array, f)
			needComma = true
		case ItemInteger:
			i, err := strconv.ParseInt(string(item.Val), 10, 64)
			if err != nil {
				return array, errors.New(
					fmt.Sprintf("Error: %s, position: %d", err.Error(), item.Start))
			}
			array = append(array, i)
			needComma = true
		case ItemNull:
			array = append(array, nil)
			needComma = true
		case ItemArrayRight:
			return array, nil
		case ItemArrayLeft:
			elem, err := p.parseArray()
			if err != nil {
				return array, err
			}
			array = append(array, elem)
			needComma = true
		case ItemDictLeft:
			elem, err := p.parseDict()
			if err != nil {
				return array, err
			}
			array = append(array, elem)
			needComma = true
		case ItemComma:
			continue
		default:
			msg := fmt.Sprintf("Error: wrong element in array, position: %d", item.Start)
			err := errors.New(msg)
			return array, err
		}
	}
	return array, nil
}

func (p *Parser) parseDict() (map[string]interface{}, error) {
	dict := make(map[string]interface{})
	needComma := false
	needColon := false
	needKey := true
	fieldName := ""
	initStart := true // 刚开始，字典可以为空
	for item := range p.Lex.Items {
		if item.Type == ItemDictRight && initStart {
			return dict, nil
		}
		if item.Type == ItemDictRight && needComma {
			return dict, nil
		}
		initStart = false
		if needKey && item.Type != ItemString {
			msg := fmt.Sprintf("Error: dict field needs to be string, position: %d", item.Start)
			return dict, errors.New(msg)
		}
		if needKey && item.Type == ItemString {
			fieldName = string(item.Val)
			needKey = false
			needColon = true
			continue
		}
		if needColon && item.Type != ItemColon {
			msg := fmt.Sprintf("Error: needs colon after %s, position: %d", fieldName, item.Start)
			return dict, errors.New(msg)
		}
		if needColon && item.Type == ItemColon {
			needColon = false
			continue
		}
		if needComma && item.Type != ItemComma {
			msg := fmt.Sprintf("Error: needs comma after value, position: %d", item.Start)
			return dict, errors.New(msg)
		}
		if needComma && item.Type == ItemComma {
			needComma = false
			needKey = true
			continue
		}
		switch item.Type {
		case ItemError:
			msg := fmt.Sprintf("Error: %s, position: %d", item.String(), item.Start)
			return dict, errors.New(msg)
		case ItemString:
			dict[fieldName] = string(item.Val)
		case ItemFloat:
			f, err := strconv.ParseFloat(string(item.Val), 64)
			if err != nil {
				return dict, errors.New(
					fmt.Sprintf("Error: %s, position: %d", err.Error(), item.Start))
			}
			dict[fieldName] = f
		case ItemInteger:
			i, err := strconv.ParseInt(string(item.Val), 10, 64)
			if err != nil {
				return dict, errors.New(
					fmt.Sprintf("Error: %s, position: %d", err.Error(), item.Start))
			}
			dict[fieldName] = i
		case ItemNull:
			dict[fieldName] = nil
		case ItemArrayLeft:
			elem, err := p.parseArray()
			if err != nil {
				return dict, err
			}
			dict[fieldName] = elem
		case ItemDictLeft:
			elem, err := p.parseDict()
			if err != nil {
				return dict, err
			}
			dict[fieldName] = elem
		default:
			msg := fmt.Sprintf("Error, wrong element %s in dict, position: %d", item.String(), item.Start)
			err := errors.New(msg)
			return dict, err
		}
		needComma = true
	}
	return dict, nil
}
