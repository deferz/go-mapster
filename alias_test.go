package mapster

import (
	"testing"
	"time"
)

// 定义各种类型的别名
type MyInt int
type MyString string
type MyFloat64 float64
type MyBool bool
type MyTime time.Time

// 测试结构体
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

// 反向转换测试结构体
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

// 切片测试结构体
type SliceSourceStruct struct {
	IntSlice      []int
	MyIntSlice    []MyInt
	StringSlice   []string
	MyStringSlice []MyString
	FloatSlice    []float64
	MyFloatSlice  []MyFloat64
	BoolSlice     []bool
	MyBoolSlice   []MyBool
	TimeSlice     []time.Time
	MyTimeSlice   []MyTime
}

type SliceTargetStruct struct {
	IntSlice      []int
	MyIntSlice    []int // 从 []MyInt 转换到 []int
	StringSlice   []string
	MyStringSlice []string // 从 []MyString 转换到 []string
	FloatSlice    []float64
	MyFloatSlice  []float64 // 从 []MyFloat64 转换到 []float64
	BoolSlice     []bool
	MyBoolSlice   []bool // 从 []MyBool 转换到 []bool
	TimeSlice     []time.Time
	MyTimeSlice   []time.Time // 从 []MyTime 转换到 []time.Time
}

func TestTypeAliasConversion(t *testing.T) {
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

	// 创建指针
	timePtr := &now
	myTimePtr := &myTime

	// 创建源结构体
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

	// 测试映射
	var target TargetStruct
	MapTo(src, &target)

	// 验证基本类型转换
	if target.IntField != testInt {
		t.Errorf("IntField conversion failed: expected %v, got %v", testInt, target.IntField)
	}

	if target.MyIntField != testInt {
		t.Errorf("MyIntField conversion failed: expected %v, got %v", testInt, target.MyIntField)
	}

	if target.StringField != testString {
		t.Errorf("StringField conversion failed: expected %v, got %v", testString, target.StringField)
	}

	if target.MyStringField != testString {
		t.Errorf("MyStringField conversion failed: expected %v, got %v", testString, target.MyStringField)
	}

	if target.FloatField != testFloat {
		t.Errorf("FloatField conversion failed: expected %v, got %v", testFloat, target.FloatField)
	}

	if target.MyFloatField != testFloat {
		t.Errorf("MyFloatField conversion failed: expected %v, got %v", testFloat, target.MyFloatField)
	}

	if target.BoolField != testBool {
		t.Errorf("BoolField conversion failed: expected %v, got %v", testBool, target.BoolField)
	}

	if target.MyBoolField != testBool {
		t.Errorf("MyBoolField conversion failed: expected %v, got %v", testBool, target.MyBoolField)
	}

	if !target.TimeField.Equal(now) {
		t.Errorf("TimeField conversion failed: expected %v, got %v", now, target.TimeField)
	}

	if !target.MyTimeField.Equal(now) {
		t.Errorf("MyTimeField conversion failed: expected %v, got %v", now, target.MyTimeField)
	}

	// 验证指针字段转换
	if target.TimePtr == nil || !target.TimePtr.Equal(now) {
		t.Errorf("TimePtr conversion failed: expected %v, got %v", now, target.TimePtr)
	}

	if target.MyTimePtr == nil || !target.MyTimePtr.Equal(now) {
		t.Errorf("MyTimePtr conversion failed: expected %v, got %v", now, target.MyTimePtr)
	}

	t.Logf("All type alias conversions passed successfully!")
}

func TestReverseTypeAliasConversion(t *testing.T) {
	// 创建测试数据
	testInt := 42
	testString := "hello world"
	testFloat := 3.14159
	testBool := true
	now := time.Now()
	timePtr := &now

	// 创建源结构体
	src := ReverseSourceStruct{
		MyIntField:    testInt,
		MyStringField: testString,
		MyFloatField:  testFloat,
		MyBoolField:   testBool,
		MyTimeField:   now,
		MyTimePtr:     timePtr,
	}

	// 测试映射
	var target ReverseTargetStruct
	MapTo(src, &target)

	// 验证反向转换
	if int(target.MyIntField) != testInt {
		t.Errorf("Reverse MyIntField conversion failed: expected %v, got %v", testInt, int(target.MyIntField))
	}

	if string(target.MyStringField) != testString {
		t.Errorf("Reverse MyStringField conversion failed: expected %v, got %v", testString, string(target.MyStringField))
	}

	if float64(target.MyFloatField) != testFloat {
		t.Errorf("Reverse MyFloatField conversion failed: expected %v, got %v", testFloat, float64(target.MyFloatField))
	}

	if bool(target.MyBoolField) != testBool {
		t.Errorf("Reverse MyBoolField conversion failed: expected %v, got %v", testBool, bool(target.MyBoolField))
	}

	if !time.Time(target.MyTimeField).Equal(now) {
		t.Errorf("Reverse MyTimeField conversion failed: expected %v, got %v", now, time.Time(target.MyTimeField))
	}

	if target.MyTimePtr == nil || !time.Time(*target.MyTimePtr).Equal(now) {
		t.Errorf("Reverse MyTimePtr conversion failed: expected %v, got %v", now, target.MyTimePtr)
	}

	t.Logf("All reverse type alias conversions passed successfully!")
}

