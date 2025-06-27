package tagx

import (
	"reflect"
	"strconv"
	"strings"
)

// MultiTagParser 多标签解析器
type MultiTagParser struct {
	supportedTags []string
}

// FieldConfig 字段配置信息
type FieldConfig struct {
	FieldName string
	FieldType reflect.Type
	JsonTag   string
	FormTag   string

	// 各标签解析结果
	Runner     *RunnerConfig
	Widget     *WidgetConfig
	Data       *DataConfig
	Permission *PermissionConfig
	Callbacks  []*CallbackConfig
	Validation string // 简化为字符串
}

// RunnerConfig runner标签配置 - 移除type字段
type RunnerConfig struct {
	Code string
	Name string
	Desc string
}

// WidgetConfig widget标签配置
type WidgetConfig struct {
	Type   string                 `json:"type"`   // 组件类型
	Config map[string]interface{} `json:"config"` // 灵活的配置参数
}

// DataConfig data标签配置 - 新增type字段
type DataConfig struct {
	Type         string // 数据类型：string, number, boolean, []string等
	Example      string
	DefaultValue string
	Source       string
	Format       string
}

// PermissionConfig permission标签配置 - 简化为三权限
type PermissionConfig struct {
	Read   bool // 可读权限
	Update bool // 可更新权限
	Create bool // 可创建权限
}

// CallbackConfig callback配置
type CallbackConfig struct {
	Event  string
	Params map[string]string
}

// NewMultiTagParser 创建多标签解析器
func NewMultiTagParser() *MultiTagParser {
	return &MultiTagParser{
		supportedTags: []string{"runner", "widget", "data", "permission", "callback", "validate"},
	}
}

// ParseStruct 解析结构体
func (p *MultiTagParser) ParseStruct(structType reflect.Type) ([]*FieldConfig, error) {
	var fields []*FieldConfig

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// 跳过未导出的字段
		if !field.IsExported() {
			continue
		}

		// 检查runner标签，如果是"-"则跳过该字段
		if runnerTag := field.Tag.Get("runner"); runnerTag == "-" {
			continue
		}

		config := &FieldConfig{
			FieldName: field.Name,
			FieldType: field.Type,
			JsonTag:   field.Tag.Get("json"),
			FormTag:   field.Tag.Get("form"),
		}

		// 解析各个标签
		if runnerTag := field.Tag.Get("runner"); runnerTag != "" {
			config.Runner = p.parseRunnerTag(runnerTag)
		}

		if widgetTag := field.Tag.Get("widget"); widgetTag != "" {
			config.Widget = p.parseWidgetTag(widgetTag)
		}

		if dataTag := field.Tag.Get("data"); dataTag != "" {
			config.Data = p.parseDataTag(dataTag)
		}

		if permissionTag := field.Tag.Get("permission"); permissionTag != "" {
			config.Permission = p.parsePermissionTag(permissionTag)
		}

		if callbackTag := field.Tag.Get("callback"); callbackTag != "" {
			config.Callbacks = p.parseCallbackTag(callbackTag)
		}

		if validateTag := field.Tag.Get("validate"); validateTag != "" {
			config.Validation = validateTag // 直接返回字符串
		}

		fields = append(fields, config)
	}

	return fields, nil
}

// parseRunnerTag 解析runner标签 - 移除type字段
// 格式: code:title;name:标题;desc:描述
func (p *MultiTagParser) parseRunnerTag(tag string) *RunnerConfig {
	config := &RunnerConfig{}
	parts := strings.Split(tag, ";")

	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			switch strings.TrimSpace(kv[0]) {
			case "code":
				config.Code = strings.TrimSpace(kv[1])
			case "name":
				config.Name = strings.TrimSpace(kv[1])
			case "desc":
				config.Desc = strings.TrimSpace(kv[1])
			}
		}
	}

	return config
}

// parseWidgetTag 解析widget标签 - 支持灵活配置
// 格式: type:input;placeholder:请输入;mode:text;options:类型1,类型2,类型3;accept:.jpg,.png;max_size:2MB
func (p *MultiTagParser) parseWidgetTag(tag string) *WidgetConfig {
	config := &WidgetConfig{
		Config: make(map[string]interface{}),
	}
	parts := strings.Split(tag, ";")

	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			if key == "type" {
				config.Type = value
			} else {
				// 其他所有配置都放入Config map中
				config.Config[key] = p.parseConfigValue(config.Type, key, value)
			}
		}
	}

	return config
}

