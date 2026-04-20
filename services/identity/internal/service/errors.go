package service

import (
	"errors"
	"fmt"
)

// ErrorKind describes a stable service-layer error category.
// ErrorKind 描述稳定的 service 层错误类别。
type ErrorKind string

const (
	// ErrorKindInvalidArgument marks input validation failures.
	// ErrorKindInvalidArgument 标记输入校验失败。
	ErrorKindInvalidArgument ErrorKind = "invalid_argument"
	// ErrorKindUnauthenticated marks authentication failures.
	// ErrorKindUnauthenticated 标记认证失败。
	ErrorKindUnauthenticated ErrorKind = "unauthenticated"
	// ErrorKindAlreadyExists marks unique resource conflicts.
	// ErrorKindAlreadyExists 标记唯一资源冲突。
	ErrorKindAlreadyExists ErrorKind = "already_exists"
	// ErrorKindFailedPrecondition marks state-dependent failures.
	// ErrorKindFailedPrecondition 标记依赖状态的失败。
	ErrorKindFailedPrecondition ErrorKind = "failed_precondition"
	// ErrorKindNotFound marks missing resources.
	// ErrorKindNotFound 标记资源不存在。
	ErrorKindNotFound ErrorKind = "not_found"
	// ErrorKindUnimplemented marks intentionally unavailable behavior.
	// ErrorKindUnimplemented 标记有意未开放的行为。
	ErrorKindUnimplemented ErrorKind = "unimplemented"
)

// Error represents a transport-agnostic business error.
// Error 表示与传输协议无关的业务错误。
type Error struct {
	Kind    ErrorKind
	Message string
	Cause   error
}

// Error implements the error interface.
// Error 实现 error 接口。
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

// Unwrap exposes the wrapped cause.
// Unwrap 暴露被包装的底层错误。
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

// NewError creates a service-layer business error.
// NewError 创建 service 层业务错误。
func NewError(kind ErrorKind, message string, cause error) error {
	return &Error{
		Kind:    kind,
		Message: message,
		Cause:   cause,
	}
}

// IsKind reports whether an error belongs to the given kind.
// IsKind 判断错误是否属于指定类别。
func IsKind(err error, kind ErrorKind) bool {
	var serviceErr *Error
	if !errors.As(err, &serviceErr) {
		return false
	}

	return serviceErr.Kind == kind
}
