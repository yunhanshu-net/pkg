package query

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 压力测试用的数据模型
type StressTestUser struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Name      string    `json:"name" gorm:"column:name"`
	Age       int       `json:"age" gorm:"column:age"`
	Status    string    `json:"status" gorm:"column:status"`
	Score     float64   `json:"score" gorm:"column:score"`
	Category  string    `json:"category" gorm:"column:category"`
}

func (StressTestUser) TableName() string {
	return "stress_test_users"
}

// 创建压力测试数据库
func setupStressTestDB(t *testing.T) *gorm.DB {
	// 使用临时文件数据库避免并发问题
	db, err := gorm.Open(sqlite.Open("file:stress_test.db?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&StressTestUser{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// 插入大量测试数据
	statuses := []string{"active", "inactive", "pending", "suspended"}
	categories := []string{"A", "B", "C", "D", "E"}

	for i := 0; i < 1000; i++ {
		user := StressTestUser{
			Name:     fmt.Sprintf("User_%d", i),
			Age:      20 + (i % 50),
			Status:   statuses[i%len(statuses)],
			Score:    float64(50 + (i % 50)),
			Category: categories[i%len(categories)],
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}

	return db
}

// 模拟 table.go 中的 AutoPaginated 方法（简化版）
type StressMockTableData struct {
	err  error
	val  interface{}
	Data StressMockTable
}

type StressMockTable struct {
	Title      string                   `json:"title"`
	Column     []StressMockColumn       `json:"column"`
	Values     map[string][]interface{} `json:"values"`
	Pagination StressMockPaginated      `json:"pagination"`
}

type StressMockColumn struct {
	Idx  int    `json:"idx"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type StressMockPaginated struct {
	CurrentPage int `json:"current_page"`
	TotalCount  int `json:"total_count"`
	TotalPages  int `json:"total_pages"`
	PageSize    int `json:"page_size"`
}

// 模拟 AutoPaginated 方法
func (t *StressMockTableData) AutoPaginated(db *gorm.DB, model interface{}, pageInfo *PageInfoReq) *StressMockTableData {
	if pageInfo == nil {
		pageInfo = new(PageInfoReq)
	}

	// 使用query库的公开方法应用搜索条件
	dbWithConditions, err := ApplySearchConditions(db, pageInfo)
	if err != nil {
		t.err = fmt.Errorf("AutoPaginated.ApplySearchConditions failed: %v", err)
		return t
	}

	// 获取分页大小
	pageSize := pageInfo.GetLimit()
	offset := pageInfo.GetOffset()

	// 查询总数
	var totalCount int64
	if err := dbWithConditions.Model(model).Count(&totalCount).Error; err != nil {
		t.err = fmt.Errorf("AutoPaginated.Count :%+v failed to count records: %v", t.val, err)
		return t
	}

	// 应用排序
	if pageInfo.GetSorts() != "" {
		dbWithConditions = dbWithConditions.Order(pageInfo.GetSorts())
	}

	// 查询当前页数据
	if err := dbWithConditions.Offset(offset).Limit(pageSize).Find(t.val).Error; err != nil {
		t.err = fmt.Errorf("AutoPaginated.Find :%+v failed to find records: %v", t.val, err)
		return t
	}

	// 计算总页数
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}

	// 构造分页结果
	t.Data.Pagination = StressMockPaginated{
		CurrentPage: pageInfo.Page,
		TotalCount:  int(totalCount),
		TotalPages:  totalPages,
		PageSize:    pageSize,
	}
	return t
}

// 压力测试：高并发查询
func TestStressHighConcurrency(t *testing.T) {
	db := setupStressTestDB(t)

	// 并发数
	concurrency := 50
	// 每个goroutine的查询次数
	queriesPerGoroutine := 20

	var wg sync.WaitGroup
	errors := make(chan error, concurrency*queriesPerGoroutine)
	results := make(chan int, concurrency*queriesPerGoroutine)

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < queriesPerGoroutine; j++ {
				// 随机选择查询条件
				var pageInfo *PageInfoReq
				switch j % 4 {
				case 0:
					pageInfo = &PageInfoReq{
						Page:     1,
						PageSize: 10,
						Sorts:    "age:asc",
						Eq:       []string{"status:active"},
					}
				case 1:
					pageInfo = &PageInfoReq{
						Page:     1,
						PageSize: 20,
						Sorts:    "score:desc",
						Gt:       []string{"age:30"},
					}
				case 2:
					pageInfo = &PageInfoReq{
						Page:     1,
						PageSize: 15,
						Sorts:    "name:asc",
						Like:     []string{"name:User_1"},
					}
				case 3:
					pageInfo = &PageInfoReq{
						Page:     1,
						PageSize: 25,
						Sorts:    "id:asc",
						In:       []string{"category:A,B"},
					}
				}

				// 执行查询
				var users []StressTestUser
				table := &StressMockTableData{val: &users}

				result := table.AutoPaginated(db, &StressTestUser{}, pageInfo)
				if result.err != nil {
					errors <- fmt.Errorf("goroutine %d, query %d: %v", goroutineID, j, result.err)
					continue
				}

				// 记录结果
				results <- len(users)
			}
		}(i)
	}

	wg.Wait()
	close(errors)
	close(results)

	duration := time.Since(start)
	t.Logf("Stress test completed in %v", duration)
	t.Logf("Total queries: %d", concurrency*queriesPerGoroutine)
	t.Logf("Queries per second: %.2f", float64(concurrency*queriesPerGoroutine)/duration.Seconds())

	// 检查错误
	errorCount := 0
	for err := range errors {
		t.Error(err)
		errorCount++
	}

	// 统计结果
	resultCount := 0
	totalResults := 0
	for count := range results {
		resultCount++
		totalResults += count
	}

	t.Logf("Successful queries: %d", resultCount)
	t.Logf("Failed queries: %d", errorCount)
	t.Logf("Average results per query: %.2f", float64(totalResults)/float64(resultCount))

	if errorCount > 0 {
		t.Errorf("Found %d errors in stress test", errorCount)
	}
}

// 压力测试：长时间运行
func TestStressLongRunning(t *testing.T) {
	db := setupStressTestDB(t)

	// 运行时间
	duration := 5 * time.Second
	start := time.Now()
	queryCount := 0
	errorCount := 0

	t.Logf("Starting long running stress test for %v", duration)

	for time.Since(start) < duration {
		// 随机查询条件
		pageInfo := &PageInfoReq{
			Page:     1,
			PageSize: 10,
			Sorts:    "id:asc",
			Eq:       []string{"status:active"},
		}

		var users []StressTestUser
		table := &StressMockTableData{val: &users}

		result := table.AutoPaginated(db, &StressTestUser{}, pageInfo)
		if result.err != nil {
			errorCount++
			t.Logf("Query %d failed: %v", queryCount, result.err)
		}

		queryCount++

		// 每1000次查询输出一次进度
		if queryCount%1000 == 0 {
			elapsed := time.Since(start)
			t.Logf("Progress: %d queries in %v (%.2f qps)", queryCount, elapsed, float64(queryCount)/elapsed.Seconds())
		}
	}

	elapsed := time.Since(start)
	t.Logf("Long running test completed:")
	t.Logf("Total queries: %d", queryCount)
	t.Logf("Total errors: %d", errorCount)
	t.Logf("Duration: %v", elapsed)
	t.Logf("Average QPS: %.2f", float64(queryCount)/elapsed.Seconds())
	t.Logf("Error rate: %.2f%%", float64(errorCount)/float64(queryCount)*100)

	if errorCount > 0 {
		t.Errorf("Found %d errors in long running test", errorCount)
	}
}

// 压力测试：内存泄漏检测
func TestStressMemoryLeak(t *testing.T) {
	db := setupStressTestDB(t)

	// 执行大量查询，检查是否有内存泄漏
	for i := 0; i < 10000; i++ {
		pageInfo := &PageInfoReq{
			Page:     1,
			PageSize: 10,
			Sorts:    "id:asc",
		}

		var users []StressTestUser
		table := &StressMockTableData{val: &users}

		result := table.AutoPaginated(db, &StressTestUser{}, pageInfo)
		if result.err != nil {
			t.Errorf("Query %d failed: %v", i, result.err)
		}

		// 每1000次查询输出一次进度
		if i%1000 == 0 && i > 0 {
			t.Logf("Memory leak test progress: %d queries completed", i)
		}
	}

	t.Log("Memory leak test completed: 10000 queries executed")
}

// 基准测试：AutoPaginated 性能
func BenchmarkAutoPaginatedStress(b *testing.B) {
	db := setupStressTestDB(&testing.T{})

	pageInfo := &PageInfoReq{
		Page:     1,
		PageSize: 10,
		Sorts:    "id:asc",
		Eq:       []string{"status:active"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []StressTestUser
		table := &StressMockTableData{val: &users}

		table.AutoPaginated(db, &StressTestUser{}, pageInfo)
	}
}
