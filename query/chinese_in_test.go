package query

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 测试用的数据模型
type ChineseTestUser struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Name      string    `json:"name" gorm:"column:name"`
	Status    string    `json:"status" gorm:"column:status"`
	Category  string    `json:"category" gorm:"column:category"`
	City      string    `json:"city" gorm:"column:city"`
	Score     int       `json:"score" gorm:"column:score"`
}

func (ChineseTestUser) TableName() string {
	return "chinese_test_users"
}

// 创建测试数据库
func setupChineseTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启SQL日志
	})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&ChineseTestUser{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// 插入中文测试数据
	testUsers := []ChineseTestUser{
		{Name: "张三", Status: "活跃", Category: "高级", City: "北京", Score: 85},
		{Name: "李四", Status: "非活跃", Category: "中级", City: "上海", Score: 92},
		{Name: "王五", Status: "活跃", Category: "高级", City: "广州", Score: 78},
		{Name: "赵六", Status: "待审核", Category: "初级", City: "深圳", Score: 88},
		{Name: "钱七", Status: "活跃", Category: "中级", City: "杭州", Score: 95},
		{Name: "孙八", Status: "非活跃", Category: "高级", City: "成都", Score: 82},
		{Name: "周九", Status: "活跃", Category: "初级", City: "武汉", Score: 90},
		{Name: "吴十", Status: "待审核", Category: "中级", City: "西安", Score: 87},
	}

	for _, user := range testUsers {
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	return db
}

// 测试中文状态查询
func TestChineseStatusInQuery(t *testing.T) {
	db := setupChineseTestDB(t)

	// 测试 in=status:活跃,待审核
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Sorts:    "id:asc",
		In:       []string{"status:活跃,待审核"},
	}

	t.Logf("=== 测试中文状态查询 in=status:活跃,待审核 ===")
	t.Logf("PageInfo: %+v", pageInfo)

	// 使用 ApplySearchConditions 方法
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 查询数据
	var users []ChineseTestUser
	if err := dbWithConditions.Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	t.Logf("Found %d users:", len(users))
	for i, user := range users {
		t.Logf("  [%d] %s - %s - %s - %s - %d", i+1, user.Name, user.Status, user.Category, user.City, user.Score)
	}

	// 验证结果
	expectedCount := 6 // 张三, 王五, 赵六, 钱七, 周九, 吴十 (活跃 + 待审核)
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}

	// 验证所有结果都是 活跃 或 待审核
	for _, user := range users {
		if user.Status != "活跃" && user.Status != "待审核" {
			t.Errorf("User %s has unexpected status: %s", user.Name, user.Status)
		}
	}
}

// 测试中文分类查询
func TestChineseCategoryInQuery(t *testing.T) {
	db := setupChineseTestDB(t)

	// 测试 in=category:高级,中级
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Sorts:    "id:asc",
		In:       []string{"category:高级,中级"},
	}

	t.Logf("=== 测试中文分类查询 in=category:高级,中级 ===")
	t.Logf("PageInfo: %+v", pageInfo)

	// 使用 ApplySearchConditions 方法
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 查询数据
	var users []ChineseTestUser
	if err := dbWithConditions.Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	t.Logf("Found %d users:", len(users))
	for i, user := range users {
		t.Logf("  [%d] %s - %s - %s - %s - %d", i+1, user.Name, user.Status, user.Category, user.City, user.Score)
	}

	// 验证结果
	expectedCount := 6 // 张三, 李四, 王五, 钱七, 孙八, 吴十 (高级 + 中级)
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}

	// 验证所有结果都是 高级 或 中级
	for _, user := range users {
		if user.Category != "高级" && user.Category != "中级" {
			t.Errorf("User %s has unexpected category: %s", user.Name, user.Category)
		}
	}
}

// 测试中文城市查询
func TestChineseCityInQuery(t *testing.T) {
	db := setupChineseTestDB(t)

	// 测试 in=city:北京,上海,广州
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Sorts:    "id:asc",
		In:       []string{"city:北京,上海,广州"},
	}

	t.Logf("=== 测试中文城市查询 in=city:北京,上海,广州 ===")
	t.Logf("PageInfo: %+v", pageInfo)

	// 使用 ApplySearchConditions 方法
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 查询数据
	var users []ChineseTestUser
	if err := dbWithConditions.Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	t.Logf("Found %d users:", len(users))
	for i, user := range users {
		t.Logf("  [%d] %s - %s - %s - %s - %d", i+1, user.Name, user.Status, user.Category, user.City, user.Score)
	}

	// 验证结果
	expectedCount := 3 // 张三(北京), 李四(上海), 王五(广州)
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}

	// 验证所有结果都是指定的城市
	validCities := map[string]bool{"北京": true, "上海": true, "广州": true}
	for _, user := range users {
		if !validCities[user.City] {
			t.Errorf("User %s has unexpected city: %s", user.Name, user.City)
		}
	}
}

