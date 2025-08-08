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

// 嵌套结构体
type NestedStruct struct {
	Value     int
	MyValue   MyInt
	Message   string
	MyMessage MyString
	Time      time.Time
	MyTime    MyTime
}

// 源结构体 - 包含嵌套结构体
type SourceStruct struct {
	Name      string
	Age       int
	MyAge     MyInt
	Address   NestedStruct
	MyAddress NestedStruct
	Tags      []string
	MyTags    []MyString
	Times     []time.Time
	MyTimes   []MyTime
}

// 目标结构体 - 包含嵌套结构体
type TargetStruct struct {
	Name      string
	Age       int
	MyAge     int // 从 MyInt 转换到 int
	Address   NestedStruct
	MyAddress NestedStruct // 嵌套结构体中的类型别名也应该转换
	Tags      []string
	MyTags    []string // 从 []MyString 转换到 []string
	Times     []time.Time
	MyTimes   []time.Time // 从 []MyTime 转换到 []time.Time
}

// 深度嵌套结构体
type DeepNestedStruct struct {
	Level1 struct {
		Level2 struct {
			Level3 struct {
				Value     int
				MyValue   MyInt
				Message   string
				MyMessage MyString
				Time      time.Time
				MyTime    MyTime
			}
		}
	}
}

type DeepNestedTarget struct {
	Level1 struct {
		Level2 struct {
			Level3 struct {
				Value     int
				MyValue   int // 从 MyInt 转换到 int
				Message   string
				MyMessage string // 从 MyString 转换到 string
				Time      time.Time
				MyTime    time.Time // 从 MyTime 转换到 time.Time
			}
		}
	}
}

