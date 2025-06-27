package typex

import (
	"context"

	"github.com/google/uuid"
	"github.com/yunhanshu-net/pkg/constants"
)

// getTraceIDFromContext 从context中获取trace_id
func getTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(constants.TraceID).(string); ok && traceID != "" {
		return traceID
	}
	// 如果是gin.Context，尝试使用GetString方法
	if ginCtx, ok := ctx.(interface{ GetString(string) string }); ok {
		if traceID := ginCtx.GetString(constants.TraceID); traceID != "" {
			return traceID
		}
	}
	// 如果没有trace_id，生成一个新的UUID作为fallback
	return uuid.New().String()
}

// getUserFromContext 从context中获取用户ID
func getUserFromContext(ctx context.Context) string {
	if user, ok := ctx.Value("user").(string); ok {
		return user
	}
	// 如果是gin.Context，尝试使用GetString方法
	if ginCtx, ok := ctx.(interface{ GetString(string) string }); ok {
		if user := ginCtx.GetString("user"); user != "" {
			return user
		}
	}
	return "anonymous" // 默认值
}

// getRunnerFromContext 从context中获取运行器ID
func getRunnerFromContext(ctx context.Context) string {
	if runner, ok := ctx.Value("runner").(string); ok {
		return runner
	}
	// 如果是gin.Context，尝试使用GetString方法
	if ginCtx, ok := ctx.(interface{ GetString(string) string }); ok {
		if runner := ginCtx.GetString("runner"); runner != "" {
			return runner
		}
	}
	return "default" // 默认值
}