func TestTypeAliasSliceConversion(t *testing.T) {
	// 创建测试数据
	ints := []int{1, 2, 3, 4, 5}
	myInts := make([]MyInt, len(ints))
	for i, v := range ints {
		myInts[i] = MyInt(v)
	}

	strings := []string{"hello", "world", "test"}
	myStrings := make([]MyString, len(strings))
	for i, v := range strings {
		myStrings[i] = MyString(v)
	}

	floats := []float64{1.1, 2.2, 3.3}
	myFloats := make([]MyFloat64, len(floats))
	for i, v := range floats {
		myFloats[i] = MyFloat64(v)
	}

	bools := []bool{true, false, true}
	myBools := make([]MyBool, len(bools))
	for i, v := range bools {
		myBools[i] = MyBool(v)
	}

	now := time.Now()
	times := []time.Time{now, now.Add(time.Hour), now.Add(2 * time.Hour)}
	myTimes := make([]MyTime, len(times))
	for i, v := range times {
		myTimes[i] = MyTime(v)
	}

	// 创建源结构体
	src := SliceSourceStruct{
		IntSlice:      ints,
		MyIntSlice:    myInts,
		StringSlice:   strings,
		MyStringSlice: myStrings,
		FloatSlice:    floats,
		MyFloatSlice:  myFloats,
		BoolSlice:     bools,
		MyBoolSlice:   myBools,
		TimeSlice:     times,
		MyTimeSlice:   myTimes,
	}

	// 测试映射
	var target SliceTargetStruct
	MapTo(src, &target)

	// 验证切片长度
	if len(target.IntSlice) != len(ints) {
		t.Errorf("IntSlice length mismatch: expected %d, got %d", len(ints), len(target.IntSlice))
	}

	if len(target.MyIntSlice) != len(ints) {
		t.Errorf("MyIntSlice length mismatch: expected %d, got %d", len(ints), len(target.MyIntSlice))
	}

	if len(target.StringSlice) != len(strings) {
		t.Errorf("StringSlice length mismatch: expected %d, got %d", len(strings), len(target.StringSlice))
	}

	if len(target.MyStringSlice) != len(strings) {
		t.Errorf("MyStringSlice length mismatch: expected %d, got %d", len(strings), len(target.MyStringSlice))
	}

	if len(target.FloatSlice) != len(floats) {
		t.Errorf("FloatSlice length mismatch: expected %d, got %d", len(floats), len(target.FloatSlice))
	}

	if len(target.MyFloatSlice) != len(floats) {
		t.Errorf("MyFloatSlice length mismatch: expected %d, got %d", len(floats), len(target.MyFloatSlice))
	}

	if len(target.BoolSlice) != len(bools) {
		t.Errorf("BoolSlice length mismatch: expected %d, got %d", len(bools), len(target.BoolSlice))
	}

	if len(target.MyBoolSlice) != len(bools) {
		t.Errorf("MyBoolSlice length mismatch: expected %d, got %d", len(bools), len(target.MyBoolSlice))
	}

	if len(target.TimeSlice) != len(times) {
		t.Errorf("TimeSlice length mismatch: expected %d, got %d", len(times), len(target.TimeSlice))
	}

	if len(target.MyTimeSlice) != len(times) {
		t.Errorf("MyTimeSlice length mismatch: expected %d, got %d", len(times), len(target.MyTimeSlice))
	}

	// 验证每个元素
	for i, expected := range ints {
		if target.IntSlice[i] != expected {
			t.Errorf("IntSlice[%d] conversion failed: expected %v, got %v", i, expected, target.IntSlice[i])
		}

		if target.MyIntSlice[i] != expected {
			t.Errorf("MyIntSlice[%d] conversion failed: expected %v, got %v", i, expected, target.MyIntSlice[i])
		}
	}

	for i, expected := range strings {
		if target.StringSlice[i] != expected {
			t.Errorf("StringSlice[%d] conversion failed: expected %v, got %v", i, expected, target.StringSlice[i])
		}

		if target.MyStringSlice[i] != expected {
			t.Errorf("MyStringSlice[%d] conversion failed: expected %v, got %v", i, expected, target.MyStringSlice[i])
		}
	}

	for i, expected := range floats {
		if target.FloatSlice[i] != expected {
			t.Errorf("FloatSlice[%d] conversion failed: expected %v, got %v", i, expected, target.FloatSlice[i])
		}

		if target.MyFloatSlice[i] != expected {
			t.Errorf("MyFloatSlice[%d] conversion failed: expected %v, got %v", i, expected, target.MyFloatSlice[i])
		}
	}

	for i, expected := range bools {
		if target.BoolSlice[i] != expected {
			t.Errorf("BoolSlice[%d] conversion failed: expected %v, got %v", i, expected, target.BoolSlice[i])
		}

		if target.MyBoolSlice[i] != expected {
			t.Errorf("MyBoolSlice[%d] conversion failed: expected %v, got %v", i, expected, target.MyBoolSlice[i])
		}
	}

	for i, expected := range times {
		if !target.TimeSlice[i].Equal(expected) {
			t.Errorf("TimeSlice[%d] conversion failed: expected %v, got %v", i, expected, target.TimeSlice[i])
		}

		if !target.MyTimeSlice[i].Equal(expected) {
			t.Errorf("MyTimeSlice[%d] conversion failed: expected %v, got %v", i, expected, target.MyTimeSlice[i])
		}
	}

	t.Logf("All type alias slice conversions passed successfully!")
}