func main() {
	fmt.Println("=== 嵌套结构体类型别名转换示例 ===\n")

	// 创建测试数据
	now := time.Now()
	myTime := MyTime(now)

	nested := NestedStruct{
		Value:     42,
		MyValue:   MyInt(42),
		Message:   "nested message",
		MyMessage: MyString("nested message"),
		Time:      now,
		MyTime:    myTime,
	}

	fmt.Printf("嵌套结构体数据:\n")
	fmt.Printf("  Value: %v\n", nested.Value)
	fmt.Printf("  MyValue: %v\n", nested.MyValue)
	fmt.Printf("  Message: %v\n", nested.Message)
	fmt.Printf("  MyMessage: %v\n", nested.MyMessage)
	fmt.Printf("  Time: %v\n", nested.Time)
	fmt.Printf("  MyTime: %v\n", nested.MyTime)
	fmt.Println()

	// 示例1: 基本嵌套结构体转换
	fmt.Println("1. 基本嵌套结构体转换:")
	src := SourceStruct{
		Name:      "John",
		Age:       30,
		MyAge:     MyInt(30),
		Address:   nested,
		MyAddress: nested,
		Tags:      []string{"tag1", "tag2"},
		MyTags:    []MyString{MyString("tag1"), MyString("tag2")},
		Times:     []time.Time{now, now.Add(time.Hour)},
		MyTimes:   []MyTime{myTime, MyTime(now.Add(time.Hour))},
	}

	var target TargetStruct
	mapster.MapTo(src, &target)

	fmt.Printf("   基本字段:\n")
	fmt.Printf("     Name: %v\n", target.Name)
	fmt.Printf("     Age: %v\n", target.Age)
	fmt.Printf("     MyAge (转换后): %v\n", target.MyAge)

	fmt.Printf("   嵌套结构体 Address:\n")
	fmt.Printf("     Value: %v\n", target.Address.Value)
	fmt.Printf("     MyValue: %v\n", target.Address.MyValue)
	fmt.Printf("     Message: %v\n", target.Address.Message)
	fmt.Printf("     MyMessage: %v\n", target.Address.MyMessage)
	fmt.Printf("     Time: %v\n", target.Address.Time)
	fmt.Printf("     MyTime: %v\n", time.Time(target.Address.MyTime))

	fmt.Printf("   嵌套结构体 MyAddress:\n")
	fmt.Printf("     Value: %v\n", target.MyAddress.Value)
	fmt.Printf("     MyValue: %v\n", target.MyAddress.MyValue)
	fmt.Printf("     Message: %v\n", target.MyAddress.Message)
	fmt.Printf("     MyMessage: %v\n", target.MyAddress.MyMessage)
	fmt.Printf("     Time: %v\n", target.MyAddress.Time)
	fmt.Printf("     MyTime: %v\n", time.Time(target.MyAddress.MyTime))

	fmt.Printf("   切片字段:\n")
	fmt.Printf("     Tags: %v\n", target.Tags)
	fmt.Printf("     MyTags (转换后): %v\n", target.MyTags)
	fmt.Printf("     Times: %v\n", target.Times)
	fmt.Printf("     MyTimes (转换后): %v\n", target.MyTimes)
	fmt.Println()

	// 示例2: 深度嵌套结构体转换
	fmt.Println("2. 深度嵌套结构体转换:")
	var deepSrc DeepNestedStruct
	deepSrc.Level1.Level2.Level3.Value = 42
	deepSrc.Level1.Level2.Level3.MyValue = MyInt(42)
	deepSrc.Level1.Level2.Level3.Message = "deep nested message"
	deepSrc.Level1.Level2.Level3.MyMessage = MyString("deep nested message")
	deepSrc.Level1.Level2.Level3.Time = now
	deepSrc.Level1.Level2.Level3.MyTime = myTime

	var deepTarget DeepNestedTarget
	mapster.MapTo(deepSrc, &deepTarget)

	fmt.Printf("   Level1.Level2.Level3:\n")
	fmt.Printf("     Value: %v\n", deepTarget.Level1.Level2.Level3.Value)
	fmt.Printf("     MyValue (转换后): %v\n", deepTarget.Level1.Level2.Level3.MyValue)
	fmt.Printf("     Message: %v\n", deepTarget.Level1.Level2.Level3.Message)
	fmt.Printf("     MyMessage (转换后): %v\n", deepTarget.Level1.Level2.Level3.MyMessage)
	fmt.Printf("     Time: %v\n", deepTarget.Level1.Level2.Level3.Time)
	fmt.Printf("     MyTime (转换后): %v\n", deepTarget.Level1.Level2.Level3.MyTime)
	fmt.Println()

	// 示例3: 验证转换的正确性
	fmt.Println("3. 验证转换的正确性:")

	// 验证基本字段
	if target.MyAge == 30 {
		fmt.Println("   ✓ MyAge -> int 转换正确")
	} else {
		fmt.Println("   ✗ MyAge -> int 转换错误")
	}

	// 验证嵌套结构体字段
	if target.Address.MyValue == 42 {
		fmt.Println("   ✓ Address.MyValue -> int 转换正确")
	} else {
		fmt.Println("   ✗ Address.MyValue -> int 转换错误")
	}

	if target.Address.MyMessage == "nested message" {
		fmt.Println("   ✓ Address.MyMessage -> string 转换正确")
	} else {
		fmt.Println("   ✗ Address.MyMessage -> string 转换错误")
	}

	if time.Time(target.Address.MyTime).Equal(now) {
		fmt.Println("   ✓ Address.MyTime -> time.Time 转换正确")
	} else {
		fmt.Println("   ✗ Address.MyTime -> time.Time 转换错误")
	}

	// 验证 MyAddress 字段
	if target.MyAddress.MyValue == 42 {
		fmt.Println("   ✓ MyAddress.MyValue -> int 转换正确")
	} else {
		fmt.Println("   ✗ MyAddress.MyValue -> int 转换错误")
	}

	if target.MyAddress.MyMessage == "nested message" {
		fmt.Println("   ✓ MyAddress.MyMessage -> string 转换正确")
	} else {
		fmt.Println("   ✗ MyAddress.MyMessage -> string 转换错误")
	}

	if time.Time(target.MyAddress.MyTime).Equal(now) {
		fmt.Println("   ✓ MyAddress.MyTime -> time.Time 转换正确")
	} else {
		fmt.Println("   ✗ MyAddress.MyTime -> time.Time 转换错误")
	}

	// 验证切片字段
	if len(target.MyTags) == 2 && target.MyTags[0] == "tag1" && target.MyTags[1] == "tag2" {
		fmt.Println("   ✓ MyTags -> []string 转换正确")
	} else {
		fmt.Println("   ✗ MyTags -> []string 转换错误")
	}

	if len(target.MyTimes) == 2 && target.MyTimes[0].Equal(now) && target.MyTimes[1].Equal(now.Add(time.Hour)) {
		fmt.Println("   ✓ MyTimes -> []time.Time 转换正确")
	} else {
		fmt.Println("   ✗ MyTimes -> []time.Time 转换错误")
	}

	// 验证深度嵌套字段
	if deepTarget.Level1.Level2.Level3.MyValue == 42 {
		fmt.Println("   ✓ Deep nested MyValue -> int 转换正确")
	} else {
		fmt.Println("   ✗ Deep nested MyValue -> int 转换错误")
	}

	if deepTarget.Level1.Level2.Level3.MyMessage == "deep nested message" {
		fmt.Println("   ✓ Deep nested MyMessage -> string 转换正确")
	} else {
		fmt.Println("   ✗ Deep nested MyMessage -> string 转换错误")
	}

	if deepTarget.Level1.Level2.Level3.MyTime.Equal(now) {
		fmt.Println("   ✓ Deep nested MyTime -> time.Time 转换正确")
	} else {
		fmt.Println("   ✗ Deep nested MyTime -> time.Time 转换错误")
	}

	fmt.Println("\n=== 嵌套结构体类型别名转换测试完成 ===")
}