// WidgetConfigParser 组件配置解析器接口
type WidgetConfigParser interface {
	ParseConfig(key, value string) (interface{}, bool)
	GetSupportedKeys() []string
}

// 各组件的配置解析器
var widgetParsers = map[string]WidgetConfigParser{
	"file_upload":  &FileUploadParser{},
	"color":        &ColorParser{},
	"datetime":     &DateTimeParser{},
	"slider":       &SliderParser{},
	"multiselect":  &MultiSelectParser{},
	"select":       &SelectParser{},
	"input":        &InputParser{},
	"switch":       &SwitchParser{},
	"radio":        &RadioParser{},
	"checkbox":     &CheckboxParser{},
	"file_display": &FileDisplayParser{},
	"table":        &TableParser{},
}

// parseConfigValue 解析配置值 - 使用组件特定的解析器
func (p *MultiTagParser) parseConfigValue(widgetType, key, value string) interface{} {
	// 查找对应组件的解析器
	if parser, exists := widgetParsers[widgetType]; exists {
		if parsedValue, handled := parser.ParseConfig(key, value); handled {
			return parsedValue
		}
	}

	// 通用解析器作为后备
	return p.parseGenericValue(key, value)
}

// parseGenericValue 通用配置解析器
func (p *MultiTagParser) parseGenericValue(key, value string) interface{} {
	// 通用的布尔值解析
	if value == "true" || value == "false" {
		return value == "true"
	}

	// 通用的数字解析
	if num, err := strconv.Atoi(value); err == nil {
		return num
	}

	// 默认返回字符串
	return value
}

// parseDataTag 解析data标签 - 新增type字段
// 格式: type:string;example:示例值;default_value:默认值;source:api://users;format:YYYY-MM-DD
func (p *MultiTagParser) parseDataTag(tag string) *DataConfig {
	config := &DataConfig{}
	parts := strings.Split(tag, ";")

	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			switch key {
			case "type":
				config.Type = value
			case "example":
				config.Example = value
			case "default_value":
				config.DefaultValue = value
			case "source":
				config.Source = value
			case "format":
				config.Format = value
			}
		}
	}

	return config
}

// parsePermissionTag 解析permission标签 - 简化为三权限
// 格式: read,update,create 或 read,update 或 create 等
func (p *MultiTagParser) parsePermissionTag(tag string) *PermissionConfig {
	config := &PermissionConfig{
		Read:   true, // 默认全部权限
		Update: true,
		Create: true,
	}

	if tag == "" {
		return config // 空标签表示全部权限
	}

	// 重置为false，根据标签内容设置
	config.Read = false
	config.Update = false
	config.Create = false

	permissions := strings.Split(tag, ",")
	for _, perm := range permissions {
		perm = strings.TrimSpace(perm)
		switch perm {
		case "read":
			config.Read = true
		case "update":
			config.Update = true
		case "create":
			config.Create = true
		}
	}

	return config
}

