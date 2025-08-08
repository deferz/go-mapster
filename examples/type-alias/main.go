package main

import (
	"fmt"
	"time"

	"github.com/deferz/go-mapster"
)

// 定义各种类型的别名
type MyInt int
type MyString string
type MyFloat64 float64
type MyBool bool
type MyTime time.Time

// 源结构体
type SourceStruct struct {
	IntField      int
	MyIntField    MyInt
	StringField   string
	MyStringField MyString
	FloatField    float64
	MyFloatField  MyFloat64
	BoolField     bool
	MyBoolField   MyBool
	TimeField     time.Time
	MyTimeField   MyTime
	TimePtr       *time.Time
	MyTimePtr     *MyTime
}

// 目标结构体
type TargetStruct struct {
	IntField      int
	MyIntField    int // 从 MyInt 转换到 int
	StringField   string
	MyStringField string // 从 MyString 转换到 string
	FloatField    float64
	MyFloatField  float64 // 从 MyFloat64 转换到 float64
	BoolField     bool
	MyBoolField   bool // 从 MyBool 转换到 bool
	TimeField     time.Time
	MyTimeField   time.Time // 从 MyTime 转换到 time.Time
	TimePtr       *time.Time
	MyTimePtr     *time.Time // 从 *MyTime 转换到 *time.Time
}

// 反向转换示例
type ReverseSourceStruct struct {
	MyIntField    int        // 源字段名与目标字段名匹配
	MyStringField string     // 源字段名与目标字段名匹配
	MyFloatField  float64    // 源字段名与目标字段名匹配
	MyBoolField   bool       // 源字段名与目标字段名匹配
	MyTimeField   time.Time  // 源字段名与目标字段名匹配
	MyTimePtr     *time.Time // 源字段名与目标字段名匹配
}

type ReverseTargetStruct struct {
	MyIntField    MyInt     // 从 int 转换到 MyInt
	MyStringField MyString  // 从 string 转换到 MyString
	MyFloatField  MyFloat64 // 从 float64 转换到 MyFloat64
	MyBoolField   MyBool    // 从 bool 转换到 MyBool
	MyTimeField   MyTime    // 从 time.Time 转换到 MyTime
	MyTimePtr     *MyTime   // 从 *time.Time 转换到 *MyTime
}

