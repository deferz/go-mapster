package mapper

import (
	"fmt"
	"reflect"
)

// circularRefTracker 用于跟踪循环引用
type circularRefTracker struct {
	visited map[uintptr]bool
}

// newCircularRefTracker 创建新的循环引用跟踪器
func newCircularRefTracker() *circularRefTracker {
	return &circularRefTracker{
		visited: make(map[uintptr]bool),
	}
}

// checkAndMark 检查并标记访问过的对象
// 如果对象已经被访问过，返回 true（表示存在循环引用）
func (c *circularRefTracker) checkAndMark(ptr uintptr) bool {
	if c.visited[ptr] {
		return true
	}
	c.visited[ptr] = true
	return false
}

// unmark 取消标记（用于退出递归时）
func (c *circularRefTracker) unmark(ptr uintptr) {
	delete(c.visited, ptr)
}

// MapValueWithCircularCheck 带循环引用检测的映射函数
func MapValueWithCircularCheck(src, dst reflect.Value, tracker *circularRefTracker) error {
	// 对于可能包含循环引用的类型（指针、切片、Map、接口）
	if src.Kind() == reflect.Ptr || src.Kind() == reflect.Slice ||
		src.Kind() == reflect.Map || src.Kind() == reflect.Interface {
		if src.CanAddr() {
			ptr := src.UnsafeAddr()
			if tracker.checkAndMark(ptr) {
				// 检测到循环引用，返回 nil（不继续映射）
				return nil
			}
			defer tracker.unmark(ptr)
		}
	}

	// 调用原始的 MapValue 函数
	return MapValue(src, dst)
}

// detectCircularReference 检测结构体字段中的循环引用
func detectCircularReference(v reflect.Value, path []reflect.Type) error {
	if !v.IsValid() || v.IsZero() {
		return nil
	}

	vType := v.Type()

	// 检查是否在路径中已经存在
	for _, t := range path {
		if t == vType {
			return fmt.Errorf("检测到循环引用: %v", vType)
		}
	}

	// 将当前类型添加到路径
	newPath := append(path, vType)

	// 根据类型递归检查
	switch v.Kind() {
	case reflect.Ptr:
		if !v.IsNil() {
			return detectCircularReference(v.Elem(), newPath)
		}
	case reflect.Struct:
		// 检查所有字段
		for i := 0; i < v.NumField(); i++ {
			if err := detectCircularReference(v.Field(i), newPath); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		// 检查元素类型
		if v.Len() > 0 {
			return detectCircularReference(v.Index(0), newPath)
		}
	case reflect.Map:
		// 检查键和值类型
		for _, key := range v.MapKeys() {
			if err := detectCircularReference(v.MapIndex(key), newPath); err != nil {
				return err
			}
		}
	}

	return nil
}
