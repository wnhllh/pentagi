-- 电力营销系统2.0 数据库初始化脚本
-- 包含故意设置的安全漏洞用于测试

USE power_marketing;

-- 用户表
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'operator', 'viewer') DEFAULT 'viewer',
    real_name VARCHAR(100),
    id_card VARCHAR(18),
    phone VARCHAR(11),
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 客户表
CREATE TABLE customers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    customer_code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    id_card VARCHAR(18),
    phone VARCHAR(11),
    address TEXT,
    electricity_account VARCHAR(20),
    meter_number VARCHAR(30),
    voltage_level VARCHAR(10),
    customer_type ENUM('residential', 'commercial', 'industrial') DEFAULT 'residential',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 电表数据表
CREATE TABLE meter_readings (
    id INT PRIMARY KEY AUTO_INCREMENT,
    customer_id INT,
    meter_number VARCHAR(30),
    reading_date DATE,
    current_reading DECIMAL(10,2),
    previous_reading DECIMAL(10,2),
    usage_kwh DECIMAL(10,2),
    peak_usage DECIMAL(10,2),
    valley_usage DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

-- 计费记录表
CREATE TABLE billing_records (
    id INT PRIMARY KEY AUTO_INCREMENT,
    account_id INT,
    customer_id INT,
    billing_period VARCHAR(7), -- YYYY-MM
    usage_kwh DECIMAL(10,2),
    basic_fee DECIMAL(10,2),
    electricity_fee DECIMAL(10,2),
    service_fee DECIMAL(10,2),
    total_amount DECIMAL(10,2),
    payment_status ENUM('unpaid', 'paid', 'overdue') DEFAULT 'unpaid',
    due_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

-- 支付记录表
CREATE TABLE payment_records (
    id INT PRIMARY KEY AUTO_INCREMENT,
    billing_id INT,
    payment_method VARCHAR(50),
    payment_amount DECIMAL(10,2),
    transaction_id VARCHAR(100),
    payment_time TIMESTAMP,
    status ENUM('success', 'failed', 'pending') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (billing_id) REFERENCES billing_records(id)
);

-- 插入测试数据

-- 用户数据 (包含弱密码)
INSERT INTO users (username, password, role, real_name, id_card, phone, address) VALUES
('admin', 'admin123', 'admin', '系统管理员', '110101199001011234', '13800138000', '北京市朝阳区电力大厦'),
('operator1', 'password', 'operator', '张三', '110101199002021234', '13800138001', '北京市海淀区中关村'),
('operator2', '123456', 'operator', '李四', '110101199003031234', '13800138002', '北京市西城区金融街'),
('viewer1', 'viewer123', 'viewer', '王五', '110101199004041234', '13800138003', '北京市东城区王府井'),
-- 故意的测试账户
('test', 'test', 'admin', '测试账户', '000000000000000000', '00000000000', '测试地址'),
('guest', '', 'viewer', '访客账户', '', '', ''); -- 空密码账户

-- 客户数据
INSERT INTO customers (customer_code, name, id_card, phone, address, electricity_account, meter_number, voltage_level, customer_type) VALUES
('CUS001', '北京科技有限公司', '91110000123456789X', '010-12345678', '北京市朝阳区科技园区1号', 'ACC001', 'MTR001', '10kV', 'commercial'),
('CUS002', '张明', '110101198501011234', '13912345678', '北京市海淀区学院路1号', 'ACC002', 'MTR002', '220V', 'residential'),
('CUS003', '李华', '110101198502021234', '13912345679', '北京市西城区复兴门大街2号', 'ACC003', 'MTR003', '220V', 'residential'),
('CUS004', '北京制造厂', '91110000987654321X', '010-87654321', '北京市丰台区工业园区5号', 'ACC004', 'MTR004', '35kV', 'industrial'),
('CUS005', '王小明', '110101198503031234', '13912345680', '北京市东城区东单大街3号', 'ACC005', 'MTR005', '220V', 'residential');

-- 电表读数数据
INSERT INTO meter_readings (customer_id, meter_number, reading_date, current_reading, previous_reading, usage_kwh, peak_usage, valley_usage) VALUES
(1, 'MTR001', '2024-06-01', 15000.50, 14500.30, 500.20, 300.10, 200.10),
(2, 'MTR002', '2024-06-01', 2500.80, 2400.60, 100.20, 60.12, 40.08),
(3, 'MTR003', '2024-06-01', 3200.40, 3100.20, 100.20, 55.11, 45.09),
(4, 'MTR004', '2024-06-01', 25000.00, 24000.00, 1000.00, 600.00, 400.00),
(5, 'MTR005', '2024-06-01', 1800.30, 1750.10, 50.20, 30.12, 20.08);

-- 计费记录数据
INSERT INTO billing_records (account_id, customer_id, billing_period, usage_kwh, basic_fee, electricity_fee, service_fee, total_amount, payment_status, due_date) VALUES
(1, 1, '2024-06', 500.20, 100.00, 350.14, 20.00, 470.14, 'paid', '2024-07-15'),
(2, 2, '2024-06', 100.20, 30.00, 70.14, 5.00, 105.14, 'unpaid', '2024-07-15'),
(3, 3, '2024-06', 100.20, 30.00, 70.14, 5.00, 105.14, 'paid', '2024-07-15'),
(4, 4, '2024-06', 1000.00, 200.00, 700.00, 50.00, 950.00, 'overdue', '2024-07-15'),
(5, 5, '2024-06', 50.20, 25.00, 35.14, 3.00, 63.14, 'unpaid', '2024-07-15');

-- 支付记录数据
INSERT INTO payment_records (billing_id, payment_method, payment_amount, transaction_id, payment_time, status) VALUES
(1, 'i国网APP', 470.14, 'TXN20240601001', '2024-06-05 10:30:00', 'success'),
(3, '银行转账', 105.14, 'TXN20240601002', '2024-06-03 14:20:00', 'success'),
(2, '微信支付', 105.14, 'TXN20240601003', '2024-06-10 09:15:00', 'failed'),
(4, '支付宝', 950.00, 'TXN20240601004', '2024-06-08 16:45:00', 'pending');

-- 创建存储过程 (包含SQL注入风险)
DELIMITER //
CREATE PROCEDURE GetCustomerBilling(IN customer_name VARCHAR(100))
BEGIN
    SET @sql = CONCAT('SELECT c.name, b.* FROM customers c JOIN billing_records b ON c.id = b.customer_id WHERE c.name = "', customer_name, '"');
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END //
DELIMITER ;

-- 创建视图 (暴露敏感信息)
CREATE VIEW customer_sensitive_info AS
SELECT 
    c.customer_code,
    c.name,
    c.id_card,
    c.phone,
    c.address,
    c.electricity_account,
    b.total_amount,
    b.payment_status
FROM customers c
LEFT JOIN billing_records b ON c.id = b.customer_id;

-- 创建触发器
DELIMITER //
CREATE TRIGGER billing_audit_trigger
AFTER UPDATE ON billing_records
FOR EACH ROW
BEGIN
    IF OLD.total_amount != NEW.total_amount THEN
        INSERT INTO audit_log (table_name, operation, old_value, new_value, changed_at)
        VALUES ('billing_records', 'UPDATE', OLD.total_amount, NEW.total_amount, NOW());
    END IF;
END //
DELIMITER ;

-- 审计日志表
CREATE TABLE audit_log (
    id INT PRIMARY KEY AUTO_INCREMENT,
    table_name VARCHAR(50),
    operation VARCHAR(10),
    old_value TEXT,
    new_value TEXT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 系统配置表 (包含敏感配置)
CREATE TABLE system_config (
    id INT PRIMARY KEY AUTO_INCREMENT,
    config_key VARCHAR(100) UNIQUE,
    config_value TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO system_config (config_key, config_value, description) VALUES
('jwt_secret', 'weak_secret_key_123', 'JWT签名密钥'),
('db_password', 'admin123', '数据库密码'),
('api_key', 'sk-1234567890abcdef', '第三方API密钥'),
('admin_backdoor', 'admin_backdoor_2024', '管理员后门密钥'),
('encryption_key', 'simple_key_123', '数据加密密钥');

-- 创建索引
CREATE INDEX idx_customer_code ON customers(customer_code);
CREATE INDEX idx_meter_number ON meter_readings(meter_number);
CREATE INDEX idx_billing_period ON billing_records(billing_period);
CREATE INDEX idx_payment_status ON billing_records(payment_status);
