package workflow

import (
	"testing"
)

func TestStepDescriptionExtraction(t *testing.T) {
	// 测试代码，包含带描述的步骤定义
	code := `
//desc: 用户验证步骤
step1 = validateUser(username: string "用户名", password: string "密码") -> (isValid: bool "验证结果", err: string "错误信息");

//desc: 发送欢迎邮件
step2 = sendWelcomeEmail(email: string "邮箱地址") -> (success: bool "发送结果");

func main() {
	//desc: 执行用户验证
	valid, err := step1("john", "password123")
	
	//desc: 如果验证成功，发送欢迎邮件
	if valid {
		//desc: 发送邮件
		success := step2("john@example.com")
	}
}
`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 检查步骤数量
	if len(result.Steps) != 2 {
		t.Fatalf("期望2个步骤，实际得到%d个", len(result.Steps))
	}

	// 检查第一个步骤的描述
	step1 := result.Steps[0]
	if step1.Name != "step1" {
		t.Errorf("步骤1名称错误: 期望 step1，实际 %s", step1.Name)
	}
	if step1.Desc != "用户验证步骤" {
		t.Errorf("步骤1描述错误: 期望 '用户验证步骤'，实际 '%s'", step1.Desc)
	}

	// 检查第二个步骤的描述
	step2 := result.Steps[1]
	if step2.Name != "step2" {
		t.Errorf("步骤2名称错误: 期望 step2，实际 %s", step2.Name)
	}
	if step2.Desc != "发送欢迎邮件" {
		t.Errorf("步骤2描述错误: 期望 '发送欢迎邮件'，实际 '%s'", step2.Desc)
	}

	// 检查主函数中的语句描述
	if result.MainFunc == nil || len(result.MainFunc.Statements) == 0 {
		t.Fatal("主函数为空")
	}

	// 检查第一个语句的描述
	stmt1 := result.MainFunc.Statements[0]
	if stmt1.Desc != "执行用户验证" {
		t.Errorf("语句1描述错误: 期望 '执行用户验证'，实际 '%s'", stmt1.Desc)
	}

	// 检查if语句的描述
	ifStmt := result.MainFunc.Statements[1]
	if ifStmt.Type != "if" {
		t.Errorf("期望if语句，实际得到 %s", ifStmt.Type)
	}
	if ifStmt.Desc != "如果验证成功，发送欢迎邮件" {
		t.Errorf("if语句描述错误: 期望 '如果验证成功，发送欢迎邮件'，实际 '%s'", ifStmt.Desc)
	}

	// 检查if语句内部的描述
	if len(ifStmt.Children) == 0 {
		t.Fatal("if语句内部为空")
	}
	innerStmt := ifStmt.Children[0]
	if innerStmt.Desc != "发送邮件" {
		t.Errorf("if内部语句描述错误: 期望 '发送邮件'，实际 '%s'", innerStmt.Desc)
	}

	t.Log("✅ 所有描述信息提取测试通过")
}
