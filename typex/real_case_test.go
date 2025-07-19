package typex

import (
	"encoding/json"
	"testing"
)

// 模拟实际的Product结构体（简化版本）
type TestProduct struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Tags        string  `json:"tags"`
	CreatedBy   string  `json:"created_by"`
	CreatedAt   Time    `json:"created_at"`
	UpdatedAt   Time    `json:"updated_at"`
}

func TestRealCaseJSON(t *testing.T) {
	// 这是用户实际遇到问题的JSON数据
	jsonData := `{
		"name": "111",
		"category": "其他",
		"price": 11,
		"stock": 11,
		"description": "",
		"status": "启用",
		"tags": "",
		"created_by": "",
		"created_at": "",
		"updated_at": ""
	}`

	var product TestProduct
	err := json.Unmarshal([]byte(jsonData), &product)
	if err != nil {
		t.Errorf("反序列化失败: %v", err)
		return
	}

	// 验证基本字段
	if product.Name != "111" {
		t.Errorf("期望Name: 111, 实际: %s", product.Name)
	}

	if product.Category != "其他" {
		t.Errorf("期望Category: 其他, 实际: %s", product.Category)
	}

	if product.Price != 11 {
		t.Errorf("期望Price: 11, 实际: %f", product.Price)
	}

	if product.Stock != 11 {
		t.Errorf("期望Stock: 11, 实际: %d", product.Stock)
	}

	if product.Status != "启用" {
		t.Errorf("期望Status: 启用, 实际: %s", product.Status)
	}

	// 验证空字符串时间字段被正确处理为零值
	if !product.CreatedAt.IsZero() {
		t.Errorf("期望CreatedAt是零值时间，实际: %v", product.CreatedAt)
	}

	if !product.UpdatedAt.IsZero() {
		t.Errorf("期望UpdatedAt是零值时间，实际: %v", product.UpdatedAt)
	}

	t.Logf("测试通过！产品: %+v", product)
}

// 为Time类型添加IsZero方法以便测试
func (t Time) IsZero() bool {
	return t.GetUnix() == 0 || t.GetUnix() == -62135596800 // Unix时间戳的零值
}

func TestCompleteOnTableAddRowsJSON(t *testing.T) {
	// 完整的OnTableAddRows请求JSON（简化版本，只包含时间相关字段）
	jsonData := `{
		"method": "GET",
		"router": "ex4/product_list1/",
		"type": "OnTableAddRows",
		"body": {
			"rows": [
				{
					"name": "111",
					"category": "其他",
					"price": 11,
					"stock": 11,
					"description": "",
					"status": "启用",
					"tags": "",
					"created_by": "",
					"created_at": "",
					"updated_at": ""
				}
			]
		}
	}`

	// 定义请求结构体
	type OnTableAddRowsRequest struct {
		Method string `json:"method"`
		Router string `json:"router"`
		Type   string `json:"type"`
		Body   struct {
			Rows []TestProduct `json:"rows"`
		} `json:"body"`
	}

	var request OnTableAddRowsRequest
	err := json.Unmarshal([]byte(jsonData), &request)
	if err != nil {
		t.Errorf("反序列化完整请求失败: %v", err)
		return
	}

	// 验证请求基本信息
	if request.Method != "GET" {
		t.Errorf("期望Method: GET, 实际: %s", request.Method)
	}

	if request.Type != "OnTableAddRows" {
		t.Errorf("期望Type: OnTableAddRows, 实际: %s", request.Type)
	}

	// 验证行数据
	if len(request.Body.Rows) != 1 {
		t.Errorf("期望1行数据，实际: %d", len(request.Body.Rows))
		return
	}

	product := request.Body.Rows[0]

	// 验证时间字段被正确处理
	if !product.CreatedAt.IsZero() {
		t.Errorf("期望CreatedAt是零值时间，实际: %v", product.CreatedAt)
	}

	if !product.UpdatedAt.IsZero() {
		t.Errorf("期望UpdatedAt是零值时间，实际: %v", product.UpdatedAt)
	}

	t.Logf("完整请求测试通过！产品: %+v", product)
}
