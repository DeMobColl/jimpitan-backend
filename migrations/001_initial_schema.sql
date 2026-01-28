-- Migration: Create initial schema for Jimpitan
-- This migration sets up all tables based on the Google Sheets structure

-- Users Table
CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(20) PRIMARY KEY COMMENT 'USR-001, USR-002, ...',
  name VARCHAR(255) NOT NULL,
  role ENUM('admin', 'petugas') NOT NULL,
  username VARCHAR(100) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL COMMENT 'SHA-256 hash',
  token VARCHAR(255) UNIQUE,
  token_expiry DATETIME,
  last_login DATETIME,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME,
  INDEX idx_username (username),
  INDEX idx_token (token),
  INDEX idx_role (role),
  INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Customers Table
CREATE TABLE IF NOT EXISTS customers (
  id VARCHAR(20) PRIMARY KEY COMMENT 'CUST-001, CUST-002, ...',
  blok VARCHAR(50) NOT NULL COMMENT 'Block/ID Number',
  nama VARCHAR(255) NOT NULL COMMENT 'Full Name',
  qr_hash VARCHAR(10) UNIQUE NOT NULL COMMENT '10-character QR identifier',
  total_setoran DECIMAL(12, 2) DEFAULT 0 COMMENT 'Total deposits sum',
  last_transaction DATETIME,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME,
  INDEX idx_blok (blok),
  INDEX idx_qr_hash (qr_hash),
  INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Transactions Table
CREATE TABLE IF NOT EXISTS transactions (
  id VARCHAR(20) PRIMARY KEY COMMENT 'Transaction ID',
  timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  customer_id VARCHAR(20) NOT NULL,
  blok VARCHAR(50) NOT NULL COMMENT 'Denormalized customer blok',
  nama VARCHAR(255) NOT NULL COMMENT 'Denormalized customer name',
  nominal DECIMAL(12, 2) NOT NULL,
  user_id VARCHAR(20) NOT NULL,
  petugas VARCHAR(255) NOT NULL COMMENT 'Staff name',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME,
  FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
  INDEX idx_customer_id (customer_id),
  INDEX idx_user_id (user_id),
  INDEX idx_timestamp (timestamp),
  INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Sessions Table
CREATE TABLE IF NOT EXISTS sessions (
  id VARCHAR(255) PRIMARY KEY,
  user_id VARCHAR(20) NOT NULL,
  token VARCHAR(255) NOT NULL UNIQUE,
  expires_at DATETIME NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  INDEX idx_user_id (user_id),
  INDEX idx_token (token),
  INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Config Table
CREATE TABLE IF NOT EXISTS config (
  id VARCHAR(50) PRIMARY KEY,
  petugas_web_login_enabled BOOLEAN DEFAULT true,
  mobile_app_version VARCHAR(20),
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert default config
INSERT IGNORE INTO config (id, petugas_web_login_enabled, mobile_app_version) 
VALUES ('default', true, '1.0.0');