// 测试中文姓名查询
func TestChineseNameInQuery(t *testing.T) {
	db := setupChineseTestDB(t)

	// 测试 in=name:张三,李四,王五
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Sorts:    "id:asc",
		In:       []string{"name:张三,李四,王五"},
	}

	t.Logf("=== 测试中文姓名查询 in=name:张三,李四,王五 ===")
	t.Logf("PageInfo: %+v", pageInfo)

	// 使用 ApplySearchConditions 方法
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 查询数据
	var users []ChineseTestUser
	if err := dbWithConditions.Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	t.Logf("Found %d users:", len(users))
	for i, user := range users {
		t.Logf("  [%d] %s - %s - %s - %s - %d", i+1, user.Name, user.Status, user.Category, user.City, user.Score)
	}

	// 验证结果
	expectedCount := 3 // 张三, 李四, 王五
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}

	// 验证所有结果都是指定的姓名
	validNames := map[string]bool{"张三": true, "李四": true, "王五": true}
	for _, user := range users {
		if !validNames[user.Name] {
			t.Errorf("User %s is not in expected names", user.Name)
		}
	}
}

// 测试中文混合查询
func TestChineseMixedInQuery(t *testing.T) {
	db := setupChineseTestDB(t)

	// 测试多个中文条件
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Sorts:    "score:desc",
		In:       []string{"status:活跃", "category:高级,中级"},
	}

	t.Logf("=== 测试中文混合查询 ===")
	t.Logf("PageInfo: %+v", pageInfo)

	// 使用 ApplySearchConditions 方法
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 查询数据
	var users []ChineseTestUser
	if err := dbWithConditions.Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	t.Logf("Found %d users:", len(users))
	for i, user := range users {
		t.Logf("  [%d] %s - %s - %s - %s - %d", i+1, user.Name, user.Status, user.Category, user.City, user.Score)
	}

	// 验证结果 - 应该是 活跃 状态且 category 为 高级 或 中级 的用户
	expectedCount := 4 // 张三(活跃,高级), 王五(活跃,高级), 钱七(活跃,中级), 周九(活跃,初级) - 但周九是初级，所以应该是3个
	// 重新计算：张三(活跃,高级), 王五(活跃,高级), 钱七(活跃,中级) = 3个
	expectedCount = 3
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}

	// 验证所有结果都满足条件
	for _, user := range users {
		if user.Status != "活跃" {
			t.Errorf("User %s has unexpected status: %s", user.Name, user.Status)
		}
		if user.Category != "高级" && user.Category != "中级" {
			t.Errorf("User %s has unexpected category: %s", user.Name, user.Category)
		}
	}
}

// 测试中文完整查询流程
func TestChineseCompleteInQueryFlow(t *testing.T) {
	db := setupChineseTestDB(t)

	t.Logf("=== 测试中文完整查询流程 ===")

	// 模拟 URL: in=status:活跃,待审核&page=1&page_size=5&sorts=score:desc
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 5,
		Sorts:    "score:desc",
		In:       []string{"status:活跃,待审核"},
	}

	t.Logf("PageInfo: %+v", pageInfo)

	// 使用 SimplePaginate 方法
	var users []ChineseTestUser
	result, err := SimplePaginate(db, &ChineseTestUser{}, &users, pageInfo)
	if err != nil {
		t.Fatalf("SimplePaginate failed: %v", err)
	}

	t.Logf("Query result:")
	t.Logf("  TotalCount: %d", result.TotalCount)
	t.Logf("  CurrentPage: %d", result.CurrentPage)
	t.Logf("  PageSize: %d", result.PageSize)
	t.Logf("  TotalPages: %d", result.TotalPages)
	t.Logf("  Items count: %d", len(users))

	t.Logf("Found users:")
	for i, user := range users {
		t.Logf("  [%d] %s - %s - %s - %s - %d", i+1, user.Name, user.Status, user.Category, user.City, user.Score)
	}

	// 验证结果
	expectedCount := 6 // 总共有6个用户满足条件
	if result.TotalCount != int64(expectedCount) {
		t.Errorf("Expected total count %d but got %d", expectedCount, result.TotalCount)
	}

	// 验证分页结果
	expectedPageCount := 5 // 第一页显示5个用户
	if len(users) != expectedPageCount {
		t.Errorf("Expected %d users on first page but got %d", expectedPageCount, len(users))
	}

	// 验证排序（按分数降序）
	for i := 1; i < len(users); i++ {
		if users[i-1].Score < users[i].Score {
			t.Errorf("Sorting is incorrect: %d < %d", users[i-1].Score, users[i].Score)
		}
	}
}

// 测试中文不存在值查询
func TestChineseNonExistentInQuery(t *testing.T) {
	db := setupChineseTestDB(t)

	// 测试不存在的状态
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		In:       []string{"status:不存在"},
	}

	t.Logf("=== 测试中文不存在值查询 in=status:不存在 ===")
	t.Logf("PageInfo: %+v", pageInfo)

	// 使用 ApplySearchConditions 方法
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 查询数据
	var users []ChineseTestUser
	if err := dbWithConditions.Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	t.Logf("Found %d users:", len(users))

	// 验证结果
	expectedCount := 0
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}
}
