-- Seed: Insert test customers with transactions
-- Generate 5 sample customers with QR hashes and transactions

-- Insert customers
INSERT INTO customers (id, blok, nama, qr_hash, total_setoran, last_transaction, created_at, updated_at)
VALUES
  ('CUST-001', '1A', 'Siti Nurhaliza', 'QR001ABC', 500000, NOW() - INTERVAL 2 DAY, NOW(), NOW()),
  ('CUST-002', '1B', 'Bambang Irawan', 'QR002XYZ', 1200000, NOW() - INTERVAL 1 DAY, NOW(), NOW()),
  ('CUST-003', '2A', 'Rina Wijaya', 'QR003DEF', 750000, NOW() - INTERVAL 5 HOUR, NOW(), NOW()),
  ('CUST-004', '2B', 'Ahmad Yusuf', 'QR004GHI', 300000, NOW() - INTERVAL 3 HOUR, NOW(), NOW()),
  ('CUST-005', '3A', 'Dwi Cahyani', 'QR005JKL', 900000, NOW() - INTERVAL 1 HOUR, NOW(), NOW())
ON DUPLICATE KEY UPDATE
  updated_at = NOW();

-- Insert sample transactions
INSERT INTO transactions (id, timestamp, customer_id, blok, nama, nominal, user_id, petugas, created_at)
VALUES
  ('TXN-001', NOW() - INTERVAL 2 DAY, 'CUST-001', '1A', 'Siti Nurhaliza', 250000, 'USR-001', 'Admin', NOW()),
  ('TXN-002', NOW() - INTERVAL 2 DAY, 'CUST-001', '1A', 'Siti Nurhaliza', 250000, 'USR-002', 'Petugas', NOW()),
  ('TXN-003', NOW() - INTERVAL 1 DAY, 'CUST-002', '1B', 'Bambang Irawan', 600000, 'USR-001', 'Admin', NOW()),
  ('TXN-004', NOW() - INTERVAL 1 DAY, 'CUST-002', '1B', 'Bambang Irawan', 600000, 'USR-002', 'Petugas', NOW()),
  ('TXN-005', NOW() - INTERVAL 5 HOUR, 'CUST-003', '2A', 'Rina Wijaya', 750000, 'USR-001', 'Admin', NOW()),
  ('TXN-006', NOW() - INTERVAL 3 HOUR, 'CUST-004', '2B', 'Ahmad Yusuf', 300000, 'USR-002', 'Petugas', NOW()),
  ('TXN-007', NOW() - INTERVAL 1 HOUR, 'CUST-005', '3A', 'Dwi Cahyani', 450000, 'USR-001', 'Admin', NOW()),
  ('TXN-008', NOW() - INTERVAL 30 MINUTE, 'CUST-005', '3A', 'Dwi Cahyani', 450000, 'USR-002', 'Petugas', NOW())
ON DUPLICATE KEY UPDATE
  timestamp = VALUES(timestamp);
