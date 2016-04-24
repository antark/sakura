// sakura
package sakura

import (
	"bufio"
	"fmt"
	"io"
	"sakura/types"
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

func (src *source) unget_token(token Token) {
	src.buffer = &token
}

func (src *source) next_token() (tok Token) {
	var token_type int = types.NOTYPE
	var token []rune

	// post 处理
	defer func() {
		if tok.token_type == types.NOTYPE {
			tok.token_type = token_type
		}
		if tok.value == nil {
			switch tok.token_type {
			case types.NUMBER:
				// 数值转换
				if strings.Contains(tok.name, ".") {
					tok.value, _ = strconv.ParseFloat(tok.name, 64)
				} else {
					tok.value, _ = strconv.ParseInt(tok.name, 10, 64)
				}
			case types.IDENTIFIER:
				// 预定义标识符
				switch tok.name {
				case "true":
					tok.token_type = types.BOOLEAN
					tok.value = true
				case "false":
					tok.token_type = types.BOOLEAN
					tok.value = false
				}
			default:
			}
		}
	}()

	// buffer 中有，直接返回
	if src.buffer != nil {
		tok = *src.buffer
		src.buffer = nil
		return tok
	}
	for {
		ch, _, err := src.ReadRune()
		if err != nil {
			return tok
		}

		switch token_type {
		case types.NOTYPE: // 无类型
			switch ch {
			case '\t', '\n', '\v', '\f', '\r', ' ': // 空白字符
			case '(', ')', '{', '}', ';': // 分隔符
				token_type = types.DELIMITER
				tok.name = string(append(token, ch))
				return tok
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
				token_type = types.NUMBER
				token = append(token, ch)
			case '"': // "，字符串
				token_type = types.CHARS
			case '+', '-', '*', '/', '%', '^': // + - * / %
				token_type = types.OPERATOR
				tok.name = string(append(token, ch))
				return tok
			case '&': // &, && 操作符
				ch, _, _ := src.ReadRune()
				switch ch {
				case '&':
					tok.name = "&&"
				default:
					src.UnreadRune()
					tok.name = "&"
				}
				token_type = types.OPERATOR
				return tok
			case '|': // |, || 操作符
				ch, _, _ := src.ReadRune()
				switch ch {
				case '|':
					tok.name = "||"
				default:
					src.UnreadRune()
					tok.name = "|"
				}
				token_type = types.OPERATOR
				return tok
			case '!': // !, != 操作符
				ch, _, _ := src.ReadRune()
				switch ch {
				case '=':
					tok.name = "!="
				default:
					src.UnreadRune()
					tok.name = "!"
				}
				token_type = types.OPERATOR
				return tok
			case '=': // =, == 操作符
				ch, _, _ := src.ReadRune()
				switch ch {
				case '=':
					tok.name = "=="
				default:
					src.UnreadRune()
					tok.name = "="
				}
				token_type = types.OPERATOR
				return tok
			case '>': // >, >=, >> 操作符
				ch, _, _ := src.ReadRune()
				switch ch {
				case '=':
					tok.name = ">="
				case '>':
					tok.name = ">>"
				default:
					src.UnreadRune()
					tok.name = ">"
				}
				token_type = types.OPERATOR
				return tok
			case '<': // <, <=, << 操作符
				ch, _, _ := src.ReadRune()
				switch ch {
				case '=':
					tok.name = "<="
				case '<':
					tok.name = "<<"
				default:
					src.UnreadRune()
					tok.name = "<"
				}
				token_type = types.OPERATOR
				return tok
			default:
				if unicode.IsLetter(ch) {
					token_type = types.IDENTIFIER
					token = append(token, ch)
				} else {
					panic("some special char encountered.")
				}
			}

		case types.CHARS:
			switch ch {
			case '"':
				tok.name = string(token)
				return tok
			default:
				token = append(token, ch)
			}
		case types.NUMBER:
			switch ch {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
				token = append(token, ch)
			default:
				src.UnreadRune() // 回退一个字符
				tok.name = string(token)
				return tok
			}
		case types.IDENTIFIER:
			if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
				token = append(token, ch)
			} else {
				src.UnreadRune() // 回退一个字符
				tok.name = string(token)
				return tok
			}
		}
	}
}

