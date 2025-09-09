package workflow

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// WorkflowState 工作流状态存储结构
type WorkflowState struct {
	FlowID    string    `json:"flow_id"`
	Status    string    `json:"status"` // pending, running, completed, failed
	Data      string    `json:"data"`   // SimpleParseResult的JSON序列化
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SQLitePersistence SQLite持久化实现
type SQLitePersistence struct {
	db *sql.DB
}

// NewSQLitePersistence 创建SQLite持久化实例
func NewSQLitePersistence(dbPath string) (*SQLitePersistence, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// 创建表
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS workflow_states (
		flow_id TEXT PRIMARY KEY,
		status TEXT NOT NULL,
		data TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(createTableSQL); err != nil {
		return nil, err
	}

	return &SQLitePersistence{db: db}, nil
}

// SaveWorkflow 保存工作流状态
func (p *SQLitePersistence) SaveWorkflow(workflow *SimpleParseResult) error {
	// 序列化工作流数据
	data, err := json.Marshal(workflow)
	if err != nil {
		return fmt.Errorf("序列化工作流失败: %v", err)
	}

	// 确定状态
	status := "running"
	if workflow.MainFunc != nil {
		allCompleted := true
		hasFailed := false
		for _, stmt := range workflow.MainFunc.Statements {
			if stmt.Status == "failed" {
				hasFailed = true
				break
			}
			if stmt.Status != "completed" {
				allCompleted = false
			}
		}
		if hasFailed {
			status = "failed"
		} else if allCompleted {
			status = "completed"
		}
	}

	// 插入或更新
	upsertSQL := `
	INSERT INTO workflow_states (flow_id, status, data, created_at, updated_at)
	VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT(flow_id) DO UPDATE SET
		status = excluded.status,
		data = excluded.data,
		updated_at = CURRENT_TIMESTAMP
	`

	_, err = p.db.Exec(upsertSQL, workflow.FlowID, status, string(data))
	return err
}

// LoadWorkflow 加载工作流状态
func (p *SQLitePersistence) LoadWorkflow(flowID string) (*SimpleParseResult, error) {
	var data string
	var status string
	var createdAt, updatedAt time.Time

	querySQL := `SELECT status, data, created_at, updated_at FROM workflow_states WHERE flow_id = ?`
	err := p.db.QueryRow(querySQL, flowID).Scan(&status, &data, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	// 反序列化工作流数据
	var workflow SimpleParseResult
	if err := json.Unmarshal([]byte(data), &workflow); err != nil {
		return nil, fmt.Errorf("反序列化工作流失败: %v", err)
	}

	return &workflow, nil
}

// ListWorkflows 列出所有工作流
func (p *SQLitePersistence) ListWorkflows() ([]WorkflowState, error) {
	querySQL := `SELECT flow_id, status, created_at, updated_at FROM workflow_states ORDER BY updated_at DESC`
	rows, err := p.db.Query(querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []WorkflowState
	for rows.Next() {
		var wf WorkflowState
		err := rows.Scan(&wf.FlowID, &wf.Status, &wf.CreatedAt, &wf.UpdatedAt)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, wf)
	}

	return workflows, nil
}

// Close 关闭数据库连接
func (p *SQLitePersistence) Close() error {
	return p.db.Close()
}

func TestSQLitePersistence(t *testing.T) {
	// 1. 创建临时数据库文件
	dbPath := "test_workflow.db"
	defer os.Remove(dbPath) // 测试完成后删除

	// 2. 创建持久化实例
	persistence, err := NewSQLitePersistence(dbPath)
	if err != nil {
		t.Fatalf("创建持久化实例失败: %v", err)
	}
	defer persistence.Close()

	// 3. 创建工作流代码
	workflowCode := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000
}

step1 = beiluo.test1.user.create_user(
    username: string "用户名",
    phone: int "手机号"
) -> (
    userId: string "用户ID",
    err: error "是否失败"
);

func main() {
    用户ID, step1Err := step1(input["用户名"], input["手机号"]){retry:3, timeout:5000}
    if step1Err != nil {
        step1.Printf("❌ 用户创建失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，用户ID: %s", 用户ID)
    
    通知信息 := "欢迎 {{用户名}} 加入我们！"
    fmt.Printf("通知: %s", 通知信息)
}`

	// 4. 解析工作流
	parser := NewSimpleParser()
	parseResult := parser.ParseWorkflow(workflowCode)
	if !parseResult.Success {
		t.Fatalf("解析失败: %s", parseResult.Error)
	}

	// 5. 设置FlowID
	parseResult.FlowID = "test-persistence-" + fmt.Sprintf("%d", time.Now().Unix())

	// 6. 创建执行器
	executor := NewExecutor()

	// 7. 设置回调函数 - 每次状态更新都保存到数据库
	executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
		t.Logf("执行步骤: %s - %s", step.Name, in.StepDesc)
		t.Logf("输入参数: %+v", in.RealInput)

		// 模拟用户创建
		time.Sleep(50 * time.Millisecond)
		return &ExecutorOut{
			Success: true,
			WantOutput: map[string]interface{}{
				"userId": "USER_" + fmt.Sprintf("%d", time.Now().Unix()),
				"err":    nil,
			},
			Error: "",
			Logs:  []string{"用户创建成功"},
		}, nil
	}

	// 状态更新回调 - 每次状态变化都保存到数据库
	executor.OnWorkFlowUpdate = func(ctx context.Context, current *SimpleParseResult) error {
		t.Logf("状态更新: FlowID=%s, 变量数量=%d", current.FlowID, len(current.Variables))

		// 保存到数据库
		if err := persistence.SaveWorkflow(current); err != nil {
			t.Errorf("保存工作流状态失败: %v", err)
			return err
		}

		t.Logf("✅ 工作流状态已保存到数据库")
		return nil
	}

	executor.OnWorkFlowExit = func(ctx context.Context, current *SimpleParseResult) error {
		t.Log("工作流正常结束")
		return nil
	}

	// 8. 执行工作流
	ctx := context.Background()
	if err := executor.Start(ctx, parseResult); err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	// 9. 验证数据库中的状态
	loadedWorkflow, err := persistence.LoadWorkflow(parseResult.FlowID)
	if err != nil {
		t.Fatalf("加载工作流失败: %v", err)
	}

	// 验证数据完整性
	if loadedWorkflow.FlowID != parseResult.FlowID {
		t.Errorf("FlowID不匹配: 期望 %s, 实际 %s", parseResult.FlowID, loadedWorkflow.FlowID)
	}

	if len(loadedWorkflow.Variables) != len(parseResult.Variables) {
		t.Errorf("变量数量不匹配: 期望 %d, 实际 %d", len(parseResult.Variables), len(loadedWorkflow.Variables))
	}

	// 验证语句状态
	if loadedWorkflow.MainFunc != nil {
		for i, stmt := range loadedWorkflow.MainFunc.Statements {
			if stmt.Status != "completed" {
				t.Errorf("语句 %d 状态不正确: 期望 completed, 实际 %s", i, stmt.Status)
			}
		}
	}

	// 10. 列出所有工作流
	workflows, err := persistence.ListWorkflows()
	if err != nil {
		t.Fatalf("列出工作流失败: %v", err)
	}

	if len(workflows) != 1 {
		t.Errorf("工作流数量不正确: 期望 1, 实际 %d", len(workflows))
	}

	t.Logf("✅ 持久化测试通过")
	t.Logf("数据库中的工作流: %+v", workflows[0])
}

func TestPersistenceWithResume(t *testing.T) {
	// 测试断点续传功能
	dbPath := "test_resume.db"
	defer os.Remove(dbPath)

	persistence, err := NewSQLitePersistence(dbPath)
	if err != nil {
		t.Fatalf("创建持久化实例失败: %v", err)
	}
	defer persistence.Close()

	// 1. 创建并执行工作流到一半
	workflowCode := `var input = map[string]interface{}{
    "用户名": "李四",
    "手机号": 13900139000
}

step1 = beiluo.test1.user.create_user(
    username: string "用户名",
    phone: int "手机号"
) -> (
    userId: string "用户ID",
    err: error "是否失败"
);

step2 = beiluo.test1.user.send_notification(
    userId: string "用户ID",
    message: string "消息"
) -> (
    success: bool "是否成功",
    err: error "是否失败"
);

func main() {
    用户ID, step1Err := step1(input["用户名"], input["手机号"]){retry:3, timeout:5000}
    if step1Err != nil {
        return
    }
    
    成功, step2Err := step2(用户ID, "欢迎消息"){retry:2, timeout:3000}
    if step2Err != nil {
        return
    }
}`

	parser := NewSimpleParser()
	parseResult := parser.ParseWorkflow(workflowCode)
	if !parseResult.Success {
		t.Fatalf("解析失败: %s", parseResult.Error)
	}

	parseResult.FlowID = "test-resume-" + fmt.Sprintf("%d", time.Now().Unix())

	// 模拟执行到一半时中断
	executor := NewExecutor()
	executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
		if step.Name == "step1" {
			// step1 成功
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"userId": "USER_123",
					"err":    nil,
				},
			}, nil
		}
		// step2 失败，模拟中断
		return &ExecutorOut{
			Success: false,
			Error:   "网络错误",
		}, nil
	}

	executor.OnWorkFlowUpdate = func(ctx context.Context, current *SimpleParseResult) error {
		return persistence.SaveWorkflow(current)
	}

	// 执行工作流（会失败）
	ctx := context.Background()
	executor.Start(ctx, parseResult)

	// 2. 从数据库恢复状态
	loadedWorkflow, err := persistence.LoadWorkflow(parseResult.FlowID)
	if err != nil {
		t.Fatalf("加载工作流失败: %v", err)
	}

	// 验证恢复的状态
	if loadedWorkflow.Variables["用户ID"].Value != "USER_123" {
		t.Errorf("用户ID未正确恢复")
	}

	// 验证语句状态
	step1Completed := false
	for _, stmt := range loadedWorkflow.MainFunc.Statements {
		if stmt.Function == "step1" && stmt.Status == "completed" {
			step1Completed = true
		}
	}

	if !step1Completed {
		t.Errorf("step1应该已完成")
	}

	t.Logf("✅ 断点续传测试通过")
}
