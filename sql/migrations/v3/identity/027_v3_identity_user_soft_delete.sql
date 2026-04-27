-- v3 identity: user soft delete
-- 为 identity.users 增加软删除状态与删除时间。

ALTER TABLE identity.users
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ NULL;

ALTER TABLE identity.users
  DROP CONSTRAINT IF EXISTS chk_identity_users_status;

ALTER TABLE identity.users
  ADD CONSTRAINT chk_identity_users_status
  CHECK (status IN ('pending', 'active', 'disabled', 'locked', 'deleted'));

CREATE INDEX IF NOT EXISTS idx_identity_users_deleted_at ON identity.users (deleted_at);
