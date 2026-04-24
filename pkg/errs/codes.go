package errs

// Code defines a stable project-level business error code.
// Code 定义稳定的项目级业务错误码。
type Code int

const (
	// CodeGatewayBadRequest identifies gateway request validation failures.
	// CodeGatewayBadRequest 标识 gateway 请求参数校验失败。
	CodeGatewayBadRequest Code = 100101
	// CodeGatewayAuthorizationRequired identifies missing authorization headers.
	// CodeGatewayAuthorizationRequired 标识缺失认证头。
	CodeGatewayAuthorizationRequired Code = 100201
	// CodeGatewayInvalidAuthorizationScheme identifies malformed authorization headers.
	// CodeGatewayInvalidAuthorizationScheme 标识格式错误的认证头。
	CodeGatewayInvalidAuthorizationScheme Code = 100202
	// CodeGatewayAccessTokenInvalid identifies invalid access tokens.
	// CodeGatewayAccessTokenInvalid 标识非法 access token。
	CodeGatewayAccessTokenInvalid Code = 100203
	// CodeGatewayAccessTokenInactive identifies inactive access tokens.
	// CodeGatewayAccessTokenInactive 标识未激活或失效的 access token。
	CodeGatewayAccessTokenInactive Code = 100204
	// CodeGatewayAccessForbidden identifies forbidden access.
	// CodeGatewayAccessForbidden 标识无权限访问。
	CodeGatewayAccessForbidden Code = 100301
	// CodeGatewayNotReady identifies gateway readiness failures.
	// CodeGatewayNotReady 标识 gateway 未就绪。
	CodeGatewayNotReady Code = 100401
	// CodeGatewayAuthServiceUnavailable identifies identity auth dependency failures.
	// CodeGatewayAuthServiceUnavailable 标识 identity 认证依赖不可用。
	CodeGatewayAuthServiceUnavailable Code = 100601
	// CodeGatewayUpstreamTimeout identifies upstream timeout failures.
	// CodeGatewayUpstreamTimeout 标识上游超时。
	CodeGatewayUpstreamTimeout Code = 100602
	// CodeGatewayInternal identifies gateway internal failures.
	// CodeGatewayInternal 标识 gateway 内部错误。
	CodeGatewayInternal Code = 109901

	// CodeIdentityInvalidArgument identifies identity validation failures.
	// CodeIdentityInvalidArgument 标识 identity 参数校验失败。
	CodeIdentityInvalidArgument Code = 110101
	// CodeIdentityInvalidCredentials identifies invalid credentials.
	// CodeIdentityInvalidCredentials 标识非法凭证。
	CodeIdentityInvalidCredentials Code = 110201
	// CodeIdentityInvalidRefreshToken identifies invalid refresh tokens.
	// CodeIdentityInvalidRefreshToken 标识非法 refresh token。
	CodeIdentityInvalidRefreshToken Code = 110202
	// CodeIdentityRefreshTokenExpired identifies expired refresh tokens.
	// CodeIdentityRefreshTokenExpired 标识已过期 refresh token。
	CodeIdentityRefreshTokenExpired Code = 110203
	// CodeIdentitySessionRevoked identifies revoked sessions.
	// CodeIdentitySessionRevoked 标识已吊销会话。
	CodeIdentitySessionRevoked Code = 110204
	// CodeIdentityAccountPending identifies pending accounts.
	// CodeIdentityAccountPending 标识待激活账号。
	CodeIdentityAccountPending Code = 110205
	// CodeIdentityAccountDisabled identifies disabled accounts.
	// CodeIdentityAccountDisabled 标识已禁用账号。
	CodeIdentityAccountDisabled Code = 110206
	// CodeIdentityAccountLocked identifies locked accounts.
	// CodeIdentityAccountLocked 标识已锁定账号。
	CodeIdentityAccountLocked Code = 110207
	// CodeIdentityInternalCallerUnauthorized identifies unauthorized internal callers.
	// CodeIdentityInternalCallerUnauthorized 标识未通过认证的内部调用方。
	CodeIdentityInternalCallerUnauthorized Code = 110208
	// CodeIdentitySSOProviderDisabled identifies disabled SSO providers.
	// CodeIdentitySSOProviderDisabled 标识已禁用的 SSO provider。
	CodeIdentitySSOProviderDisabled Code = 110401
	// CodeIdentitySSOProviderNotReady identifies unsupported or not-ready SSO providers.
	// CodeIdentitySSOProviderNotReady 标识未就绪的 SSO provider。
	CodeIdentitySSOProviderNotReady Code = 110402
	// CodeIdentitySSOStateInvalid identifies invalid SSO state.
	// CodeIdentitySSOStateInvalid 标识非法 SSO state。
	CodeIdentitySSOStateInvalid Code = 110403
	// CodeIdentityUserNotFound identifies missing users.
	// CodeIdentityUserNotFound 标识用户不存在。
	CodeIdentityUserNotFound Code = 110501
	// CodeIdentitySessionNotFound identifies missing sessions.
	// CodeIdentitySessionNotFound 标识会话不存在。
	CodeIdentitySessionNotFound Code = 110504
	// CodeIdentityAccountNotFound identifies missing accounts.
	// CodeIdentityAccountNotFound 标识账号不存在。
	CodeIdentityAccountNotFound Code = 110505
	// CodeIdentityUsernameAlreadyExists identifies username conflicts.
	// CodeIdentityUsernameAlreadyExists 标识用户名冲突。
	CodeIdentityUsernameAlreadyExists Code = 110502
	// CodeIdentityEmailAlreadyExists identifies email conflicts.
	// CodeIdentityEmailAlreadyExists 标识邮箱冲突。
	CodeIdentityEmailAlreadyExists Code = 110503
	// CodeIdentityDependencyUnavailable identifies identity dependency readiness failures.
	// CodeIdentityDependencyUnavailable 标识 identity 依赖不可用。
	CodeIdentityDependencyUnavailable Code = 110601
	// CodeIdentityInternal identifies identity internal failures.
	// CodeIdentityInternal 标识 identity 内部错误。
	CodeIdentityInternal Code = 119901

	// CodeContentInvalidArgument identifies content validation failures.
	// CodeContentInvalidArgument 标识 content 参数校验失败。
	CodeContentInvalidArgument Code = 120101
	// CodeContentInvalidType identifies invalid content types.
	// CodeContentInvalidType 标识非法内容类型。
	CodeContentInvalidType Code = 120102
	// CodeContentInvalidStatus identifies invalid content status values.
	// CodeContentInvalidStatus 标识非法内容状态。
	CodeContentInvalidStatus Code = 120103
	// CodeContentInvalidVisibility identifies invalid visibility values.
	// CodeContentInvalidVisibility 标识非法可见性。
	CodeContentInvalidVisibility Code = 120104
	// CodeContentInvalidAIAccess identifies invalid AI access values.
	// CodeContentInvalidAIAccess 标识非法 AI 访问策略。
	CodeContentInvalidAIAccess Code = 120105
	// CodeContentInternalCallerUnauthorized identifies unauthorized internal callers.
	// CodeContentInternalCallerUnauthorized 标识未通过认证的内部调用方。
	CodeContentInternalCallerUnauthorized Code = 120201
	// CodeContentAccessForbidden identifies forbidden content access.
	// CodeContentAccessForbidden 标识内容访问被拒绝。
	CodeContentAccessForbidden Code = 120301
	// CodeContentInvalidTransition identifies invalid content state transitions.
	// CodeContentInvalidTransition 标识非法内容状态流转。
	CodeContentInvalidTransition Code = 120401
	// CodeContentNotFound identifies missing content items.
	// CodeContentNotFound 标识内容不存在。
	CodeContentNotFound Code = 120501
	// CodeContentSlugAlreadyExists identifies content slug conflicts.
	// CodeContentSlugAlreadyExists 标识内容 slug 冲突。
	CodeContentSlugAlreadyExists Code = 120502
	// CodeContentTagNotFound identifies missing content tags.
	// CodeContentTagNotFound 标识内容标签不存在。
	CodeContentTagNotFound Code = 120503
	// CodeContentTagAlreadyExists identifies tag name or slug conflicts.
	// CodeContentTagAlreadyExists 标识标签名称或 slug 冲突。
	CodeContentTagAlreadyExists Code = 120504
	// CodeContentRevisionNotFound identifies missing content revisions.
	// CodeContentRevisionNotFound 标识内容版本不存在。
	CodeContentRevisionNotFound Code = 120505
	// CodeContentTagInUse identifies tags that are still bound to content.
	// CodeContentTagInUse 标识仍被内容绑定的标签。
	CodeContentTagInUse Code = 120506
	// CodeContentInternal identifies content internal failures.
	// CodeContentInternal 标识 content 内部错误。
	CodeContentInternal Code = 129901
)

// String returns the decimal business code string.
// String 返回十进制业务错误码字符串。
func (c Code) String() string {
	return itoa(int(c))
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}

	var buf [16]byte
	index := len(buf)
	for value > 0 {
		index--
		buf[index] = byte('0' + value%10)
		value /= 10
	}

	return string(buf[index:])
}
