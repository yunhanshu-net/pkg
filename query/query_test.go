package query

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestProduct 测试用的产品模型
type TestProduct struct {
	ID          int     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string  `gorm:"column:name;comment:产品名称" json:"name"`
	Category    string  `gorm:"column:category;comment:产品分类" json:"category"`
	Price       float64 `gorm:"column:price;comment:产品价格" json:"price"`
	Stock       int     `gorm:"column:stock;comment:库存数量" json:"stock"`
	Description string  `gorm:"column:description;comment:产品描述" json:"description"`
	Status      string  `gorm:"column:status;comment:产品状态" json:"status"`
	Tags        string  `gorm:"column:tags;comment:产品标签" json:"tags"`
	CreatedBy   string  `gorm:"column:created_by;comment:创建人" json:"created_by"`
}

func (TestProduct) TableName() string {
	return "test_products"
}

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	// 使用内存SQLite数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启SQL日志
	})
	if err != nil {
		t.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&TestProduct{})
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	// 插入测试数据
	testData := []TestProduct{
		{Name: "iPhone 15", Category: "手机", Price: 5999.00, Stock: 100, Status: "启用", Tags: "苹果,智能手机", CreatedBy: "admin"},
		{Name: "MacBook Pro", Category: "笔记本", Price: 12999.00, Stock: 50, Status: "启用", Tags: "苹果,笔记本", CreatedBy: "admin"},
		{Name: "小米13", Category: "手机", Price: 3999.00, Stock: 200, Status: "启用", Tags: "小米,智能手机", CreatedBy: "user1"},
		{Name: "华为P60", Category: "手机", Price: 4999.00, Stock: 150, Status: "禁用", Tags: "华为,智能手机", CreatedBy: "user2"},
		{Name: "联想ThinkPad", Category: "笔记本", Price: 8999.00, Stock: 80, Status: "启用", Tags: "联想,商务", CreatedBy: "admin"},
		{Name: "iPad Air", Category: "平板", Price: 4599.00, Stock: 120, Status: "启用", Tags: "苹果,平板", CreatedBy: "user1"},
		{Name: "Surface Pro", Category: "平板", Price: 7999.00, Stock: 60, Status: "禁用", Tags: "微软,平板", CreatedBy: "user2"},
		{Name: "AirPods Pro", Category: "耳机", Price: 1999.00, Stock: 300, Status: "启用", Tags: "苹果,耳机", CreatedBy: "admin"},
	}

	for _, product := range testData {
		if err := db.Create(&product).Error; err != nil {
			t.Fatalf("插入测试数据失败: %v", err)
		}
	}

	return db
}

// TestBasicPagination 测试基础分页功能
func TestBasicPagination(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	t.Log("=== 测试基础分页功能 ===")

	// 测试第一页
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 3,
	}

	var products []TestProduct
	result, err := AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("分页查询失败: %v", err)
	}

	t.Logf("第一页结果: 当前页=%d, 总数=%d, 总页数=%d, 每页数量=%d",
		result.CurrentPage, result.TotalCount, result.TotalPages, result.PageSize)

	if result.CurrentPage != 1 {
		t.Errorf("期望当前页为1, 实际: %d", result.CurrentPage)
	}
	if result.TotalCount != 8 {
		t.Errorf("期望总数为8, 实际: %d", result.TotalCount)
	}
	if len(products) != 3 {
		t.Errorf("期望返回3条记录, 实际: %d", len(products))
	}

	// 打印第一页数据
	for i, product := range products {
		t.Logf("第一页数据[%d]: %s - %s - ¥%.2f", i+1, product.Name, product.Category, product.Price)
	}
}

// TestSorting 测试排序功能
func TestSorting(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	t.Log("=== 测试排序功能 ===")

	// 测试按价格降序排序
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 5,
		Sorts:    "price:desc",
	}

	var products []TestProduct
	_, err := AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("排序查询失败: %v", err)
	}

	t.Logf("按价格降序排序结果:")
	for i, product := range products {
		t.Logf("[%d] %s - ¥%.2f", i+1, product.Name, product.Price)
	}

	// 验证排序是否正确（价格应该是降序）
	for i := 1; i < len(products); i++ {
		if products[i-1].Price < products[i].Price {
			t.Errorf("排序错误: 第%d个产品价格(%.2f) < 第%d个产品价格(%.2f)",
				i, products[i-1].Price, i+1, products[i].Price)
		}
	}

	// 测试多字段排序
	pageInfo.Sorts = "category:asc,price:desc"
	products = []TestProduct{} // 重置切片
	_, err = AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("多字段排序查询失败: %v", err)
	}

	t.Logf("按分类升序、价格降序排序结果:")
	for i, product := range products {
		t.Logf("[%d] %s - %s - ¥%.2f", i+1, product.Name, product.Category, product.Price)
	}
}

