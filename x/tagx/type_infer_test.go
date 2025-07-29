package tagx

import (
	"reflect"
	"testing"
)

// TestTypeInferrer 测试类型推断器
func TestTypeInferrer(t *testing.T) {
	inferrer := NewTypeInferrer()

	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		// 基础类型测试
		{"string", "hello", "string"},
		{"int", 123, "number"},
		{"int8", int8(123), "number"},
		{"int16", int16(123), "number"},
		{"int32", int32(123), "number"},
		{"int64", int64(123), "number"},
		{"uint", uint(123), "number"},
		{"uint8", uint8(123), "number"},
		{"uint16", uint16(123), "number"},
		{"uint32", uint32(123), "number"},
		{"uint64", uint64(123), "number"},
		{"float32", float32(123.45), "float"},
		{"float64", float64(123.45), "float"},
		{"bool", true, "boolean"},

		// 切片类型测试
		{"[]string", []string{"a", "b"}, "[]string"},
		{"[]int", []int{1, 2, 3}, "[]number"},
		{"[]int8", []int8{1, 2, 3}, "[]number"},
		{"[]int16", []int16{1, 2, 3}, "[]number"},
		{"[]int32", []int32{1, 2, 3}, "[]number"},
		{"[]int64", []int64{1, 2, 3}, "[]number"},
		{"[]uint", []uint{1, 2, 3}, "[]number"},
		{"[]uint8", []uint8{1, 2, 3}, "[]number"},
		{"[]uint16", []uint16{1, 2, 3}, "[]number"},
		{"[]uint32", []uint32{1, 2, 3}, "[]number"},
		{"[]uint64", []uint64{1, 2, 3}, "[]number"},
		{"[]float32", []float32{1.1, 2.2}, "[]float"},
		{"[]float64", []float64{1.1, 2.2}, "[]float"},
		{"[]bool", []bool{true, false}, "[]boolean"},
		{"[]struct", []struct{}{}, "[]struct"},

		// 指针类型测试
		{"*string", new(string), "string"},
		{"*int", new(int), "number"},
		{"*float64", new(float64), "float"},
		{"*bool", new(bool), "boolean"},

		// 复杂类型测试
		{"struct", struct{}{}, "struct"},
		// 注意：map、interface 等动态类型不支持，因为无法在程序启动时确定具体类型

		// 特殊类型测试
		{"nil", nil, "string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inferrer.InferTypeFromValue(tt.value)
			if result != tt.expected {
				t.Errorf("InferTypeFromValue(%v) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

// TestInferType 测试全局类型推断函数
func TestInferType(t *testing.T) {
	tests := []struct {
		name     string
		typ      reflect.Type
		expected string
	}{
		{"string", reflect.TypeOf(""), "string"},
		{"int", reflect.TypeOf(0), "number"},
		{"float64", reflect.TypeOf(0.0), "float"},
		{"bool", reflect.TypeOf(false), "boolean"},
		{"[]string", reflect.TypeOf([]string{}), "[]string"},
		{"[]int", reflect.TypeOf([]int{}), "[]number"},
		{"struct", reflect.TypeOf(struct{}{}), "struct"},
		{"*string", reflect.TypeOf(new(string)), "string"},
		// 注意：map、interface 等动态类型不支持，因为无法在程序启动时确定具体类型
		{"map", reflect.TypeOf(map[string]interface{}{}), ""},
		{"interface", reflect.TypeOf((*interface{})(nil)).Elem(), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InferType(tt.typ)
			if result != tt.expected {
				t.Errorf("InferType(%v) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

// TestGetTypeFromFieldConfig 测试从FieldConfig获取类型
func TestGetTypeFromFieldConfig(t *testing.T) {
	// 测试有data标签的情况
	fieldWithData := &FieldConfig{
		FieldType: reflect.TypeOf(""),
		Data: &DataConfig{
			Type: "custom_type",
		},
	}
	result := GetTypeFromFieldConfig(fieldWithData)
	if result != "custom_type" {
		t.Errorf("GetTypeFromFieldConfig with data tag = %v, want custom_type", result)
	}

	// 测试没有data标签的情况
	fieldWithoutData := &FieldConfig{
		FieldType: reflect.TypeOf(0),
		Data:      nil,
	}
	result = GetTypeFromFieldConfig(fieldWithoutData)
	if result != "number" {
		t.Errorf("GetTypeFromFieldConfig without data tag = %v, want number", result)
	}
}

// TestTypeMappingTable 测试类型映射表
func TestTypeMappingTable(t *testing.T) {
	expectedMappings := map[string]string{
		"string":   "string",
		"int":      "number",
		"float64":  "float",
		"bool":     "boolean",
		"[]string": "[]string",
		"[]int":    "[]number",
		"[]float64": "[]float",
		"[]bool":   "[]boolean",
		"[]struct": "[]struct",
	}

	for goType, expectedType := range expectedMappings {
		if mappedType, exists := TypeMappingTable[goType]; !exists {
			t.Errorf("TypeMappingTable missing mapping for %v", goType)
		} else if mappedType != expectedType {
			t.Errorf("TypeMappingTable[%v] = %v, want %v", goType, mappedType, expectedType)
		}
	}
}

// TestTypeConstants 测试类型常量
func TestTypeConstants(t *testing.T) {
	constants := map[string]string{
		"TypeString":   TypeString,
		"TypeNumber":   TypeNumber,
		"TypeFloat":    TypeFloat,
		"TypeBoolean":  TypeBoolean,
		"TypeStruct":   TypeStruct,
		"TypeFiles":    TypeFiles,
		"TypeStrings":  TypeStrings,
		"TypeNumbers":  TypeNumbers,
		"TypeFloats":   TypeFloats,
		"TypeBooleans": TypeBooleans,
		"TypeStructs":  TypeStructs,
	}

	for name, value := range constants {
		if value == "" {
			t.Errorf("Type constant %v is empty", name)
		}
	}
} 