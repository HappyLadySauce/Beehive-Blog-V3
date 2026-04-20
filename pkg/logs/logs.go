package logs

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/zeromicro/go-zero/core/logx"
)

type contextKey string

const contextFieldsKey contextKey = "beehive.logs.fields"

var sensitiveFragments = []string{
	"password",
	"token",
	"secret",
	"authorization",
	"cookie",
}

// Field describes one structured log field.
// Field 描述一个结构化日志字段。
type Field struct {
	Key   string
	Value any
}

// Logger wraps the project logger entry point.
// Logger 封装项目统一日志入口。
type Logger struct {
	ctx context.Context
}

// WithFields stores structured fields in context for downstream loggers.
// WithFields 将结构化字段写入上下文供后续日志使用。
func WithFields(ctx context.Context, fields ...Field) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(fields) == 0 {
		return ctx
	}

	merged := append([]Field{}, fieldsFromContext(ctx)...)
	merged = append(merged, fields...)
	return context.WithValue(ctx, contextFieldsKey, merged)
}

// WithRequestID stores a request identifier in the log context.
// WithRequestID 在日志上下文中写入请求标识。
func WithRequestID(ctx context.Context, requestID string) context.Context {
	if strings.TrimSpace(requestID) == "" {
		return ctx
	}

	return WithFields(ctx, String("request_id", requestID))
}

// Ctx creates a contextual logger.
// Ctx 创建带上下文的日志记录器。
func Ctx(ctx context.Context) *Logger {
	if ctx == nil {
		ctx = context.Background()
	}

	return &Logger{ctx: ctx}
}

// With adds structured fields to the logger.
// With 为日志记录器附加结构化字段。
func (l *Logger) With(fields ...Field) *Logger {
	if l == nil {
		return Ctx(context.Background())
	}

	return &Logger{ctx: WithFields(l.ctx, fields...)}
}

// Debug writes a structured debug log.
// Debug 写入结构化调试日志。
func (l *Logger) Debug(action string, fields ...Field) {
	l.log(debugLevel, action, nil, fields...)
}

// Info writes a structured info log.
// Info 写入结构化信息日志。
func (l *Logger) Info(action string, fields ...Field) {
	l.log(infoLevel, action, nil, fields...)
}

// Warn writes a structured warn log.
// Warn 写入结构化警告日志。
func (l *Logger) Warn(action string, fields ...Field) {
	l.log(warnLevel, action, nil, fields...)
}

// Error writes a structured error log.
// Error 写入结构化错误日志。
func (l *Logger) Error(action string, err error, fields ...Field) {
	l.log(errorLevel, action, err, fields...)
}

// Debugf writes a formatted debug log.
// Debugf 写入格式化调试日志。
func (l *Logger) Debugf(format string, args ...any) {
	l.printf(debugLevel, format, args...)
}

// Infof writes a formatted info log.
// Infof 写入格式化信息日志。
func (l *Logger) Infof(format string, args ...any) {
	l.printf(infoLevel, format, args...)
}

// Warnf writes a formatted warn log.
// Warnf 写入格式化警告日志。
func (l *Logger) Warnf(format string, args ...any) {
	l.printf(warnLevel, format, args...)
}

// Errorf writes a formatted error log.
// Errorf 写入格式化错误日志。
func (l *Logger) Errorf(format string, args ...any) {
	l.printf(errorLevel, format, args...)
}

// String creates a string field.
// String 创建字符串字段。
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Int64 creates an int64 field.
// Int64 创建 int64 字段。
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Int creates an int field.
// Int 创建 int 字段。
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Bool creates a bool field.
// Bool 创建 bool 字段。
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// Any creates a generic field.
// Any 创建通用字段。
func Any(key string, value any) Field {
	return Field{Key: key, Value: value}
}

// RequestID creates a request identifier field.
// RequestID 创建请求标识字段。
func RequestID(value string) Field {
	return String("request_id", value)
}

// UserID creates a user identifier field.
// UserID 创建用户标识字段。
func UserID(value string) Field {
	return String("user_id", value)
}

// SessionID creates a session identifier field.
// SessionID 创建会话标识字段。
func SessionID(value string) Field {
	return String("session_id", value)
}

// CodeField creates a business code field.
// CodeField 创建业务错误码字段。
func CodeField(err error) Field {
	code, ok := errs.CodeOf(err)
	if !ok {
		return Field{}
	}

	return Int("code", int(code))
}

// Err creates a safe error summary field.
// Err 创建安全错误摘要字段。
func Err(err error) Field {
	if err == nil {
		return Field{}
	}

	return String("error", safeErrorMessage(err))
}

type level int

const (
	debugLevel level = iota
	infoLevel
	warnLevel
	errorLevel
)

