package errs

import "errors"

// Coder exposes a stable business code from an error value.
// Coder 暴露错误值上的稳定业务错误码。
type Coder interface {
	ErrorCode() Code
}

// Error defines the unified project-level domain error.
// Error 定义统一的项目级领域错误。
type Error struct {
	Code      Code
	Message   string
	Cause     error
	Reference string
	Meta      map[string]any
}

// Error implements the error interface.
// Error 实现 error 接口。
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e.Message == "" {
		return e.Code.String()
	}

	return e.Message
}

// Unwrap exposes the wrapped cause.
// Unwrap 暴露被包装的底层原因。
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Cause
}

// ErrorCode returns the stable business code.
// ErrorCode 返回稳定业务错误码。
func (e *Error) ErrorCode() Code {
	if e == nil {
		return 0
	}

	return e.Code
}

// Is reports whether the target error matches by business code.
// Is 按业务错误码判断目标错误是否匹配。
func (e *Error) Is(target error) bool {
	if e == nil || target == nil {
		return false
	}

	targetCode, ok := CodeOf(target)
	return ok && e.Code == targetCode
}

// Option customizes an Error during construction.
// Option 在构造错误时自定义 Error。
type Option func(target *Error)

// WithCause attaches a wrapped cause to the error.
// WithCause 为错误附加底层原因。
func WithCause(cause error) Option {
	return func(target *Error) {
		target.Cause = cause
	}
}

// WithReference attaches a traceable reference token.
// WithReference 为错误附加可追踪引用标识。
func WithReference(reference string) Option {
	return func(target *Error) {
		target.Reference = reference
	}
}

// WithMeta attaches safe structured metadata.
// WithMeta 为错误附加安全的结构化元数据。
func WithMeta(meta map[string]any) Option {
	return func(target *Error) {
		if len(meta) == 0 {
			return
		}

		cloned := make(map[string]any, len(meta))
		for key, value := range meta {
			cloned[key] = value
		}
		target.Meta = cloned
	}
}

// New creates a unified business error.
// New 创建统一业务错误。
func New(code Code, message string, opts ...Option) error {
	err := &Error{
		Code:    code,
		Message: message,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(err)
		}
	}

	return err
}

// Wrap wraps an existing error with a stable business code and message.
// Wrap 使用稳定业务码和消息包装现有错误。
func Wrap(err error, code Code, message string, opts ...Option) error {
	if err == nil {
		return New(code, message, opts...)
	}

	options := make([]Option, 0, len(opts)+1)
	options = append(options, WithCause(err))
	options = append(options, opts...)
	return New(code, message, options...)
}

type sentinel struct {
	code Code
}

func (s sentinel) Error() string {
	return s.code.String()
}

func (s sentinel) ErrorCode() Code {
	return s.code
}

// E builds a sentinel error for errors.Is matching by business code.
// E 构建用于 errors.Is 按业务错误码匹配的哨兵错误。
func E(code Code) error {
	return sentinel{code: code}
}

// Parse extracts the first Error from a wrapped error chain.
// Parse 从错误链中提取第一个 Error。
func Parse(err error) *Error {
	if err == nil {
		return nil
	}

	var target *Error
	if errors.As(err, &target) {
		return target
	}

	return nil
}

// CodeOf extracts the first business code from an error chain.
// CodeOf 从错误链中提取第一个业务错误码。
func CodeOf(err error) (Code, bool) {
	if err == nil {
		return 0, false
	}

	var target Coder
	if errors.As(err, &target) {
		return target.ErrorCode(), true
	}

	return 0, false
}

// MessageOf extracts the first stable client-facing message from an error chain.
// MessageOf 从错误链中提取第一个稳定对外消息。
func MessageOf(err error) string {
	parsed := Parse(err)
	if parsed == nil {
		return ""
	}

	return parsed.Message
}

// IsCode reports whether an error chain contains the given business code.
// IsCode 判断错误链中是否包含指定业务码。
func IsCode(err error, code Code) bool {
	return errors.Is(err, E(code))
}
