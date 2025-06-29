-- Battery ERP System Database Initialization
-- This script will create the database, tables, and initial data

-- Create database
CREATE DATABASE IF NOT EXISTS battery_erp CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Source the schema and seed files
SOURCE schema.sql;
SOURCE seed.sql;

-- Display initialization completion message
SELECT '========================================' as '';
SELECT 'Battery ERP Database Initialized Successfully!' as '';
SELECT '========================================' as '';
SELECT 'Default Login Credentials:' as '';
SELECT '  Admin: admin / admin123' as '';
SELECT '  Operator: operator / admin123' as '';
SELECT '========================================' as '';