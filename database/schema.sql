-- Battery ERP System Database Schema
-- MySQL 8.0+ Required

CREATE DATABASE IF NOT EXISTS battery_erp CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE battery_erp;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    real_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'normal' COMMENT 'super_admin or normal',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_username (username),
    INDEX idx_role (role),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB COMMENT='System users';

-- Battery categories table
CREATE TABLE IF NOT EXISTS battery_categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    unit_price DECIMAL(10,2) NOT NULL COMMENT 'Price per kg',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_name (name),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB COMMENT='Battery categories/types';

-- Inbound orders table (purchase/receipt)
CREATE TABLE IF NOT EXISTS inbound_orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(50) NOT NULL UNIQUE,
    supplier_name VARCHAR(100) NOT NULL,
    total_amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'completed' COMMENT 'completed, cancelled',
    notes TEXT,
    created_by BIGINT UNSIGNED NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_order_no (order_no),
    INDEX idx_supplier_name (supplier_name),
    INDEX idx_status (status),
    INDEX idx_created_by (created_by),
    INDEX idx_created_at (created_at),
    
    FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB COMMENT='Inbound purchase orders';

-- Inbound order items table
CREATE TABLE IF NOT EXISTS inbound_order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    category_id BIGINT UNSIGNED NOT NULL,
    gross_weight DECIMAL(10,3) NOT NULL COMMENT 'Gross weight in kg',
    tare_weight DECIMAL(10,3) NOT NULL COMMENT 'Tare weight in kg',
    net_weight DECIMAL(10,3) NOT NULL COMMENT 'Net weight in kg',
    unit_price DECIMAL(10,2) NOT NULL COMMENT 'Price per kg',
    sub_total DECIMAL(15,2) NOT NULL COMMENT 'Net weight * unit price',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_order_id (order_id),
    INDEX idx_category_id (category_id),
    
    FOREIGN KEY (order_id) REFERENCES inbound_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES battery_categories(id)
) ENGINE=InnoDB COMMENT='Inbound order line items';

-- Outbound orders table (sales)
CREATE TABLE IF NOT EXISTS outbound_orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(50) NOT NULL UNIQUE,
    customer_name VARCHAR(100) NOT NULL,
    total_amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'completed' COMMENT 'completed, cancelled',
    notes TEXT,
    created_by BIGINT UNSIGNED NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_order_no (order_no),
    INDEX idx_customer_name (customer_name),
    INDEX idx_status (status),
    INDEX idx_created_by (created_by),
    INDEX idx_created_at (created_at),
    
    FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB COMMENT='Outbound sales orders';

-- Outbound order items table
CREATE TABLE IF NOT EXISTS outbound_order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    category_id BIGINT UNSIGNED NOT NULL,
    weight DECIMAL(10,3) NOT NULL COMMENT 'Weight in kg',
    unit_price DECIMAL(10,2) NOT NULL COMMENT 'Price per kg',
    sub_total DECIMAL(15,2) NOT NULL COMMENT 'Weight * unit price',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_order_id (order_id),
    INDEX idx_category_id (category_id),
    
    FOREIGN KEY (order_id) REFERENCES outbound_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES battery_categories(id)
) ENGINE=InnoDB COMMENT='Outbound order line items';

-- Inventory table
CREATE TABLE IF NOT EXISTS inventory (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    category_id BIGINT UNSIGNED NOT NULL UNIQUE,
    current_weight_kg DECIMAL(12,3) NOT NULL DEFAULT 0 COMMENT 'Current inventory weight in kg',
    last_inbound_at DATETIME NULL,
    last_outbound_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_category_id (category_id),
    INDEX idx_current_weight_kg (current_weight_kg),
    
    FOREIGN KEY (category_id) REFERENCES battery_categories(id)
) ENGINE=InnoDB COMMENT='Real-time inventory tracking';