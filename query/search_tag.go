package query

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// SearchFieldConfig 搜索字段配置
type SearchFieldConfig struct {
	Field      string        `json:"field"`      // 字段名
	Name       string        `json:"name"`       // 显示名称
	DataType   string        `json:"data_type"`  // 数据类型
	Operators  []string      `json:"operators"`  // 支持的操作符
	Widget     *WidgetConfig `json:"widget"`     // 组件配置
	Permission string        `json:"permission"` // 权限配置
}

// WidgetConfig 组件配置
type WidgetConfig struct {
	Type        string            `json:"type"`        // 组件类型
	Options     []string          `json:"options"`     // 选项列表
	Placeholder string            `json:"placeholder"` // 占位符
	Prefix      string            `json:"prefix"`      // 前缀
	Suffix      string            `json:"suffix"`      // 后缀
	Format      string            `json:"format"`      // 格式
	TrueValue   string            `json:"true_value"`  // 开关真值
	FalseValue  string            `json:"false_value"` // 开关假值
	Extra       map[string]string `json:"extra"`       // 额外配置
}

// SearchFormConfig 搜索表单配置
type SearchFormConfig struct {
	Fields []SearchFieldConfig `json:"fields"` // 搜索字段列表
}

// BuildQueryConfigFromModel 根据模型的search标签构建QueryConfig
func BuildQueryConfigFromModel(model interface{}) (*QueryConfig, error) {
	config := NewQueryConfig()

	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// 获取search标签
		searchTag := field.Tag.Get("search")
		if searchTag == "" {
			continue
		}

		// 获取字段名（优先使用runner标签中的code）
		fieldName := getFieldCode(field)
		if fieldName == "" {
			fieldName = getGormColumnName(field)
		}
		if fieldName == "" {
			fieldName = strings.ToLower(field.Name)
		}

		// 解析支持的操作符
		operators := parseOperators(searchTag)
		if len(operators) == 0 {
			continue
		}

		// 检查权限标签
		permissionTag := field.Tag.Get("permission")
		if permissionTag == "write" {
			// 只写权限的字段不允许查询
			config.DenyField(fieldName)
			continue
		}

		// 添加到白名单
		config.AllowField(fieldName, operators...)
	}

	return config, nil
}

// ValidateSearchRequest 根据模型的search标签验证查询请求
func ValidateSearchRequest(model interface{}, pageInfo *PageInfoReq) error {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// 构建字段映射
	fieldMap := make(map[string][]string)

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		searchTag := field.Tag.Get("search")
		if searchTag == "" {
			continue
		}

		fieldName := getFieldCode(field)
		if fieldName == "" {
			fieldName = getGormColumnName(field)
		}
		if fieldName == "" {
			fieldName = strings.ToLower(field.Name)
		}

		operators := parseOperators(searchTag)
		fieldMap[fieldName] = operators
	}

	// 验证查询参数
	if err := validateOperators(pageInfo.Eq, "eq", fieldMap); err != nil {
		return err
	}
	if err := validateOperators(pageInfo.Like, "like", fieldMap); err != nil {
		return err
	}
	if err := validateOperators(pageInfo.In, "in", fieldMap); err != nil {
		return err
	}
	if err := validateOperators(pageInfo.Gt, "gt", fieldMap); err != nil {
		return err
	}
	if err := validateOperators(pageInfo.Gte, "gte", fieldMap); err != nil {
		return err
	}
	if err := validateOperators(pageInfo.Lt, "lt", fieldMap); err != nil {
		return err
	}
	if err := validateOperators(pageInfo.Lte, "lte", fieldMap); err != nil {
		return err
	}

	return nil
}

// GenerateSearchFormConfig 根据模型标签生成搜索表单配置
func GenerateSearchFormConfig(model interface{}) (*SearchFormConfig, error) {
	config := &SearchFormConfig{
		Fields: []SearchFieldConfig{},
	}

	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// 检查search标签
		searchTag := field.Tag.Get("search")
		if searchTag == "" {
			continue
		}

		// 检查权限
		permissionTag := field.Tag.Get("permission")
		if permissionTag == "write" {
			continue
		}

		// 构建字段配置
		fieldConfig := SearchFieldConfig{
			Field:      getFieldCode(field),
			Name:       getFieldName(field),
			DataType:   getDataType(field),
			Operators:  parseOperators(searchTag),
			Widget:     parseWidgetConfig(field),
			Permission: permissionTag,
		}

		if fieldConfig.Field == "" {
			fieldConfig.Field = getGormColumnName(field)
		}
		if fieldConfig.Field == "" {
			fieldConfig.Field = strings.ToLower(field.Name)
		}

		if fieldConfig.Name == "" {
			fieldConfig.Name = field.Name
		}

		config.Fields = append(config.Fields, fieldConfig)
	}

	return config, nil
}

