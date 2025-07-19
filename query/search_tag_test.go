package query

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestSearchProduct 测试用的产品模型（带search标签）
type TestSearchProduct struct {
	ID int `json:"id" gorm:"primaryKey" runner:"code:id;name:产品ID" data:"type:number"`

	// 文本搜索字段
	Name string `json:"name" gorm:"column:name;comment:产品名称" runner:"code:name;name:产品名称" widget:"type:input;placeholder:请输入产品名称" data:"type:string" search:"like,eq" validate:"required"`

	// 分类选择字段 - 支持精确匹配和多选
	Category string `json:"category" gorm:"column:category;comment:产品分类" runner:"code:category;name:产品分类" widget:"type:select;options:手机,笔记本,平板,耳机" data:"type:string" search:"eq,in" validate:"required"`

	// 价格范围搜索
	Price float64 `json:"price" gorm:"column:price;comment:产品价格" runner:"code:price;name:产品价格" widget:"type:input;prefix:¥;precision:2" data:"type:float" search:"gte,lte,eq" validate:"required,min=0"`

	// 库存数量搜索
	Stock int `json:"stock" gorm:"column:stock;comment:库存数量" runner:"code:stock;name:库存数量" widget:"type:input;suffix:件" data:"type:number" search:"gte,lte,eq,gt,lt" validate:"required,min=0"`

	// 状态开关搜索
	Status bool `json:"status" gorm:"column:status;comment:产品状态" runner:"code:status;name:产品状态" widget:"type:switch;true_label:启用;false_label:禁用" data:"type:boolean;default_value:true" search:"eq" validate:"required"`

	// 标签模糊搜索
	Tags string `json:"tags" gorm:"column:tags;comment:产品标签" runner:"code:tags;name:产品标签" widget:"type:tag;separator:,;max_tags:5" data:"type:string" search:"like,in"`

	// 只读字段 - 不可搜索
	CreatedAt string `json:"created_at" gorm:"column:created_at" runner:"code:created_at;name:创建时间" widget:"type:datetime;format:datetime" data:"type:string" permission:"read"`

	// 密码字段 - 不可搜索
	SecretKey string `json:"-" gorm:"column:secret_key;comment:密钥" runner:"code:secret_key;name:密钥" widget:"type:input;mode:password" data:"type:string" search:"eq" permission:"write" validate:"required"`
}

func (TestSearchProduct) TableName() string {
	return "test_search_products"
}

