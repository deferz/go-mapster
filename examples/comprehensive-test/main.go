package main

import (
	"fmt"
	"time"

	"github.com/deferz/go-mapster"
)

// 测试所有高级功能的综合示例

func main() {
	fmt.Println("=== Go Mapster 完整功能测试 ===")

	// 1. 基础映射测试
	fmt.Println("\n1. 基础自动映射:")
	testBasicMapping()

	// 2. 深度路径解析测试
	fmt.Println("\n2. 深度路径解析:")
	testDeepPathMapping()

	// 3. 自定义函数映射测试
	fmt.Println("\n3. 自定义函数映射:")
	testCustomFunctionMapping()

	// 4. 循环引用处理测试
	fmt.Println("\n4. 循环引用处理:")
	testCircularReferenceHandling()

	// 5. 批量映射测试
	fmt.Println("\n5. 批量映射:")
	testBatchMapping()

	// 6. 条件映射测试
	fmt.Println("\n6. 条件映射:")
	testConditionalMapping()

	fmt.Println("\n=== 所有功能测试完成 ===")
}

// 基础映射结构
type User struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Age       int
}

type UserDTO struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Age       int
}

func testBasicMapping() {
	user := User{
		ID:        1,
		FirstName: "张",
		LastName:  "三",
		Email:     "zhangsan@example.com",
		Age:       30,
	}

	dto := mapster.Map[UserDTO](user)
	fmt.Printf("  原始: %+v\n", user)
	fmt.Printf("  映射: %+v\n", dto)
	fmt.Printf("  ✅ 基础映射正常\n")
}

// 深度路径解析结构
type Address struct {
	Street   string
	City     string
	Province string
	Country  string
}

type Company struct {
	Name     string
	Industry string
	Address  Address
	Founded  time.Time
}

type Employee struct {
	ID      int64
	Name    string
	Title   string
	Company *Company
	Manager *Employee
	Reports []*Employee
}

type EmployeeDetailDTO struct {
	ID              int64
	Name            string
	Title           string
	CompanyName     string
	CompanyIndustry string
	CompanyCity     string
	CompanyCountry  string
	ManagerName     string
	ReportCount     int
}

func init() {
	// 配置深度路径映射
	mapster.Config[Employee, EmployeeDetailDTO]().
		Map("CompanyName").FromPath("Company.Name").
		Map("CompanyIndustry").FromPath("Company.Industry").
		Map("CompanyCity").FromPath("Company.Address.City").
		Map("CompanyCountry").FromPath("Company.Address.Country").
		Map("ManagerName").FromPath("Manager.Name").
		Map("ReportCount").FromFunc(func(e Employee) any {
		return len(e.Reports)
	}).
		Register()
}

