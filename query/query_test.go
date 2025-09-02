package query

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 测试用的数据模型
type TestUser struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Name      string    `json:"name" gorm:"column:name"`
	Age       int       `json:"age" gorm:"column:age"`
	Status    string    `json:"status" gorm:"column:status"`
	Score     float64   `json:"score" gorm:"column:score"`
}

func (TestUser) TableName() string {
	return "test_users"
}

// 创建测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&TestUser{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// 插入测试数据
	testUsers := []TestUser{
		{Name: "Alice", Age: 25, Status: "active", Score: 85.5},
		{Name: "Bob", Age: 30, Status: "inactive", Score: 92.0},
		{Name: "Charlie", Age: 35, Status: "active", Score: 78.5},
		{Name: "David", Age: 28, Status: "pending", Score: 88.0},
		{Name: "Eve", Age: 32, Status: "active", Score: 95.5},
		{Name: "Frank", Age: 27, Status: "inactive", Score: 82.0},
		{Name: "Grace", Age: 29, Status: "active", Score: 90.0},
		{Name: "Henry", Age: 31, Status: "pending", Score: 87.5},
		{Name: "Ivy", Age: 26, Status: "active", Score: 93.0},
		{Name: "Jack", Age: 33, Status: "inactive", Score: 79.5},
	}

	for _, user := range testUsers {
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	return db
}

// 测试 ApplySearchConditions 方法的一致性
func TestApplySearchConditionsConsistency(t *testing.T) {
	db := setupTestDB(t)

	// 测试多次查询的一致性
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("ConsistencyTest_%d", i), func(t *testing.T) {
			// 创建分页请求
			pageInfo := &PageInfoReq{
				Page:     1,
				PageSize: 5,
				Sorts:    "age:asc",
				Eq:       []string{"status:active"}, // 使用精确匹配而不是模糊匹配
			}

			// 使用 ApplySearchConditions 方法
			dbWithConditions, err := ApplySearchConditions(db, pageInfo)
			if err != nil {
				t.Fatalf("ApplySearchConditions failed: %v", err)
			}

			// 查询总数
			var totalCount int64
			if err := dbWithConditions.Model(&TestUser{}).Count(&totalCount).Error; err != nil {
				t.Fatalf("Count failed: %v", err)
			}

			// 应用排序
			if pageInfo.GetSorts() != "" {
				dbWithConditions = dbWithConditions.Order(pageInfo.GetSorts())
			}

			// 查询数据
			var users []TestUser
			pageSize := pageInfo.GetLimit()
			offset := pageInfo.GetOffset()
			if err := dbWithConditions.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
				t.Fatalf("Find failed: %v", err)
			}

			// 检查结果
			if len(users) == 0 {
				t.Error("Expected users but got empty slice")
			}

			// 验证排序是否正确（按年龄升序）
			for j := 1; j < len(users); j++ {
				if users[j-1].Age > users[j].Age {
					t.Errorf("Sorting is incorrect: %d > %d", users[j-1].Age, users[j].Age)
				}
			}

			// 验证所有用户状态都是 active
			for _, user := range users {
				if user.Status != "active" {
					t.Errorf("Expected status 'active' but got '%s'", user.Status)
				}
			}

			t.Logf("Test %d: Found %d users, total count: %d", i, len(users), totalCount)
		})
	}
}

// 测试 AutoPaginateTable 方法
func TestAutoPaginateTable(t *testing.T) {
	db := setupTestDB(t)

	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 3,
		Sorts:    "score:desc",
		Eq:       []string{"status:active"},
	}

	var users []TestUser
	result, err := AutoPaginateTable(context.Background(), db, &TestUser{}, &users, pageInfo)
	if err != nil {
		t.Fatalf("AutoPaginateTable failed: %v", err)
	}

	// 验证结果
	if result == nil {
		t.Fatal("Result is nil")
	}

	if len(users) == 0 {
		t.Error("Expected users but got empty slice")
	}

	// 验证排序（按分数降序）
	for i := 1; i < len(users); i++ {
		if users[i-1].Score < users[i].Score {
			t.Errorf("Sorting is incorrect: %f < %f", users[i-1].Score, users[i].Score)
		}
	}

	// 验证状态过滤
	for _, user := range users {
		if user.Status != "active" {
			t.Errorf("Expected status 'active' but got '%s'", user.Status)
		}
	}

	t.Logf("Found %d users, total count: %d", len(users), result.TotalCount)
}

// 测试并发查询的一致性
func TestConcurrentQueries(t *testing.T) {
	// 为每个goroutine创建独立的数据库连接
	done := make(chan bool, 10)
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// 每个goroutine创建自己的数据库
			db := setupTestDB(t)

			pageInfo := &PageInfoReq{
				Page:     1,
				PageSize: 5,
				Sorts:    "name:asc",
				Like:     []string{"name:a"}, // 查找名字包含 'a' 的用户
			}

			// 使用 ApplySearchConditions 方法
			dbWithConditions, err := ApplySearchConditions(db, pageInfo)
			if err != nil {
				errors <- fmt.Errorf("goroutine %d: ApplySearchConditions failed: %v", id, err)
				return
			}

			// 查询总数
			var totalCount int64
			if err := dbWithConditions.Model(&TestUser{}).Count(&totalCount).Error; err != nil {
				errors <- fmt.Errorf("goroutine %d: Count failed: %v", id, err)
				return
			}

			// 应用排序
			if pageInfo.GetSorts() != "" {
				dbWithConditions = dbWithConditions.Order(pageInfo.GetSorts())
			}

			// 查询数据
			var users []TestUser
			pageSize := pageInfo.GetLimit()
			offset := pageInfo.GetOffset()
			if err := dbWithConditions.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
				errors <- fmt.Errorf("goroutine %d: Find failed: %v", id, err)
				return
			}

			// 验证结果不为空
			if len(users) == 0 {
				errors <- fmt.Errorf("goroutine %d: expected users but got empty slice", id)
				return
			}

			// 验证排序
			for j := 1; j < len(users); j++ {
				if users[j-1].Name > users[j].Name {
					errors <- fmt.Errorf("goroutine %d: sorting is incorrect", id)
					return
				}
			}

			t.Logf("Goroutine %d: Found %d users", id, len(users))
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 检查错误
	close(errors)
	for err := range errors {
		if err != nil {
			t.Error(err)
		}
	}
}

