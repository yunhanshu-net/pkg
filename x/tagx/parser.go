package tagx

import (
	"fmt"
	"github.com/yunhanshu-net/function-go/view/widget/types"
	"github.com/yunhanshu-net/pkg/x/slicesx"
	"reflect"
	"strings"
)

type RunnerTag map[string]string

func ParserKv(tag string) map[string]string {
	if tag == "" {
		return make(map[string]string)
	}
	mp := make(map[string]string)
	split := strings.Split(tag, ";")
	for _, s := range split {
		vals := strings.Split(s, ":")
		key := vals[0]
		value := vals[1]
		mp[key] = value
	}
	return mp
}

type RunnerFieldInfo struct {
	Name  string              // 字段名（含匿名字段层级，如"User.ID"）
	Type  reflect.StructField // 字段类型
	Value reflect.Value
	Tags  RunnerTag // 解析后的标签键值对
	Index []int     // 字段索引路径（用于反射）
}

func (r *RunnerFieldInfo) IsSearchCond() bool {
	_, ok := r.Tags["search_cond"]
	return ok
}

func (r *RunnerFieldInfo) GetCode() string {
	get := r.Tags["code"]
	if get != "" {
		return get
	}
	get = r.Type.Tag.Get("json")
	if get != "" {
		return strings.Split(get, ",")[0]
	}
	return r.Type.Name
}

func (r *RunnerFieldInfo) GetName() string {
	get := r.Tags["name"]
	if get != "" {
		return get
	}
	get = r.Type.Tag.Get("json")
	if get != "" {
		return strings.Split(get, ",")[0]
	}
	return r.Type.Name
}

func (r *RunnerFieldInfo) GetDesc() string {
	return r.Tags["desc"]
}
func (r *RunnerFieldInfo) GetDefaultValue() interface{} {
	return r.Tags["default_value"]
}

func (r *RunnerFieldInfo) GetCallbacks() string {
	return r.Tags["callback"]
}
func (r *RunnerFieldInfo) GetExample() string {
	return r.Tags["example"]
}

func (r *RunnerFieldInfo) GetShow() string {
	return r.Tags["show"]
}
func (r *RunnerFieldInfo) GetHidden() string {
	return r.Tags["hidden"]
}

func (r *RunnerFieldInfo) GetRequired() bool {
	validate := r.Type.Tag.Get("validate")
	split := strings.Split(validate, ",")
	return slicesx.ContainsString(split, "required")
}
func (r *RunnerFieldInfo) GetValidates() string {
	validate := r.Type.Tag.Get("validate")
	split := strings.Split(validate, ",")
	return strings.Join(slicesx.RemoveString(split, "required"), ",")
}

func (r *RunnerFieldInfo) GetValueType() string {
	valueType, err := r.getValueType()
	if err != nil {
		return types.ValueString
	}
	return types.UseValueType(r.Tags["type"], valueType)
}
func (i *RunnerFieldInfo) getValueType() (string, error) {
	switch i.Type.Type.Kind() {
	case reflect.Struct:
		return types.ValueObject, nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
		return types.ValueNumber, nil
	case reflect.Float32, reflect.Float64:
		return types.ValueFloat, nil
	case reflect.String:
		return types.ValueString, nil
	case reflect.Bool:
		return types.ValueBoolean, nil
	case reflect.Slice, reflect.Array:
		return types.ValueArray, nil
	case reflect.Map:
		return types.ValueObject, nil
	case reflect.Interface:
		return types.ValueObject, nil
	default:
		return "", fmt.Errorf("unsupported field type: %v", i.Type.Type.Kind())
	}
}

func ParseStructFieldsTypeOf(obj reflect.Type, tagKey string) ([]*RunnerFieldInfo, error) {
	return parseFields(obj, tagKey, nil, nil), nil
}

func GetSliceElementType(slice interface{}) (tp reflect.Type, err error) {
	t := reflect.TypeOf(slice)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice")
	}

	// 获取切片元素的类型
	elementType := t.Elem()

	// 如果是指针类型，可以进一步获取指向的类型
	if elementType.Kind() == reflect.Ptr {
		elementType = elementType.Elem()
	}
	//if elementType.Kind() != reflect.Struct {
	//	return nil, fmt.Errorf("input must be a struct")
	//}
	return elementType, nil
}

// 递归解析字段，处理匿名字段
func parseFields(t reflect.Type, tagKey string, parentIndex []int, parentNames []string) []*RunnerFieldInfo {
	var fields []*RunnerFieldInfo

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		currentIndex := append(parentIndex, i)
		currentNames := append(parentNames, field.Name)

		// 处理匿名字段（结构体类型）
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			embeddedFields := parseFields(field.Type, tagKey, currentIndex, parentNames)
			fields = append(fields, embeddedFields...)
			continue
		}
		get := field.Tag.Get(tagKey)
		if get == "-" {
			continue
		}
		// 普通字段：生成FieldInfo
		info := &RunnerFieldInfo{
			Value: reflect.New(field.Type).Elem(),
			Name:  strings.Join(currentNames, "."), // 生成层级字段名
			Type:  field,
			Tags:  ParseTagToMap(get),
			Index: currentIndex,
		}
		fields = append(fields, info)
	}
	return fields
}
func parseFieldsTag(t reflect.Type, tagKey string, parentIndex []int, parentNames []string) []*RunnerFieldInfo {
	var fields []*RunnerFieldInfo

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		currentIndex := append(parentIndex, i)
		currentNames := append(parentNames, field.Name)

		// 处理匿名字段（结构体类型）
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			embeddedFields := parseFields(field.Type, tagKey, currentIndex, parentNames)
			fields = append(fields, embeddedFields...)
			continue
		}

		// 普通字段：生成FieldInfo
		info := &RunnerFieldInfo{
			Name:  strings.Join(currentNames, "."), // 生成层级字段名
			Type:  field,
			Tags:  ParseTagToMap(field.Tag.Get(tagKey)),
			Index: currentIndex,
		}
		fields = append(fields, info)
	}
	return fields
}

// ParseTagToMap 解析标签字符串为键值对（支持GORM风格分号分隔）
func ParseTagToMap(tag string) map[string]string {
	result := make(map[string]string)
	for _, pair := range strings.Split(tag, ";") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, ":", 2)
		key := strings.TrimSpace(kv[0])
		if key == "" {
			continue
		}
		value := ""
		if len(kv) > 1 {
			value = strings.TrimSpace(kv[1])
		}
		result[key] = value
	}
	return result
}
