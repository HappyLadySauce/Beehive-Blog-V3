-- Add scoped OAuth state support for sensitive reauthentication flows.
-- 为敏感操作重验流程增加带作用域的 OAuth state 支持。
ALTER TABLE identity.oauth_login_states
  ADD COLUMN IF NOT EXISTS purpose VARCHAR(32) NOT NULL DEFAULT 'login',
  ADD COLUMN IF NOT EXISTS subject_user_id BIGINT NULL;

CREATE INDEX IF NOT EXISTS idx_identity_oauth_login_states_purpose_subject
  ON identity.oauth_login_states (purpose, subject_user_id);