// 测试复杂查询条件
func TestComplexSearchConditions(t *testing.T) {
	db := setupTestDB(t)

	// 复杂查询：年龄大于25，状态为active，分数大于80，按分数降序
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Sorts:    "score:desc",
		Gt:       []string{"age:25", "score:80"},
		Eq:       []string{"status:active"},
	}

	// 使用 ApplySearchConditions 方法
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 查询总数
	var totalCount int64
	if err := dbWithConditions.Model(&TestUser{}).Count(&totalCount).Error; err != nil {
		t.Fatalf("Count failed: %v", err)
	}

	// 应用排序
	if pageInfo.GetSorts() != "" {
		dbWithConditions = dbWithConditions.Order(pageInfo.GetSorts())
	}

	// 查询数据
	var users []TestUser
	pageSize := pageInfo.GetLimit()
	offset := pageInfo.GetOffset()
	if err := dbWithConditions.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	// 验证所有条件都满足
	for _, user := range users {
		if user.Age <= 25 {
			t.Errorf("Expected age > 25 but got %d", user.Age)
		}
		if user.Status != "active" {
			t.Errorf("Expected status 'active' but got '%s'", user.Status)
		}
		if user.Score <= 80 {
			t.Errorf("Expected score > 80 but got %f", user.Score)
		}
	}

	// 验证排序
	for i := 1; i < len(users); i++ {
		if users[i-1].Score < users[i].Score {
			t.Errorf("Sorting is incorrect: %f < %f", users[i-1].Score, users[i].Score)
		}
	}

	t.Logf("Complex query found %d users, total count: %d", len(users), totalCount)
}

// 测试分页边界情况
func TestPaginationBoundaries(t *testing.T) {
	db := setupTestDB(t)

	// 测试第一页
	pageInfo1 := &PageInfoReq{
		Page:     1,
		PageSize: 3,
		Sorts:    "id:asc",
	}

	// 使用 ApplySearchConditions 方法
	dbWithConditions1, err := ApplySearchConditions(db, pageInfo1)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 应用排序
	if pageInfo1.GetSorts() != "" {
		dbWithConditions1 = dbWithConditions1.Order(pageInfo1.GetSorts())
	}

	// 查询第一页数据
	var users1 []TestUser
	pageSize1 := pageInfo1.GetLimit()
	offset1 := pageInfo1.GetOffset()
	if err := dbWithConditions1.Offset(offset1).Limit(pageSize1).Find(&users1).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	// 测试第二页
	pageInfo2 := &PageInfoReq{
		Page:     2,
		PageSize: 3,
		Sorts:    "id:asc",
	}

	// 使用 ApplySearchConditions 方法
	dbWithConditions2, err := ApplySearchConditions(db, pageInfo2)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 应用排序
	if pageInfo2.GetSorts() != "" {
		dbWithConditions2 = dbWithConditions2.Order(pageInfo2.GetSorts())
	}

	// 查询第二页数据
	var users2 []TestUser
	pageSize2 := pageInfo2.GetLimit()
	offset2 := pageInfo2.GetOffset()
	if err := dbWithConditions2.Offset(offset2).Limit(pageSize2).Find(&users2).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	// 验证分页结果不重复
	if len(users1) > 0 && len(users2) > 0 {
		// 检查是否有重复的ID
		idMap := make(map[int]bool)
		for _, user := range users1 {
			idMap[user.ID] = true
		}
		for _, user := range users2 {
			if idMap[user.ID] {
				t.Errorf("Found duplicate ID %d in different pages", user.ID)
			}
		}
	}

	t.Logf("Page 1: %d users, Page 2: %d users", len(users1), len(users2))
}

// 测试空结果
func TestEmptyResults(t *testing.T) {
	db := setupTestDB(t)

	// 查询不存在的条件
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Eq:       []string{"status:nonexistent"},
	}

	// 使用 ApplySearchConditions 方法
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 查询数据
	var users []TestUser
	pageSize := pageInfo.GetLimit()
	offset := pageInfo.GetOffset()
	if err := dbWithConditions.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	// 应该返回空结果
	if len(users) != 0 {
		t.Errorf("Expected empty result but got %d users", len(users))
	}

	t.Log("Empty result test passed")
}

// 基准测试
func BenchmarkApplySearchConditions(b *testing.B) {
	db := setupTestDB(&testing.T{})

	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Sorts:    "id:asc",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 使用 ApplySearchConditions 方法
		dbWithConditions, err := ApplySearchConditions(db, pageInfo)
		if err != nil {
			b.Fatalf("ApplySearchConditions failed: %v", err)
		}

		// 应用排序
		if pageInfo.GetSorts() != "" {
			dbWithConditions = dbWithConditions.Order(pageInfo.GetSorts())
		}

		// 查询数据
		var users []TestUser
		pageSize := pageInfo.GetLimit()
		offset := pageInfo.GetOffset()
		if err := dbWithConditions.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
			b.Fatalf("Find failed: %v", err)
		}
	}
}
