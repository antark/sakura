// types
package types

func Op_values(op string, a interface{}, b interface{}) (value interface{}) {
	if _, yes := a.(int64); yes {
		if b == nil {
			return -a.(int64)
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
			return uint64(ia) << uint64(ib)
		case ">>":
			return uint64(ia) >> uint64(ib)
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
			return -a.(float64)
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
