package service

const (
	maxAvatarURLLength  = 2048
	maxAuditEventLength = 64
)

// UserManagementService handles identity user administration and self-service mutations.
// UserManagementService 处理 identity 用户管理与用户自助修改。
type UserManagementService struct {
	deps Dependencies
}

// NewUserManagementService creates a UserManagementService.
// NewUserManagementService 创建 UserManagementService。
func NewUserManagementService(deps Dependencies) *UserManagementService {
	return &UserManagementService{deps: deps}
}