// TestEqualQuery 测试等于查询
func TestEqualQuery(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	t.Log("=== 测试等于查询(eq) ===")

	// 测试单个等于条件
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Eq:       []string{"category:手机"},
	}

	var products []TestProduct
	result, err := AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("等于查询失败: %v", err)
	}

	t.Logf("分类=手机的产品数量: %d", result.TotalCount)
	for i, product := range products {
		t.Logf("[%d] %s - %s", i+1, product.Name, product.Category)
		if product.Category != "手机" {
			t.Errorf("查询结果错误: 期望分类为'手机', 实际: %s", product.Category)
		}
	}

	// 测试多个等于条件
	pageInfo.Eq = []string{"category:手机", "status:启用"}
	products = []TestProduct{} // 重置切片
	result, err = AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("多条件等于查询失败: %v", err)
	}

	t.Logf("分类=手机且状态=启用的产品数量: %d", result.TotalCount)
	for i, product := range products {
		t.Logf("[%d] %s - %s - %s", i+1, product.Name, product.Category, product.Status)
		if product.Category != "手机" || product.Status != "启用" {
			t.Errorf("查询结果错误: 期望分类='手机'且状态='启用', 实际: %s/%s",
				product.Category, product.Status)
		}
	}
}

// TestLikeQuery 测试模糊查询
func TestLikeQuery(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	t.Log("=== 测试模糊查询(like) ===")

	// 测试单个模糊查询
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Like:     []string{"name:iPhone"},
	}

	var products []TestProduct
	result, err := AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("模糊查询失败: %v", err)
	}

	t.Logf("名称包含'iPhone'的产品数量: %d", result.TotalCount)
	for i, product := range products {
		t.Logf("[%d] %s", i+1, product.Name)
	}

	// 测试多个模糊查询
	pageInfo.Like = []string{"tags:苹果", "name:Pro"}
	products = []TestProduct{} // 重置切片
	result, err = AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("多条件模糊查询失败: %v", err)
	}

	t.Logf("标签包含'苹果'且名称包含'Pro'的产品数量: %d", result.TotalCount)
	for i, product := range products {
		t.Logf("[%d] %s - %s", i+1, product.Name, product.Tags)
	}
}

// TestNumericComparison 测试数值比较查询
func TestNumericComparison(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	t.Log("=== 测试数值比较查询 ===")

	// 测试大于查询(gt)
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Gt:       []string{"price:5000"},
	}

	var products []TestProduct
	result, err := AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("大于查询失败: %v", err)
	}

	t.Logf("价格>5000的产品数量: %d", result.TotalCount)
	for i, product := range products {
		t.Logf("[%d] %s - ¥%.2f", i+1, product.Name, product.Price)
		if product.Price <= 5000 {
			t.Errorf("查询结果错误: 产品价格%.2f不大于5000", product.Price)
		}
	}

	// 测试大于等于查询(gte)
	pageInfo = &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Gte:      []string{"stock:100"},
	}

	products = []TestProduct{} // 重置切片
	result, err = AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("大于等于查询失败: %v", err)
	}

	t.Logf("库存>=100的产品数量: %d", result.TotalCount)
	for i, product := range products {
		t.Logf("[%d] %s - 库存:%d", i+1, product.Name, product.Stock)
		if product.Stock < 100 {
			t.Errorf("查询结果错误: 产品库存%d不大于等于100", product.Stock)
		}
	}

	// 测试小于查询(lt)
	pageInfo = &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Lt:       []string{"price:3000"},
	}

	products = []TestProduct{} // 重置切片
	result, err = AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("小于查询失败: %v", err)
	}

	t.Logf("价格<3000的产品数量: %d", result.TotalCount)
	for i, product := range products {
		t.Logf("[%d] %s - ¥%.2f", i+1, product.Name, product.Price)
		if product.Price >= 3000 {
			t.Errorf("查询结果错误: 产品价格%.2f不小于3000", product.Price)
		}
	}

	// 测试小于等于查询(lte)
	pageInfo = &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Lte:      []string{"stock:100"},
	}

	products = []TestProduct{} // 重置切片
	result, err = AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("小于等于查询失败: %v", err)
	}

	t.Logf("库存<=100的产品数量: %d", result.TotalCount)
	for i, product := range products {
		t.Logf("[%d] %s - 库存:%d", i+1, product.Name, product.Stock)
		if product.Stock > 100 {
			t.Errorf("查询结果错误: 产品库存%d不小于等于100", product.Stock)
		}
	}
}

