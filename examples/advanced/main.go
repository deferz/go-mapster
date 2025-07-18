package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/deferz/go-mapster"
)

// Domain models
type Customer struct {
	ID          int64
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	DateOfBirth time.Time
	IsActive    bool
}

type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
	Country string
}

type Order struct {
	ID              int64
	CustomerID      int64
	Customer        Customer
	Items           []OrderItem
	TotalAmount     float64
	CreatedAt       time.Time
	Status          string
	ShippingAddress Address
}

type OrderItem struct {
	ID       int64
	OrderID  int64
	Name     string
	Quantity int
	Price    float64
}

// DTOs
type CustomerDTO struct {
	ID       int64
	FullName string
	Email    string
	Phone    string
	Age      int
	Status   string
}

type OrderDTO struct {
	ID           int64
	CustomerName string
	ItemCount    int
	TotalAmount  string
	CreatedDate  string
	Status       string
	ShippingInfo string
}

type OrderSummaryDTO struct {
	OrderNumber    string
	CustomerInfo   string
	Summary        string
	FormattedTotal string
}

func init() {
	// Configure Customer to CustomerDTO mapping
	mapster.Config[Customer, CustomerDTO]().
		Map("FullName").FromFunc(func(c Customer) interface{} {
		return c.FirstName + " " + c.LastName
	}).
		Map("Phone").FromField("PhoneNumber").
		Map("Age").FromFunc(func(c Customer) interface{} {
		return int(time.Since(c.DateOfBirth).Hours() / 24 / 365)
	}).
		Map("Status").FromFunc(func(c Customer) interface{} {
		if c.IsActive {
			return "Active"
		}
		return "Inactive"
	}).
		Register()

	// Configure Order to OrderDTO mapping
	mapster.Config[Order, OrderDTO]().
		Map("CustomerName").FromFunc(func(o Order) interface{} {
		return o.Customer.FirstName + " " + o.Customer.LastName
	}).
		Map("ItemCount").FromFunc(func(o Order) interface{} {
		return len(o.Items)
	}).
		Map("TotalAmount").FromFunc(func(o Order) interface{} {
		return fmt.Sprintf("$%.2f", o.TotalAmount)
	}).
		Map("CreatedDate").FromFunc(func(o Order) interface{} {
		return o.CreatedAt.Format("2006-01-02")
	}).
		Map("ShippingInfo").FromFunc(func(o Order) interface{} {
		addr := o.ShippingAddress
		return fmt.Sprintf("%s, %s, %s %s", addr.Street, addr.City, addr.State, addr.ZipCode)
	}).
		Register()

	// Configure Order to OrderSummaryDTO mapping with more complex logic
	mapster.Config[Order, OrderSummaryDTO]().
		Map("OrderNumber").FromFunc(func(o Order) interface{} {
		return fmt.Sprintf("ORD-%06d", o.ID)
	}).
		Map("CustomerInfo").FromFunc(func(o Order) interface{} {
		return fmt.Sprintf("%s (%s)",
			o.Customer.FirstName+" "+o.Customer.LastName,
			o.Customer.Email)
	}).
		Map("Summary").FromFunc(func(o Order) interface{} {
		itemNames := make([]string, len(o.Items))
		for i, item := range o.Items {
			itemNames[i] = fmt.Sprintf("%s x%d", item.Name, item.Quantity)
		}
		return strings.Join(itemNames, ", ")
	}).
		Map("FormattedTotal").FromFunc(func(o Order) interface{} {
		return fmt.Sprintf("Total: $%.2f (%s)", o.TotalAmount, strings.ToUpper(o.Status))
	}).
		Register()
}

func main() {
	// Sample data
	customer := Customer{
		ID:          1,
		FirstName:   "John",
		LastName:    "Smith",
		Email:       "john.smith@email.com",
		PhoneNumber: "+1-555-0123",
		DateOfBirth: time.Date(1985, 3, 15, 0, 0, 0, 0, time.UTC),
		IsActive:    true,
	}

	address := Address{
		Street:  "123 Main St",
		City:    "Seattle",
		State:   "WA",
		ZipCode: "98101",
		Country: "USA",
	}

	order := Order{
		ID:         12345,
		CustomerID: customer.ID,
		Customer:   customer,
		Items: []OrderItem{
			{ID: 1, Name: "Laptop", Quantity: 1, Price: 1299.99},
			{ID: 2, Name: "Mouse", Quantity: 2, Price: 29.99},
			{ID: 3, Name: "Keyboard", Quantity: 1, Price: 79.99},
		},
		TotalAmount:     1439.96,
		CreatedAt:       time.Now(),
		Status:          "pending",
		ShippingAddress: address,
	}

	fmt.Println("=== Advanced Mapster Example ===")

	// Customer mapping
	customerDTO := mapster.Map[CustomerDTO](customer)
	fmt.Println("Customer DTO:")
	fmt.Printf("  ID: %d\n", customerDTO.ID)
	fmt.Printf("  Full Name: %s\n", customerDTO.FullName)
	fmt.Printf("  Email: %s\n", customerDTO.Email)
	fmt.Printf("  Phone: %s\n", customerDTO.Phone)
	fmt.Printf("  Age: %d\n", customerDTO.Age)
	fmt.Printf("  Status: %s\n", customerDTO.Status)
	fmt.Println()

	// Order mapping
	orderDTO := mapster.Map[OrderDTO](order)
	fmt.Println("Order DTO:")
	fmt.Printf("  ID: %d\n", orderDTO.ID)
	fmt.Printf("  Customer: %s\n", orderDTO.CustomerName)
	fmt.Printf("  Items: %d\n", orderDTO.ItemCount)
	fmt.Printf("  Total: %s\n", orderDTO.TotalAmount)
	fmt.Printf("  Date: %s\n", orderDTO.CreatedDate)
	fmt.Printf("  Status: %s\n", orderDTO.Status)
	fmt.Printf("  Shipping: %s\n", orderDTO.ShippingInfo)
	fmt.Println()

	// Order summary mapping
	summaryDTO := mapster.Map[OrderSummaryDTO](order)
	fmt.Println("Order Summary DTO:")
	fmt.Printf("  Order Number: %s\n", summaryDTO.OrderNumber)
	fmt.Printf("  Customer: %s\n", summaryDTO.CustomerInfo)
	fmt.Printf("  Items: %s\n", summaryDTO.Summary)
	fmt.Printf("  %s\n", summaryDTO.FormattedTotal)
	fmt.Println()

	// Batch mapping
	customers := []Customer{customer, {
		ID:          2,
		FirstName:   "Jane",
		LastName:    "Doe",
		Email:       "jane.doe@email.com",
		PhoneNumber: "+1-555-0456",
		DateOfBirth: time.Date(1990, 7, 22, 0, 0, 0, 0, time.UTC),
		IsActive:    false,
	}}

	customerDTOs := mapster.MapSlice[CustomerDTO](customers)
	fmt.Println("Batch Customer Mapping:")
	for i, dto := range customerDTOs {
		fmt.Printf("  %d. %s (%s) - %s\n", i+1, dto.FullName, dto.Email, dto.Status)
	}
	fmt.Println()

	fmt.Println("=== Performance Info ===")
	fmt.Println("This example demonstrates:")
	fmt.Println("• Complex field mapping with custom functions")
	fmt.Println("• Date/time processing")
	fmt.Println("• String formatting and transformation")
	fmt.Println("• Nested object field access")
	fmt.Println("• Conditional logic in mapping")
	fmt.Println("• Batch processing of slices")
}
