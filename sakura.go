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
	var token_type int = NOTYPE
	var token []rune

	// post 处理
	defer func() {
		if tok.token_type == NOTYPE {
			tok.token_type = token_type
		}
		if tok.value == nil {
			switch tok.token_type {
			case NUMBER:
				// 数值转换
				if strings.Contains(tok.name, ".") {
					tok.value, _ = strconv.ParseFloat(tok.name, 64)
				} else {
					tok.value, _ = strconv.ParseInt(tok.name, 10, 64)
				}
			case IDENTIFIER:
				// 预定义标识符
				switch tok.name {
				case "true":
					tok.token_type = BOOLEAN
					tok.value = true
				case "false":
					tok.token_type = BOOLEAN
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
		case NOTYPE: // 无类型
			switch ch {
			case '\t', '\n', '\v', '\f', '\r', ' ': // 空白字符
			case '(', ')', '{', '}', '+', '-', '*', '/', ';': // 分隔符
				token_type = DELIMITER
				tok.name = string(append(token, ch))
				return tok
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
				token_type = NUMBER
				token = append(token, ch)
			case '"': // "，字符串
				token_type = STRING
			case '&': // &, && 操作符
				ch, _, _ := src.ReadRune()
				switch ch {
				case '&':
					tok.name = "&&"
				default:
					src.UnreadRune()
					tok.name = "&"
				}
				token_type = OPERATOR
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
				token_type = OPERATOR
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
				token_type = OPERATOR
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
				token_type = OPERATOR
				return tok
			case '>': // !, != 操作符
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
				token_type = OPERATOR
				return tok
			case '<': // !, != 操作符
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
				token_type = OPERATOR
				return tok
			default:
				if unicode.IsLetter(ch) {
					token_type = IDENTIFIER
					token = append(token, ch)
				} else {
					panic("some special char encountered.")
				}
			}

		case STRING:
			switch ch {
			case '"':
				tok.name = string(token)
				return tok
			default:
				token = append(token, ch)
			}
		case NUMBER:
			switch ch {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
				token = append(token, ch)
			default:
				src.UnreadRune() // 回退一个字符
				tok.name = string(token)
				return tok
			}
		case IDENTIFIER:
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

var src *source = &source{bufio.NewReader(os.Stdin), nil} // source
var symbol_table map[string]*Value = make(map[string]*Value, 20)

func run() {
	for {
		token := src.next_token()
		if token.token_type == NOTYPE {
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
	if token.token_type == NOTYPE {
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
	if token.token_type != DELIMITER || token.name != ";" {
		// fmt.Println("statement: expect ; ")
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

func expression(level int) (value Value) {
	var target func(int) Value
	if level < 5 {
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
			sum := op_values(op.name, left.value, right.value)
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
			value = primary(level)
			value.value = op_values("-", value.value, nil)
		case "(":
			value = expression(1)
			token = src.next_token()
			if token.name != ")" {
				panic(") expected")
			}
		case "!":
		}
	}
	return value
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

		if op == "%" {
			ia := a.(int64)
			ib := b.(int64)
			return ia % ib
		}

		if op == "<<" {
			ia := uint64(a.(int64))
			ib := uint64(b.(int64))
			return ia << ib
		}

		if op == ">>" {
			ia := uint64(a.(int64))
			ib := uint64(b.(int64))
			return ia >> ib
		}

		if op == "&" {
			ia := a.(int64)
			ib := b.(int64)
			return ia & ib
		}

		if op == "^" {
			ia := a.(int64)
			ib := b.(int64)
			return ia ^ ib
		}

		if op == "|" {
			ia := a.(int64)
			ib := b.(int64)
			return ia | ib
		}

		if op == ">" {
			ia := a.(int64)
			ib := b.(int64)
			return ia > ib
		}

		if op == "<" {
			ia := a.(int64)
			ib := b.(int64)
			return ia < ib
		}

		if op == ">=" {
			ia := a.(int64)
			ib := b.(int64)
			return ia >= ib
		}

		if op == "<=" {
			ia := a.(int64)
			ib := b.(int64)
			return ia <= ib
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

		if op == ">" {
			ia := a.(float64)
			ib := b.(float64)
			return ia > ib
		}

		if op == ">=" {
			ia := a.(float64)
			ib := b.(float64)
			return ia >= ib
		}

		if op == "<" {
			ia := a.(float64)
			ib := b.(float64)
			return ia < ib
		}
		if op == "<=" {
			ia := a.(float64)
			ib := b.(float64)
			return ia <= ib
		}
	}

	if _, yes := a.(string); yes {
		if op == "+" {
			ia := a.(string)
			ib := b.(string)
			return ia + ib
		}
		if op == ">" {
			ia := a.(string)
			ib := b.(string)
			return ia > ib
		}
		if op == ">=" {
			ia := a.(string)
			ib := b.(string)
			return ia >= ib
		}
		if op == "<" {
			ia := a.(string)
			ib := b.(string)
			return ia < ib
		}
		if op == "<=" {
			ia := a.(string)
			ib := b.(string)
			return ia <= ib
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

	if op == "==" {
		return a == b
	}
	if op == "!=" {
		return a != b
	}
	return nil
}

func main() {
	run()
	/*for {
		token := src.next_token()
		fmt.Println(token)
		if token.token_type == NOTYPE || token.token_type == UNKNOWN {
			break
		}
	}*/
}
