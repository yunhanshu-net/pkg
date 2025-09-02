package query

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 测试用的数据模型
type LikeDebugUser struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Name      string    `json:"name" gorm:"column:name"`
	Phone     string    `json:"phone" gorm:"column:phone"`
}

func (LikeDebugUser) TableName() string {
	return "like_debug_users"
}

// 创建测试数据库
func setupLikeDebugDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启SQL日志
	})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&LikeDebugUser{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// 插入测试数据
	testUsers := []LikeDebugUser{
		{Name: "Alice", Phone: "13812345678"},
		{Name: "Bob", Phone: "13987654321"},
		{Name: "Charlie", Phone: "13827641234"},
		{Name: "David", Phone: "13727645678"},
	}

	for _, user := range testUsers {
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	return db
}

// 测试 parseFieldValues 函数
func TestParseFieldValues(t *testing.T) {
	t.Logf("=== 测试 parseFieldValues 函数 ===")

	// 测试单个条件
	input := "phone:2764"
	conditions, err := parseFieldValues(input)
	if err != nil {
		t.Fatalf("parseFieldValues failed: %v", err)
	}

	t.Logf("Input: %s", input)
	t.Logf("Parsed conditions: %+v", conditions)

	if len(conditions) != 1 {
		t.Errorf("Expected 1 condition but got %d", len(conditions))
	}

	if field, ok := conditions["phone"]; !ok || field != "2764" {
		t.Errorf("Expected phone=2764 but got %+v", conditions)
	}
}

// 测试 validateAndBuildCondition 函数
func TestValidateAndBuildCondition(t *testing.T) {
	db := setupLikeDebugDB(t)

	t.Logf("=== 测试 validateAndBuildCondition 函数 ===")

	// 测试 like 条件
	inputs := []string{"phone:2764"}
	operator := "like"

	t.Logf("Inputs: %+v", inputs)
	t.Logf("Operator: %s", operator)

	// 调用 validateAndBuildCondition
	var dbPtr *gorm.DB = db
	err := validateAndBuildCondition(&dbPtr, inputs, operator, nil)
	if err != nil {
		t.Fatalf("validateAndBuildCondition failed: %v", err)
	}

	// 查询数据
	var users []LikeDebugUser
	if err := dbPtr.Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	t.Logf("Found %d users:", len(users))
	for i, user := range users {
		t.Logf("  [%d] %s - %s", i+1, user.Name, user.Phone)
	}

	// 验证结果
	expectedCount := 2 // Charlie 和 David 包含 2764
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}
}

// 测试 buildWhereConditionsWithoutConfig 函数
func TestBuildWhereConditionsWithoutConfig(t *testing.T) {
	db := setupLikeDebugDB(t)

	t.Logf("=== 测试 buildWhereConditionsWithoutConfig 函数 ===")

	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Like:     []string{"phone:2764"},
	}

	t.Logf("PageInfo: %+v", pageInfo)

	// 调用 buildWhereConditionsWithoutConfig
	var dbPtr *gorm.DB = db
	err := buildWhereConditionsWithoutConfig(&dbPtr, pageInfo)
	if err != nil {
		t.Fatalf("buildWhereConditionsWithoutConfig failed: %v", err)
	}

	// 查询数据
	var users []LikeDebugUser
	if err := dbPtr.Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	t.Logf("Found %d users:", len(users))
	for i, user := range users {
		t.Logf("  [%d] %s - %s", i+1, user.Name, user.Phone)
	}

	// 验证结果
	expectedCount := 2
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}
}

// 测试 ApplySearchConditions 函数
func TestApplySearchConditions(t *testing.T) {
	db := setupLikeDebugDB(t)

	t.Logf("=== 测试 ApplySearchConditions 函数 ===")

	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Like:     []string{"phone:2764"},
	}

	t.Logf("PageInfo: %+v", pageInfo)

	// 调用 ApplySearchConditions
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 查询数据
	var users []LikeDebugUser
	if err := dbWithConditions.Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	t.Logf("Found %d users:", len(users))
	for i, user := range users {
		t.Logf("  [%d] %s - %s", i+1, user.Name, user.Phone)
	}

	// 验证结果
	expectedCount := 2
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}
}

// 测试手动构建的查询
func TestManualLikeQuery(t *testing.T) {
	db := setupLikeDebugDB(t)

	t.Logf("=== 测试手动构建的 like 查询 ===")

	// 手动构建查询
	var users []LikeDebugUser
	if err := db.Where("phone LIKE ?", "%2764%").Find(&users).Error; err != nil {
		t.Fatalf("Manual query failed: %v", err)
	}

	t.Logf("Manual query found %d users:", len(users))
	for i, user := range users {
		t.Logf("  [%d] %s - %s", i+1, user.Name, user.Phone)
	}

	// 验证结果
	expectedCount := 2
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}
}

// 测试不同的 like 条件
func TestDifferentLikeConditions(t *testing.T) {
	db := setupLikeDebugDB(t)

	testCases := []struct {
		name     string
		like     []string
		expected int
	}{
		{"phone:2764", []string{"phone:2764"}, 2},
		{"phone:138", []string{"phone:138"}, 2},
		{"name:Alice", []string{"name:Alice"}, 1},
		{"phone:999", []string{"phone:999"}, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pageInfo := &PageInfoReq{
				Page:     1,
				PageSize: 10,
				Like:     tc.like,
			}

			t.Logf("=== 测试 %s ===", tc.name)
			t.Logf("PageInfo: %+v", pageInfo)

			// 使用 ApplySearchConditions 方法
			dbWithConditions, err := ApplySearchConditions(db, pageInfo)
			if err != nil {
				t.Fatalf("ApplySearchConditions failed: %v", err)
			}

			// 查询数据
			var users []LikeDebugUser
			if err := dbWithConditions.Find(&users).Error; err != nil {
				t.Fatalf("Find failed: %v", err)
			}

			t.Logf("Found %d users:", len(users))
			for i, user := range users {
				t.Logf("  [%d] %s - %s", i+1, user.Name, user.Phone)
			}

			if len(users) != tc.expected {
				t.Errorf("Expected %d users but got %d", tc.expected, len(users))
			}
		})
	}
}
