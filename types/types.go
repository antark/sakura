// types
package types

import (
	"strconv"
	"strings"
)

// scanner 中 token 类型
const (
	NOTYPE     = iota // 无类型
	DELIMITER         // 分隔符
	NUMBER            // 整数
	BOOLEAN           // 布尔类型
	CHARS             // 字符串
	OPERATOR          // 操作符
	IDENTIFIER        // 标识符
	UNKNOWN           // 未知类型
)

// parser 中 value 类型
const (
	INT_64 = iota
	FLOAT_64
	BOOL
	STRING
)

func Op_values(op string, a interface{}, b interface{}) (value interface{}) {
	if _, yes := a.(int64); yes {
		if b == nil {
			switch op {
			case "+":
				return a
			case "-":
				return -a.(int64)
			case "^":
				return ^a.(int64)
			default:
				return nil
			}
		}
		var ia, ib int64
		ia = a.(int64)
		ib = b.(int64)
		switch op {
		case "+":
			return ia + ib
		case "-":
			return ia - ib
		case "*":
			return ia * ib
		case "/":
			return ia * ib
		case "%":
			return ia / ib
		case "<<":
			return int64(uint64(ia) << uint64(ib))
		case ">>":
			return int64(uint64(ia) >> uint64(ib))
		case "&":
			return ia & ib
		case "|":
			return ia | ib
		case "^":
			return ia ^ ib
		case ">":
			return ia > ib
		case "<":
			return ia < ib
		case ">=":
			return ia >= ib
		case "<=":
			return ia <= ib
		}
	}

	if _, yes := a.(float64); yes {
		if b == nil {
			switch op {
			case "+":
				return a
			case "-":
				return -a.(float64)
			default:
				return nil
			}
		}
		var ia, ib float64
		ia = a.(float64)
		ib = b.(float64)
		switch op {
		case "+":
			return ia + ib
		case "-":
			return ia - ib
		case "*":
			return ia * ib
		case "/":
			return ia * ib
		case ">":
			return ia > ib
		case "<":
			return ia < ib
		case ">=":
			return ia >= ib
		case "<=":
			return ia <= ib
		}
	}

	if _, yes := a.(string); yes {
		if b == nil {
			switch op {
			case "+":
				if strings.Contains(a.(string), ".") {
					value, _ := strconv.ParseFloat(a.(string), 64)
					return value
				} else {
					value, _ := strconv.ParseInt(a.(string), 10, 64)
					return value
				}
			default:
				return nil
			}
		}
		var ia, ib string
		ia = a.(string)
		ib = b.(string)
		switch op {
		case "+":
			return ia + ib
		case ">":
			return ia > ib
		case "<":
			return ia < ib
		case ">=":
			return ia >= ib
		case "<=":
			return ia <= ib
		}
	}

	if _, yes := a.(bool); yes {
		if b == nil {
			switch op {
			case "!":
				return !a.(bool)
			default:
				return nil
			}
		}
		var ia, ib bool
		ia = a.(bool)
		ib = b.(bool)
		switch op {
		case "&&":
			return ia && ib
		case "||":
			return ia || ib
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