// parseCallbackTag 解析callback标签 - 字段级别的回调
// 格式: OnInputFuzzy(delay:300,min:2);OnBlur(validate:true) (不包含OnPageLoad等函数级别回调)
func (p *MultiTagParser) parseCallbackTag(tag string) []*CallbackConfig {
	var callbacks []*CallbackConfig

	parts := strings.Split(tag, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 跳过函数级别的回调
		if strings.HasPrefix(part, "OnPageLoad") {
			continue // OnPageLoad是函数级别回调，不在这里处理
		}

		if strings.Contains(part, "(") {
			// 带参数的回调: OnInputFuzzy(delay:300,min:2)
			eventEnd := strings.Index(part, "(")
			event := strings.TrimSpace(part[:eventEnd])
			paramsStr := part[eventEnd+1 : len(part)-1] // 去掉括号

			params := make(map[string]string)
			if paramsStr != "" {
				paramPairs := strings.Split(paramsStr, ",")
				for _, pair := range paramPairs {
					kv := strings.SplitN(pair, ":", 2)
					if len(kv) == 2 {
						params[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
					}
				}
			}

			callbacks = append(callbacks, &CallbackConfig{
				Event:  event,
				Params: params,
			})
		} else {
			// 无参数的回调: OnValueChange
			callbacks = append(callbacks, &CallbackConfig{
				Event:  part,
				Params: make(map[string]string),
			})
		}
	}

	return callbacks
}

// GetCode 获取字段代码
func (f *FieldConfig) GetCode() string {
	if f.Runner != nil && f.Runner.Code != "" {
		return f.Runner.Code
	}

	// 从json标签获取
	if f.JsonTag != "" {
		parts := strings.Split(f.JsonTag, ",")
		if len(parts) > 0 && parts[0] != "" {
			return parts[0]
		}
	}

	// 使用字段名的小写形式
	return strings.ToLower(f.FieldName)
}

// GetName 获取字段显示名称
func (f *FieldConfig) GetName() string {
	if f.Runner != nil && f.Runner.Name != "" {
		return f.Runner.Name
	}
	return f.FieldName
}

// GetType 获取字段类型 - 支持复合类型，参考widget_value_type.go
func (f *FieldConfig) GetType() string {
	// 优先使用data标签中的type
	if f.Data != nil && f.Data.Type != "" {
		return f.Data.Type
	}

	// 根据Go类型推断，参考widget_value_type.go中的类型定义
	typeStr := f.FieldType.String()

	// 特殊类型检查（文件类型）
	if typeStr == "*files.Files" || typeStr == "files.Files" || typeStr == "files.Writer" {
		return "files"
	}

	switch f.FieldType.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "number"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice:
		// 支持切片类型
		elemType := f.FieldType.Elem()
		switch elemType.Kind() {
		case reflect.String:
			return "[]string" // 支持[]string类型
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return "[]number"
		default:
			return "array"
		}
	case reflect.Map:
		return "object"
	case reflect.Struct:
		// 检查是否是时间类型
		if typeStr == "time.Time" {
			return "time"
		}
		// 检查是否是files相关类型
		if typeStr == "files.Writer" || typeStr == "files.Files" {
			return "files"
		}
		return "object"
	case reflect.Ptr:
		// 处理指针类型
		elemType := f.FieldType.Elem()
		if elemType.String() == "files.Files" {
			return "files"
		}
		// 递归处理指针指向的类型
		return f.getTypeFromReflectType(elemType)
	default:
		return "string"
	}
}

// getTypeFromReflectType 辅助方法，处理reflect.Type的类型推断
func (f *FieldConfig) getTypeFromReflectType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "number"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Struct:
		if t.String() == "files.Files" || t.String() == "files.Writer" {
			return "files"
		}
		return "object"
	default:
		return "string"
	}
}

// FileUploadParser 文件上传组件配置解析器
type FileUploadParser struct{}

// ParseConfig 解析文件上传组件配置
func (p *FileUploadParser) ParseConfig(key, value string) (interface{}, bool) {
	switch key {
	case "accept":
		// 文件类型限制，字符串格式：.jpg,.png,.pdf
		return value, true
	case "max_size":
		// 最大文件大小，字符串格式：2MB, 10GB等
		return value, true
	case "max_count":
		// 最大文件数量，数字格式
		if num, err := strconv.Atoi(value); err == nil {
			return num, true
		}
		return value, true
	case "preview", "drag_drop":
		// 布尔值配置
		return value == "true", true
	case "placeholder":
		// 占位符文本
		return value, true
	}
	return nil, false
}

// GetSupportedKeys 获取支持的配置键
func (p *FileUploadParser) GetSupportedKeys() []string {
	return []string{"accept", "max_size", "max_count", "preview", "drag_drop", "placeholder"}
}

// ColorParser 颜色选择器组件配置解析器
type ColorParser struct{}

// ParseConfig 解析颜色选择器组件配置
func (p *ColorParser) ParseConfig(key, value string) (interface{}, bool) {
	switch key {
	case "format":
		// 颜色格式：hex, rgb, rgba, hsl, hsla
		return value, true
	case "predefine":
		// 预定义颜色列表，逗号分隔
		if value != "" {
			colors := strings.Split(value, ",")
			for i, color := range colors {
				colors[i] = strings.TrimSpace(color)
			}
			return colors, true
		}
		return []string{}, true
	case "show_alpha", "show_swatches", "allow_empty":
		// 布尔值配置
		return value == "true", true
	case "placeholder":
		// 占位符文本
		return value, true
	}
	return nil, false
}

