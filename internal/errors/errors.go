package errors

import (
	"fmt"
	"strings"
)

// ErrorType categorizes hull errors for structured handling.
type ErrorType string

const (
	ErrParse            ErrorType = "PARSE"
	ErrExpression       ErrorType = "EXPRESSION"
	ErrIncludeCycle     ErrorType = "INCLUDE_CYCLE"
	ErrIncludeNotFound  ErrorType = "INCLUDE_NOT_FOUND"
	ErrUndefinedVar     ErrorType = "UNDEFINED_VAR"
	ErrFunction         ErrorType = "FUNCTION_ERROR"
	ErrType             ErrorType = "TYPE_ERROR"
	ErrSchemaValidation ErrorType = "SCHEMA_VALIDATION"
	ErrKube             ErrorType = "KUBE_ERROR"
	ErrRelease          ErrorType = "RELEASE_ERROR"
	ErrPackageInvalid   ErrorType = "PACKAGE_INVALID"
	ErrDependency       ErrorType = "DEPENDENCY_ERROR"
	ErrCLIFlag          ErrorType = "CLI_FLAG"
	ErrCLIValidation    ErrorType = "CLI_VALIDATION"
	ErrRepo             ErrorType = "REPO_ERROR"
	ErrArchive          ErrorType = "ARCHIVE_ERROR"
	ErrRegistry         ErrorType = "REGISTRY_ERROR"
	ErrAuth             ErrorType = "AUTH_ERROR"
	ErrConflict         ErrorType = "DEPENDENCY_CONFLICT"
	ErrCycle            ErrorType = "DEPENDENCY_CYCLE"
	ErrDigest           ErrorType = "DIGEST_MISMATCH"
	ErrSignature        ErrorType = "SIGNATURE_ERROR"
	ErrLockFile         ErrorType = "LOCKFILE_ERROR"
	ErrRateLimit        ErrorType = "RATE_LIMIT"
	ErrReleaseNotFound  ErrorType = "RELEASE_NOT_FOUND"
	ErrInternal         ErrorType = "INTERNAL"
)

// HullError is the structured error type used throughout hull.
type HullError struct {
	Type       ErrorType
	Message    string
	FilePath   string
	Line       int
	Column     int
	Expression string
	Cause      error
	Context    map[string]string
}

// Error implements the error interface.
func (e *HullError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "[%s] %s", e.Type, e.Message)

	if e.FilePath != "" {
		fmt.Fprintf(&b, " (file: %s", e.FilePath)
		if e.Line > 0 {
			fmt.Fprintf(&b, ":%d", e.Line)
			if e.Column > 0 {
				fmt.Fprintf(&b, ":%d", e.Column)
			}
		}
		b.WriteString(")")
	}

	if e.Expression != "" {
		fmt.Fprintf(&b, " expr: %q", e.Expression)
	}

	if e.Cause != nil {
		fmt.Fprintf(&b, ": %s", e.Cause.Error())
	}

	return b.String()
}

// Unwrap returns the underlying cause for errors.Is/As support.
func (e *HullError) Unwrap() error {
	return e.Cause
}

// NewError creates a HullError with a type and message.
func NewError(errType ErrorType, message string) *HullError {
	return &HullError{
		Type:    errType,
		Message: message,
	}
}

// NewErrorf creates a HullError with a formatted message.
func NewErrorf(errType ErrorType, format string, args ...any) *HullError {
	return &HullError{
		Type:    errType,
		Message: fmt.Sprintf(format, args...),
	}
}

// WrapError wraps an existing error into a HullError.
func WrapError(errType ErrorType, message string, cause error) *HullError {
	return &HullError{
		Type:    errType,
		Message: message,
		Cause:   cause,
	}
}

// WrapErrorf wraps an existing error with a formatted message.
func WrapErrorf(errType ErrorType, cause error, format string, args ...any) *HullError {
	return &HullError{
		Type:    errType,
		Message: fmt.Sprintf(format, args...),
		Cause:   cause,
	}
}

// WithFile attaches file location info and returns the same error for chaining.
func (e *HullError) WithFile(path string, line, column int) *HullError {
	e.FilePath = path
	e.Line = line
	e.Column = column
	return e
}

// WithExpression attaches an expression string and returns the same error.
func (e *HullError) WithExpression(expr string) *HullError {
	e.Expression = expr
	return e
}

// WithContext attaches a key-value pair to the context map and returns the same error.
func (e *HullError) WithContext(key, value string) *HullError {
	if nil == e.Context {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// --- Domain-specific constructors ---

// ParseError creates a parse error with file location.
func ParseError(message, filePath string, line, column int) *HullError {
	return NewError(ErrParse, message).WithFile(filePath, line, column)
}

// ExpressionError creates an expression evaluation error.
func ExpressionError(message, expr string, cause error) *HullError {
	he := WrapError(ErrExpression, message, cause)
	he.Expression = expr
	return he
}

// PackageError creates a package-related error.
func PackageError(message, filePath string, cause error) *HullError {
	he := WrapError(ErrPackageInvalid, message, cause)
	he.FilePath = filePath
	return he
}

// KubeError creates a Kubernetes-related error.
func KubeError(message string, cause error) *HullError {
	return WrapError(ErrKube, message, cause)
}

// CLIError creates a CLI flag or validation error.
func CLIError(errType ErrorType, message string) *HullError {
	return NewError(errType, message)
}

// InternalError creates an internal/unexpected error.
func InternalError(message string, cause error) *HullError {
	return WrapError(ErrInternal, message, cause)
}

// FormatUserFriendly returns a human-readable error suitable for CLI output.
func FormatUserFriendly(err error) string {
	he, ok := err.(*HullError)
	if !ok {
		return fmt.Sprintf("Error: %s", err.Error())
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Error: %s\n", he.Message)

	if he.FilePath != "" {
		fmt.Fprintf(&b, "  File: %s", he.FilePath)
		if he.Line > 0 {
			fmt.Fprintf(&b, ", line %d", he.Line)
			if he.Column > 0 {
				fmt.Fprintf(&b, ", column %d", he.Column)
			}
		}
		b.WriteString("\n")
	}

	if he.Expression != "" {
		fmt.Fprintf(&b, "  Expression: %s\n", he.Expression)
	}

	if he.Cause != nil {
		fmt.Fprintf(&b, "  Caused by: %s\n", he.Cause.Error())
	}

	for k, v := range he.Context {
		fmt.Fprintf(&b, "  %s: %s\n", k, v)
	}

	return b.String()
}