// TestInQuery 测试IN查询
func TestInQuery(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	t.Log("=== 测试IN查询 ===")

	// 测试单字段IN查询
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		In:       []string{"category:手机,category:笔记本"},
	}

	var products []TestProduct
	result, err := AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("IN查询失败: %v", err)
	}

	t.Logf("分类在[手机,笔记本]中的产品数量: %d", result.TotalCount)
	for i, product := range products {
		t.Logf("[%d] %s - %s", i+1, product.Name, product.Category)
		if product.Category != "手机" && product.Category != "笔记本" {
			t.Errorf("查询结果错误: 产品分类'%s'不在[手机,笔记本]中", product.Category)
		}
	}

	// 测试多字段IN查询
	pageInfo.In = []string{"category:手机,category:笔记本", "status:启用"}
	products = []TestProduct{} // 重置切片
	result, err = AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("多字段IN查询失败: %v", err)
	}

	t.Logf("分类在[手机,笔记本]且状态=启用的产品数量: %d", result.TotalCount)
	for i, product := range products {
		t.Logf("[%d] %s - %s - %s", i+1, product.Name, product.Category, product.Status)
		if (product.Category != "手机" && product.Category != "笔记本") || product.Status != "启用" {
			t.Errorf("查询结果错误: 产品'%s'不符合条件", product.Name)
		}
	}
}

// TestComplexQuery 测试复合查询
func TestComplexQuery(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	t.Log("=== 测试复合查询 ===")

	// 测试多种条件组合
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Eq:       []string{"status:启用"},
		Like:     []string{"tags:苹果"},
		Gt:       []string{"price:2000"},
		Lt:       []string{"price:10000"},
		In:       []string{"category:手机,category:笔记本,category:平板"},
		Sorts:    "price:desc",
	}

	var products []TestProduct
	result, err := AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("复合查询失败: %v", err)
	}

	t.Logf("复合查询结果数量: %d", result.TotalCount)
	t.Logf("查询条件: 状态=启用 AND 标签包含'苹果' AND 价格>2000 AND 价格<10000 AND 分类在[手机,笔记本,平板]中")

	for i, product := range products {
		t.Logf("[%d] %s - %s - ¥%.2f - %s - %s",
			i+1, product.Name, product.Category, product.Price, product.Status, product.Tags)

		// 验证查询结果
		if product.Status != "启用" {
			t.Errorf("查询结果错误: 产品状态'%s'不等于'启用'", product.Status)
		}
		if product.Price <= 2000 || product.Price >= 10000 {
			t.Errorf("查询结果错误: 产品价格%.2f不在(2000,10000)范围内", product.Price)
		}
		if product.Category != "手机" && product.Category != "笔记本" && product.Category != "平板" {
			t.Errorf("查询结果错误: 产品分类'%s'不在允许范围内", product.Category)
		}
	}
}

