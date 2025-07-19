// Package httpx 提供链式调用的HTTP请求工具
package httpx

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ConnectivityResult 连通性检查结果
type ConnectivityResult struct {
	// 基础网络状态
	Reachable  bool   `json:"reachable"`   // 是否可达
	StatusCode int    `json:"status_code"` // HTTP状态码
	Error      string `json:"error"`       // 错误信息

	// 性能指标
	Latency      time.Duration `json:"latency"`       // 网络延迟
	ResponseTime time.Duration `json:"response_time"` // 响应时间

	// 网络诊断
	Timeout     bool `json:"timeout"`      // 是否超时
	DNSResolved bool `json:"dns_resolved"` // DNS是否解析成功

	// 服务器响应头信息
	Server        string `json:"server"`         // 服务器标识
	ContentType   string `json:"content_type"`   // 内容类型
	ContentLength int64  `json:"content_length"` // 内容长度

	// SSL信息
	SSLValid bool `json:"ssl_valid"` // SSL证书是否有效

	// 请求信息（新增）
	RequestMethod  string            `json:"request_method"`  // 请求方法
	RequestURL     string            `json:"request_url"`     // 请求URL
	RequestHeaders map[string]string `json:"request_headers"` // 请求头

	// 扩展数据
	Metadata map[string]interface{} `json:"metadata"` // 扩展信息
}

// HttpRequest 链式调用的HTTP请求构建器
type HttpRequest struct {
	method  string
	url     string
	headers map[string]string
	body    interface{}
	timeout time.Duration

	// 连通性检查结果作为内部状态
	connectivity *ConnectivityResult
}

// Delete 创建DELETE请求
func Delete(url string) *HttpRequest {
	return &HttpRequest{
		method:  "DELETE",
		url:     url,
		headers: make(map[string]string),
		timeout: 60 * time.Second,
	}
}

// Get 创建GET请求
func Get(url string) *HttpRequest {
	return &HttpRequest{
		method:  "GET",
		url:     url,
		headers: make(map[string]string),
		timeout: 60 * time.Second,
	}
}

// Post 创建POST请求
func Post(url string) *HttpRequest {
	return &HttpRequest{
		method:  "POST",
		url:     url,
		headers: make(map[string]string),
		timeout: 60 * time.Second,
	}
}

// Put 创建PUT请求
func Put(url string) *HttpRequest {
	return &HttpRequest{
		method:  "PUT",
		url:     url,
		headers: make(map[string]string),
		timeout: 60 * time.Second,
	}
}

// Head 创建HEAD请求
func Head(url string) *HttpRequest {
	return &HttpRequest{
		method:  "HEAD",
		url:     url,
		headers: make(map[string]string),
		timeout: 60 * time.Second,
	}
}

// Header 设置请求头
func (h *HttpRequest) Header(key, value string) *HttpRequest {
	h.headers[key] = value
	return h
}

// Headers 批量设置请求头
func (h *HttpRequest) Headers(headers map[string]string) *HttpRequest {
	for k, v := range headers {
		h.headers[k] = v
	}
	return h
}

// Body 设置请求体
func (h *HttpRequest) Body(body interface{}) *HttpRequest {
	h.body = body
	return h
}

// Timeout 设置超时时间
func (h *HttpRequest) Timeout(timeout time.Duration) *HttpRequest {
	h.timeout = timeout
	return h
}

