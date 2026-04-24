package grpcx

import (
	"strconv"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errorDomain = "beehive.blog.v3"

// CodeToGRPC maps a business code to a gRPC status code.
// CodeToGRPC 将业务错误码映射为 gRPC 状态码。
func CodeToGRPC(code errs.Code) codes.Code {
	switch code {
	case errs.CodeGatewayNotReady, errs.CodeGatewayAuthServiceUnavailable, errs.CodeIdentityDependencyUnavailable:
		return codes.Unavailable
	case errs.CodeGatewayUpstreamTimeout:
		return codes.DeadlineExceeded
	case errs.CodeContentNotFound, errs.CodeContentTagNotFound, errs.CodeContentRevisionNotFound:
		return codes.NotFound
	case errs.CodeContentTagInUse:
		return codes.FailedPrecondition
	case errs.CodeIdentityAccountPending, errs.CodeIdentityAccountDisabled, errs.CodeIdentityAccountLocked,
		errs.CodeIdentitySSOProviderDisabled, errs.CodeIdentitySSOProviderNotReady:
		return codes.FailedPrecondition
	case errs.CodeIdentitySSOStateInvalid:
		return codes.Unauthenticated
	}

	category := (int(code) / 100) % 100
	switch category {
	case 1:
		return codes.InvalidArgument
	case 2:
		return codes.Unauthenticated
	case 3:
		return codes.PermissionDenied
	case 4:
		return codes.FailedPrecondition
	case 5:
		if int(code)%100 == 1 {
			return codes.NotFound
		}
		return codes.AlreadyExists
	case 6:
		return codes.Unavailable
	case 99:
		return codes.Internal
	default:
		return codes.Internal
	}
}

// ToStatus converts a domain error into a gRPC status error with structured details.
// ToStatus 将领域错误转换为带结构化明细的 gRPC status 错误。
func ToStatus(err error, fallbackMessage string) error {
	if err == nil {
		return nil
	}
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	parsed := errs.Parse(err)
	if parsed == nil {
		if fallbackMessage == "" {
			fallbackMessage = "internal error"
		}
		parsed = &errs.Error{
			Code:    errs.CodeIdentityInternal,
			Message: fallbackMessage,
			Cause:   err,
		}
	}

	st := status.New(CodeToGRPC(parsed.Code), parsed.Message)
	detail := &errdetails.ErrorInfo{
		Reason: strconv.Itoa(int(parsed.Code)),
		Domain: errorDomain,
		Metadata: map[string]string{
			"message": parsed.Message,
		},
	}
	if parsed.Reference != "" {
		detail.Metadata["reference"] = parsed.Reference
	}

	withDetails, detailErr := st.WithDetails(detail)
	if detailErr != nil {
		return st.Err()
	}

	return withDetails.Err()
}

// ParseStatus extracts a domain error from a gRPC status error when structured details exist.
// ParseStatus 在存在结构化明细时从 gRPC status 错误中提取领域错误。
func ParseStatus(err error) (*errs.Error, bool) {
	if err == nil {
		return nil, false
	}

	st, ok := status.FromError(err)
	if !ok {
		return nil, false
	}

	for _, detail := range st.Details() {
		errorInfo, ok := detail.(*errdetails.ErrorInfo)
		if !ok || errorInfo.GetDomain() != errorDomain {
			continue
		}

		codeValue, parseErr := strconv.Atoi(errorInfo.GetReason())
		if parseErr != nil {
			continue
		}

		parsed := &errs.Error{
			Code:    errs.Code(codeValue),
			Message: st.Message(),
		}
		if metadata := errorInfo.GetMetadata(); metadata != nil {
			if reference := metadata["reference"]; reference != "" {
				parsed.Reference = reference
			}
		}

		return parsed, true
	}

	return nil, false
}

// FromStatus rebuilds a domain error from a gRPC status error.
// FromStatus 从 gRPC status 错误重建领域错误。
func FromStatus(err error, fallbackCode errs.Code, fallbackMessage string) error {
	if err == nil {
		return nil
	}

	if parsed, ok := ParseStatus(err); ok {
		return parsed
	}

	if fallbackCode == 0 {
		return nil
	}
	if fallbackMessage == "" {
		fallbackMessage = "internal server error"
	}

	return errs.New(fallbackCode, fallbackMessage)
}
