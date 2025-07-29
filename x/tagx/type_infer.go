package tagx

import (
	"reflect"
)

// TypeInferrer 类型推断器
type TypeInferrer struct {
	// 特殊类型映射表
	specialTypeMap map[string]string
	// 基础类型映射表
	basicTypeMap map[reflect.Kind]string
	// 切片类型映射表
	sliceTypeMap map[reflect.Kind]string
}

// NewTypeInferrer 创建类型推断器
func NewTypeInferrer() *TypeInferrer {
	return &TypeInferrer{
		// 特殊类型映射表
		specialTypeMap: map[string]string{
			"*files.Files": "files",
			"files.Files":  "files",
			"files.Writer": "files",
		},
		// 基础类型映射表
		basicTypeMap: map[reflect.Kind]string{
			reflect.String: "string",
			reflect.Int:    "number",
			reflect.Int8:   "number",
			reflect.Int16:  "number",
			reflect.Int32:  "number",
			reflect.Int64:  "number",
			reflect.Uint:   "number",
			reflect.Uint8:  "number",
			reflect.Uint16: "number",
			reflect.Uint32: "number",
			reflect.Uint64: "number",
			reflect.Float32: "float",
			reflect.Float64: "float",
			reflect.Bool:    "boolean",
			reflect.Struct:  "struct",
		},
		// 切片类型映射表 - 只支持定义好的具体类型
		sliceTypeMap: map[reflect.Kind]string{
			reflect.String: "[]string",
			reflect.Int:    "[]number",
			reflect.Int8:   "[]number",
			reflect.Int16:  "[]number",
			reflect.Int32:  "[]number",
			reflect.Int64:  "[]number",
			reflect.Uint:   "[]number",
			reflect.Uint8:  "[]number",
			reflect.Uint16: "[]number",
			reflect.Uint32: "[]number",
			reflect.Uint64: "[]number",
			reflect.Float32: "[]float",
			reflect.Float64: "[]float",
			reflect.Bool:    "[]boolean",
			reflect.Struct:  "[]struct",
		},
	}
}

// InferType 推断字段类型
func (t *TypeInferrer) InferType(fieldType reflect.Type) string {
	// 处理指针类型
	if fieldType.Kind() == reflect.Ptr {
		return t.InferType(fieldType.Elem())
	}

	typeStr := fieldType.String()

	// 1. 检查特殊类型映射表
	if mappedType, exists := t.specialTypeMap[typeStr]; exists {
		return mappedType
	}

	// 2. 处理切片类型
	if fieldType.Kind() == reflect.Slice {
		elemType := fieldType.Elem()
		// 处理指针切片
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		
		// 检查切片元素类型
		if mappedType, exists := t.sliceTypeMap[elemType.Kind()]; exists {
			return mappedType
		}
		
		// 特殊切片类型检查
		if elemType.String() == "files.Files" {
			return "files"
		}
		
		// 不支持抽象数组类型，返回空字符串表示忽略该字段
		// 注意：只有定义好的具体类型才支持，避免抽象类型导致渲染问题
		return ""
	}

	// 3. 处理基础类型
	if mappedType, exists := t.basicTypeMap[fieldType.Kind()]; exists {
		return mappedType
	}

	// 4. 不支持的类型返回空字符串（表示忽略该字段）
	// 注意：map、interface 等动态类型不支持，因为无法在程序启动时确定具体类型
	return ""
}

// InferTypeFromValue 从值推断类型
func (t *TypeInferrer) InferTypeFromValue(value interface{}) string {
	if value == nil {
		return "string"
	}
	return t.InferType(reflect.TypeOf(value))
}

// InferTypeFromField 从结构体字段推断类型
func (t *TypeInferrer) InferTypeFromField(field reflect.StructField) string {
	return t.InferType(field.Type)
}

// GetTypeFromFieldConfig 从 FieldConfig 获取类型（兼容现有代码）
func (t *TypeInferrer) GetTypeFromFieldConfig(field *FieldConfig) string {
	// 优先使用data标签中的type
	if field.Data != nil && field.Data.Type != "" {
		return field.Data.Type
	}

	// 使用类型推断器推断类型
	return t.InferType(field.FieldType)
}

// 全局类型推断器实例
var globalTypeInferrer = NewTypeInferrer()

// InferType 全局类型推断函数
func InferType(fieldType reflect.Type) string {
	return globalTypeInferrer.InferType(fieldType)
}

// InferTypeFromValue 全局从值推断类型函数
func InferTypeFromValue(value interface{}) string {
	return globalTypeInferrer.InferTypeFromValue(value)
}

// InferTypeFromField 全局从结构体字段推断类型函数
func InferTypeFromField(field reflect.StructField) string {
	return globalTypeInferrer.InferTypeFromField(field)
}

// GetTypeFromFieldConfig 全局从 FieldConfig 获取类型函数
func GetTypeFromFieldConfig(field *FieldConfig) string {
	return globalTypeInferrer.GetTypeFromFieldConfig(field)
}

	// 类型映射常量
	const (
		TypeString   = "string"
		TypeNumber   = "number"
		TypeFloat    = "float"
		TypeBoolean  = "boolean"
		TypeStruct   = "struct"
		TypeFiles    = "files"
		TypeStrings  = "[]string"
		TypeNumbers  = "[]number"
		TypeFloats   = "[]float"
		TypeBooleans = "[]boolean"
		TypeStructs  = "[]struct"
	)

	// 类型映射表（用于文档和测试）
	var TypeMappingTable = map[string]string{
		// 基础类型
		"string":   TypeString,
		"int":      TypeNumber,
		"int8":     TypeNumber,
		"int16":    TypeNumber,
		"int32":    TypeNumber,
		"int64":    TypeNumber,
		"uint":     TypeNumber,
		"uint8":    TypeNumber,
		"uint16":   TypeNumber,
		"uint32":   TypeNumber,
		"uint64":   TypeNumber,
		"float32":  TypeFloat,
		"float64":  TypeFloat,
		"bool":     TypeBoolean,
		
		// 切片类型 - 只支持定义好的具体类型
		"[]string":   TypeStrings,
		"[]int":      TypeNumbers,
		"[]int8":     TypeNumbers,
		"[]int16":    TypeNumbers,
		"[]int32":    TypeNumbers,
		"[]int64":    TypeNumbers,
		"[]uint":     TypeNumbers,
		"[]uint8":    TypeNumbers,
		"[]uint16":   TypeNumbers,
		"[]uint32":   TypeNumbers,
		"[]uint64":   TypeNumbers,
		"[]float32":  TypeFloats,
		"[]float64":  TypeFloats,
		"[]bool":     TypeBooleans,
		"[]struct":   TypeStructs,
		
		// 特殊类型
		"*files.Files": TypeFiles,
		"files.Files":  TypeFiles,
		"files.Writer": TypeFiles,
		
		// 复杂类型
		"struct": TypeStruct,
		// 注意：map、interface 等动态类型不支持，因为无法在程序启动时确定具体类型
		// 注意：抽象数组类型不支持，只允许定义好的具体类型
		// 注意：时间类型统一使用 int64 时间戳存储，前端通过 datetime 组件展示
	} 