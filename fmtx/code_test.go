package fmtx

import (
	"testing"
)

func TestDeleteVar(t *testing.T) {
	tests := []struct {
		name     string
		variable string
		code     string
		expected string
	}{
		{
			name:     "删除var定义的RouterGroup",
			variable: "RouterGroup",
			code: `package main

var RouterGroup = "/api"
var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
			expected: `package main

var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
		},
		{
			name:     "删除:=定义的RouterGroup",
			variable: "RouterGroup",
			code: `package main

RouterGroup := "/api"
var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
			expected: `package main

var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
		},
		{
			name:     "删除const定义的RouterGroup",
			variable: "RouterGroup",
			code: `package main

const RouterGroup = "/api"
var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
			expected: `package main

var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
		},
		{
			name:     "删除换行定义的RouterGroup",
			variable: "RouterGroup",
			code: `package main

RouterGroup
    := "/api"
var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
			expected: `package main

var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
		},
		{
			name:     "删除带空格换行的RouterGroup",
			variable: "RouterGroup",
			code: `package main

RouterGroup 
    := "/api"
var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
			expected: `package main

var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
		},
		{
			name:     "删除var换行定义的RouterGroup",
			variable: "RouterGroup",
			code: `package main

var RouterGroup
    = "/api"
var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
			expected: `package main

var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
		},
		{
			name:     "删除const换行定义的RouterGroup",
			variable: "RouterGroup",
			code: `package main

const RouterGroup
    = "/api"
var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
			expected: `package main

var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
		},
		{
			name:     "删除多个RouterGroup定义",
			variable: "RouterGroup",
			code: `package main

var RouterGroup = "/api"
RouterGroup := "/api2"
const RouterGroup = "/api3"
var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
			expected: `package main

var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`,
		},
		{
			name:     "空变量名",
			variable: "",
			code:     "var RouterGroup = \"/api\"",
			expected: "var RouterGroup = \"/api\"",
		},
		{
			name:     "空代码",
			variable: "RouterGroup",
			code:     "",
			expected: "",
		},
		{
			name:     "不存在的变量",
			variable: "NonExistentVar",
			code:     "var RouterGroup = \"/api\"",
			expected: "var RouterGroup = \"/api\"",
		},
		{
			name:     "变量名部分匹配",
			variable: "Router",
			code:     "var RouterGroup = \"/api\"",
			expected: "var RouterGroup = \"/api\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeleteVar(tt.variable, tt.code)
			if result != tt.expected {
				t.Errorf("DeleteVar() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecodeFileName(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "从函数组EnName提取文件名",
			code: `package retail

var RetailPurchaseGroup = &runner.FunctionGroup{
	CnName: "采购管理",
	EnName: "retail_purchase",
}

func init() {
	runner.Get(RouterGroup+"/retail_purchase_list", RetailPurchaseList, RetailPurchaseListOption)
}`,
			expected: "retail_purchase.go",
		},
		{
			name: "从路由注册提取文件名（Post）",
			code: `package demo

func init() {
	// 注册Form函数：汇率计算
	runner.Post(RouterGroup+"/exchange_rate_calculate", ExchangeRateCalculate, ExchangeRateCalculateOption)
}`,
			expected: "exchange_rate_calculate.go",
		},
		{
			name: "从路由注册提取文件名（Get）",
			code: `package demo

func init() {
	// 注册Table函数：汇率配置管理
	runner.Get(RouterGroup+"/exchange_rate_config_list", ExchangeRateConfigList, ExchangeRateConfigListOption)
}`,
			expected: "exchange_rate_config_list.go",
		},
		{
			name: "从文件名注释提取",
			code: `// 文件名：user_management.go
package user

func someFunction() {
	fmt.Println("hello")
}`,
			expected: "user_management.go",
		},
		{
			name: "从文件注释提取（中文冒号）",
			code: `// 文件：order_system.go
package order

func someFunction() {
	fmt.Println("hello")
}`,
			expected: "order_system.go",
		},
		{
			name: "注释优先级高于函数组",
			code: `// 文件名：user_management.go
package retail

var RetailPurchaseGroup = &runner.FunctionGroup{
	CnName: "采购管理",
	EnName: "retail_purchase",
}

func init() {
	runner.Get(RouterGroup+"/some_other_route", SomeFunction, SomeOption)
}`,
			expected: "user_management.go",
		},
		{
			name: "函数组优先级高于路由注册",
			code: `package retail

var RetailPurchaseGroup = &runner.FunctionGroup{
	CnName: "采购管理",
	EnName: "retail_purchase",
}

func init() {
	runner.Get(RouterGroup+"/some_other_route", SomeFunction, SomeOption)
}`,
			expected: "retail_purchase.go",
		},
		{
			name: "路由注册优先级高于注释（无注释时）",
			code: `package demo

func init() {
	runner.Post(RouterGroup+"/new_function_name", NewFunction, NewOption)
}`,
			expected: "new_function_name.go",
		},
		{
			name: "多个路由注册取第一个",
			code: `package demo

func init() {
	runner.Post(RouterGroup+"/first_route", FirstFunction, FirstOption)
	runner.Get(RouterGroup+"/second_route", SecondFunction, SecondOption)
}`,
			expected: "first_route.go",
		},
		{
			name:     "空代码",
			code:     "",
			expected: "",
		},
		{
			name: "无法解析的代码",
			code: `package demo

func someFunction() {
	fmt.Println("hello")
}`,
			expected: "",
		},
		{
			name: "复杂函数组定义",
			code: `package complex

var ComplexSystemGroup = &runner.FunctionGroup{
	CnName: "复杂系统管理",
	EnName: "complex_system",
	Description: "这是一个复杂的系统",
	Version: "1.0.0",
}

func init() {
	runner.Get(RouterGroup+"/complex_list", ComplexList, ComplexListOption)
}`,
			expected: "complex_system.go",
		},
		{
			name: "带空格的函数组定义",
			code: `package retail

var RetailPurchaseGroup = &runner.FunctionGroup {
	CnName: "采购管理" ,
	EnName: "retail_purchase" ,
}

func init() {
	runner.Get(RouterGroup+"/retail_purchase_list", RetailPurchaseList, RetailPurchaseListOption)
}`,
			expected: "retail_purchase.go",
		},
		{
			name: "多行函数组定义",
			code: `package pdf

var PdfEncryptGroup = &runner.FunctionGroup{
    CnName: "PDF工具-加密解密",
    EnName: "pdf_encrypt",
}

func init() {
    runner.Post(RouterGroup+"/pdf_encrypt", PdfEncrypt, PdfEncryptOption)
}`,
			expected: "pdf_encrypt.go",
		},
		{
			name: "复杂多行函数组定义",
			code: `package image

var ImageConvertGroup = &runner.FunctionGroup{
    CnName: "图片格式转换工具",
    EnName: "image_convert",
    Description: "支持PNG、JPG、WebP格式互转",
}

func init() {
    runner.Post(RouterGroup+"/image_convert", ImageConvert, ImageConvertOption)
}`,
			expected: "image_convert.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecodeFileName(tt.code)
			if result != tt.expected {
				t.Errorf("DecodeFileName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// 基准测试
func BenchmarkDeleteVar(b *testing.B) {
	code := `package main

var RouterGroup = "/api"
var OtherVar = "test"

func main() {
    fmt.Println(RouterGroup)
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeleteVar("RouterGroup", code)
	}
}

func BenchmarkDecodeFileName(b *testing.B) {
	code := `package retail

var RetailPurchaseGroup = &runner.FunctionGroup{
	CnName: "采购管理",
	EnName: "retail_purchase",
}

func init() {
	runner.Get(RouterGroup+"/retail_purchase_list", RetailPurchaseList, RetailPurchaseListOption)
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DecodeFileName(code)
	}
}

func TestConvertToRouterGroup(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "转换Post路由（api路径）",
			code: `func init() {
	runner.Post("/api/healthcare/new_function_name", NewFunction, NewOption)
}`,
			expected: `func init() {
	runner.Post(RouterGroup+"/new_function_name", NewFunction, NewOption)
}`,
		},
		{
			name: "转换Get路由（api路径）",
			code: `func init() {
	runner.Get("/api/retail/medicine_list", MedicineList, MedicineListOption)
}`,
			expected: `func init() {
	runner.Get(RouterGroup+"/medicine_list", MedicineList, MedicineListOption)
}`,
		},
		{
			name: "转换任意路径（非api开头）",
			code: `func init() {
	runner.Post("/healthcare/new_function_name", NewFunction, NewOption)
}`,
			expected: `func init() {
	runner.Post(RouterGroup+"/new_function_name", NewFunction, NewOption)
}`,
		},
		{
			name: "转换复杂路径（多级目录）",
			code: `func init() {
	runner.Get("/v1/admin/user/management", UserManagement, UserManagementOption)
}`,
			expected: `func init() {
	runner.Get(RouterGroup+"/management", UserManagement, UserManagementOption)
}`,
		},
		{
			name: "转换简单路径（只有两级）",
			code: `func init() {
	runner.Post("/user/create", UserCreate, UserCreateOption)
}`,
			expected: `func init() {
	runner.Post(RouterGroup+"/create", UserCreate, UserCreateOption)
}`,
		},
		{
			name: "转换单路径段",
			code: `func init() {
	runner.Get("/single_path", SingleFunction, SingleOption)
}`,
			expected: `func init() {
	runner.Get(RouterGroup+"/single_path", SingleFunction, SingleOption)
}`,
		},
		{
			name: "转换多个路由",
			code: `func init() {
	runner.Post("/api/healthcare/prescription_create", PrescriptionCreate, PrescriptionCreateOption)
	runner.Get("/api/healthcare/medicine_list", MedicineList, MedicineListOption)
	runner.Post("/api/healthcare/registration_create", RegistrationCreate, RegistrationCreateOption)
}`,
			expected: `func init() {
	runner.Post(RouterGroup+"/prescription_create", PrescriptionCreate, PrescriptionCreateOption)
	runner.Get(RouterGroup+"/medicine_list", MedicineList, MedicineListOption)
	runner.Post(RouterGroup+"/registration_create", RegistrationCreate, RegistrationCreateOption)
}`,
		},
		{
			name: "转换单引号路由",
			code: `func init() {
	runner.Post('/api/healthcare/new_function_name', NewFunction, NewOption)
}`,
			expected: `func init() {
	runner.Post(RouterGroup+"/new_function_name", NewFunction, NewOption)
}`,
		},
		{
			name: "转换复杂路径",
			code: `func init() {
	runner.Post("/api/healthcare/very_long_function_name_with_underscores", LongFunction, LongOption)
}`,
			expected: `func init() {
	runner.Post(RouterGroup+"/very_long_function_name_with_underscores", LongFunction, LongOption)
}`,
		},
		{
			name: "混合已转换和未转换的路由",
			code: `func init() {
	runner.Post("/api/healthcare/new_function", NewFunction, NewOption)
	runner.Get(RouterGroup+"/existing_function", ExistingFunction, ExistingOption)
	runner.Post("/api/retail/another_function", AnotherFunction, AnotherOption)
}`,
			expected: `func init() {
	runner.Post(RouterGroup+"/new_function", NewFunction, NewOption)
	runner.Get(RouterGroup+"/existing_function", ExistingFunction, ExistingOption)
	runner.Post(RouterGroup+"/another_function", AnotherFunction, AnotherOption)
}`,
		},
		{
			name:     "空代码",
			code:     "",
			expected: "",
		},
		{
			name: "转换单路径段",
			code: `func init() {
	runner.Post(RouterGroup+"/already_converted", Function, Option)
	runner.Get("/single_path", OtherFunction, OtherOption)
}`,
			expected: `func init() {
	runner.Post(RouterGroup+"/already_converted", Function, Option)
	runner.Get(RouterGroup+"/single_path", OtherFunction, OtherOption)
}`,
		},
		{
			name: "带注释的代码",
			code: `func init() {
	// 药品管理 - Table函数
	runner.Get("/api/healthcare/medicine_list", MedicineList, MedicineListOption)
	
	// 医生管理 - Table函数
	runner.Get("/api/healthcare/doctor_manage", DoctorManage, DoctorManageOption)
}`,
			expected: `func init() {
	// 药品管理 - Table函数
	runner.Get(RouterGroup+"/medicine_list", MedicineList, MedicineListOption)
	
	// 医生管理 - Table函数
	runner.Get(RouterGroup+"/doctor_manage", DoctorManage, DoctorManageOption)
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToRouterGroup(tt.code)
			if result != tt.expected {
				t.Errorf("ConvertToRouterGroup() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func BenchmarkConvertToRouterGroup(b *testing.B) {
	code := `func init() {
	// 药品管理 - Table函数
	runner.Get("/api/healthcare/medicine_list", MedicineList, MedicineListOption)

	// 医生管理 - Table函数
	runner.Get("/api/healthcare/doctor_manage", DoctorManage, DoctorManageOption)

	// 挂号管理 - Table函数
	runner.Get("/api/healthcare/registration_list", RegistrationList, RegistrationListOption)

	// 科室管理 - Table函数
	runner.Get("/api/healthcare/department_list", DepartmentList, DepartmentListOption)

	// 处方管理 - Form函数（类似收银台）
	runner.Post("/api/healthcare/prescription_create", PrescriptionCreate, PrescriptionCreateOption)

	// 挂号创建 - Form函数
	runner.Post("/api/healthcare/registration_create", RegistrationCreate, RegistrationCreateOption)
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConvertToRouterGroup(code)
	}
}
