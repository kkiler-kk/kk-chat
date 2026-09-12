package middleware

import (
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/labstack/echo/v4"
	"server-go/internal/usecase/apperror"
)

// RequestLogger 访问日志：method/path/status/latency/request_id/ip。
func RequestLogger(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			logger.Info("http",
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"latency_ms", time.Since(start).Milliseconds(),
				"request_id", reqID,
				"ip", c.RealIP(),
			)
			return err
		}
	}
}

// Recover 捕获 panic → 记录堆栈 → 以命名返回值把 apperror 传回外层 ErrorHandler 统一渲染。
// 注意：不能用 c.Error()——那会绕过自定义 ErrorHandler 直达 Echo 默认处理器。
func Recover(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("panic recovered",
						"panic", r, "stack", string(debug.Stack()),
						"path", c.Request().URL.Path)
					err = apperror.New(apperror.CodeInternal, "服务器内部错误")
				}
			}()
			return next(c)
		}
	}
}