// TestBuildQueryConfigFromModel 测试从模型构建查询配置
func TestBuildQueryConfigFromModel(t *testing.T) {
	t.Log("=== 测试从模型构建查询配置 ===")

	config, err := BuildQueryConfigFromModel(&TestSearchProduct{})
	if err != nil {
		t.Fatalf("构建查询配置失败: %v", err)
	}

	// 检查白名单配置
	expectedFields := map[string][]string{
		"name":     {"like", "eq"},
		"category": {"eq", "in"},
		"price":    {"gte", "lte", "eq"},
		"stock":    {"gte", "lte", "eq", "gt", "lt"},
		"status":   {"eq"},
		"tags":     {"like", "in"},
	}

	for field, expectedOps := range expectedFields {
		actualOps, exists := config.Fields[field]
		if !exists {
			t.Errorf("字段 %s 未在白名单中找到", field)
			continue
		}

		t.Logf("字段 %s 支持的操作符: %v", field, actualOps)

		// 检查操作符是否匹配
		for _, expectedOp := range expectedOps {
			found := false
			for _, actualOp := range actualOps {
				if actualOp == expectedOp {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("字段 %s 缺少操作符 %s", field, expectedOp)
			}
		}
	}

	// 检查黑名单配置
	if _, exists := config.Blacklist["secret_key"]; !exists {
		t.Errorf("secret_key 字段应该在黑名单中")
	}

	t.Logf("查询配置构建成功，白名单字段数: %d，黑名单字段数: %d",
		len(config.Fields), len(config.Blacklist))
}

// TestValidateSearchRequest 测试搜索请求验证
func TestValidateSearchRequest(t *testing.T) {
	t.Log("=== 测试搜索请求验证 ===")

	// 测试合法的搜索请求
	validPageInfo := &PageInfoReq{
		Eq:   []string{"category:手机", "status:启用"},
		Like: []string{"name:苹果", "tags:智能"},
		Gte:  []string{"price:1000", "stock:10"},
		Lte:  []string{"price:5000", "stock:100"},
	}

	err := ValidateSearchRequest(&TestSearchProduct{}, validPageInfo)
	if err != nil {
		t.Errorf("合法请求验证失败: %v", err)
	} else {
		t.Log("合法请求验证通过")
	}

	// 测试不支持的字段
	invalidPageInfo1 := &PageInfoReq{
		Eq: []string{"unknown_field:test"}, // unknown_field不支持搜索
	}

	err = ValidateSearchRequest(&TestSearchProduct{}, invalidPageInfo1)
	if err == nil {
		t.Error("应该拒绝不支持搜索的字段")
	} else {
		t.Logf("正确拒绝了不支持的字段: %v", err)
	}

	// 测试不支持的操作符
	invalidPageInfo2 := &PageInfoReq{
		Like: []string{"status:启用"}, // status不支持like操作
	}

	err = ValidateSearchRequest(&TestSearchProduct{}, invalidPageInfo2)
	if err == nil {
		t.Error("应该拒绝不支持的操作符")
	} else {
		t.Logf("正确拒绝了不支持的操作符: %v", err)
	}
}

// TestGenerateSearchFormConfig 测试生成搜索表单配置
func TestGenerateSearchFormConfig(t *testing.T) {
	t.Log("=== 测试生成搜索表单配置 ===")

	config, err := GenerateSearchFormConfig(&TestSearchProduct{})
	if err != nil {
		t.Fatalf("生成搜索表单配置失败: %v", err)
	}

	// 打印配置JSON
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	t.Logf("搜索表单配置:\n%s", string(configJSON))

	// 验证字段数量
	expectedFieldCount := 6 // name, category, price, stock, status, tags
	if len(config.Fields) != expectedFieldCount {
		t.Errorf("期望 %d 个搜索字段，实际: %d", expectedFieldCount, len(config.Fields))
	}

	// 验证特定字段配置
	for _, field := range config.Fields {
		t.Logf("字段: %s, 名称: %s, 类型: %s, 操作符: %v",
			field.Field, field.Name, field.DataType, field.Operators)

		switch field.Field {
		case "name":
			if field.DataType != "string" {
				t.Errorf("name字段数据类型应为string，实际: %s", field.DataType)
			}
			if field.Widget.Type != "input" {
				t.Errorf("name字段组件类型应为input，实际: %s", field.Widget.Type)
			}

		case "category":
			if len(field.Widget.Options) != 4 {
				t.Errorf("category字段应有4个选项，实际: %d", len(field.Widget.Options))
			}

		case "price":
			if field.Widget.Prefix != "¥" {
				t.Errorf("price字段前缀应为¥，实际: %s", field.Widget.Prefix)
			}

		case "status":
			if field.Widget.TrueValue != "启用" {
				t.Errorf("status字段真值应为启用，实际: %s", field.Widget.TrueValue)
			}
		}
	}
}

// TestAutoSearchPaginated 测试自动搜索分页查询
func TestAutoSearchPaginated(t *testing.T) {
	t.Log("=== 测试自动搜索分页查询 ===")

	// 设置测试数据库
	db := setupSearchTestDB(t)

	// 测试基础搜索
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Eq:       []string{"category:手机"},
		Like:     []string{"name:iPhone"},
		Gte:      []string{"price:1000"},
		Lte:      []string{"price:8000"},
	}

	var products []TestSearchProduct
	result, err := AutoSearchPaginated(db, &TestSearchProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("自动搜索分页查询失败: %v", err)
	}

	t.Logf("搜索结果: 总数=%d, 当前页数据量=%d", result.TotalCount, len(products))

	// 验证搜索结果
	for i, product := range products {
		t.Logf("[%d] %s - %s - ¥%.2f", i+1, product.Name, product.Category, product.Price)

		// 验证搜索条件
		if product.Category != "手机" {
			t.Errorf("搜索结果分类错误: 期望手机，实际: %s", product.Category)
		}
		if product.Price < 1000 || product.Price > 8000 {
			t.Errorf("搜索结果价格超出范围: %.2f", product.Price)
		}
		if !strings.Contains(product.Name, "iPhone") {
			t.Errorf("搜索结果名称不匹配: %s", product.Name)
		}
	}

	// 测试无效搜索请求
	invalidPageInfo := &PageInfoReq{
		Like: []string{"secret_key:test"}, // 不支持的字段
	}

	_, err = AutoSearchPaginated(db, &TestSearchProduct{}, &products, invalidPageInfo)
	if err == nil {
		t.Error("应该拒绝无效的搜索请求")
	} else {
		t.Logf("正确拒绝了无效请求: %v", err)
	}
}

