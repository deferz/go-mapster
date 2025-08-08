package mapster

import (
	"testing"
	"time"
)

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
type NestedSourceStruct struct {
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
type NestedTargetStruct struct {
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

func TestNestedTypeAliasConversion(t *testing.T) {
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

	// 创建源结构体
	src := NestedSourceStruct{
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

	// 测试映射
	var target NestedTargetStruct
	MapTo(src, &target)

	// 验证基本字段转换
	if target.Name != "John" {
		t.Errorf("Name conversion failed: expected John, got %v", target.Name)
	}

	if target.Age != 30 {
		t.Errorf("Age conversion failed: expected 30, got %v", target.Age)
	}

	if target.MyAge != 30 {
		t.Errorf("MyAge conversion failed: expected 30, got %v", target.MyAge)
	}

	// 验证嵌套结构体字段转换
	if target.Address.Value != 42 {
		t.Errorf("Address.Value conversion failed: expected 42, got %v", target.Address.Value)
	}

	if target.Address.MyValue != 42 {
		t.Errorf("Address.MyValue conversion failed: expected 42, got %v", target.Address.MyValue)
	}

	if target.Address.Message != "nested message" {
		t.Errorf("Address.Message conversion failed: expected 'nested message', got %v", target.Address.Message)
	}

	if target.Address.MyMessage != "nested message" {
		t.Errorf("Address.MyMessage conversion failed: expected 'nested message', got %v", target.Address.MyMessage)
	}

	if !target.Address.Time.Equal(now) {
		t.Errorf("Address.Time conversion failed: expected %v, got %v", now, target.Address.Time)
	}

	if !time.Time(target.Address.MyTime).Equal(now) {
		t.Errorf("Address.MyTime conversion failed: expected %v, got %v", now, time.Time(target.Address.MyTime))
	}

	// 验证 MyAddress 字段转换（嵌套结构体中的类型别名）
	if target.MyAddress.Value != 42 {
		t.Errorf("MyAddress.Value conversion failed: expected 42, got %v", target.MyAddress.Value)
	}

	if target.MyAddress.MyValue != 42 {
		t.Errorf("MyAddress.MyValue conversion failed: expected 42, got %v", target.MyAddress.MyValue)
	}

	if target.MyAddress.Message != "nested message" {
		t.Errorf("MyAddress.Message conversion failed: expected 'nested message', got %v", target.MyAddress.Message)
	}

	if target.MyAddress.MyMessage != "nested message" {
		t.Errorf("MyAddress.MyMessage conversion failed: expected 'nested message', got %v", target.MyAddress.MyMessage)
	}

	if !target.MyAddress.Time.Equal(now) {
		t.Errorf("MyAddress.Time conversion failed: expected %v, got %v", now, target.MyAddress.Time)
	}

	if !time.Time(target.MyAddress.MyTime).Equal(now) {
		t.Errorf("MyAddress.MyTime conversion failed: expected %v, got %v", now, time.Time(target.MyAddress.MyTime))
	}

	// 验证切片字段转换
	if len(target.Tags) != 2 {
		t.Errorf("Tags length mismatch: expected 2, got %d", len(target.Tags))
	}

	if len(target.MyTags) != 2 {
		t.Errorf("MyTags length mismatch: expected 2, got %d", len(target.MyTags))
	}

	if target.Tags[0] != "tag1" || target.Tags[1] != "tag2" {
		t.Errorf("Tags conversion failed: expected [tag1 tag2], got %v", target.Tags)
	}

	if target.MyTags[0] != "tag1" || target.MyTags[1] != "tag2" {
		t.Errorf("MyTags conversion failed: expected [tag1 tag2], got %v", target.MyTags)
	}

	if len(target.Times) != 2 {
		t.Errorf("Times length mismatch: expected 2, got %d", len(target.Times))
	}

	if len(target.MyTimes) != 2 {
		t.Errorf("MyTimes length mismatch: expected 2, got %d", len(target.MyTimes))
	}

	if !target.Times[0].Equal(now) || !target.Times[1].Equal(now.Add(time.Hour)) {
		t.Errorf("Times conversion failed")
	}

	if !target.MyTimes[0].Equal(now) || !target.MyTimes[1].Equal(now.Add(time.Hour)) {
		t.Errorf("MyTimes conversion failed")
	}

	t.Logf("All nested type alias conversions passed successfully!")
}

func TestDeepNestedTypeAliasConversion(t *testing.T) {
	// 创建测试数据
	now := time.Now()
	myTime := MyTime(now)

	// 创建深度嵌套的源结构体
	var src DeepNestedStruct
	src.Level1.Level2.Level3.Value = 42
	src.Level1.Level2.Level3.MyValue = MyInt(42)
	src.Level1.Level2.Level3.Message = "deep nested message"
	src.Level1.Level2.Level3.MyMessage = MyString("deep nested message")
	src.Level1.Level2.Level3.Time = now
	src.Level1.Level2.Level3.MyTime = myTime

	// 测试映射
	var target DeepNestedTarget
	MapTo(src, &target)

	// 验证深度嵌套结构体中的类型别名转换
	if target.Level1.Level2.Level3.Value != 42 {
		t.Errorf("Deep nested Value conversion failed: expected 42, got %v", target.Level1.Level2.Level3.Value)
	}

	if target.Level1.Level2.Level3.MyValue != 42 {
		t.Errorf("Deep nested MyValue conversion failed: expected 42, got %v", target.Level1.Level2.Level3.MyValue)
	}

	if target.Level1.Level2.Level3.Message != "deep nested message" {
		t.Errorf("Deep nested Message conversion failed: expected 'deep nested message', got %v", target.Level1.Level2.Level3.Message)
	}

	if target.Level1.Level2.Level3.MyMessage != "deep nested message" {
		t.Errorf("Deep nested MyMessage conversion failed: expected 'deep nested message', got %v", target.Level1.Level2.Level3.MyMessage)
	}

	if !target.Level1.Level2.Level3.Time.Equal(now) {
		t.Errorf("Deep nested Time conversion failed: expected %v, got %v", now, target.Level1.Level2.Level3.Time)
	}

	if !target.Level1.Level2.Level3.MyTime.Equal(now) {
		t.Errorf("Deep nested MyTime conversion failed: expected %v, got %v", now, target.Level1.Level2.Level3.MyTime)
	}

	t.Logf("All deep nested type alias conversions passed successfully!")
}

// 测试包含指针的嵌套结构体
type PointerNestedStruct struct {
	Value     *int
	MyValue   *MyInt
	Message   *string
	MyMessage *MyString
	Time      *time.Time
	MyTime    *MyTime
}

type PointerNestedSource struct {
	Data   PointerNestedStruct
	MyData PointerNestedStruct
}

type PointerNestedTarget struct {
	Data   PointerNestedStruct
	MyData PointerNestedStruct
}

func TestPointerNestedTypeAliasConversion(t *testing.T) {
	// 创建测试数据
	value := 42
	myValue := MyInt(42)
	message := "pointer message"
	myMessage := MyString("pointer message")
	now := time.Now()
	myTime := MyTime(now)

	// 创建源结构体
	src := PointerNestedSource{
		Data: PointerNestedStruct{
			Value:     &value,
			MyValue:   &myValue,
			Message:   &message,
			MyMessage: &myMessage,
			Time:      &now,
			MyTime:    &myTime,
		},
		MyData: PointerNestedStruct{
			Value:     &value,
			MyValue:   &myValue,
			Message:   &message,
			MyMessage: &myMessage,
			Time:      &now,
			MyTime:    &myTime,
		},
	}

	// 测试映射
	var target PointerNestedTarget
	MapTo(src, &target)

	// 验证指针嵌套结构体中的类型别名转换
	if *target.Data.Value != 42 {
		t.Errorf("Data.Value conversion failed: expected 42, got %v", *target.Data.Value)
	}

	if *target.Data.MyValue != 42 {
		t.Errorf("Data.MyValue conversion failed: expected 42, got %v", *target.Data.MyValue)
	}

	if *target.Data.Message != "pointer message" {
		t.Errorf("Data.Message conversion failed: expected 'pointer message', got %v", *target.Data.Message)
	}

	if *target.Data.MyMessage != "pointer message" {
		t.Errorf("Data.MyMessage conversion failed: expected 'pointer message', got %v", *target.Data.MyMessage)
	}

	if !target.Data.Time.Equal(now) {
		t.Errorf("Data.Time conversion failed: expected %v, got %v", now, target.Data.Time)
	}

	if !time.Time(*target.Data.MyTime).Equal(now) {
		t.Errorf("Data.MyTime conversion failed: expected %v, got %v", now, time.Time(*target.Data.MyTime))
	}

	// 验证 MyData 字段转换
	if *target.MyData.Value != 42 {
		t.Errorf("MyData.Value conversion failed: expected 42, got %v", *target.MyData.Value)
	}

	if *target.MyData.MyValue != 42 {
		t.Errorf("MyData.MyValue conversion failed: expected 42, got %v", *target.MyData.MyValue)
	}

	if *target.MyData.Message != "pointer message" {
		t.Errorf("MyData.Message conversion failed: expected 'pointer message', got %v", *target.MyData.Message)
	}

	if *target.MyData.MyMessage != "pointer message" {
		t.Errorf("MyData.MyMessage conversion failed: expected 'pointer message', got %v", *target.MyData.MyMessage)
	}

	if !target.MyData.Time.Equal(now) {
		t.Errorf("MyData.Time conversion failed: expected %v, got %v", now, target.MyData.Time)
	}

	if !time.Time(*target.MyData.MyTime).Equal(now) {
		t.Errorf("MyData.MyTime conversion failed: expected %v, got %v", now, time.Time(*target.MyData.MyTime))
	}

	t.Logf("All pointer nested type alias conversions passed successfully!")
}

// 基准测试：嵌套结构体性能测试
func BenchmarkNestedTypeAliasConversion(b *testing.B) {
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

	src := NestedSourceStruct{
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var target NestedTargetStruct
		MapTo(src, &target)
	}
}