// GetSupportedKeys 获取支持的配置键
func (p *ColorParser) GetSupportedKeys() []string {
	return []string{"format", "predefine", "show_alpha", "show_swatches", "allow_empty", "placeholder"}
}

// DateTimeParser 日期时间组件解析器
type DateTimeParser struct{}

func (p *DateTimeParser) ParseConfig(key, value string) (interface{}, bool) {
	switch key {
	case "format":
		// 验证格式是否有效
		validFormats := map[string]bool{
			"date": true, "datetime": true, "time": true,
			"daterange": true, "datetimerange": true,
			"month": true, "year": true, "week": true,
		}
		if validFormats[value] {
			return value, true
		}
		return "date", true // 默认格式
	case "placeholder", "start_placeholder", "end_placeholder":
		return value, true
	case "default_value", "default_time", "min_date", "max_date", "separator":
		return value, true
	case "shortcuts":
		return value, true // 直接返回字符串，由前端解析
	case "disabled":
		return strings.ToLower(value) == "true", true
	default:
		return nil, false
	}
}

func (p *DateTimeParser) GetSupportedKeys() []string {
	return []string{"format", "placeholder", "start_placeholder", "end_placeholder",
		"default_value", "default_time", "min_date", "max_date", "separator", "shortcuts", "disabled"}
}

// SliderParser 滑块组件配置解析器
type SliderParser struct{}

// ParseConfig 解析滑块组件配置
func (p *SliderParser) ParseConfig(key, value string) (interface{}, bool) {
	switch key {
	case "mode":
		// 滑块模式：range等
		return value, true
	case "min", "max", "step":
		// 数值配置，优先解析为数字
		if num, err := strconv.Atoi(value); err == nil {
			return num, true
		}
		return value, true
	case "show_stops", "show_tooltip":
		// 布尔值配置
		return value == "true", true
	}
	return nil, false
}

// GetSupportedKeys 获取支持的配置键
func (p *SliderParser) GetSupportedKeys() []string {
	return []string{"mode", "min", "max", "step", "show_stops", "show_tooltip"}
}

// MultiSelectParser 多选组件配置解析器
type MultiSelectParser struct{}

// ParseConfig 解析多选组件配置
func (p *MultiSelectParser) ParseConfig(key, value string) (interface{}, bool) {
	switch key {
	case "options":
		// 选项列表，逗号分隔
		if value != "" {
			options := strings.Split(value, ",")
			for i, option := range options {
				options[i] = strings.TrimSpace(option)
			}
			return options, true
		}
		return []string{}, true
	case "placeholder":
		// 占位符文本
		return value, true
	case "multiple_limit":
		// 最大选择数量
		if num, err := strconv.Atoi(value); err == nil {
			return num, true
		}
		return value, true
	case "filterable", "clearable", "collapse_tags":
		// 布尔值配置
		return value == "true", true
	}
	return nil, false
}

// GetSupportedKeys 获取支持的配置键
func (p *MultiSelectParser) GetSupportedKeys() []string {
	return []string{"options", "placeholder", "multiple_limit", "filterable", "clearable", "collapse_tags"}
}

// SelectParser 选择器组件配置解析器
type SelectParser struct{}

// ParseConfig 解析选择器组件配置
func (p *SelectParser) ParseConfig(key, value string) (interface{}, bool) {
	switch key {
	case "options":
		// 选项列表，逗号分隔
		if value != "" {
			options := strings.Split(value, ",")
			for i, option := range options {
				options[i] = strings.TrimSpace(option)
			}
			return options, true
		}
		return []string{}, true
	case "placeholder":
		// 占位符文本
		return value, true
	case "filterable", "clearable":
		// 布尔值配置
		return value == "true", true
	}
	return nil, false
}

// GetSupportedKeys 获取支持的配置键
func (p *SelectParser) GetSupportedKeys() []string {
	return []string{"options", "placeholder", "filterable", "clearable"}
}