// TestQueryConfig 测试查询配置
func TestQueryConfig(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	t.Log("=== 测试查询配置 ===")

	// 创建查询配置，只允许特定字段查询
	config := NewQueryConfig()
	config.AllowField("name", "like", "eq")
	config.AllowField("category", "eq", "in")
	config.AllowField("price", "gt", "gte", "lt", "lte")
	config.DenyField("description") // 禁止查询描述字段

	// 测试允许的查询
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Eq:       []string{"category:手机"},
		Like:     []string{"name:iPhone"},
		Gt:       []string{"price:3000"},
	}

	var products []TestProduct
	result, err := AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo, config)
	if err != nil {
		t.Fatalf("配置查询失败: %v", err)
	}

	t.Logf("配置查询成功，结果数量: %d", result.TotalCount)

	// 测试被禁止的字段查询
	pageInfo.Like = []string{"description:测试"}
	products = []TestProduct{} // 重置切片
	_, err = AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo, config)
	if err == nil {
		t.Error("期望查询被禁止字段时返回错误，但实际成功了")
	} else {
		t.Logf("正确拦截了被禁止的字段查询: %v", err)
	}

	// 测试不允许的操作符
	pageInfo.Like = []string{"name:测试"}
	pageInfo.Gt = []string{"category:测试"} // category不允许gt操作
	_, err = AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo, config)
	if err == nil {
		t.Error("期望使用不允许的操作符时返回错误，但实际成功了")
	} else {
		t.Logf("正确拦截了不允许的操作符: %v", err)
	}
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	t.Log("=== 测试边界情况 ===")

	// 测试空查询条件
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
	}

	var products []TestProduct
	result, err := AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("空查询条件失败: %v", err)
	}
	t.Logf("空查询条件结果数量: %d", result.TotalCount)

	// 测试无效的查询格式
	pageInfo = &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Eq:       []string{"invalid_format"}, // 缺少冒号
	}

	products = []TestProduct{} // 重置切片
	_, err = AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err == nil {
		t.Error("期望无效查询格式返回错误，但实际成功了")
	} else {
		t.Logf("正确拦截了无效查询格式: %v", err)
	}

	// 测试超大页码
	pageInfo = &PageInfoReq{
		Page:     9999,
		PageSize: 10,
	}

	products = []TestProduct{} // 重置切片
	result, err = AutoPaginateTable(ctx, db, &TestProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("超大页码查询失败: %v", err)
	}
	t.Logf("超大页码查询结果数量: %d", len(products))
	if len(products) != 0 {
		t.Error("超大页码应该返回空结果")
	}
}

// containsSubstring 辅助函数：检查字符串是否包含子串
func containsSubstring(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str == substr ||
			len(str) > len(substr) &&
				(str[:len(substr)] == substr ||
					str[len(str)-len(substr):] == substr ||
					findSubstring(str, substr)))
}

func findSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestPageInfoReq_WithSorts(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		add      string
		expected string
	}{
		{
			name:     "第一次设置排序",
			initial:  "",
			add:      "category:asc,price:desc",
			expected: "category:asc,price:desc",
		},
		{
			name:     "添加第二个排序条件",
			initial:  "category:asc",
			add:      "price:desc",
			expected: "category:asc,price:desc",
		},
		{
			name:     "添加多个排序条件",
			initial:  "category:asc",
			add:      "price:desc,created_at:asc",
			expected: "category:asc,price:desc,created_at:asc",
		},
		{
			name:     "使用减号前缀格式",
			initial:  "category:asc",
			add:      "-price,-created_at",
			expected: "category:asc,price:desc,created_at:desc",
		},
		{
			name:     "不覆盖已存在的排序条件",
			initial:  "category:asc,price:desc",
			add:      "category:desc,stock:desc",
			expected: "category:asc,price:desc,stock:desc",
		},
		{
			name:     "混合格式且不覆盖",
			initial:  "id:asc",
			add:      "id:desc,price:desc,-stock",
			expected: "id:asc,price:desc,stock:desc",
		},
		{
			name:     "前端排序优先级",
			initial:  "created_at:desc,price:asc",
			add:      "id:desc,created_at:asc",
			expected: "created_at:desc,price:asc,id:desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PageInfoReq{
				Sorts: tt.initial,
			}

			result := req.WithSorts(tt.add)

			if result.Sorts != tt.expected {
				t.Errorf("WithSorts() = %v, want %v", result.Sorts, tt.expected)
			}
		})
	}
}

