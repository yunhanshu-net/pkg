package query

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 测试用的数据模型
type InTestUser struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Name      string    `json:"name" gorm:"column:name"`
	Status    string    `json:"status" gorm:"column:status"`
	Category  string    `json:"category" gorm:"column:category"`
	Score     int       `json:"score" gorm:"column:score"`
}

func (InTestUser) TableName() string {
	return "in_test_users"
}

// 创建测试数据库
func setupInTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启SQL日志
	})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&InTestUser{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// 插入测试数据
	testUsers := []InTestUser{
		{Name: "Alice", Status: "active", Category: "A", Score: 85},
		{Name: "Bob", Status: "inactive", Category: "B", Score: 92},
		{Name: "Charlie", Status: "active", Category: "A", Score: 78},
		{Name: "David", Status: "pending", Category: "C", Score: 88},
		{Name: "Eve", Status: "active", Category: "B", Score: 95},
		{Name: "Frank", Status: "inactive", Category: "A", Score: 82},
		{Name: "Grace", Status: "active", Category: "C", Score: 90},
		{Name: "Henry", Status: "pending", Category: "B", Score: 87},
	}

	for _, user := range testUsers {
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	return db
}

// 测试 in 查询条件
func TestInQueryCondition(t *testing.T) {
	db := setupInTestDB(t)

	// 测试 in=status:active,pending
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Sorts:    "id:asc",
		In:       []string{"status:active,pending"},
	}

	t.Logf("=== 测试 in=status:active,pending ===")
	t.Logf("PageInfo: %+v", pageInfo)

	// 使用 ApplySearchConditions 方法
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 查询数据
	var users []InTestUser
	if err := dbWithConditions.Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	t.Logf("Found %d users:", len(users))
	for i, user := range users {
		t.Logf("  [%d] %s - %s - %s - %d", i+1, user.Name, user.Status, user.Category, user.Score)
	}

	// 验证结果
	expectedCount := 6 // Alice, Charlie, David, Eve, Grace, Henry (active + pending)
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}

	// 验证所有结果都是 active 或 pending
	for _, user := range users {
		if user.Status != "active" && user.Status != "pending" {
			t.Errorf("User %s has unexpected status: %s", user.Name, user.Status)
		}
	}
}

// 测试不同的 in 查询条件
func TestVariousInConditions(t *testing.T) {
	db := setupInTestDB(t)

	testCases := []struct {
		name     string
		in       []string
		expected int
		check    func(user InTestUser) bool
	}{
		{
			name:     "status:active,pending",
			in:       []string{"status:active,pending"},
			expected: 6,
			check: func(user InTestUser) bool {
				return user.Status == "active" || user.Status == "pending"
			},
		},
		{
			name:     "category:A,B",
			in:       []string{"category:A,B"},
			expected: 6,
			check: func(user InTestUser) bool {
				return user.Category == "A" || user.Category == "B"
			},
		},
		{
			name:     "score:85,92,95",
			in:       []string{"score:85,92,95"},
			expected: 3,
			check: func(user InTestUser) bool {
				return user.Score == 85 || user.Score == 92 || user.Score == 95
			},
		},
		{
			name:     "name:Alice,Bob,Charlie",
			in:       []string{"name:Alice,Bob,Charlie"},
			expected: 3,
			check: func(user InTestUser) bool {
				return user.Name == "Alice" || user.Name == "Bob" || user.Name == "Charlie"
			},
		},
		{
			name:     "status:nonexistent",
			in:       []string{"status:nonexistent"},
			expected: 0,
			check: func(user InTestUser) bool {
				return false // 不应该有任何用户
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pageInfo := &PageInfoReq{
				Page:     1,
				PageSize: 10,
				In:       tc.in,
			}

			t.Logf("=== 测试 %s ===", tc.name)
			t.Logf("PageInfo: %+v", pageInfo)

			// 使用 ApplySearchConditions 方法
			dbWithConditions, err := ApplySearchConditions(db, pageInfo)
			if err != nil {
				t.Fatalf("ApplySearchConditions failed: %v", err)
			}

			// 查询数据
			var users []InTestUser
			if err := dbWithConditions.Find(&users).Error; err != nil {
				t.Fatalf("Find failed: %v", err)
			}

			t.Logf("Found %d users:", len(users))
			for i, user := range users {
				t.Logf("  [%d] %s - %s - %s - %d", i+1, user.Name, user.Status, user.Category, user.Score)
			}

			if len(users) != tc.expected {
				t.Errorf("Expected %d users but got %d", tc.expected, len(users))
			}

			// 验证每个用户都满足条件
			for _, user := range users {
				if !tc.check(user) {
					t.Errorf("User %s does not satisfy the condition", user.Name)
				}
			}
		})
	}
}