// Do 执行请求并可选绑定响应体到结构体
func (h *HttpRequest) Do(resBindBody interface{}) (*RequestContext, error) {
	var bodyStr string
	if h.body != nil {
		bodyBytes, err := json.Marshal(h.body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		bodyStr = string(bodyBytes)
	}

	return h.request(bodyStr, resBindBody)
}

// DoString 执行请求并返回字符串响应
func (h *HttpRequest) DoString() (*RequestContext, error) {
	return h.Do(nil)
}

// ConnectivityCheck 执行连通性检查，返回HttpRequest支持链式调用
func (h *HttpRequest) ConnectivityCheck() *HttpRequest {
	h.connectivity = h.performConnectivityCheck()
	return h
}

// GetConnectivity 获取连通性检查结果
func (h *HttpRequest) GetConnectivity() *ConnectivityResult {
	return h.connectivity
}

// performConnectivityCheck 执行连通性检查的具体实现
func (h *HttpRequest) performConnectivityCheck() *ConnectivityResult {
	result := &ConnectivityResult{}

	// 创建HEAD请求
	req, err := http.NewRequest("HEAD", h.url, nil)
	if err != nil {
		result.Error = fmt.Sprintf("创建请求失败: %v", err)
		return result
	}

	// 设置请求头
	for k, v := range h.headers {
		req.Header.Set(k, v)
	}

	// 执行请求
	start := time.Now()
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	result.ResponseTime = time.Since(start)

	if err != nil {
		result.Error = err.Error()

		// 基础错误分析
		if strings.Contains(err.Error(), "timeout") {
			result.Timeout = true
		}
		if strings.Contains(err.Error(), "no such host") {
			result.DNSResolved = false
		}

		return result
	}
	defer resp.Body.Close()

	// 设置基础信息
	result.Reachable = true
	result.StatusCode = resp.StatusCode
	result.DNSResolved = true

	// 服务器信息
	result.Server = resp.Header.Get("Server")
	result.ContentType = resp.Header.Get("Content-Type")
	result.ContentLength = resp.ContentLength

	// SSL信息 - 验证SSL证书有效性
	if strings.HasPrefix(h.url, "https://") {
		result.SSLValid = h.verifySSL(h.url)
	}

	return result
}

// verifySSL 验证SSL证书有效性
func (h *HttpRequest) verifySSL(url string) bool {
	// 解析URL获取主机名
	if !strings.HasPrefix(url, "https://") {
		return false
	}

	// 提取主机名和端口
	host := strings.TrimPrefix(url, "https://")
	if strings.Contains(host, "/") {
		host = strings.Split(host, "/")[0]
	}

	// 确保有端口号
	if !strings.Contains(host, ":") {
		host += ":443"
	}

	// 创建TLS配置
	config := &tls.Config{
		InsecureSkipVerify: false, // 严格验证SSL证书
		ServerName:         strings.Split(host, ":")[0],
	}

	// 创建连接
	conn, err := tls.Dial("tcp", host, config)
	if err != nil {
		return false
	}
	defer conn.Close()

	// 验证证书链
	err = conn.VerifyHostname(strings.Split(host, ":")[0])
	return err == nil
}

// getBodyString 获取请求体字符串
func (h *HttpRequest) getBodyString() string {
	if h.body == nil {
		return ""
	}
	bodyBytes, _ := json.Marshal(h.body)
	return string(bodyBytes)
}

// DryRun 返回DryRun案例，不实际执行
func (h *HttpRequest) DryRun() *HttpDryRunCase {
	dryRunCase := &HttpDryRunCase{
		CaseType:    "http",
		Method:      h.method,
		Url:         h.url,
		Headers:     h.headers,
		Body:        h.getBodyString(),
		Description: fmt.Sprintf("HTTP %s 调用", h.method),
	}

	// 如果已执行连通性检查，则包含结果
	if h.connectivity != nil {
		dryRunCase.Connectivity = h.connectivity
	}

	return dryRunCase
}

// RequestContext 请求上下文
type RequestContext struct {
	StartTime     time.Time
	EndTime       time.Time
	Cost          time.Duration
	Url           string
	ReqBody       string
	Method        string
	ResBodyString string
	resetBody     string
	Code          int
	Headers       map[string]string
}

// ResetLogBody 重置日志体
func (r *RequestContext) ResetLogBody(resetLogBody string) *RequestContext {
	r.resetBody = resetLogBody
	return r
}

// String 字符串表示
func (r *RequestContext) String(messages ...string) string {
	if r == nil {
		if len(messages) > 0 {
			return fmt.Sprintf("%s nil", messages[0])
		}
		return "nil"
	}

	body := r.ResBodyString
	if r.resetBody != "" {
		body = r.resetBody
	}

	msg := ""
	if len(messages) > 0 {
		msg = messages[0]
	}

	return fmt.Sprintf("msg:%s request method:%v cost:%v url:%v body:%v "+
		",res_code: %v res_body: %v", msg, r.Method, r.Cost, r.Url, r.ReqBody, r.Code, body)
}

// OK 检查是否成功
func (r *RequestContext) OK() bool {
	return r != nil && r.Code == 200
}

// request 执行HTTP请求的核心方法
func (h *HttpRequest) request(bodyStr string, resBindBody interface{}) (*RequestContext, error) {
	start := time.Now()

	// 创建请求体
	var body io.Reader
	if bodyStr != "" {
		body = strings.NewReader(bodyStr)
	}

	// 创建HTTP请求
	req, err := http.NewRequest(h.method, h.url, body)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	for k, v := range h.headers {
		req.Header.Set(k, v)
	}

	// 如果没有设置Content-Type且有请求体，默认设置为application/json; charset=utf-8
	if bodyStr != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: h.timeout,
	}

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		end := time.Now()
		return &RequestContext{
			Url:           h.url,
			Method:        h.method,
			StartTime:     start,
			EndTime:       end,
			Cost:          end.Sub(start),
			ReqBody:       bodyStr,
			Code:          0,
			ResBodyString: "",
		}, fmt.Errorf("HTTP请求执行失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		end := time.Now()
		return &RequestContext{
			Url:           h.url,
			Method:        h.method,
			StartTime:     start,
			EndTime:       end,
			Cost:          end.Sub(start),
			ReqBody:       bodyStr,
			Code:          resp.StatusCode,
			ResBodyString: "",
		}, fmt.Errorf("读取响应体失败: %w", err)
	}

	end := time.Now()
	reqCtx := &RequestContext{
		Url:           h.url,
		Method:        h.method,
		StartTime:     start,
		EndTime:       end,
		Cost:          end.Sub(start),
		ReqBody:       bodyStr,
		Code:          resp.StatusCode,
		ResBodyString: string(respBody),
		Headers:       h.headers,
	}

	// 如果提供了绑定目标，尝试解析JSON响应
	if resBindBody != nil {
		if err := json.Unmarshal(respBody, resBindBody); err != nil {
			return reqCtx, fmt.Errorf("解析响应JSON失败: %w", err)
		}
	}

	return reqCtx, nil
}

// HttpDryRunCase HTTP 操作的 DryRun 实现
// 字段全部大写，便于外部直接访问
type HttpDryRunCase struct {
	CaseType    string                 `json:"case_type"`
	Method      string                 `json:"method"`
	Url         string                 `json:"url"`
	Headers     map[string]string      `json:"headers"`
	Body        string                 `json:"body"`
	Description string                 `json:"description"`
	Meta        map[string]interface{} `json:"meta"`

	// 连通性检查结果
	Connectivity *ConnectivityResult `json:"connectivity,omitempty"`
}

func (h *HttpDryRunCase) Type() string {
	return "http"
}

func (h *HttpDryRunCase) Map() map[string]interface{} {
	return map[string]interface{}{
		"method":      h.Method,
		"url":         h.Url,
		"headers":     h.Headers,
		"body":        h.Body,
		"description": h.Description,
	}
}

func (h *HttpDryRunCase) Metadata() map[string]interface{} {
	return h.Meta
}

// getRiskLevel 根据 HTTP 方法判断风险等级
func getRiskLevel(method string) string {
	switch method {
	case "DELETE":
		return "high"
	case "PUT", "PATCH":
		return "medium"
	case "POST":
		return "medium"
	default:
		return "low"
	}
}
