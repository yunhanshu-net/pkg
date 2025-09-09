package workflow

//
//import (
//	"strings"
//	"testing"
//)
//
//// 测试步骤级别日志记录
//func TestStepLogging(t *testing.T) {
//	code := `var input = map[string]interface{}{
//    "用户名": "张三",
//    "手机号": 13800138000,
//}
//
//step1 = beiluo.test1.devops.devops_script_create(
//    username: string "用户名",
//    phone: int "手机号"
//) -> (
//    workId: string "工号",
//    username: string "用户名",
//    err: error "是否失败"
//);
//
//step2 = beiluo.test1.crm.crm_interview_schedule(
//    username: string "用户名"
//) -> (
//    interviewTime: string "面试时间",
//    interviewer: string "面试官名称",
//    err: error "是否失败"
//);
//
//func main() {
//    //desc: 开始用户注册和面试安排流程
//    fmt.Println("开始用户注册和面试安排流程...")
//
//    //desc: 创建用户账号，获取工号
//    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"]){retry:3, timeout:5000, priority:"high"}
//
//    //desc: 检查用户创建是否成功
//    if step1Err != nil {
//        //desc: 用户创建失败，记录错误并退出
//        step1.Printf("创建用户失败: %v", step1Err)
//        return
//    }
//
//    //desc: 用户创建成功，记录成功日志
//    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
//
//    //desc: 安排面试时间，联系面试官
//    面试时间, 面试官名称, step2Err := step2(用户名){retry:2, timeout:3000, priority:"normal"}
//
//    //desc: 检查面试安排是否成功
//    if step2Err != nil {
//        //desc: 面试安排失败，记录错误并退出
//        step2.Printf("安排面试失败: %v", step2Err)
//        return
//    }
//
//    //desc: 面试安排成功，记录详细信息
//    step2.Printf("✅ 面试安排成功，时间: %s，面试官: %s", 面试时间, 面试官名称)
//
//    fmt.Printf("流程完成，工号: %s，面试时间: %s\n", 工号, 面试时间)
//}`
//
//	// 解析工作流
//	parser := NewSimpleParser()
//	result := parser.ParseWorkflow(code)
//
//	if !result.Success {
//		t.Fatalf("解析失败: %s", result.Error)
//	}
//
//	// 执行工作流
//	executor := NewWorkflowExecutor()
//	execResult, updatedParseResult := executor.ExecuteWorkflowWithResult(code)
//
//	if !execResult.Success {
//		t.Fatalf("执行失败: %s", execResult.Error)
//	}
//
//	// 验证步骤状态
//	if len(result.MainFunc.Statements) < 8 {
//		t.Fatalf("期望至少8个语句，实际%d个", len(result.MainFunc.Statements))
//	}
//
//	// 检查步骤日志（使用更新后的解析结果）
//	var step1Logs, step2Logs int
//	for _, step := range updatedParseResult.Steps {
//		if step.Name == "step1" {
//			step1Logs = len(step.Logs)
//		} else if step.Name == "step2" {
//			step2Logs = len(step.Logs)
//		}
//	}
//
//	// 检查是否有步骤级别的日志记录（通过内容判断）
//	hasStep1Log := false
//	hasStep2Log := false
//	for _, stmt := range result.MainFunc.Statements {
//		if stmt.Type == "print" && strings.Contains(stmt.Content, "step1.Printf") {
//			hasStep1Log = true
//		}
//		if stmt.Type == "print" && strings.Contains(stmt.Content, "step2.Printf") {
//			hasStep2Log = true
//		}
//	}
//
//	if !hasStep1Log {
//		t.Error("没有找到 step1.Printf 日志记录")
//	}
//	if !hasStep2Log {
//		t.Error("没有找到 step2.Printf 日志记录")
//	}
//
//	t.Logf("step1 日志数量: %d", step1Logs)
//	t.Logf("step2 日志数量: %d", step2Logs)
//	t.Logf("找到 step1.Printf: %v", hasStep1Log)
//	t.Logf("找到 step2.Printf: %v", hasStep2Log)
//
//	// 验证步骤日志是否正确添加
//	if step1Logs == 0 {
//		t.Error("step1 没有日志记录")
//	}
//	if step2Logs == 0 {
//		t.Error("step2 没有日志记录")
//	}
//
//	// 验证日志内容
//	for _, step := range updatedParseResult.Steps {
//		if step.Name == "step1" {
//			t.Logf("step1 日志内容:")
//			for i, log := range step.Logs {
//				t.Logf("  [%d] %s: %s", i, log.Source, log.Message)
//			}
//		} else if step.Name == "step2" {
//			t.Logf("step2 日志内容:")
//			for i, log := range step.Logs {
//				t.Logf("  [%d] %s: %s", i, log.Source, log.Message)
//			}
//		}
//	}
//}
//
//// 测试语句状态管理
//func TestStatementStatus(t *testing.T) {
//	code := `var input = map[string]interface{}{
//    "订单号": "ORD-001",
//    "金额": 1000,
//}
//
//step1 = beiluo.test1.order.validate_order(string 订单号, int 金额) -> (bool 验证结果, err 是否失败);
//
//func main() {
//    fmt.Println("开始订单处理流程...")
//
//    验证结果, step1Err := step1(input["订单号"], input["金额"]){retry:2, timeout:3000}
//    if step1Err != nil {
//        step1.Printf("订单验证失败: %v", step1Err)
//        return
//    }
//
//    if 验证结果 {
//        step1.Printf("订单验证通过")
//    } else {
//        step1.Printf("订单验证失败")
//    }
//}`
//
//	// 解析工作流
//	parser := NewSimpleParser()
//	result := parser.ParseWorkflow(code)
//
//	if !result.Success {
//		t.Fatalf("解析失败: %s", result.Error)
//	}
//
//	// 执行工作流
//	executor := NewWorkflowExecutor()
//	execResult := executor.ExecuteWorkflow(code)
//
//	if !execResult.Success {
//		t.Fatalf("执行失败: %s", execResult.Error)
//	}
//
//	// 验证语句状态 - 检查执行结果中的语句状态
//	for _, stmt := range result.MainFunc.Statements {
//		// 只检查实际执行的语句
//		if stmt.Status == StatusPending && stmt.Type != "other" {
//			t.Logf("语句 %s 状态为 pending", stmt.Content)
//		}
//	}
//
//	// 检查执行结果
//	if !execResult.Success {
//		t.Errorf("执行失败: %s", execResult.Error)
//	}
//
//	t.Logf("执行结果: 成功=%v, 步骤数=%d", execResult.Success, len(execResult.Steps))
//	t.Log("语句状态验证完成")
//}
//
//// 测试重试次数管理
//func TestRetryCount(t *testing.T) {
//	stmt := &SimpleStatement{
//		Type:       "function-call",
//		Content:    "test()",
//		Status:     StatusPending,
//		RetryCount: 0,
//	}
//
//	// 测试重试次数增加
//	if stmt.RetryCount != 0 {
//		t.Errorf("期望重试次数为0，实际为%d", stmt.RetryCount)
//	}
//
//	stmt.IncrementRetry()
//	if stmt.RetryCount != 1 {
//		t.Errorf("期望重试次数为1，实际为%d", stmt.RetryCount)
//	}
//
//	stmt.IncrementRetry()
//	if stmt.RetryCount != 2 {
//		t.Errorf("期望重试次数为2，实际为%d", stmt.RetryCount)
//	}
//
//	// 测试重试次数重置
//	stmt.ResetRetry()
//	if stmt.RetryCount != 0 {
//		t.Errorf("期望重试次数为0，实际为%d", stmt.RetryCount)
//	}
//
//	t.Log("重试次数管理测试通过")
//}
