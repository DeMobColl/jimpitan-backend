-- Seed: Insert first admin user
-- Username: admin
-- Password: admin123 (SHA-256 hash)
-- Role: admin

INSERT INTO users (
  id,
  name,
  role,
  username,
  password_hash,
  created_at,
  updated_at
) VALUES (
  'USR-001',
  'Administrator',
  'admin',
  'admin',
  '240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9',
  NOW(),
  NOW()
) ON DUPLICATE KEY UPDATE
  updated_at = NOW();

-- Insert petugas user for testing
INSERT INTO users (
  id,
  name,
  role,
  username,
  password_hash,
  created_at,
  updated_at
) VALUES (
  'USR-002',
  'Petugas Default',
  'petugas',
  'petugas',
  '2dad904f71aa0dcf6ea1addaa084a5865ffe448e4d3f900668e1cc7e7b6153d7',
  NOW(),
  NOW()
) ON DUPLICATE KEY UPDATE
  updated_at = NOW();
