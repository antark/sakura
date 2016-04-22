// sakura
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Token struct {
	token_type int
	name       string
	value      interface{}
}
type Value struct {
	value_type int
	value      interface{}
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
	DELIMITER         // 分隔符
	NUMBER            // 整数
	BOOLEAN           // 布尔类型
	STRING            // 字符串
	OPERATOR          // 操作符
	IDENTIFIER        // 标识符
	UNKNOWN           // 未知类型
)

const (
	INT_64 = iota
	FLOAT_64
	BOOL
	CHARS
)

func (src *source) unget_token(token Token) {
	src.buffer = &token
}

func (src *source) next_token() (tok Token) {
	// predeclared identifier
	defer func() {
		switch tok.name {
		case "true":
			tok.token_type = BOOLEAN
			tok.value = true
		case "false":
			tok.token_type = BOOLEAN
			tok.value = false
		}
	}()

	// buffer 中有，直接返回
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

		// 字符串
		if token_type == NOTYPE && ch == '"' {
			token_type = STRING
			continue
		}
		if token_type == STRING {
			if ch != '"' {
				token = append(token, ch)
				continue
			} else {
				tok.token_type = token_type
				tok.name = string(token)
				return tok
			}
		}

		// 空白字符
		if unicode.IsSpace(ch) && token_type == NOTYPE {
			continue
		}

		// && ||
		if ch == '&' && token_type == NOTYPE {
			ch, _, _ := src.ReadRune()
			if ch == '&' {
				tok.token_type = OPERATOR
				tok.name = string(append(token, '&', '&'))
				return tok
			} else {
				src.UnreadRune()
				tok.token_type = OPERATOR
				tok.name = string(append(token, '&'))
				return tok
			}
		}

		if ch == '|' && token_type == NOTYPE {
			ch, _, _ := src.ReadRune()
			if ch == '|' {
				tok.token_type = OPERATOR
				tok.name = string(append(token, '|', '|'))
				return tok
			} else {
				src.UnreadRune()
				tok.token_type = OPERATOR
				tok.name = string(append(token, '|'))
				return tok
			}
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
		if unicode.IsLetter(ch) && token_type == NOTYPE ||
			(unicode.IsLetter(ch) || unicode.IsDigit(ch) && token_type == IDENTIFIER) {
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

		// 整数、浮点数
		if tok.token_type == NUMBER {
			if strings.Contains(tok.name, ".") {
				tok.value, _ = strconv.ParseFloat(tok.name, 64)
			} else {
				tok.value, _ = strconv.ParseInt(tok.name, 10, 64)
			}
		}
		return tok
	}
	return tok
}

var src *source = &source{bufio.NewReader(os.Stdin), nil} // source
var symbol_table map[string]*Value = make(map[string]*Value, 20)

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
		fmt.Println(value.value)
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
	if symbol_table[key.name] != nil {
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
	symbol_table[key.name] = &value
	return nil
}

// Expression = Term | Term + Term | Term - Term
func expression() (value Value) {
	left := term()
	// fmt.Println("----", left.value) // todo

	var ops []string
	switch left.value_type {
	case INT_64:
		ops = append(ops, "+", "-")
	case FLOAT_64:
		ops = append(ops, "+", "-")
	case BOOL:
		ops = append(ops, "||")
	case CHARS:
		ops = append(ops, "+")
	}

NEXT:
	for {
		op := src.next_token()

		for _, op_item := range ops {
			// fmt.Println("----", op_item) // todo

			if op.name != op_item {
				continue
			}
			right := term()
			// fmt.Println("----", right.value) // todo

			sum := op_values(op.name, left.value, right.value)
			// fmt.Println("----", sum) // todo

			left.value = sum
			continue NEXT
		}
		src.unget_token(op)
		return left
	}
}

func op_values(op string, a interface{}, b interface{}) (value interface{}) {
	// fmt.Println(op, a, b)
	if _, yes := a.(int64); yes {
		if b == nil {
			return -a.(int64)
		}
		if op == "+" {
			ia := a.(int64)
			ib := b.(int64)
			return ia + ib
		}

		if op == "-" {
			ia := a.(int64)
			ib := b.(int64)
			return ia - ib
		}
		if op == "*" {
			ia := a.(int64)
			ib := b.(int64)
			return ia * ib
		}

		if op == "/" {
			ia := a.(int64)
			ib := b.(int64)
			return ia / ib
		}
	}

	if _, yes := a.(float64); yes {
		if b == nil {
			return -a.(float64)
		}
		if op == "+" {
			ia := a.(float64)
			ib := b.(float64)
			return ia + ib
		}

		if op == "-" {
			ia := a.(float64)
			ib := b.(float64)
			return ia - ib
		}
		if op == "*" {
			ia := a.(float64)
			ib := b.(float64)
			return ia * ib
		}

		if op == "/" {
			ia := a.(float64)
			ib := b.(float64)
			return ia / ib
		}
	}

	if _, yes := a.(string); yes {
		if op == "+" {
			ia := a.(string)
			ib := b.(string)
			return ia + ib
		}
	}

	if _, yes := a.(bool); yes {
		if op == "||" {
			ia := a.(bool)
			ib := b.(bool)
			return ia || ib
		}
		if op == "&&" {
			ia := a.(bool)
			ib := b.(bool)
			return ia && ib
		}
	}
	return nil
}

func term() (value Value) {
	left := primary()
	// fmt.Println("--------", left.value) // todo

	var ops []string
	switch left.value_type {
	case INT_64:
		ops = append(ops, "*", "/")
	case FLOAT_64:
		ops = append(ops, "*", "/")
	case BOOL:
		ops = append(ops, "&&")
	}

NEXT:
	for {
		op := src.next_token()

		for _, op_item := range ops {
			// fmt.Println("--------", op_item) // todo

			if op.name != op_item {
				continue
			}
			right := primary()
			// fmt.Println("--------", right.value) // todo

			mul := op_values(op.name, left.value, right.value)
			left.value = mul
			continue NEXT
		}
		src.unget_token(op)
		return left
	}
}

func primary() (value Value) {
	token := src.next_token()
	// fmt.Println("------------", token.name)

	switch token.token_type {
	// 标识符
	case IDENTIFIER:
		value = *symbol_table[token.name]
	// 数值
	case NUMBER:
		switch token.value.(type) {
		case uint64:
			value.value_type = INT_64
		case float64:
			value.value_type = FLOAT_64
		}
		value.value = token.value
	// 布尔
	case BOOLEAN:
		value.value_type = BOOL
		value.value = token.value
	// 字符串
	case STRING:
		value.value_type = CHARS
		value.value = token.name
	default:
		switch token.name {
		case "-":
			value = primary()
			value.value = op_values("-", value.value, nil)
		case "(":
			value = expression()
			token = src.next_token()
			if token.name != ")" {
				panic(") expected")
			}
		case "!":
		}
	}
	return value
}

func main() {
	run()
}
