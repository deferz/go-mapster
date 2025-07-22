package main

import (
	"fmt"

	"github.com/deferz/go-mapster"
)

// 深度路径解析和循环引用处理示例

// 地址结构
type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}

// 公司结构
type Company struct {
	Name      string
	Address   Address
	Employees []*Employee // 员工列表，可能造成循环引用
	CEO       *Employee   // CEO，可能造成循环引用
}

// 员工结构
type Employee struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Company   *Company  // 可能造成循环引用
	Manager   *Employee // 可能造成自我引用
}

// 扁平化DTO（演示深度路径映射）
type EmployeeFlatDTO struct {
	ID               int64
	FirstName        string
	LastName         string
	Email            string
	CompanyName      string
	CompanyStreet    string
	CompanyCity      string
	CompanyState     string
	ManagerFirstName string
	ManagerLastName  string
}

// 安全DTO（避免循环引用）
type EmployeeSafeDTO struct {
	ID          int64
	FirstName   string
	LastName    string
	Email       string
	CompanyName string
	ManagerName string
}

func init() {
	// 配置深度路径映射
	mapster.Config[Employee, EmployeeFlatDTO]().
		Map("CompanyName").FromPath("Company.Name").
		Map("CompanyStreet").FromPath("Company.Address.Street").
		Map("CompanyCity").FromPath("Company.Address.City").
		Map("CompanyState").FromPath("Company.Address.State").
		Map("ManagerFirstName").FromPath("Manager.FirstName").
		Map("ManagerLastName").FromPath("Manager.LastName").
		Register()

	// 配置安全映射（避免循环引用）- 使用不同的映射避免冲突
	mapster.Config[Employee, EmployeeSafeDTO]().
		Map("CompanyName").FromFunc(func(e Employee) any {
		if e.Company != nil {
			return e.Company.Name
		}
		return ""
	}).
		Map("ManagerName").FromFunc(func(e Employee) any {
		if e.Manager != nil {
			return e.Manager.FirstName + " " + e.Manager.LastName
		}
		return ""
	}).
		Register()
}

func main() {
	fmt.Println("=== 深度路径解析和循环引用处理示例 ===")

	// 创建公司
	company := &Company{
		Name: "科技创新有限公司",
		Address: Address{
			Street:  "中关村软件园2号楼",
			City:    "北京",
			State:   "北京市",
			ZipCode: "100190",
		},
	}

	// 创建CEO
	ceo := &Employee{
		ID:        1,
		FirstName: "张",
		LastName:  "总",
		Email:     "ceo@company.com",
		Company:   company,
	}

	// 创建经理
	manager := &Employee{
		ID:        2,
		FirstName: "李",
		LastName:  "经理",
		Email:     "manager@company.com",
		Company:   company,
		Manager:   ceo, // 经理的上级是CEO
	}

	// 创建员工
	employee := &Employee{
		ID:        3,
		FirstName: "王",
		LastName:  "工程师",
		Email:     "engineer@company.com",
		Company:   company,
		Manager:   manager, // 员工的上级是经理
	}

	// 设置循环引用
	company.CEO = ceo
	company.Employees = []*Employee{ceo, manager, employee}
	ceo.Manager = ceo // 自我引用（CEO没有上级，但为了演示循环引用）

	fmt.Println("\n1. 深度路径映射示例:")

	// 测试深度路径映射
	flatDTO := mapster.Map[EmployeeFlatDTO](employee)
	fmt.Printf("员工: %s %s\n", flatDTO.FirstName, flatDTO.LastName)
	fmt.Printf("邮箱: %s\n", flatDTO.Email)
	fmt.Printf("公司名称: %s\n", flatDTO.CompanyName)
	fmt.Printf("公司地址: %s, %s, %s\n", flatDTO.CompanyStreet, flatDTO.CompanyCity, flatDTO.CompanyState)
	fmt.Printf("经理: %s %s\n", flatDTO.ManagerFirstName, flatDTO.ManagerLastName)

	fmt.Println("\n2. 安全映射示例（避免循环引用）:")

	// 测试安全映射（避免循环引用）
	safeDTO := mapster.Map[EmployeeSafeDTO](employee)
	fmt.Printf("员工: %s %s\n", safeDTO.FirstName, safeDTO.LastName)
	fmt.Printf("邮箱: %s\n", safeDTO.Email)
	fmt.Printf("公司名称: %s\n", safeDTO.CompanyName)
	fmt.Printf("经理: %s\n", safeDTO.ManagerName)

	fmt.Println("\n3. 批量安全映射:")

	// 批量映射测试
	safeDTOs := mapster.MapSlice[EmployeeSafeDTO](company.Employees)
	for i, dto := range safeDTOs {
		fmt.Printf("员工 %d: %s %s (%s)\n", i+1, dto.FirstName, dto.LastName, dto.CompanyName)
	}

	fmt.Println("\n=== 功能特性说明 ===")
	fmt.Println("✅ 深度路径解析:")
	fmt.Println("  • Company.Name - 访问嵌套对象属性")
	fmt.Println("  • Company.Address.City - 多层嵌套访问")
	fmt.Println("  • Manager.FirstName - 指针对象属性")
	fmt.Println("  • 支持 nil 安全检查")
	fmt.Println("  • 支持接口和 map 类型")
	fmt.Println()
	fmt.Println("✅ 循环引用处理:")
	fmt.Println("  • 自动检测指针循环引用")
	fmt.Println("  • 最大深度限制防止栈溢出")
	fmt.Println("  • 优雅处理自我引用")
	fmt.Println("  • 通过自定义函数避免循环")
	fmt.Println()
	fmt.Println("🔧 实现策略:")
	fmt.Println("  • 路径解析: 使用反射和字符串分割")
	fmt.Println("  • 循环检测: 指针地址追踪")
	fmt.Println("  • 安全映射: 自定义函数控制映射逻辑")
}
