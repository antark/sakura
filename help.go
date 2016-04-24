package sakura

import (
	"fmt"
)

func help() {
	fmt.Println("支持数据类型：整数、浮点数、bool类型，和字符串的：")
	fmt.Println("四则运算：+  -  *  /  %")
	fmt.Println("按位运算：<<  >>  &  |  ^")
	fmt.Println("关系比较：==  !=  >  <  >=  <=")
	fmt.Println("逻辑运算：!  &&  ||")
	fmt.Println("字符串拼接：+")
	fmt.Println("字符串数值转换：+")
	fmt.Println()

	fmt.Println("优先级：一元运算符 + - ^ ! 最高，其次 * / 一级，其次 + - 一级，再者 == != 一级，然后是 && 和 || ")
	fmt.Println()

	fmt.Println("examples as:")
	fmt.Println("1+2;  3*4;  (1+2)*(3-4)+10-20;  1<<10;  1234&4567;")
	fmt.Println("true && false || true; ")
	fmt.Println("\"hello\"+\",\"+\"world\";")
	fmt.Println("help")
	fmt.Println("quit")
}