func TestPageInfoReq_GetSorts(t *testing.T) {
	tests := []struct {
		name     string
		sorts    string
		expected string
	}{
		{
			name:     "单个排序条件",
			sorts:    "category:asc",
			expected: "category ASC",
		},
		{
			name:     "多个排序条件",
			sorts:    "category:asc,price:desc",
			expected: "category ASC, price DESC",
		},
		{
			name:     "带空格的排序条件",
			sorts:    "category : asc , price : desc",
			expected: "category ASC, price DESC",
		},
		{
			name:     "空排序条件",
			sorts:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &PageInfoReq{
				Sorts: tt.sorts,
			}

			result := req.GetSorts()

			if result != tt.expected {
				t.Errorf("GetSorts() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestBooleanProduct 测试布尔值查询的产品模型
type TestBooleanProduct struct {
	ID       int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"column:name;comment:产品名称" json:"name"`
	IsActive bool   `gorm:"column:is_active;comment:是否激活" json:"is_active"`
	IsPublic bool   `gorm:"column:is_public;comment:是否公开" json:"is_public"`
}

func (TestBooleanProduct) TableName() string {
	return "test_boolean_products"
}

// TestBooleanQuery 测试布尔值查询
func TestBooleanQuery(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	t.Log("=== 测试布尔值查询 ===")

	// 自动迁移
	err := db.AutoMigrate(&TestBooleanProduct{})
	if err != nil {
		t.Fatalf("数据库迁移失败: %v", err)
	}

	// 插入测试数据
	testData := []TestBooleanProduct{
		{Name: "产品A", IsActive: true, IsPublic: true},
		{Name: "产品B", IsActive: true, IsPublic: false},
		{Name: "产品C", IsActive: false, IsPublic: true},
		{Name: "产品D", IsActive: false, IsPublic: false},
	}

	for _, product := range testData {
		if err := db.Create(&product).Error; err != nil {
			t.Fatalf("插入测试数据失败: %v", err)
		}
	}

	// 测试布尔值等于查询
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Eq:       []string{"is_active:true"},
	}

	var products []TestBooleanProduct
	result, err := AutoPaginateTable(ctx, db, &TestBooleanProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("布尔值查询失败: %v", err)
	}

	t.Logf("is_active=true 的产品数量: %d", result.TotalCount)
	for i, product := range products {
		t.Logf("[%d] %s - IsActive: %v", i+1, product.Name, product.IsActive)
		if !product.IsActive {
			t.Errorf("查询结果错误: 期望 IsActive 为 true, 实际: %v", product.IsActive)
		}
	}

	// 测试布尔值不等于查询
	pageInfo.Eq = []string{"is_active:false"}
	products = []TestBooleanProduct{} // 重置切片
	result, err = AutoPaginateTable(ctx, db, &TestBooleanProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("布尔值不等于查询失败: %v", err)
	}

	t.Logf("is_active=false 的产品数量: %d", result.TotalCount)
	for i, product := range products {
		t.Logf("[%d] %s - IsActive: %v", i+1, product.Name, product.IsActive)
		if product.IsActive {
			t.Errorf("查询结果错误: 期望 IsActive 为 false, 实际: %v", product.IsActive)
		}
	}

	// 测试多个布尔条件
	pageInfo.Eq = []string{"is_active:true", "is_public:false"}
	products = []TestBooleanProduct{} // 重置切片
	result, err = AutoPaginateTable(ctx, db, &TestBooleanProduct{}, &products, pageInfo)
	if err != nil {
		t.Fatalf("多布尔条件查询失败: %v", err)
	}

	t.Logf("is_active=true 且 is_public=false 的产品数量: %d", result.TotalCount)
	for i, product := range products {
		t.Logf("[%d] %s - IsActive: %v, IsPublic: %v", i+1, product.Name, product.IsActive, product.IsPublic)
		if !product.IsActive || product.IsPublic {
			t.Errorf("查询结果错误: 期望 IsActive=true 且 IsPublic=false, 实际: IsActive=%v, IsPublic=%v", 
				product.IsActive, product.IsPublic)
		}
	}

	// 测试IN查询 - 布尔值
	t.Log("测试IN查询 - 布尔值")
	var inResults []TestBooleanProduct
	inPageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		In:       []string{"is_active:true,false"},
	}
	
	inResult, err := AutoPaginateTable(ctx, db, &TestBooleanProduct{}, &inResults, inPageInfo)
	if err != nil {
		t.Fatalf("IN查询失败: %v", err)
	}
	
	if inResult.TotalCount != 4 {
		t.Errorf("IN查询结果数量错误: 期望4，实际%d", inResult.TotalCount)
	}
	
	// 验证所有结果都包含在查询条件中
	for _, product := range inResults {
		if product.IsActive != true && product.IsActive != false {
			t.Errorf("IN查询结果包含无效的布尔值: %v", product.IsActive)
		}
	}
	
	// 测试NOT IN查询 - 布尔值
	t.Log("测试NOT IN查询 - 布尔值")
	var notInResults []TestBooleanProduct
	notInPageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		NotIn:    []string{"is_public:false"},
	}
	
	_, err = AutoPaginateTable(ctx, db, &TestBooleanProduct{}, &notInResults, notInPageInfo)
	if err != nil {
		t.Fatalf("NOT IN查询失败: %v", err)
	}
	
	// 验证所有结果都不等于false
	for _, product := range notInResults {
		if product.IsPublic == false {
			t.Errorf("NOT IN查询结果包含被排除的值: %v", product.IsPublic)
		}
	}

	t.Log("布尔值查询测试通过")
}
