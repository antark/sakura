// sakura
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type Value struct {
	value_type int
	name       string
	value      interface{}
}
type Token struct {
	token_type int
	name       string
	value      int64
}
type Exception struct {
	msg string
}

func (e Exception) Error() string {
	return e.msg
}

type source struct {
	*bufio.Reader
	buffer *Token
}

const (
	NOTYPE     = iota // 无类型
	DELIMITER  = iota // 分隔符
	NUMBER     = iota // 数字
	IDENTIFIER = iota // 标识符
	UNKNOWN    = iota // 未知类型
)

func (src *source) unget_token(token Token) {
	src.buffer = &token
}

func (src *source) next_token() Token {
	var tok Token
	if src.buffer != nil {
		tok = *src.buffer
		src.buffer = nil
		return tok
	}

	var token_type int = NOTYPE
	var token []rune
	for {
		ch, _, err := src.ReadRune()
		if err != nil {
			return tok
		}

		// 空白字符
		if unicode.IsSpace(ch) && token_type == NOTYPE {
			continue
		}

		// 分隔符
		switch ch {
		case '(', ')', '{', '}', '+', '-', '*', '/', '=', ';':
			if token_type == NOTYPE {
				tok.token_type = DELIMITER
				tok.name = string(append(token, ch))
				return tok
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

		tok.token_type = token_type
		tok.name = string(token)
		if tok.token_type == NUMBER {
			tok.value, _ = strconv.ParseInt(tok.name, 10, 64)
		}
		return tok
	}
	return tok
}

var src *source = &source{bufio.NewReader(os.Stdin), nil} // source
var symbol_table map[string]int64 = make(map[string]int64, 20)

func run() {
	for {
		token := src.next_token()
		if token.token_type == UNKNOWN {
			break
		}
		switch token.name {
		case "quit":
			return
		case "help":
			fmt.Println("hehe")
			continue
		}
		src.unget_token(token)
		statement()
	}
}

func statement() {
	token := src.next_token()
	if token.token_type == UNKNOWN {
		return
	}
	switch string(token.name) {
	case "let":
		err := declaration()
		if err != nil {
			fmt.Println("error encounter:", err.Error())
		}
	default:
		src.unget_token(token)
		value := expression()
		fmt.Println(value)
	}
	token = src.next_token()
	if token.token_type != DELIMITER || token.name != ";" {
		fmt.Println("statement: expect ; ")
		src.unget_token(token)
	}
}

// 常量定义类似于 abc = 3.14
func declaration() error {
	var err *Exception
	key := src.next_token()
	if key.token_type != IDENTIFIER {
		err = &Exception{"declare: name expected"}
	}
	if symbol_table[key.name] != 0 {
		err = &Exception{"declare: name declare twice"}
	}

	equal := src.next_token()
	if equal.token_type != IDENTIFIER && equal.name != "=" {
		return Exception{"declare: = expected"}
	}
	value := expression()

	if err != nil {
		return *err
	}
	symbol_table[key.name] = value
	return nil
}

// Expression = Term | Term + Term | Term - Term
func expression() int64 {
	left := term()

	for {
		op := src.next_token()
		switch op.name {
		case "+":
			left += term()
		case "-":
			left -= term()
		default:
			src.unget_token(op)
			return left
		}
	}
}
func term() int64 {
	left := primary()

	for {
		op := src.next_token()
		switch op.name {
		case "*":
			left *= primary()
		case "/":
			left /= primary()
		default:
			src.unget_token(op)
			return left
		}
	}
}

func primary() int64 {
	token := src.next_token()
	if token.token_type == IDENTIFIER {
		return symbol_table[token.name]
	}

	if token.token_type == NUMBER {
		return token.value
	}
	return 0
}

func main() {
	run()
}
