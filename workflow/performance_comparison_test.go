package workflow

import (
	"testing"
)

// 性能对比测试 - 展示指针优化的效果
func TestPerformanceComparison(t *testing.T) {
	// 测试不同大小的工作流
	testCases := []struct {
		name     string
		workflow func() string
	}{
		{"Simple", getSimpleWorkflow},
		{"Medium", getMediumWorkflow},
		{"Complex", getComplexWorkflow},
		{"Large", getLargeWorkflow},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code := tc.workflow()
			parser := NewSimpleParser()

			// 预热
			for i := 0; i < 100; i++ {
				result := parser.ParseWorkflow(code)
				if !result.Success {
					t.Fatalf("解析失败: %s", result.Error)
				}
			}

			// 性能测试
			iterations := 1000
			for i := 0; i < iterations; i++ {
				result := parser.ParseWorkflow(code)
				if !result.Success {
					t.Fatalf("解析失败: %s", result.Error)
				}

				// 模拟使用解析结果
				_ = len(result.MainFunc.Statements)
				_ = len(result.Steps)
				_ = len(result.Variables)

				// 模拟访问嵌套结构
				for _, stmt := range result.MainFunc.Statements {
					_ = len(stmt.Children)
					_ = len(stmt.Args)
					_ = len(stmt.Returns)
					_ = len(stmt.Metadata)
				}
			}
		})
	}
}

// 内存使用分析
func TestMemoryUsage(t *testing.T) {
	code := getComplexWorkflow()
	parser := NewSimpleParser()

	// 解析工作流
	result := parser.ParseWorkflow(code)
	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 分析内存使用
	t.Logf("工作流步骤数量: %d", len(result.Steps))
	t.Logf("主函数语句数量: %d", len(result.MainFunc.Statements))
	t.Logf("变量映射数量: %d", len(result.Variables))

	// 统计嵌套结构
	totalStatements := 0
	totalArgs := 0
	totalReturns := 0

	var countStructures func(statements []*SimpleStatement)
	countStructures = func(statements []*SimpleStatement) {
		for _, stmt := range statements {
			totalStatements++
			totalArgs += len(stmt.Args)
			totalReturns += len(stmt.Returns)

			if len(stmt.Children) > 0 {
				countStructures(stmt.Children)
			}
		}
	}

	countStructures(result.MainFunc.Statements)

	t.Logf("总语句数量: %d", totalStatements)
	t.Logf("总参数数量: %d", totalArgs)
	t.Logf("总返回值数量: %d", totalReturns)
}

// 并发性能测试
func TestConcurrentPerformance(t *testing.T) {
	code := getMediumWorkflow()
	parser := NewSimpleParser()

	// 并发解析测试
	concurrency := 10
	iterations := 100

	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			for j := 0; j < iterations; j++ {
				result := parser.ParseWorkflow(code)
				if !result.Success {
					t.Errorf("解析失败: %s", result.Error)
					done <- false
					return
				}
			}
			done <- true
		}()
	}

	// 等待所有goroutine完成
	successCount := 0
	for i := 0; i < concurrency; i++ {
		if <-done {
			successCount++
		}
	}

	if successCount != concurrency {
		t.Errorf("并发测试失败，成功: %d/%d", successCount, concurrency)
	}
}

// 压力测试 - 连续解析大量工作流
func TestStressTest(t *testing.T) {
	code := getComplexWorkflow()
	parser := NewSimpleParser()

	// 连续解析1000次
	for i := 0; i < 1000; i++ {
		result := parser.ParseWorkflow(code)
		if !result.Success {
			t.Fatalf("第%d次解析失败: %s", i+1, result.Error)
		}

		// 模拟使用解析结果
		_ = len(result.MainFunc.Statements)
		_ = len(result.Steps)
		_ = len(result.Variables)
	}
}

// 基准测试 - 不同工作流大小的性能对比
func BenchmarkWorkflowSizes(b *testing.B) {
	tests := []struct {
		name     string
		workflow func() string
	}{
		{"Simple", getSimpleWorkflow},
		{"Medium", getMediumWorkflow},
		{"Complex", getComplexWorkflow},
		{"Large", getLargeWorkflow},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			code := tt.workflow()
			parser := NewSimpleParser()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				result := parser.ParseWorkflow(code)
				if !result.Success {
					b.Fatalf("解析失败: %s", result.Error)
				}
			}
		})
	}
}

// 内存分配分析
func BenchmarkMemoryAllocation(b *testing.B) {
	code := getComplexWorkflow()
	parser := NewSimpleParser()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result := parser.ParseWorkflow(code)
		if !result.Success {
			b.Fatalf("解析失败: %s", result.Error)
		}

		// 模拟使用解析结果
		_ = len(result.MainFunc.Statements)
		_ = len(result.Steps)
		_ = len(result.Variables)

		// 模拟访问嵌套结构
		for _, stmt := range result.MainFunc.Statements {
			_ = len(stmt.Children)
			_ = len(stmt.Args)
			_ = len(stmt.Returns)
			_ = len(stmt.Metadata)
		}
	}
}
