package tests

import (
	"testing"
	"time"

	mapster "github.com/deferz/go-mapster"
)

// 多层嵌套的基础结构体
type BaseEntity struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuditInfo struct {
	BaseEntity // 第一层嵌套
	CreatedBy  string
	UpdatedBy  string
}

type MetaInfo struct {
	Version   int
	IsDeleted bool
}

// 多层嵌套的源结构体
type MultiLevelSource struct {
	AuditInfo // 第二层嵌套（包含BaseEntity）
	MetaInfo  // 第一层嵌套
	Name      string
	Value     float64
}

// 多层嵌套的目标结构体
type MultiLevelTarget struct {
	AuditInfo // 第二层嵌套（包含BaseEntity）
	MetaInfo  // 第一层嵌套
	Name      string
	Value     float64
}

// 扁平化的目标结构体
type FlattenedTarget struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string
	Version   int
	IsDeleted bool
	Name      string
	Value     float64
}

// TestMultiLevelEmbedding 测试多层匿名嵌套结构体映射
func TestMultiLevelEmbedding(t *testing.T) {
	now := time.Now()

	// 创建源对象
	src := MultiLevelSource{
		AuditInfo: AuditInfo{
			BaseEntity: BaseEntity{
				ID:        1001,
				CreatedAt: now,
				UpdatedAt: now.Add(time.Hour),
			},
			CreatedBy: "admin",
			UpdatedBy: "system",
		},
		MetaInfo: MetaInfo{
			Version:   2,
			IsDeleted: false,
		},
		Name:  "多层嵌套测试",
		Value: 123.45,
	}

	// 测试多层嵌套结构体映射
	t.Run("多层嵌套结构体映射", func(t *testing.T) {
		dst, err := mapster.Map[MultiLevelTarget](src)
		if err != nil {
			t.Fatalf("多层嵌套结构体映射失败: %v", err)
		}

		// 验证第三层嵌套字段 (BaseEntity)
		if dst.ID != src.ID {
			t.Errorf("期望 ID=%d, 得到 %d", src.ID, dst.ID)
		}
		if !dst.CreatedAt.Equal(src.CreatedAt) {
			t.Errorf("期望 CreatedAt=%v, 得到 %v", src.CreatedAt, dst.CreatedAt)
		}
		if !dst.UpdatedAt.Equal(src.UpdatedAt) {
			t.Errorf("期望 UpdatedAt=%v, 得到 %v", src.UpdatedAt, dst.UpdatedAt)
		}

		// 验证第二层嵌套字段 (AuditInfo)
		if dst.CreatedBy != src.CreatedBy {
			t.Errorf("期望 CreatedBy=%s, 得到 %s", src.CreatedBy, dst.CreatedBy)
		}
		if dst.UpdatedBy != src.UpdatedBy {
			t.Errorf("期望 UpdatedBy=%s, 得到 %s", src.UpdatedBy, dst.UpdatedBy)
		}

		// 验证第一层嵌套字段 (MetaInfo)
		if dst.Version != src.Version {
			t.Errorf("期望 Version=%d, 得到 %d", src.Version, dst.Version)
		}
		if dst.IsDeleted != src.IsDeleted {
			t.Errorf("期望 IsDeleted=%v, 得到 %v", src.IsDeleted, dst.IsDeleted)
		}

		// 验证普通字段
		if dst.Name != src.Name {
			t.Errorf("期望 Name=%s, 得到 %s", src.Name, dst.Name)
		}
		if dst.Value != src.Value {
			t.Errorf("期望 Value=%f, 得到 %f", src.Value, dst.Value)
		}
	})

	// 测试多层嵌套结构体到扁平化结构体的映射
	t.Run("多层嵌套结构体到扁平化结构体映射", func(t *testing.T) {
		dst, err := mapster.Map[FlattenedTarget](src)
		if err != nil {
			t.Fatalf("多层嵌套结构体到扁平化结构体映射失败: %v", err)
		}

		// 验证原来在第三层 (BaseEntity) 的字段
		if dst.ID != src.ID {
			t.Errorf("期望 ID=%d, 得到 %d", src.ID, dst.ID)
		}
		if !dst.CreatedAt.Equal(src.CreatedAt) {
			t.Errorf("期望 CreatedAt=%v, 得到 %v", src.CreatedAt, dst.CreatedAt)
		}
		if !dst.UpdatedAt.Equal(src.UpdatedAt) {
			t.Errorf("期望 UpdatedAt=%v, 得到 %v", src.UpdatedAt, dst.UpdatedAt)
		}

		// 验证原来在第二层 (AuditInfo) 的字段
		if dst.CreatedBy != src.CreatedBy {
			t.Errorf("期望 CreatedBy=%s, 得到 %s", src.CreatedBy, dst.CreatedBy)
		}
		if dst.UpdatedBy != src.UpdatedBy {
			t.Errorf("期望 UpdatedBy=%s, 得到 %s", src.UpdatedBy, dst.UpdatedBy)
		}

		// 验证原来在第一层 (MetaInfo) 的字段
		if dst.Version != src.Version {
			t.Errorf("期望 Version=%d, 得到 %d", src.Version, dst.Version)
		}
		if dst.IsDeleted != src.IsDeleted {
			t.Errorf("期望 IsDeleted=%v, 得到 %v", src.IsDeleted, dst.IsDeleted)
		}

		// 验证普通字段
		if dst.Name != src.Name {
			t.Errorf("期望 Name=%s, 得到 %s", src.Name, dst.Name)
		}
		if dst.Value != src.Value {
			t.Errorf("期望 Value=%f, 得到 %f", src.Value, dst.Value)
		}
	})
}