var src *source // input source
var symbol_table map[string]*Value = make(map[string]*Value, 20)

func Run(reader io.Reader) {
	src = &source{bufio.NewReader(reader), nil} // source
	for {
		token := src.next_token()
		if token.token_type == types.NOTYPE {
			break
		}
		switch token.name {
		case "quit":
			return
		case "help":
			help()
			continue
		}
		src.unget_token(token)
		statement()
	}
}

func statement() {
	token := src.next_token()
	if token.token_type == types.NOTYPE {
		return
	}
	switch string(token.name) {
	case "let":
		err := declaration()
		if err != nil {
			fmt.Println("error encounter:", err.Error())
		}
		// fmt.Println("statement:", symbol_table)
	default:
		src.unget_token(token)
		value := expression(1)
		fmt.Println(value.value)
	}
	token = src.next_token()
	if token.token_type != types.DELIMITER || token.name != ";" {
		// fmt.Println("statement: expect ; ")
		src.unget_token(token)
	}
}

// 常量定义类似于 abc = 3.14
func declaration() error {
	var err *Exception
	key := src.next_token()
	if key.token_type != types.IDENTIFIER {
		err = &Exception{"declare: name expected"}
	}
	/*
		if symbol_table[key.name] != nil {
			err = &Exception{"declare: name declare twice"}
		}
	*/

	equal := src.next_token()
	if equal.token_type != types.IDENTIFIER && equal.name != "=" {
		return Exception{"declare: = expected"}
	}
	value := expression(1)

	if err != nil {
		return *err
	}
	symbol_table[key.name] = &value
	return nil
}

// Expression = Term | Term + Term | Term - Term

var op_levels = map[int][]string{
	1: []string{"||"},
	2: []string{"&&"},
	3: []string{"==", "!=", ">", ">=", "<", "<="},
	4: []string{"+", "-", "|", "^"},
	5: []string{"*", "/", "%", "<<", ">>", "&"},
}

var expression_max_level = 5

func expression(level int) (value Value) {
	var target func(int) Value
	if level < expression_max_level {
		target = expression
	} else {
		target = primary
	}

	left := target(level + 1)
	// fmt.Println("----", left.value) // todo

	var ops []string = op_levels[level]

NEXT:
	for {
		op := src.next_token()

		for _, op_item := range ops {
			// fmt.Println("----", op_item) // todo
			if op.name != op_item {
				continue
			}
			right := target(level + 1)
			// fmt.Println("----", right.value) // todo
			sum := types.Op_values(op.name, left.value, right.value)
			// fmt.Println("----", sum) // todo
			left.value = sum
			continue NEXT
		}
		src.unget_token(op)
		return left
	}
}

func primary(level int) (value Value) {
	token := src.next_token()
	// fmt.Println("------------", token.name)

	switch token.token_type {
	// 标识符
	case types.IDENTIFIER:
		value = *symbol_table[token.name]
	// 数值
	case types.NUMBER:
		switch token.value.(type) {
		case uint64:
			value.value_type = types.INT_64
		case float64:
			value.value_type = types.FLOAT_64
		}
		value.value = token.value
	// 布尔
	case types.BOOLEAN:
		value.value_type = types.BOOL
		value.value = token.value
	// 字符串
	case types.CHARS:
		value.value_type = types.STRING
		value.value = token.name
	default:
		switch token.name {
		case "+", "-", "!", "^":
			value = primary(level)
			value.value = types.Op_values(token.name, value.value, nil)
		case "(":
			value = expression(1)
			token = src.next_token()
			if token.name != ")" {
				panic(") expected")
			}
		}
	}
	return value
}