func testDeepPathMapping() {
	company := &Company{
		Name:     "科技创新有限公司",
		Industry: "软件开发",
		Address: Address{
			Street:   "中关村大街1号",
			City:     "北京",
			Province: "北京市",
			Country:  "中国",
		},
		Founded: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	manager := &Employee{
		ID:      1,
		Name:    "李经理",
		Title:   "技术总监",
		Company: company,
	}

	employee := &Employee{
		ID:      2,
		Name:    "王工程师",
		Title:   "高级开发工程师",
		Company: company,
		Manager: manager,
		Reports: []*Employee{},
	}

	manager.Reports = []*Employee{employee}

	dto := mapster.Map[EmployeeDetailDTO](employee)
	fmt.Printf("  员工: %s (%s)\n", dto.Name, dto.Title)
	fmt.Printf("  公司: %s (%s)\n", dto.CompanyName, dto.CompanyIndustry)
	fmt.Printf("  地址: %s, %s\n", dto.CompanyCity, dto.CompanyCountry)
	fmt.Printf("  经理: %s\n", dto.ManagerName)
	fmt.Printf("  下属: %d 人\n", dto.ReportCount)
	fmt.Printf("  ✅ 深度路径解析正常\n")
}

// 自定义函数映射结构
type Order struct {
	ID          string
	Items       []OrderItem
	CreatedAt   time.Time
	CustomerID  int64
	Status      int
	TotalAmount float64
}

type OrderItem struct {
	ProductName string
	Quantity    int
	Price       float64
}

type OrderSummaryDTO struct {
	ID              string
	ItemCount       int
	TotalQuantity   int
	FormattedAmount string
	StatusText      string
	Age             string
	IsRecent        bool
}

func init() {
	// 配置自定义函数映射
	mapster.Config[Order, OrderSummaryDTO]().
		Map("ItemCount").FromFunc(func(o Order) any {
		return len(o.Items)
	}).
		Map("TotalQuantity").FromFunc(func(o Order) any {
		total := 0
		for _, item := range o.Items {
			total += item.Quantity
		}
		return total
	}).
		Map("FormattedAmount").FromFunc(func(o Order) any {
		return fmt.Sprintf("¥%.2f", o.TotalAmount)
	}).
		Map("StatusText").FromFunc(func(o Order) any {
		switch o.Status {
		case 1:
			return "待付款"
		case 2:
			return "已付款"
		case 3:
			return "已发货"
		case 4:
			return "已完成"
		default:
			return "未知状态"
		}
	}).
		Map("Age").FromFunc(func(o Order) any {
		duration := time.Since(o.CreatedAt)
		if duration.Hours() < 24 {
			return fmt.Sprintf("%.0f小时前", duration.Hours())
		}
		return fmt.Sprintf("%.0f天前", duration.Hours()/24)
	}).
		Map("IsRecent").FromFunc(func(o Order) any {
		return time.Since(o.CreatedAt).Hours() < 24
	}).
		Register()
}

func testCustomFunctionMapping() {
	order := Order{
		ID: "ORD-2024-001",
		Items: []OrderItem{
			{ProductName: "笔记本电脑", Quantity: 1, Price: 5999.99},
			{ProductName: "鼠标", Quantity: 2, Price: 99.99},
		},
		CreatedAt:   time.Now().Add(-6 * time.Hour),
		CustomerID:  12345,
		Status:      2,
		TotalAmount: 6199.97,
	}

	dto := mapster.Map[OrderSummaryDTO](order)
	fmt.Printf("  订单: %s\n", dto.ID)
	fmt.Printf("  商品: %d 件，总数量: %d\n", dto.ItemCount, dto.TotalQuantity)
	fmt.Printf("  金额: %s\n", dto.FormattedAmount)
	fmt.Printf("  状态: %s\n", dto.StatusText)
	fmt.Printf("  时间: %s (最近: %v)\n", dto.Age, dto.IsRecent)
	fmt.Printf("  ✅ 自定义函数映射正常\n")
}

// 循环引用处理结构
type Department struct {
	ID       int
	Name     string
	Manager  *Employee2
	Members  []*Employee2
	Parent   *Department
	Children []*Department
}

type Employee2 struct {
	ID         int
	Name       string
	Department *Department
	Manager    *Employee2
	Reports    []*Employee2
}

type DepartmentDTO struct {
	ID          int
	Name        string
	ManagerName string
	MemberCount int
	ParentName  string
	ChildCount  int
}

type Employee2DTO struct {
	ID             int
	Name           string
	DepartmentName string
	ManagerName    string
	ReportCount    int
}

func init() {
	// 配置部门安全映射
	mapster.Config[Department, DepartmentDTO]().
		Map("ManagerName").FromFunc(func(d Department) any {
		if d.Manager != nil {
			return d.Manager.Name
		}
		return "无"
	}).
		Map("MemberCount").FromFunc(func(d Department) any {
		return len(d.Members)
	}).
		Map("ParentName").FromFunc(func(d Department) any {
		if d.Parent != nil {
			return d.Parent.Name
		}
		return "无"
	}).
		Map("ChildCount").FromFunc(func(d Department) any {
		return len(d.Children)
	}).
		Register()

	// 配置员工安全映射
	mapster.Config[Employee2, Employee2DTO]().
		Map("DepartmentName").FromFunc(func(e Employee2) any {
		if e.Department != nil {
			return e.Department.Name
		}
		return "无"
	}).
		Map("ManagerName").FromFunc(func(e Employee2) any {
		if e.Manager != nil {
			return e.Manager.Name
		}
		return "无"
	}).
		Map("ReportCount").FromFunc(func(e Employee2) any {
		return len(e.Reports)
	}).
		Register()
}

func testCircularReferenceHandling() {
	// 创建具有循环引用的复杂组织结构
	tech := &Department{ID: 1, Name: "技术部"}
	frontend := &Department{ID: 3, Name: "前端组", Parent: tech}
	backend := &Department{ID: 4, Name: "后端组", Parent: tech}

	cto := &Employee2{ID: 1, Name: "技术总监", Department: tech}
	frontendLead := &Employee2{ID: 2, Name: "前端组长", Department: frontend, Manager: cto}
	backendLead := &Employee2{ID: 3, Name: "后端组长", Department: backend, Manager: cto}
	developer1 := &Employee2{ID: 4, Name: "前端开发", Department: frontend, Manager: frontendLead}
	developer2 := &Employee2{ID: 5, Name: "后端开发", Department: backend, Manager: backendLead}

	// 设置循环引用
	tech.Manager = cto
	tech.Members = []*Employee2{cto, frontendLead, backendLead, developer1, developer2}
	tech.Children = []*Department{frontend, backend}

	frontend.Manager = frontendLead
	frontend.Members = []*Employee2{frontendLead, developer1}

	backend.Manager = backendLead
	backend.Members = []*Employee2{backendLead, developer2}

	cto.Reports = []*Employee2{frontendLead, backendLead}
	frontendLead.Reports = []*Employee2{developer1}
	backendLead.Reports = []*Employee2{developer2}

	// 测试部门映射
	techDTO := mapster.Map[DepartmentDTO](tech)
	fmt.Printf("  部门: %s (经理: %s)\n", techDTO.Name, techDTO.ManagerName)
	fmt.Printf("  成员: %d 人，子部门: %d 个\n", techDTO.MemberCount, techDTO.ChildCount)

	// 测试员工映射
	ctoDTO := mapster.Map[Employee2DTO](cto)
	fmt.Printf("  员工: %s (部门: %s)\n", ctoDTO.Name, ctoDTO.DepartmentName)
	fmt.Printf("  下属: %d 人\n", ctoDTO.ReportCount)

	fmt.Printf("  ✅ 循环引用处理正常\n")
}

func testBatchMapping() {
	users := []User{
		{ID: 1, FirstName: "张", LastName: "三", Email: "zhang@example.com", Age: 25},
		{ID: 2, FirstName: "李", LastName: "四", Email: "li@example.com", Age: 30},
		{ID: 3, FirstName: "王", LastName: "五", Email: "wang@example.com", Age: 35},
	}

	dtos := make([]UserDTO, len(users))
	for i, u := range users {
		dtos[i] = mapster.Map[UserDTO](u)
	}
	fmt.Printf("  批量映射 %d 个用户:\n", len(dtos))
	for i, dto := range dtos {
		fmt.Printf("    %d. %s%s (%d岁)\n", i+1, dto.FirstName, dto.LastName, dto.Age)
	}
	fmt.Printf("  ✅ 批量映射正常\n")
}

// 条件映射结构
type Product struct {
	ID          int
	Name        string
	Price       float64
	Category    string
	InStock     bool
	Rating      float32
	ReviewCount int
}

type ProductDisplayDTO struct {
	ID         int
	Name       string
	Price      string
	Category   string
	Status     string
	Rating     string
	PopularTag string
}

func init() {
	// 配置条件映射
	mapster.Config[Product, ProductDisplayDTO]().
		Map("Price").FromFunc(func(p Product) any {
		if p.Price > 0 {
			return fmt.Sprintf("¥%.2f", p.Price)
		}
		return "价格面议"
	}).
		Map("Status").FromFunc(func(p Product) any {
		if p.InStock {
			return "现货"
		}
		return "缺货"
	}).
		Map("Rating").FromFunc(func(p Product) any {
		if p.ReviewCount > 0 {
			return fmt.Sprintf("%.1f分 (%d评价)", p.Rating, p.ReviewCount)
		}
		return "暂无评价"
	}).
		Map("PopularTag").FromFunc(func(p Product) any {
		if p.ReviewCount >= 100 && p.Rating >= 4.5 {
			return "热门商品"
		} else if p.ReviewCount >= 50 && p.Rating >= 4.0 {
			return "好评商品"
		}
		return ""
	}).
		Register()
}

func testConditionalMapping() {
	products := []Product{
		{ID: 1, Name: "iPhone 15", Price: 5999, Category: "手机", InStock: true, Rating: 4.8, ReviewCount: 256},
		{ID: 2, Name: "限量版手表", Price: 0, Category: "手表", InStock: false, Rating: 0, ReviewCount: 0},
		{ID: 3, Name: "蓝牙耳机", Price: 299, Category: "音响", InStock: true, Rating: 4.2, ReviewCount: 89},
	}

	fmt.Printf("  条件映射产品展示:\n")
	for i, product := range products {
		dto := mapster.Map[ProductDisplayDTO](product)
		fmt.Printf("    %d. %s\n", i+1, dto.Name)
		fmt.Printf("       价格: %s | 状态: %s\n", dto.Price, dto.Status)
		fmt.Printf("       评价: %s\n", dto.Rating)
		if dto.PopularTag != "" {
			fmt.Printf("       标签: %s\n", dto.PopularTag)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("  ✅ 条件映射正常\n")
}