// InputParser 输入框组件配置解析器
type InputParser struct{}

// ParseConfig 解析输入框组件配置
func (p *InputParser) ParseConfig(key, value string) (interface{}, bool) {
	switch key {
	case "mode":
		// 输入框模式：text, password, textarea等
		return value, true
	case "placeholder":
		// 占位符文本
		return value, true
	case "maxlength":
		// 最大长度
		if num, err := strconv.Atoi(value); err == nil {
			return num, true
		}
		return value, true
	case "show_password", "clearable", "show_word_limit":
		// 布尔值配置
		return value == "true", true
	}
	return nil, false
}

// GetSupportedKeys 获取支持的配置键
func (p *InputParser) GetSupportedKeys() []string {
	return []string{"mode", "placeholder", "maxlength", "show_password", "clearable", "show_word_limit"}
}

// SwitchParser 开关组件配置解析器
type SwitchParser struct{}

// ParseConfig 解析开关组件配置
func (p *SwitchParser) ParseConfig(key, value string) (interface{}, bool) {
	switch key {
	case "true_label", "false_label":
		// 开关标签文本
		return value, true
	case "inline_prompt":
		// 是否显示内联提示
		return value == "true", true
	}
	return nil, false
}

// GetSupportedKeys 获取支持的配置键
func (p *SwitchParser) GetSupportedKeys() []string {
	return []string{"true_label", "false_label", "inline_prompt"}
}

// RadioParser 单选按钮组件解析器
type RadioParser struct{}

func (p *RadioParser) ParseConfig(key, value string) (interface{}, bool) {
	switch key {
	case "options":
		return value, true // 直接返回options字符串，由前端解析
	case "default_value":
		return value, true
	default:
		return nil, false
	}
}

func (p *RadioParser) GetSupportedKeys() []string {
	return []string{"options", "default_value"}
}

// CheckboxParser 复选框组件解析器
type CheckboxParser struct{}

func (p *CheckboxParser) ParseConfig(key, value string) (interface{}, bool) {
	switch key {
	case "options":
		return value, true // 直接返回options字符串，由前端解析
	case "default_value":
		return value, true
	case "multiple_limit":
		if val, err := strconv.Atoi(value); err == nil {
			return val, true
		}
		return 0, true
	case "show_select_all":
		return strings.ToLower(value) == "true", true
	default:
		return nil, false
	}
}

func (p *CheckboxParser) GetSupportedKeys() []string {
	return []string{"options", "default_value", "multiple_limit", "show_select_all"}
}

// FileDisplayParser 文件展示组件解析器
type FileDisplayParser struct{}

func (p *FileDisplayParser) ParseConfig(key, value string) (interface{}, bool) {
	switch key {
	case "display_mode":
		return value, true // card, list, grid
	case "preview", "download":
		return strings.ToLower(value) == "true", true
	case "max_preview":
		if val, err := strconv.Atoi(value); err == nil {
			return val, true
		}
		return 10, true
	default:
		return nil, false
	}
}

func (p *FileDisplayParser) GetSupportedKeys() []string {
	return []string{"display_mode", "preview", "download", "max_preview"}
}

// TableParser 表格组件解析器
type TableParser struct{}

func (p *TableParser) ParseConfig(key, value string) (interface{}, bool) {
	switch key {
	case "columns":
		return value, true // 表格列配置
	case "pagination":
		return strings.ToLower(value) == "true", true
	case "sortable", "filterable", "selectable":
		return strings.ToLower(value) == "true", true
	case "page_size":
		if val, err := strconv.Atoi(value); err == nil {
			return val, true
		}
		return 10, true
	case "height":
		return value, true // 表格高度
	default:
		return nil, false
	}
}

func (p *TableParser) GetSupportedKeys() []string {
	return []string{"columns", "pagination", "sortable", "filterable", "selectable", "page_size", "height"}
}

// TestAllWidgetParsers 测试所有组件解析器（用于验证完整性）
func TestAllWidgetParsers() map[string][]string {
	result := make(map[string][]string)
	for widgetType, parser := range widgetParsers {
		result[widgetType] = parser.GetSupportedKeys()
	}
	return result
}