// 测试更复杂的多层嵌套结构体，包含多个嵌套路径
type ComplexBase struct {
	ID   int
	Code string
}

type ComplexAudit struct {
	ComplexBase // 嵌套
	Timestamp   time.Time
}

type ComplexMeta struct {
	Tags     []string
	Category string
}

type ComplexData struct {
	ComplexMeta // 嵌套
	Content     string
}

type ComplexSource struct {
	ComplexAudit // 嵌套路径1
	ComplexData  // 嵌套路径2
	Name         string
}

type ComplexTarget struct {
	ComplexAudit // 嵌套路径1
	ComplexData  // 嵌套路径2
	Name         string
}

// 完全扁平化的目标结构体
type ComplexFlattenedTarget struct {
	ID        int
	Code      string
	Timestamp time.Time
	Tags      []string
	Category  string
	Content   string
	Name      string
}

// TestComplexMultiLevelEmbedding 测试复杂的多层匿名嵌套结构体映射
func TestComplexMultiLevelEmbedding(t *testing.T) {
	now := time.Now()

	// 创建源对象
	src := ComplexSource{
		ComplexAudit: ComplexAudit{
			ComplexBase: ComplexBase{
				ID:   2001,
				Code: "COMPLEX-001",
			},
			Timestamp: now,
		},
		ComplexData: ComplexData{
			ComplexMeta: ComplexMeta{
				Tags:     []string{"test", "complex", "embedding"},
				Category: "测试类别",
			},
			Content: "这是一个复杂的多层嵌套测试内容",
		},
		Name: "复杂嵌套测试",
	}

	// 测试复杂多层嵌套结构体映射
	t.Run("复杂多层嵌套结构体映射", func(t *testing.T) {
		dst, err := mapster.Map[ComplexTarget](src)
		if err != nil {
			t.Fatalf("复杂多层嵌套结构体映射失败: %v", err)
		}

		// 验证嵌套路径1的字段
		if dst.ID != src.ID {
			t.Errorf("期望 ID=%d, 得到 %d", src.ID, dst.ID)
		}
		if dst.Code != src.Code {
			t.Errorf("期望 Code=%s, 得到 %s", src.Code, dst.Code)
		}
		if !dst.Timestamp.Equal(src.Timestamp) {
			t.Errorf("期望 Timestamp=%v, 得到 %v", src.Timestamp, dst.Timestamp)
		}

		// 验证嵌套路径2的字段
		if len(dst.Tags) != len(src.Tags) {
			t.Errorf("期望 Tags长度=%d, 得到 %d", len(src.Tags), len(dst.Tags))
		} else {
			for i, tag := range src.Tags {
				if dst.Tags[i] != tag {
					t.Errorf("期望 Tags[%d]=%s, 得到 %s", i, tag, dst.Tags[i])
				}
			}
		}
		if dst.Category != src.Category {
			t.Errorf("期望 Category=%s, 得到 %s", src.Category, dst.Category)
		}
		if dst.Content != src.Content {
			t.Errorf("期望 Content=%s, 得到 %s", src.Content, dst.Content)
		}

		// 验证普通字段
		if dst.Name != src.Name {
			t.Errorf("期望 Name=%s, 得到 %s", src.Name, dst.Name)
		}
	})

	// 测试复杂多层嵌套结构体到扁平化结构体的映射
	t.Run("复杂多层嵌套结构体到扁平化结构体映射", func(t *testing.T) {
		dst, err := mapster.Map[ComplexFlattenedTarget](src)
		if err != nil {
			t.Fatalf("复杂多层嵌套结构体到扁平化结构体映射失败: %v", err)
		}

		// 验证原来在嵌套路径1的字段
		if dst.ID != src.ID {
			t.Errorf("期望 ID=%d, 得到 %d", src.ID, dst.ID)
		}
		if dst.Code != src.Code {
			t.Errorf("期望 Code=%s, 得到 %s", src.Code, dst.Code)
		}
		if !dst.Timestamp.Equal(src.Timestamp) {
			t.Errorf("期望 Timestamp=%v, 得到 %v", src.Timestamp, dst.Timestamp)
		}

		// 验证原来在嵌套路径2的字段
		if len(dst.Tags) != len(src.Tags) {
			t.Errorf("期望 Tags长度=%d, 得到 %d", len(src.Tags), len(dst.Tags))
		} else {
			for i, tag := range src.Tags {
				if dst.Tags[i] != tag {
					t.Errorf("期望 Tags[%d]=%s, 得到 %s", i, tag, dst.Tags[i])
				}
			}
		}
		if dst.Category != src.Category {
			t.Errorf("期望 Category=%s, 得到 %s", src.Category, dst.Category)
		}
		if dst.Content != src.Content {
			t.Errorf("期望 Content=%s, 得到 %s", src.Content, dst.Content)
		}

		// 验证普通字段
		if dst.Name != src.Name {
			t.Errorf("期望 Name=%s, 得到 %s", src.Name, dst.Name)
		}
	})
}
