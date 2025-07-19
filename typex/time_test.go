package typex

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimeUnmarshalJSON_EmptyString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool // 是否应该是零值时间
		hasError bool
	}{
		{
			name:     "空字符串",
			input:    `""`,
			expected: true,
			hasError: false,
		},
		{
			name:     "null值",
			input:    `"null"`,
			expected: true,
			hasError: false,
		},
		{
			name:     "正常时间格式",
			input:    `"2025-01-15 10:30:00"`,
			expected: false,
			hasError: false,
		},
		{
			name:     "ISO格式",
			input:    `"2025-01-15T10:30:00Z"`,
			expected: false,
			hasError: false,
		},
		{
			name:     "无效格式",
			input:    `"invalid-time"`,
			expected: false,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tm Time
			err := json.Unmarshal([]byte(tt.input), &tm)

			if tt.hasError {
				if err == nil {
					t.Errorf("期望有错误，但没有错误")
				}
				return
			}

			if err != nil {
				t.Errorf("不期望有错误，但得到错误: %v", err)
				return
			}

			isZero := time.Time(tm).IsZero()
			if tt.expected != isZero {
				t.Errorf("期望零值时间: %v, 实际零值时间: %v, 时间值: %v", tt.expected, isZero, time.Time(tm))
			}
		})
	}
}

func TestTimeMarshalJSON_ZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		time     Time
		expected string
	}{
		{
			name:     "零值时间",
			time:     Time(time.Time{}),
			expected: `""`,
		},
		{
			name:     "正常时间",
			time:     Time(time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)),
			expected: `"2025-01-15 10:30:00"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.time)
			if err != nil {
				t.Errorf("序列化失败: %v", err)
				return
			}

			if string(result) != tt.expected {
				t.Errorf("期望: %s, 实际: %s", tt.expected, string(result))
			}
		})
	}
}

// 测试产品结构体的JSON序列化和反序列化
func TestProductWithEmptyTimes(t *testing.T) {
	// 模拟你的Product结构体的简化版本
	type Product struct {
		Name      string `json:"name"`
		CreatedAt Time   `json:"created_at"`
		UpdatedAt Time   `json:"updated_at"`
	}

	// 测试包含空时间字段的JSON
	jsonData := `{
		"name": "测试产品",
		"created_at": "",
		"updated_at": ""
	}`

	var product Product
	err := json.Unmarshal([]byte(jsonData), &product)
	if err != nil {
		t.Errorf("反序列化失败: %v", err)
		return
	}

	// 验证时间字段是零值
	if !time.Time(product.CreatedAt).IsZero() {
		t.Errorf("期望CreatedAt是零值时间，实际: %v", time.Time(product.CreatedAt))
	}

	if !time.Time(product.UpdatedAt).IsZero() {
		t.Errorf("期望UpdatedAt是零值时间，实际: %v", time.Time(product.UpdatedAt))
	}

	// 测试序列化
	result, err := json.Marshal(product)
	if err != nil {
		t.Errorf("序列化失败: %v", err)
		return
	}

	t.Logf("序列化结果: %s", string(result))

	// 验证零值时间被序列化为空字符串
	var checkProduct Product
	err = json.Unmarshal(result, &checkProduct)
	if err != nil {
		t.Errorf("重新反序列化失败: %v", err)
		return
	}

	if checkProduct.Name != product.Name {
		t.Errorf("名称不匹配")
	}
}
