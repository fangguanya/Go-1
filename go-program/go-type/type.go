package go_type

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"

	api_models "github.com/angenalZZZ/Go/go-program/api-models"
)

/**
命令行参数
*/
var temperature = FlagCelsius("t", 20.0, "the temperature")

// 类型检查
func TestTypeCheck() {
	var p api_models.IPoint = &api_models.Point{X: 1, Y: 2}
	var p2 = make([]api_models.Point, 2)
	fmt.Println("-------------------------\n类型检查：")

	// 命令行参数
	fmt.Printf("  命令行参数/摄氏温度: %s\n", temperature)

	// type assertion (*指针类型)
	if p0, ok := p.(*api_models.Point); ok {
		fmt.Printf("  类型断言: %p  %p\n", &p, p0)
	}
	// interface{} 接受任意类型的变量, 不同动态类型的变量不可比较, 只能与nil比较
	var w io.Writer // zeroValue=nil, 接受实现接口: Write(p []byte) 类型的变量, 下面的动态值决定了接收者类型(*T)的不同
	fmt.Printf("  接口w io.Writer(type)：%T, (value)：%[1]v \n", w)
	w = os.Stdout
	fmt.Printf("  接口w os.Stdout(type)：%T, (value)：%[1]v \n", w)
	w = new(bytes.Buffer)
	fmt.Println("  接口w new(bytes.Buffer)(type)：", reflect.TypeOf(w), ", (value)：", w) // %T: reflect.TypeOf(w)

	//var v1 bool
	//var v2 byte   // uint8  [true 或 false]
	//var v3 rune   // uint8, uint16, uint32 [unicode 编码: 1, 2, 4 个字节]
	//var v4 int    // 32位
	//var v40 uint  // 64位
	//var v5 int8   // -128~127
	//var v50 uint8 // 0 ~ 255
	//var v6 int16
	//var v60 uint16
	//var v7 int32
	//var v70 uint32
	//var v8 int64
	//var v80 uint64
	//var v9 uintptr // 存储指针的 uint32 或 uint64
	//var f1 float32 // 小数位数精确到  7 位
	//var f2 float64 // 小数位数精确到 15 位
	//var c1 complex64
	//var c2 complex128
	//var s1 string  // readonly byte slice
	//var s2 stringS

	fmt.Println(`  格式化p：%v %+v %T %#v make(Slice::Point)`)
	fmt.Printf("  格式化p：%v %+v %T %#v [%d]Point\n", p, p, p, p, cap(p2))
	fmt.Printf("  格式化i：%c %8.1f %8.2f %8x\n", 65, 12.5, 12.509, 54349)

	// 类型检查 指针
	PtrTypeCheck()

	// 二维数组
	TwoImensionalArrays(4, 2)

	// 斐波那契数列
	new(Fibonacci).Sequence(20, 2*time.Second, func(s []int) {
		fmt.Printf("  斐波那契数列: %v", s)
	})
}