func (l *Logger) log(targetLevel level, action string, err error, fields ...Field) {
	if l == nil {
		return
	}

	line := buildLine(action, err, fieldsFromContext(l.ctx), fields)
	entry := logx.WithContext(l.ctx)
	switch targetLevel {
	case debugLevel:
		entry.Debugf("%s", line)
	case warnLevel:
		entry.Infof("level=warn %s", line)
	case errorLevel:
		entry.Errorf("%s", line)
	default:
		entry.Infof("%s", line)
	}
}

func (l *Logger) printf(targetLevel level, format string, args ...any) {
	if l == nil {
		return
	}

	message := sanitizeMessage(fmt.Sprintf(format, args...))
	entry := logx.WithContext(l.ctx)
	switch targetLevel {
	case debugLevel:
		entry.Debugf("%s", message)
	case warnLevel:
		entry.Infof("level=warn %s", message)
	case errorLevel:
		entry.Errorf("%s", message)
	default:
		entry.Infof("%s", message)
	}
}

func buildLine(action string, err error, base []Field, fields []Field) string {
	ordered := []Field{
		String("action", strings.TrimSpace(action)),
	}

	merged := mergeFields(base, fields)
	if codeField := CodeField(err); codeField.Key != "" {
		merged = append(merged, codeField)
	}
	if err != nil {
		merged = append(merged, String("cause", safeErrorMessage(err)))
	}

	ordered = append(ordered, merged...)
	return renderFields(ordered)
}

func fieldsFromContext(ctx context.Context) []Field {
	if ctx == nil {
		return nil
	}

	fields, ok := ctx.Value(contextFieldsKey).([]Field)
	if !ok || len(fields) == 0 {
		return nil
	}

	cloned := make([]Field, len(fields))
	copy(cloned, fields)
	return cloned
}

func mergeFields(base []Field, extra []Field) []Field {
	if len(base) == 0 && len(extra) == 0 {
		return nil
	}

	order := make([]string, 0, len(base)+len(extra))
	values := make(map[string]Field, len(base)+len(extra))
	appendField := func(field Field) {
		key := strings.TrimSpace(field.Key)
		if key == "" {
			return
		}
		if _, exists := values[key]; !exists {
			order = append(order, key)
		}
		values[key] = Field{Key: key, Value: field.Value}
	}

	for _, field := range base {
		appendField(field)
	}
	for _, field := range extra {
		appendField(field)
	}

	preferred := []string{"request_id", "route", "user_id", "session_id", "provider"}
	result := make([]Field, 0, len(order))
	used := make(map[string]struct{}, len(order))
	for _, key := range preferred {
		if field, ok := values[key]; ok {
			result = append(result, field)
			used[key] = struct{}{}
		}
	}

	remaining := make([]string, 0, len(order))
	for _, key := range order {
		if _, ok := used[key]; ok {
			continue
		}
		remaining = append(remaining, key)
	}
	sort.Strings(remaining)
	for _, key := range remaining {
		result = append(result, values[key])
	}

	return result
}

func renderFields(fields []Field) string {
	parts := make([]string, 0, len(fields))
	for _, field := range fields {
		key := strings.TrimSpace(field.Key)
		if key == "" {
			continue
		}

		value := sanitizeValue(key, field.Value)
		if value == "" {
			continue
		}
		parts = append(parts, key+"="+value)
	}

	return strings.Join(parts, " ")
}

func sanitizeValue(key string, value any) string {
	if isSensitiveKey(key) {
		return "[REDACTED]"
	}

	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return quoteIfNeeded(sanitizeMessage(typed))
	case fmt.Stringer:
		return quoteIfNeeded(sanitizeMessage(typed.String()))
	case error:
		return quoteIfNeeded(safeErrorMessage(typed))
	case errs.Code:
		return typed.String()
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return quoteIfNeeded(sanitizeMessage(fmt.Sprintf("%v", typed)))
	}
}

func safeErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	if parsed := errs.Parse(err); parsed != nil && parsed.Cause != nil {
		return sanitizeMessage(parsed.Cause.Error())
	}

	return sanitizeMessage(err.Error())
}

func sanitizeMessage(message string) string {
	trimmed := strings.TrimSpace(message)
	if trimmed == "" {
		return ""
	}

	lowered := strings.ToLower(trimmed)
	for _, fragment := range sensitiveFragments {
		if strings.Contains(lowered, fragment) {
			return "[REDACTED]"
		}
	}

	return trimmed
}

func isSensitiveKey(key string) bool {
	lowered := strings.ToLower(strings.TrimSpace(key))
	for _, fragment := range sensitiveFragments {
		if strings.Contains(lowered, fragment) {
			return true
		}
	}

	return false
}

func quoteIfNeeded(value string) string {
	if value == "" {
		return ""
	}

	if strings.ContainsAny(value, " \t\r\n=") {
		return strconv.Quote(value)
	}

	return value
}