// 测试多个 in 条件
func TestMultipleInConditions(t *testing.T) {
	db := setupInTestDB(t)

	// 测试多个 in 条件
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		In:       []string{"status:active", "category:A,B"},
	}

	t.Logf("=== 测试多个 in 条件 ===")
	t.Logf("PageInfo: %+v", pageInfo)

	// 使用 ApplySearchConditions 方法
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 查询数据
	var users []InTestUser
	if err := dbWithConditions.Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	t.Logf("Found %d users:", len(users))
	for i, user := range users {
		t.Logf("  [%d] %s - %s - %s - %d", i+1, user.Name, user.Status, user.Category, user.Score)
	}

	// 验证结果 - 应该是 active 状态且 category 为 A 或 B 的用户
	expectedCount := 3 // Alice (active, A), Charlie (active, A), Eve (active, B)
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}

	// 验证所有结果都满足条件
	for _, user := range users {
		if user.Status != "active" {
			t.Errorf("User %s has unexpected status: %s", user.Name, user.Status)
		}
		if user.Category != "A" && user.Category != "B" {
			t.Errorf("User %s has unexpected category: %s", user.Name, user.Category)
		}
	}
}

// 测试手动构建的 in 查询
func TestManualInQuery(t *testing.T) {
	db := setupInTestDB(t)

	t.Logf("=== 测试手动构建的 in 查询 ===")

	// 手动构建查询
	var users []InTestUser
	if err := db.Where("status IN ?", []string{"active", "pending"}).Find(&users).Error; err != nil {
		t.Fatalf("Manual query failed: %v", err)
	}

	t.Logf("Manual query found %d users:", len(users))
	for i, user := range users {
		t.Logf("  [%d] %s - %s - %s - %d", i+1, user.Name, user.Status, user.Category, user.Score)
	}

	// 验证结果
	expectedCount := 6 // 实际有6个用户满足条件
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}
}

// 测试数字类型的 in 查询
func TestNumericInQuery(t *testing.T) {
	db := setupInTestDB(t)

	// 测试数字类型的 in 查询
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		In:       []string{"score:85,92,95"},
	}

	t.Logf("=== 测试数字类型的 in 查询 ===")
	t.Logf("PageInfo: %+v", pageInfo)

	// 使用 ApplySearchConditions 方法
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 查询数据
	var users []InTestUser
	if err := dbWithConditions.Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	t.Logf("Found %d users:", len(users))
	for i, user := range users {
		t.Logf("  [%d] %s - %s - %s - %d", i+1, user.Name, user.Status, user.Category, user.Score)
	}

	// 验证结果
	expectedCount := 3 // Alice (85), Bob (92), Eve (95)
	if len(users) != expectedCount {
		t.Errorf("Expected %d users but got %d", expectedCount, len(users))
	}

	// 验证所有结果都是指定的分数
	validScores := map[int]bool{85: true, 92: true, 95: true}
	for _, user := range users {
		if !validScores[user.Score] {
			t.Errorf("User %s has unexpected score: %d", user.Name, user.Score)
		}
	}
}

// 测试完整的查询流程
func TestCompleteInQueryFlow(t *testing.T) {
	db := setupInTestDB(t)

	t.Logf("=== 测试完整的 in 查询流程 ===")

	// 模拟 URL: in=status:active,pending&page=1&page_size=5&sorts=score:desc
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 5,
		Sorts:    "score:desc",
		In:       []string{"status:active,pending"},
	}

	t.Logf("PageInfo: %+v", pageInfo)

	// 使用 SimplePaginate 方法
	var users []InTestUser
	result, err := SimplePaginate(db, &InTestUser{}, &users, pageInfo)
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
		t.Logf("  [%d] %s - %s - %s - %d", i+1, user.Name, user.Status, user.Category, user.Score)
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