// AutoSearchPaginated 自动搜索分页查询
func AutoSearchPaginated[T any](
	db *gorm.DB,
	model interface{},
	data T,
	pageInfo *PageInfoReq,
) (*PaginatedTable[T], error) {
	// 验证搜索请求
	if err := ValidateSearchRequest(model, pageInfo); err != nil {
		return nil, fmt.Errorf("搜索参数验证失败: %w", err)
	}

	// 构建查询配置
	config, err := BuildQueryConfigFromModel(model)
	if err != nil {
		return nil, fmt.Errorf("构建查询配置失败: %w", err)
	}

	// 执行分页查询
	return AutoPaginateTable(nil, db, model, data, pageInfo, config)
}

// 辅助函数

// getFieldCode 获取字段的code（从runner标签）
func getFieldCode(field reflect.StructField) string {
	runnerTag := field.Tag.Get("runner")
	if runnerTag == "" {
		return ""
	}

	pairs := strings.Split(runnerTag, ";")
	for _, pair := range pairs {
		if strings.HasPrefix(pair, "code:") {
			return strings.TrimPrefix(pair, "code:")
		}
	}

	return ""
}

// getFieldName 获取字段的显示名称（从runner标签）
func getFieldName(field reflect.StructField) string {
	runnerTag := field.Tag.Get("runner")
	if runnerTag == "" {
		return ""
	}

	pairs := strings.Split(runnerTag, ";")
	for _, pair := range pairs {
		if strings.HasPrefix(pair, "name:") {
			return strings.TrimPrefix(pair, "name:")
		}
	}

	return ""
}

// getGormColumnName 获取GORM字段名
func getGormColumnName(field reflect.StructField) string {
	gormTag := field.Tag.Get("gorm")
	if gormTag == "" {
		return ""
	}

	parts := strings.Split(gormTag, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}

	return ""
}

// getDataType 获取数据类型（从data标签）
func getDataType(field reflect.StructField) string {
	dataTag := field.Tag.Get("data")
	if dataTag == "" {
		// 根据Go类型推断
		switch field.Type.Kind() {
		case reflect.String:
			return "string"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return "number"
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return "number"
		case reflect.Float32, reflect.Float64:
			return "float"
		case reflect.Bool:
			return "boolean"
		default:
			return "string"
		}
	}

	pairs := strings.Split(dataTag, ";")
	for _, pair := range pairs {
		if strings.HasPrefix(pair, "type:") {
			return strings.TrimPrefix(pair, "type:")
		}
	}

	return "string"
}

// parseOperators 解析操作符列表
func parseOperators(searchTag string) []string {
	if searchTag == "" {
		return nil
	}

	operators := strings.Split(searchTag, ",")
	result := make([]string, 0, len(operators))

	for _, op := range operators {
		op = strings.TrimSpace(op)
		if op != "" {
			result = append(result, op)
		}
	}

	return result
}

// parseWidgetConfig 解析组件配置
func parseWidgetConfig(field reflect.StructField) *WidgetConfig {
	widgetTag := field.Tag.Get("widget")
	if widgetTag == "" {
		return nil
	}

	config := &WidgetConfig{
		Extra: make(map[string]string),
	}

	pairs := strings.Split(widgetTag, ";")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "type":
			config.Type = value
		case "options":
			config.Options = strings.Split(value, ",")
		case "placeholder":
			config.Placeholder = value
		case "prefix":
			config.Prefix = value
		case "suffix":
			config.Suffix = value
		case "format":
			config.Format = value
		case "true_value":
			config.TrueValue = value
		case "false_value":
			config.FalseValue = value
		default:
			config.Extra[key] = value
		}
	}

	return config
}

// validateOperators 验证操作符
func validateOperators(conditions []string, operator string, fieldMap map[string][]string) error {
	for _, condition := range conditions {
		parts := strings.SplitN(condition, ":", 2)
		if len(parts) != 2 {
			continue
		}

		fieldName := strings.TrimSpace(parts[0])

		// 检查字段是否有search标签
		allowedOps, exists := fieldMap[fieldName]
		if !exists {
			return fmt.Errorf("字段 %s 不支持搜索", fieldName)
		}

		// 检查操作符是否被允许
		found := false
		for _, allowedOp := range allowedOps {
			if allowedOp == operator {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("字段 %s 不支持 %s 操作", fieldName, operator)
		}
	}

	return nil
}