func main() {
	fmt.Println("=== 类型别名转换示例 ===\n")

	// 创建测试数据
	testInt := 42
	myInt := MyInt(42)
	testString := "hello world"
	myString := MyString("hello world")
	testFloat := 3.14159
	myFloat := MyFloat64(3.14159)
	testBool := true
	myBool := MyBool(true)
	now := time.Now()
	myTime := MyTime(now)
	timePtr := &now
	myTimePtr := &myTime

	fmt.Printf("原始数据:\n")
	fmt.Printf("  int: %v\n", testInt)
	fmt.Printf("  MyInt: %v\n", myInt)
	fmt.Printf("  string: %v\n", testString)
	fmt.Printf("  MyString: %v\n", myString)
	fmt.Printf("  float64: %v\n", testFloat)
	fmt.Printf("  MyFloat64: %v\n", myFloat)
	fmt.Printf("  bool: %v\n", testBool)
	fmt.Printf("  MyBool: %v\n", myBool)
	fmt.Printf("  time.Time: %v\n", now)
	fmt.Printf("  MyTime: %v\n", myTime)
	fmt.Println()

	// 示例1: 从类型别名转换到基本类型
	fmt.Println("1. 从类型别名转换到基本类型:")
	src := SourceStruct{
		IntField:      testInt,
		MyIntField:    myInt,
		StringField:   testString,
		MyStringField: myString,
		FloatField:    testFloat,
		MyFloatField:  myFloat,
		BoolField:     testBool,
		MyBoolField:   myBool,
		TimeField:     now,
		MyTimeField:   myTime,
		TimePtr:       timePtr,
		MyTimePtr:     myTimePtr,
	}

	var target TargetStruct
	mapster.MapTo(src, &target)

	fmt.Printf("   IntField: %v\n", target.IntField)
	fmt.Printf("   MyIntField (转换后): %v\n", target.MyIntField)
	fmt.Printf("   StringField: %v\n", target.StringField)
	fmt.Printf("   MyStringField (转换后): %v\n", target.MyStringField)
	fmt.Printf("   FloatField: %v\n", target.FloatField)
	fmt.Printf("   MyFloatField (转换后): %v\n", target.MyFloatField)
	fmt.Printf("   BoolField: %v\n", target.BoolField)
	fmt.Printf("   MyBoolField (转换后): %v\n", target.MyBoolField)
	fmt.Printf("   TimeField: %v\n", target.TimeField)
	fmt.Printf("   MyTimeField (转换后): %v\n", target.MyTimeField)
	fmt.Printf("   TimePtr: %v\n", target.TimePtr)
	fmt.Printf("   MyTimePtr (转换后): %v\n", target.MyTimePtr)
	fmt.Println()

	// 示例2: 从基本类型转换到类型别名
	fmt.Println("2. 从基本类型转换到类型别名:")
	reverseSrc := ReverseSourceStruct{
		MyIntField:    testInt,
		MyStringField: testString,
		MyFloatField:  testFloat,
		MyBoolField:   testBool,
		MyTimeField:   now,
		MyTimePtr:     timePtr,
	}

	var reverseTarget ReverseTargetStruct
	mapster.MapTo(reverseSrc, &reverseTarget)

	fmt.Printf("   MyIntField (转换后): %v\n", reverseTarget.MyIntField)
	fmt.Printf("   MyStringField (转换后): %v\n", reverseTarget.MyStringField)
	fmt.Printf("   MyFloatField (转换后): %v\n", reverseTarget.MyFloatField)
	fmt.Printf("   MyBoolField (转换后): %v\n", reverseTarget.MyBoolField)
	fmt.Printf("   MyTimeField (转换后): %v\n", reverseTarget.MyTimeField)
	fmt.Printf("   MyTimePtr (转换后): %v\n", reverseTarget.MyTimePtr)
	fmt.Println()

	// 示例3: 验证转换的正确性
	fmt.Println("3. 验证转换的正确性:")

	// 验证基本类型转换
	if target.MyIntField == testInt {
		fmt.Println("   ✓ MyInt -> int 转换正确")
	} else {
		fmt.Println("   ✗ MyInt -> int 转换错误")
	}

	if target.MyStringField == testString {
		fmt.Println("   ✓ MyString -> string 转换正确")
	} else {
		fmt.Println("   ✗ MyString -> string 转换错误")
	}

	if target.MyFloatField == testFloat {
		fmt.Println("   ✓ MyFloat64 -> float64 转换正确")
	} else {
		fmt.Println("   ✗ MyFloat64 -> float64 转换错误")
	}

	if target.MyBoolField == testBool {
		fmt.Println("   ✓ MyBool -> bool 转换正确")
	} else {
		fmt.Println("   ✗ MyBool -> bool 转换错误")
	}

	if target.MyTimeField.Equal(now) {
		fmt.Println("   ✓ MyTime -> time.Time 转换正确")
	} else {
		fmt.Println("   ✗ MyTime -> time.Time 转换错误")
	}

	if target.MyTimePtr != nil && target.MyTimePtr.Equal(now) {
		fmt.Println("   ✓ *MyTime -> *time.Time 转换正确")
	} else {
		fmt.Println("   ✗ *MyTime -> *time.Time 转换错误")
	}

	// 验证反向转换
	if int(reverseTarget.MyIntField) == testInt {
		fmt.Println("   ✓ int -> MyInt 转换正确")
	} else {
		fmt.Println("   ✗ int -> MyInt 转换错误")
	}

	if string(reverseTarget.MyStringField) == testString {
		fmt.Println("   ✓ string -> MyString 转换正确")
	} else {
		fmt.Println("   ✗ string -> MyString 转换错误")
	}

	if float64(reverseTarget.MyFloatField) == testFloat {
		fmt.Println("   ✓ float64 -> MyFloat64 转换正确")
	} else {
		fmt.Println("   ✗ float64 -> MyFloat64 转换错误")
	}

	if bool(reverseTarget.MyBoolField) == testBool {
		fmt.Println("   ✓ bool -> MyBool 转换正确")
	} else {
		fmt.Println("   ✗ bool -> MyBool 转换错误")
	}

	if time.Time(reverseTarget.MyTimeField).Equal(now) {
		fmt.Println("   ✓ time.Time -> MyTime 转换正确")
	} else {
		fmt.Println("   ✗ time.Time -> MyTime 转换错误")
	}

	if reverseTarget.MyTimePtr != nil && time.Time(*reverseTarget.MyTimePtr).Equal(now) {
		fmt.Println("   ✓ *time.Time -> *MyTime 转换正确")
	} else {
		fmt.Println("   ✗ *time.Time -> *MyTime 转换错误")
	}

	fmt.Println("\n=== 类型别名转换测试完成 ===")
}