// 测试空值和零值处理
func TestTypeAliasNilHandling(t *testing.T) {
	// 创建包含空指针的源结构体
	src := SourceStruct{
		IntField:      0,              // 零值
		MyIntField:    MyInt(0),       // 零值
		StringField:   "",             // 零值
		MyStringField: MyString(""),   // 零值
		FloatField:    0.0,            // 零值
		MyFloatField:  MyFloat64(0.0), // 零值
		BoolField:     false,          // 零值
		MyBoolField:   MyBool(false),  // 零值
		TimeField:     time.Time{},    // 零值
		MyTimeField:   MyTime{},       // 零值
		TimePtr:       nil,
		MyTimePtr:     nil,
	}

	// 测试映射
	var target TargetStruct
	MapTo(src, &target)

	// 验证零值处理
	if target.IntField != 0 {
		t.Errorf("IntField zero value handling failed: expected 0, got %v", target.IntField)
	}

	if target.MyIntField != 0 {
		t.Errorf("MyIntField zero value handling failed: expected 0, got %v", target.MyIntField)
	}

	if target.StringField != "" {
		t.Errorf("StringField zero value handling failed: expected empty string, got %v", target.StringField)
	}

	if target.MyStringField != "" {
		t.Errorf("MyStringField zero value handling failed: expected empty string, got %v", target.MyStringField)
	}

	if target.FloatField != 0.0 {
		t.Errorf("FloatField zero value handling failed: expected 0.0, got %v", target.FloatField)
	}

	if target.MyFloatField != 0.0 {
		t.Errorf("MyFloatField zero value handling failed: expected 0.0, got %v", target.MyFloatField)
	}

	if target.BoolField != false {
		t.Errorf("BoolField zero value handling failed: expected false, got %v", target.BoolField)
	}

	if target.MyBoolField != false {
		t.Errorf("MyBoolField zero value handling failed: expected false, got %v", target.MyBoolField)
	}

	if !target.TimeField.IsZero() {
		t.Errorf("TimeField zero value handling failed: expected zero time, got %v", target.TimeField)
	}

	if !target.MyTimeField.IsZero() {
		t.Errorf("MyTimeField zero value handling failed: expected zero time, got %v", target.MyTimeField)
	}

	// 验证空指针处理
	if target.TimePtr != nil {
		t.Errorf("TimePtr nil handling failed: expected nil, got %v", target.TimePtr)
	}

	if target.MyTimePtr != nil {
		t.Errorf("MyTimePtr nil handling failed: expected nil, got %v", target.MyTimePtr)
	}

	t.Logf("All type alias nil/zero value handling passed successfully!")
}

// 基准测试：性能测试
func BenchmarkTypeAliasConversion(b *testing.B) {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var target TargetStruct
		MapTo(src, &target)
	}
}
