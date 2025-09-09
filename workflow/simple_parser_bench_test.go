package workflow

import (
	"testing"
)

// 简单工作流 - 用于基础性能测试
func getSimpleWorkflow() string {
	return `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(string 用户名, int 手机号) -> (string 工号, string 用户名, err 是否失败);

func main() {
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    fmt.Printf("用户创建成功，工号: %s\n", 工号)
}`
}

// 中等复杂度工作流 - 用于中等性能测试
func getMediumWorkflow() string {
	return `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
    "部门": "技术部",
}

step1 = beiluo.test1.devops.devops_script_create(string 用户名, int 手机号, string 邮箱) -> (string 工号, string 用户名, string 部门, err 是否失败);
step2 = beiluo.test1.crm.crm_interview_schedule(string 用户名, string 部门) -> (string 面试时间, string 面试官名称, err 是否失败);
step3 = beiluo.test1.notification.send_email(string 邮箱, string 内容) -> (err 是否失败);

func main() {
    // 创建用户
    工号, 用户名, 部门, step1Err := step1(input["用户名"], input["手机号"], input["邮箱"]){retry:3, timeout:5000, priority:"high"}
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    
    // 安排面试
    面试时间, 面试官名称, step2Err := step2(用户名, 部门){retry:2, timeout:3000, priority:"normal"}
    if step2Err != nil {
        step2.Printf("安排面试失败: %v", step2Err)
        return
    }
    
    // 发送通知
    通知内容 := "你收到了:{{用户名}},时间：{{面试时间}}的面试安排，请关注"
    step3Err := step3(input["邮箱"], 通知内容){retry:1, timeout:2000, priority:"low"}
    if step3Err != nil {
        step3.Printf("发送通知失败: %v", step3Err)
        return
    }
    
    fmt.Printf("流程完成，工号: %s，面试时间: %s\n", 工号, 面试时间)
}`
}

// 复杂工作流 - 用于高负载性能测试
func getComplexWorkflow() string {
	return `var input = map[string]interface{}{
    "订单号": "ORD-2024-001",
    "客户姓名": "李四",
    "客户电话": 13900139000,
    "收货地址": "北京市朝阳区xxx街道xxx号",
    "商品列表": []string{"商品A", "商品B", "商品C"},
    "支付方式": "微信支付",
    "配送方式": "标准配送",
}

// 静态工作流步骤
step1 = beiluo.test1.order.validate_order[订单验证用例] -> (bool 验证结果, string 验证信息, err 是否失败);
step2 = beiluo.test1.order.create_order[订单创建用例] -> (string 订单ID, string 订单状态, err 是否失败);
step3 = beiluo.test1.payment.process_payment[支付处理用例] -> (string 支付流水号, string 支付状态, err 是否失败);

// 动态工作流步骤
step4 = beiluo.test1.inventory.check_stock(string 商品列表) -> (bool 库存充足, string 缺货商品, err 是否失败);
step5 = beiluo.test1.inventory.reserve_stock(string 订单ID, string 商品列表) -> (string 预留单号, err 是否失败);
step6 = beiluo.test1.warehouse.pick_goods(string 订单ID, string 预留单号) -> (string 拣货单号, string 已拣商品, err 是否失败);
step7 = beiluo.test1.logistics.arrange_delivery(string 订单ID, string 收货地址, string 配送方式) -> (string 物流单号, string 预计送达时间, err 是否失败);
step8 = beiluo.test1.notification.send_order_notification(string 客户电话, string 订单号, string 物流单号) -> (err 是否失败);

func main() {
    fmt.Println("🛒 开始订单处理和物流配送流程...")
    
    // 1. 验证订单信息
    验证结果, 验证信息, step1Err := step1()
    if step1Err != nil {
        step1.Printf("订单验证失败: %v", step1Err)
        return
    }
    if !验证结果 {
        fmt.Printf("订单验证不通过: %s\n", 验证信息)
        return
    }
    
    // 2. 创建订单
    订单ID, 订单状态, step2Err := step2()
    if step2Err != nil {
        step2.Printf("订单创建失败: %v", step2Err)
        return
    }
    
    // 3. 检查库存
    库存充足, 缺货商品, step4Err := step4(input["商品列表"]){retry:2, timeout:5000, priority:"high"}
    if step4Err != nil {
        step4.Printf("库存检查失败: %v", step4Err)
        return
    }
    if !库存充足 {
        fmt.Printf("库存不足，缺货商品: %v\n", 缺货商品)
        return
    }
    
    // 4. 预留库存
    预留单号, step5Err := step5(订单ID, input["商品列表"]){retry:3, timeout:8000, priority:"critical"}
    if step5Err != nil {
        fmt.Printf("库存预留失败: %v\n", step5Err)
        return
    }
    
    // 5. 处理支付
    支付流水号, 支付状态, step3Err := step3()
    if step3Err != nil {
        fmt.Printf("支付处理失败: %v\n", step3Err)
        return
    }
    
    // 6. 拣货
    拣货单号, 已拣商品, step6Err := step6(订单ID, 预留单号){retry:2, timeout:10000, priority:"high"}
    if step6Err != nil {
        fmt.Printf("拣货失败: %v\n", step6Err)
        return
    }
    
    // 7. 安排配送
    物流单号, 预计送达时间, step7Err := step7(订单ID, input["收货地址"], input["配送方式"]){retry:1, timeout:5000, priority:"normal"}
    if step7Err != nil {
        fmt.Printf("配送安排失败: %v\n", step7Err)
        return
    }
    
    // 8. 发送通知
    step8Err := step8(input["客户电话"], input["订单号"], 物流单号){retry:1, timeout:3000, priority:"low"}
    if step8Err != nil {
        fmt.Printf("通知发送失败: %v\n", step8Err)
        return
    }
    
    fmt.Printf("订单处理完成，订单ID: %s，物流单号: %s\n", 订单ID, 物流单号)
}`
}

