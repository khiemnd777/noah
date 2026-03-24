-- 1. Tạo user admin nếu chưa có
INSERT INTO users (
  email,
  password,
  name,
  active,
  avatar,
  provider,
  provider_id,
  phone,
  created_at,
  updated_at
)
SELECT 
  'admin@noah.io',
  '$2a$10$W4WXXgBmh4mKSgfRpHmEsusc2x7t7JP3E5ZcrK0k2k2liIEFjTCE2', -- bcrypt: "sa"
  'Administrator',
  true,
  'https://api.dicebear.com/9.x/initials/svg?seed=Admin',
  NULL,
  NULL,
  NULL,
  NOW(),
  NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM users WHERE email = 'admin@noah.io'
);

-- 2. Gán user admin ↔ role admin
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
JOIN roles r ON r.role_name = 'admin'
WHERE u.email = 'admin@noah.io'
ON CONFLICT DO NOTHING;
