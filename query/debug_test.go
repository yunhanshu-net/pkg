package query

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 测试用的数据模型
type DebugUser struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Name      string    `json:"name" gorm:"column:name"`
	Age       int       `json:"age" gorm:"column:age"`
	Status    string    `json:"status" gorm:"column:status"`
	Score     float64   `json:"score" gorm:"column:score"`
}

func (DebugUser) TableName() string {
	return "debug_users"
}

// 创建调试数据库
func setupDebugDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开启SQL日志
	})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&DebugUser{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// 插入测试数据
	testUsers := []DebugUser{
		{Name: "Alice", Age: 25, Status: "active", Score: 85.5},
		{Name: "Bob", Age: 30, Status: "inactive", Score: 92.0},
		{Name: "Charlie", Age: 35, Status: "active", Score: 78.5},
		{Name: "David", Age: 28, Status: "pending", Score: 88.0},
		{Name: "Eve", Age: 32, Status: "active", Score: 95.5},
	}

	for _, user := range testUsers {
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	return db
}

// 调试测试：查看生成的SQL
func TestDebugSQL(t *testing.T) {
	db := setupDebugDB(t)

	// 测试 Like 查询
	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Like:     []string{"status:active"},
	}

	fmt.Println("=== 测试 Like 查询 ===")
	fmt.Printf("PageInfo: %+v\n", pageInfo)

	// 使用 ApplySearchConditions 方法
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	// 查询数据并打印SQL
	var users []DebugUser
	if err := dbWithConditions.Find(&users).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	fmt.Printf("Found %d users:\n", len(users))
	for _, user := range users {
		fmt.Printf("  - %s (status: %s)\n", user.Name, user.Status)
	}

	// 测试 Eq 查询
	pageInfo2 := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Eq:       []string{"status:active"},
	}

	fmt.Println("\n=== 测试 Eq 查询 ===")
	fmt.Printf("PageInfo: %+v\n", pageInfo2)

	dbWithConditions2, err := ApplySearchConditions(db, pageInfo2)
	if err != nil {
		t.Fatalf("ApplySearchConditions failed: %v", err)
	}

	var users2 []DebugUser
	if err := dbWithConditions2.Find(&users2).Error; err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	fmt.Printf("Found %d users:\n", len(users2))
	for _, user := range users2 {
		fmt.Printf("  - %s (status: %s)\n", user.Name, user.Status)
	}
}

// 调试测试：手动构建查询
func TestManualQuery(t *testing.T) {
	db := setupDebugDB(t)

	fmt.Println("=== 手动构建查询 ===")

	// 手动构建 Like 查询
	var users []DebugUser
	if err := db.Where("status LIKE ?", "%active%").Find(&users).Error; err != nil {
		t.Fatalf("Manual Like query failed: %v", err)
	}

	fmt.Printf("Manual Like query found %d users:\n", len(users))
	for _, user := range users {
		fmt.Printf("  - %s (status: %s)\n", user.Name, user.Status)
	}

	// 手动构建 Eq 查询
	var users2 []DebugUser
	if err := db.Where("status = ?", "active").Find(&users2).Error; err != nil {
		t.Fatalf("Manual Eq query failed: %v", err)
	}

	fmt.Printf("Manual Eq query found %d users:\n", len(users2))
	for _, user := range users2 {
		fmt.Printf("  - %s (status: %s)\n", user.Name, user.Status)
	}
}