// 超大型工作流 - 用于极限性能测试
func getLargeWorkflow() string {
	workflow := `var input = map[string]interface{}{
    "项目名称": "大型项目",
    "版本号": "v2.0.0",
    "环境": "production",
    "团队规模": 50,
    "模块数量": 20,
}

// 定义20个步骤
`

	// 添加20个步骤定义
	for i := 1; i <= 20; i++ {
		workflow += `step` + string(rune('0'+i)) + ` = beiluo.test1.module.module` + string(rune('0'+i)) + `(string 项目名称, string 版本号, int 模块编号) -> (string 结果, err 是否失败);
`
	}

	workflow += `
func main() {
    fmt.Println("🚀 开始大型项目处理流程...")
    
    var 结果列表 []string
    var 错误列表 []string
    
`

	// 添加20个步骤调用
	for i := 1; i <= 20; i++ {
		workflow += `    // 步骤` + string(rune('0'+i)) + `
    结果` + string(rune('0'+i)) + `, step` + string(rune('0'+i)) + `Err := step` + string(rune('0'+i)) + `(input["项目名称"], input["版本号"], ` + string(rune('0'+i)) + `){retry:2, timeout:3000, priority:"normal"}
    if step` + string(rune('0'+i)) + `Err != nil {
        错误列表 = append(错误列表, "步骤` + string(rune('0'+i)) + `失败: " + step` + string(rune('0'+i)) + `Err.Error())
    } else {
        结果列表 = append(结果列表, 结果` + string(rune('0'+i)) + `)
    }
    
`
	}

	workflow += `
    fmt.Printf("处理完成，成功: %d，失败: %d\n", len(结果列表), len(错误列表))
    if len(错误列表) > 0 {
        fmt.Printf("错误详情: %v\n", 错误列表)
    }
}`

	return workflow
}

// Benchmark简单工作流解析
func BenchmarkSimpleParser_SimpleWorkflow(b *testing.B) {
	code := getSimpleWorkflow()
	parser := NewSimpleParser()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := parser.ParseWorkflow(code)
		if !result.Success {
			b.Fatalf("解析失败: %s", result.Error)
		}
	}
}

// Benchmark中等复杂度工作流解析
func BenchmarkSimpleParser_MediumWorkflow(b *testing.B) {
	code := getMediumWorkflow()
	parser := NewSimpleParser()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := parser.ParseWorkflow(code)
		if !result.Success {
			b.Fatalf("解析失败: %s", result.Error)
		}
	}
}

// Benchmark复杂工作流解析
func BenchmarkSimpleParser_ComplexWorkflow(b *testing.B) {
	code := getComplexWorkflow()
	parser := NewSimpleParser()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := parser.ParseWorkflow(code)
		if !result.Success {
			b.Fatalf("解析失败: %s", result.Error)
		}
	}
}

// Benchmark超大型工作流解析
func BenchmarkSimpleParser_LargeWorkflow(b *testing.B) {
	code := getLargeWorkflow()
	parser := NewSimpleParser()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := parser.ParseWorkflow(code)
		if !result.Success {
			b.Fatalf("解析失败: %s", result.Error)
		}
	}
}

// Benchmark内存分配测试
func BenchmarkSimpleParser_MemoryAllocation(b *testing.B) {
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
	}
}

// Benchmark并发解析测试
func BenchmarkSimpleParser_Concurrent(b *testing.B) {
	code := getMediumWorkflow()
	parser := NewSimpleParser()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			result := parser.ParseWorkflow(code)
			if !result.Success {
				b.Fatalf("解析失败: %s", result.Error)
			}
		}
	})
}

// Benchmark不同工作流大小对比
func BenchmarkSimpleParser_SizeComparison(b *testing.B) {
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
			for i := 0; i < b.N; i++ {
				result := parser.ParseWorkflow(code)
				if !result.Success {
					b.Fatalf("解析失败: %s", result.Error)
				}
			}
		})
	}
}
