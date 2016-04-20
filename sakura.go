// sakura
package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

type runtime struct {
}
type Sakura struct {
}

type source struct {
	*bufio.Reader
}

const (
	NOTYPE     = iota // 无类型
	DELIMITER  = iota // 分隔符
	NUMBER     = iota // 数字
	IDENTIFIER = iota // 标识符
)

func (src *source) next_token() []rune {
	var token []rune
	var token_type int = NOTYPE
	for {
		ch, _, err := src.ReadRune()
		if err != nil {
			return nil
		}

		// 空白字符
		if unicode.IsSpace(ch) && token_type == NOTYPE {
			continue
		}

		// 分隔符
		switch ch {
		case '(', ')', '{', '}', '+', '-', '*', '/', '=':
			if token_type == NOTYPE {
				return append(token, ch)
			}
		}

		// 数字
		if (unicode.IsDigit(ch) || ch == '.') && (token_type == NOTYPE || token_type == NUMBER) {
			token_type = NUMBER
			token = append(token, ch)
			continue
		}

		// 标识符
		if unicode.IsLetter(ch) && (token_type == NOTYPE || token_type == IDENTIFIER) {
			token_type = IDENTIFIER
			token = append(token, ch)
			continue
		}

		if token_type == NOTYPE {
			panic("some special char encountered.")
		}

		src.UnreadRune() // 回退一个字符
		return token
	}
	return nil
}

var src *source = &source{bufio.NewReader(os.Stdin)} // source
var symbol_table map[string]string = make(map[string]string, 20)

func main() {
	for {
		token := src.next_token()
		if token == nil {
			break
		}
		switch string(token) {
		case "let":
			left := string(src.next_token())
			assign := string(src.next_token())
			right := string(src.next_token())
			if left == "" || assign != "=" || right == "" {
				panic("let format invalid, let a = 10")
			}
			symbol_table[left] = right
		case "quit":
			break
		case "help":

		}
	}
	fmt.Println(symbol_table)
}