// setupSearchTestDB 设置搜索测试数据库
func setupSearchTestDB(t *testing.T) *gorm.DB {
	// 使用内存SQLite数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&TestSearchProduct{})
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	// 插入测试数据
	testData := []TestSearchProduct{
		{Name: "iPhone 15", Category: "手机", Price: 5999.00, Stock: 100, Status: true, Tags: "苹果,智能手机", CreatedAt: "2024-01-01", SecretKey: "secret1"},
		{Name: "iPhone 14", Category: "手机", Price: 4999.00, Stock: 150, Status: true, Tags: "苹果,智能手机", CreatedAt: "2024-01-02", SecretKey: "secret2"},
		{Name: "MacBook Pro", Category: "笔记本", Price: 12999.00, Stock: 50, Status: true, Tags: "苹果,笔记本", CreatedAt: "2024-01-03", SecretKey: "secret3"},
		{Name: "小米13", Category: "手机", Price: 3999.00, Stock: 200, Status: true, Tags: "小米,智能手机", CreatedAt: "2024-01-04", SecretKey: "secret4"},
		{Name: "华为P60", Category: "手机", Price: 4999.00, Stock: 150, Status: false, Tags: "华为,智能手机", CreatedAt: "2024-01-05", SecretKey: "secret5"},
		{Name: "iPad Air", Category: "平板", Price: 4599.00, Stock: 120, Status: true, Tags: "苹果,平板", CreatedAt: "2024-01-06", SecretKey: "secret6"},
		{Name: "AirPods Pro", Category: "耳机", Price: 1999.00, Stock: 300, Status: true, Tags: "苹果,耳机", CreatedAt: "2024-01-07", SecretKey: "secret7"},
	}

	for _, product := range testData {
		if err := db.Create(&product).Error; err != nil {
			t.Fatalf("插入测试数据失败: %v", err)
		}
	}

	return db
}

// TestTagParsing 测试标签解析功能
func TestTagParsing(t *testing.T) {
	t.Log("=== 测试标签解析功能 ===")

	// 测试search标签解析
	operators := parseOperators("like,eq,gte,lte")
	expectedOps := []string{"like", "eq", "gte", "lte"}

	if len(operators) != len(expectedOps) {
		t.Errorf("操作符数量不匹配: 期望%d，实际%d", len(expectedOps), len(operators))
	}

	for i, expected := range expectedOps {
		if operators[i] != expected {
			t.Errorf("操作符[%d]不匹配: 期望%s，实际%s", i, expected, operators[i])
		}
	}

	// 测试widget标签解析
	field := TestSearchProduct{}
	modelType := reflect.TypeOf(field)
	categoryField, _ := modelType.FieldByName("Category")

	widgetConfig := parseWidgetConfig(categoryField)
	if widgetConfig.Type != "select" {
		t.Errorf("组件类型解析错误: 期望select，实际%s", widgetConfig.Type)
	}

	if len(widgetConfig.Options) != 4 {
		t.Errorf("选项数量解析错误: 期望4，实际%d", len(widgetConfig.Options))
	}

	t.Logf("标签解析测试通过")
}
