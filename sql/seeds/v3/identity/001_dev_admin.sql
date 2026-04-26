-- v3 identity development admin seed.
-- v3 identity 开发环境管理员种子数据。
--
-- Account: admin@beehive.local / Admin@123456
-- 账号：admin@beehive.local / Admin@123456

WITH matched_user AS (
  SELECT id
  FROM identity.users
  WHERE username = 'admin' OR LOWER(email) = 'admin@beehive.local'
  ORDER BY
    CASE WHEN LOWER(email) = 'admin@beehive.local' THEN 0 ELSE 1 END,
    id
  LIMIT 1
),
updated_user AS (
  UPDATE identity.users
  SET
    username = 'admin',
    email = 'admin@beehive.local',
    nickname = 'Admin',
    role = 'admin',
    status = 'active',
    updated_at = NOW()
  WHERE id IN (SELECT id FROM matched_user)
  RETURNING id
),
inserted_user AS (
  INSERT INTO identity.users (
    username,
    email,
    nickname,
    role,
    status,
    created_at,
    updated_at
  )
  SELECT
    'admin',
    'admin@beehive.local',
    'Admin',
    'admin',
    'active',
    NOW(),
    NOW()
  WHERE NOT EXISTS (SELECT 1 FROM updated_user)
  RETURNING id
),
admin_user AS (
  SELECT id FROM updated_user
  UNION ALL
  SELECT id FROM inserted_user
)
INSERT INTO identity.credential_locals (
  user_id,
  password_hash,
  password_updated_at,
  created_at,
  updated_at
)
SELECT
  id,
  '$2a$12$o3oPCKiKhrvlRn3zmibIDedhEWndxc2jiU.OQkOxDcBUWJ.SsuchW',
  NOW(),
  NOW(),
  NOW()
FROM admin_user
ON CONFLICT (user_id) DO UPDATE SET
  password_hash = EXCLUDED.password_hash,
  password_updated_at = NOW(),
  updated_at = NOW();
