package workflow

import (
	"testing"
)

// 测试步骤描述功能
func TestStepDescription(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(
    username: string "用户名",
    phone: int "手机号"
) -> (
    workId: string "工号",
    username: string "用户名",
    err: error "是否失败"
);

step2 = beiluo.test1.crm.crm_interview_schedule(
    username: string "用户名"
) -> (
    interviewTime: string "面试时间",
    interviewer: string "面试官名称",
    err: error "是否失败"
);

func main() {
    //desc: 开始用户注册流程
    fmt.Println("开始用户注册和面试安排流程...")
    
    //desc: 创建用户账号，获取工号
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"]){retry:3, timeout:5000, priority:"high"}
    
    //desc: 检查用户创建是否成功
    if step1Err != nil {
        //desc: 用户创建失败，记录错误并退出
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    
    //desc: 用户创建成功，记录日志
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
    
    //desc: 安排面试时间，联系面试官
    面试时间, 面试官名称, step2Err := step2(用户名){retry:2, timeout:3000, priority:"normal"}
    
    //desc: 检查面试安排是否成功
    if step2Err != nil {
        //desc: 面试安排失败，记录错误并退出
        step2.Printf("安排面试失败: %v", step2Err)
        return
    }
    
    //desc: 面试安排成功，记录详细信息
    step2.Printf("✅ 面试安排成功，时间: %s，面试官: %s", 面试时间, 面试官名称)
    
    //desc: 流程完成，输出最终结果
    fmt.Printf("流程完成，工号: %s，面试时间: %s\n", 工号, 面试时间)
}`

	// 解析工作流
	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 验证步骤数量
	if len(result.MainFunc.Statements) < 8 {
		t.Fatalf("期望至少8个语句，实际%d个", len(result.MainFunc.Statements))
	}

	// 验证描述信息
	expectedDescs := map[string]bool{
		"开始用户注册流程":       true,
		"创建用户账号，获取工号":    true,
		"检查用户创建是否成功":     true,
		"用户创建失败，记录错误并退出": true,
		"用户创建成功，记录日志":    true,
		"安排面试时间，联系面试官":   true,
		"检查面试安排是否成功":     true,
		"面试安排失败，记录错误并退出": true,
		"面试安排成功，记录详细信息":  true,
		"流程完成，输出最终结果":    true,
	}

	foundDescs := make(map[string]bool)

	// 递归收集所有语句的描述信息
	var collectDescs func(stmts []*SimpleStatement)
	collectDescs = func(stmts []*SimpleStatement) {
		for _, stmt := range stmts {
			if stmt.Desc != "" {
				foundDescs[stmt.Desc] = true
			}
			// 递归处理子语句
			if len(stmt.Children) > 0 {
				collectDescs(stmt.Children)
			}
		}
	}

	collectDescs(result.MainFunc.Statements)

	// 验证所有期望的描述都被找到
	for expectedDesc := range expectedDescs {
		if !foundDescs[expectedDesc] {
			t.Errorf("未找到期望的描述: %s", expectedDesc)
		}
	}

	// 验证没有意外的描述
	for foundDesc := range foundDescs {
		if !expectedDescs[foundDesc] {
			t.Errorf("意外的描述: %s", foundDesc)
		}
	}

	// 打印所有语句的描述信息
	t.Log("=== 语句描述信息 ===")
	for i, stmt := range result.MainFunc.Statements {
		if stmt.Desc != "" {
			t.Logf("语句 %d: %s - %s", i+1, stmt.Type, stmt.Desc)
		} else {
			t.Logf("语句 %d: %s - (无描述)", i+1, stmt.Type)
		}
	}
}

// 测试嵌套语句的描述功能
func TestNestedStatementDescription(t *testing.T) {
	code := `func main() {
    //desc: 开始订单处理流程
    fmt.Println("开始订单处理流程...")
    
    //desc: 验证订单信息
    验证结果, step1Err := step1(input["订单号"], input["金额"]){retry:2, timeout:3000}
    
    //desc: 检查订单验证结果
    if step1Err != nil {
        //desc: 订单验证失败，记录错误
        step1.Printf("订单验证失败: %v", step1Err)
        return
    }
    
    //desc: 根据验证结果决定后续流程
    if 验证结果 {
        //desc: 订单验证通过，开始处理支付
        fmt.Println("订单验证通过，开始处理支付...")
        
        //desc: 处理支付流程
        支付流水号, step2Err := step2(input["订单号"], input["金额"]){retry:3, timeout:5000, priority:"high"}
        
        //desc: 检查支付是否成功
        if step2Err != nil {
            //desc: 支付失败，记录错误
            step2.Printf("支付失败: %v", step2Err)
            return
        }
        
        //desc: 支付成功，记录流水号
        fmt.Printf("订单处理完成，支付流水号: %s\n", 支付流水号)
    } else {
        //desc: 订单验证失败，流程结束
        step1.Printf("订单验证失败")
        fmt.Println("订单验证失败，流程结束")
    }
}`

	// 解析工作流
	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 验证描述信息
	t.Log("=== 嵌套语句描述信息 ===")
	for i, stmt := range result.MainFunc.Statements {
		if stmt.Desc != "" {
			t.Logf("语句 %d: %s - %s", i+1, stmt.Type, stmt.Desc)
		} else {
			t.Logf("语句 %d: %s - (无描述)", i+1, stmt.Type)
		}

		// 检查嵌套语句的描述
		if len(stmt.Children) > 0 {
			t.Logf("  └─ 嵌套语句:")
			for j, child := range stmt.Children {
				if child.Desc != "" {
					t.Logf("     %d. %s - %s", j+1, child.Type, child.Desc)
				} else {
					t.Logf("     %d. %s - (无描述)", j+1, child.Type)
				}
			}
		}
	}
}
