package main

import (
	"fmt"
	"testing"
	// "time"
)

// func TestArray(t *testing.T) {
// 	data := `["a", "", "\"", 12, 0, 0.23, 4.090, -12, -0.45, -19, 0]`
// 	lex := NewLexer([]byte(data))
// 	go func() {
// 		for item := range lex.Items {
// 			fmt.Printf("%+v\n", item)
// 		}
// 	}()
// 	lex.Run()
// 	time.Sleep(12)
// }

func TestParser(t *testing.T) {
	data := []string{
		`[] `,
		`{}`,
		`["a", "", "\"", 12, 0, 0.23, 4.090, -12, -0.45, -19, 0]`,
		`{"a": 1}`,
		`[1, "a", 1.90, {"c": 1, "b": [1,2]}]`,
		`[1, "a", [1, 3], [{"a": 1}]]`,
		`{"a": ["1", "a", {"c": 2}], "b": {"e": 1, "d": 3}, "c": 5}`,
	}
	for i, item := range data {
		fmt.Println("++++++++++++")
		fmt.Println(item)
		parser := NewParser([]byte(item))
		value, err := parser.Parse()
		if err != nil {
			panic(err)
		}
		fmt.Printf("item: %d, %+v\n", i, value)
		fmt.Println("===========\n")
	}
}
