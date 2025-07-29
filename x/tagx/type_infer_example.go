package tagx

import (
	"fmt"
	"reflect"
)

// ExampleTypeInference 类型推断示例
func ExampleTypeInference() {
	// 创建类型推断器
	inferrer := NewTypeInferrer()

	// 示例1：基础类型推断
	fmt.Println("=== 基础类型推断 ===")
	fmt.Printf("string -> %s\n", inferrer.InferType(reflect.TypeOf("")))
	fmt.Printf("int -> %s\n", inferrer.InferType(reflect.TypeOf(0)))
	fmt.Printf("float64 -> %s\n", inferrer.InferType(reflect.TypeOf(0.0)))
	fmt.Printf("bool -> %s\n", inferrer.InferType(reflect.TypeOf(false)))

	// 示例2：切片类型推断
	fmt.Println("\n=== 切片类型推断 ===")
	fmt.Printf("[]string -> %s\n", inferrer.InferType(reflect.TypeOf([]string{})))
	fmt.Printf("[]int -> %s\n", inferrer.InferType(reflect.TypeOf([]int{})))
	fmt.Printf("[]float64 -> %s\n", inferrer.InferType(reflect.TypeOf([]float64{})))
	fmt.Printf("[]bool -> %s\n", inferrer.InferType(reflect.TypeOf([]bool{})))
	fmt.Printf("[]struct -> %s\n", inferrer.InferType(reflect.TypeOf([]struct{}{})))

	// 示例3：指针类型推断
	fmt.Println("\n=== 指针类型推断 ===")
	fmt.Printf("*string -> %s\n", inferrer.InferType(reflect.TypeOf(new(string))))
	fmt.Printf("*int -> %s\n", inferrer.InferType(reflect.TypeOf(new(int))))
	fmt.Printf("*float64 -> %s\n", inferrer.InferType(reflect.TypeOf(new(float64))))

	// 示例4：复杂类型推断
	fmt.Println("\n=== 复杂类型推断 ===")
	fmt.Printf("struct{} -> %s\n", inferrer.InferType(reflect.TypeOf(struct{}{})))
	
	// 示例5：不支持的类型（会被忽略）
	fmt.Println("\n=== 不支持的类型（会被忽略） ===")
	fmt.Printf("map[string]interface{} -> %s\n", inferrer.InferType(reflect.TypeOf(map[string]interface{}{})))
	fmt.Printf("interface{} -> %s\n", inferrer.InferType(reflect.TypeOf((*interface{})(nil)).Elem()))
	// 注意：map、interface 等动态类型不支持，因为无法在程序启动时确定具体类型

	// 示例6：从值推断类型
	fmt.Println("\n=== 从值推断类型 ===")
	fmt.Printf("值 \"hello\" -> %s\n", inferrer.InferTypeFromValue("hello"))
	fmt.Printf("值 123 -> %s\n", inferrer.InferTypeFromValue(123))
	fmt.Printf("值 123.45 -> %s\n", inferrer.InferTypeFromValue(123.45))
	fmt.Printf("值 true -> %s\n", inferrer.InferTypeFromValue(true))
	fmt.Printf("值 []string{\"a\", \"b\"} -> %s\n", inferrer.InferTypeFromValue([]string{"a", "b"}))

	// 示例7：全局函数使用
	fmt.Println("\n=== 全局函数使用 ===")
	fmt.Printf("InferType(string) -> %s\n", InferType(reflect.TypeOf("")))
	fmt.Printf("InferTypeFromValue(123) -> %s\n", InferTypeFromValue(123))

	// 示例8：FieldConfig 类型推断
	fmt.Println("\n=== FieldConfig 类型推断 ===")
	
	// 有 data 标签的情况（现在会自动填充type）
	fieldWithData := &FieldConfig{
		FieldType: reflect.TypeOf(""),
		Data: &DataConfig{
			Type: "string", // 现在会自动推断，无需手动指定
		},
	}
	fmt.Printf("有data标签的字段 -> %s\n", GetTypeFromFieldConfig(fieldWithData))

	// 没有 data 标签的情况
	fieldWithoutData := &FieldConfig{
		FieldType: reflect.TypeOf(0),
		Data:      nil,
	}
	fmt.Printf("没有data标签的字段 -> %s\n", GetTypeFromFieldConfig(fieldWithoutData))
}

// ExampleTypeMapping 类型映射表示例
func ExampleTypeMapping() {
	fmt.Println("=== 类型映射表 ===")
	
	// 显示部分类型映射
	examples := []string{"string", "int", "float64", "bool", "[]string", "[]int", "[]float64", "[]bool", "[]struct"}
	
	for _, goType := range examples {
		if mappedType, exists := TypeMappingTable[goType]; exists {
			fmt.Printf("%s -> %s\n", goType, mappedType)
		}
	}
}

// ExampleConstants 类型常量示例
func ExampleConstants() {
	fmt.Println("=== 类型常量 ===")
	fmt.Printf("TypeString: %s\n", TypeString)
	fmt.Printf("TypeNumber: %s\n", TypeNumber)
	fmt.Printf("TypeFloat: %s\n", TypeFloat)
	fmt.Printf("TypeBoolean: %s\n", TypeBoolean)
	fmt.Printf("TypeStruct: %s\n", TypeStruct)
	fmt.Printf("TypeFiles: %s\n", TypeFiles)
	fmt.Printf("TypeStrings: %s\n", TypeStrings)
	fmt.Printf("TypeNumbers: %s\n", TypeNumbers)
	fmt.Printf("TypeFloats: %s\n", TypeFloats)
	fmt.Printf("TypeBooleans: %s\n", TypeBooleans)
	fmt.Printf("TypeStructs: %s\n", TypeStructs)
	
	// 添加时间类型说明
	fmt.Println("\n=== 时间类型说明 ===")
	fmt.Println("时间类型统一使用 int64 时间戳存储，前端通过 datetime 组件展示")
	fmt.Println("不支持 time.Time 类型，因为后端统一使用时间戳存储")
} 