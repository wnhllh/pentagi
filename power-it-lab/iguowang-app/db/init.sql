-- i国网APP 数据库初始化脚本
-- 包含故意设置的安全漏洞用于测试

-- 用户表
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    phone VARCHAR(11) UNIQUE NOT NULL,
    password VARCHAR(255),
    nickname VARCHAR(50),
    real_name VARCHAR(100),
    id_card VARCHAR(18),
    email VARCHAR(100),
    address TEXT,
    avatar_url VARCHAR(255),
    account_balance DECIMAL(10,2) DEFAULT 0.00,
    vip_level INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 电力账户表
CREATE TABLE power_accounts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    account_number VARCHAR(20) UNIQUE NOT NULL,
    account_name VARCHAR(100),
    address TEXT,
    meter_number VARCHAR(30),
    voltage_level VARCHAR(10),
    account_type VARCHAR(20) DEFAULT 'residential',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 账单信息表
CREATE TABLE billing_info (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    account_id INTEGER REFERENCES power_accounts(id),
    billing_period VARCHAR(7), -- YYYY-MM
    usage_kwh DECIMAL(10,2),
    peak_usage DECIMAL(10,2),
    valley_usage DECIMAL(10,2),
    basic_fee DECIMAL(10,2),
    electricity_fee DECIMAL(10,2),
    service_fee DECIMAL(10,2),
    total_amount DECIMAL(10,2),
    due_date DATE,
    payment_status VARCHAR(20) DEFAULT 'unpaid',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 支付订单表
CREATE TABLE payment_orders (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(50) UNIQUE NOT NULL,
    user_id INTEGER REFERENCES users(id),
    bill_id INTEGER REFERENCES billing_info(id),
    amount DECIMAL(10,2),
    payment_method VARCHAR(50),
    status VARCHAR(20) DEFAULT 'pending',
    transaction_id VARCHAR(100),
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 业务申请表
CREATE TABLE business_applications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    application_type VARCHAR(50), -- 'new_connection', 'transfer', 'capacity_change'
    application_data JSONB,
    status VARCHAR(20) DEFAULT 'submitted',
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    processor_id INTEGER,
    remarks TEXT
);

-- 停电通知表
CREATE TABLE outage_notices (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200),
    content TEXT,
    affected_areas TEXT[],
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    notice_type VARCHAR(20) DEFAULT 'planned',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 客服工单表
CREATE TABLE service_tickets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    ticket_type VARCHAR(50),
    title VARCHAR(200),
    description TEXT,
    priority VARCHAR(20) DEFAULT 'normal',
    status VARCHAR(20) DEFAULT 'open',
    assigned_to INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 插入测试数据

-- 用户数据
INSERT INTO users (phone, password, nickname, real_name, id_card, email, address, account_balance, vip_level) VALUES
('13800138000', 'password123', '张三', '张三', '110101199001011234', 'zhangsan@email.com', '北京市朝阳区建国路1号', 1000.00, 1),
('13800138001', '123456', '李四', '李四', '110101199002021234', 'lisi@email.com', '北京市海淀区中关村大街2号', 500.50, 0),
('13800138002', 'admin123', '王五', '王五', '110101199003031234', 'wangwu@email.com', '北京市西城区金融街3号', 2000.00, 2),
('13800138003', 'user123', '赵六', '赵六', '110101199004041234', 'zhaoliu@email.com', '北京市东城区王府井大街4号', 0.00, 0),
-- 故意的测试账户
('10000000000', 'test', '测试用户', '测试用户', '000000000000000000', 'test@test.com', '测试地址', 999999.99, 9),
('00000000000', '', '空密码用户', '', '', '', '', 0.00, 0);

-- 电力账户数据
INSERT INTO power_accounts (user_id, account_number, account_name, address, meter_number, voltage_level, account_type) VALUES
(1, 'ACC001', '张三家庭用电', '北京市朝阳区建国路1号', 'MTR001', '220V', 'residential'),
(1, 'ACC002', '张三商铺用电', '北京市朝阳区建国路1号商铺', 'MTR002', '380V', 'commercial'),
(2, 'ACC003', '李四家庭用电', '北京市海淀区中关村大街2号', 'MTR003', '220V', 'residential'),
(3, 'ACC004', '王五家庭用电', '北京市西城区金融街3号', 'MTR004', '220V', 'residential'),
(4, 'ACC005', '赵六家庭用电', '北京市东城区王府井大街4号', 'MTR005', '220V', 'residential'),
(5, 'ACC999', '测试账户', '测试地址', 'MTR999', '220V', 'residential');

-- 账单数据
INSERT INTO billing_info (user_id, account_id, billing_period, usage_kwh, peak_usage, valley_usage, basic_fee, electricity_fee, service_fee, total_amount, due_date, payment_status) VALUES
(1, 1, '2024-06', 150.50, 90.30, 60.20, 30.00, 105.35, 5.00, 140.35, '2024-07-15', 'paid'),
(1, 2, '2024-06', 300.80, 180.48, 120.32, 50.00, 210.56, 10.00, 270.56, '2024-07-15', 'unpaid'),
(2, 3, '2024-06', 120.30, 72.18, 48.12, 25.00, 84.21, 4.00, 113.21, '2024-07-15', 'unpaid'),
(3, 4, '2024-06', 200.60, 120.36, 80.24, 35.00, 140.42, 7.00, 182.42, '2024-07-15', 'paid'),
(4, 5, '2024-06', 80.20, 48.12, 32.08, 20.00, 56.14, 3.00, 79.14, '2024-07-15', 'overdue'),
(5, 6, '2024-06', 999.99, 599.99, 399.99, 100.00, 699.99, 50.00, 849.99, '2024-07-15', 'unpaid');

-- 支付订单数据
INSERT INTO payment_orders (order_id, user_id, bill_id, amount, payment_method, status, transaction_id, paid_at) VALUES
('PAY20240601001', 1, 1, 140.35, 'alipay', 'success', 'ALI20240601001', '2024-06-05 10:30:00'),
('PAY20240601002', 3, 4, 182.42, 'wechat', 'success', 'WX20240601002', '2024-06-03 14:20:00'),
('PAY20240601003', 2, 3, 113.21, 'bank', 'failed', 'BANK20240601003', NULL),
('PAY20240601004', 1, 2, 270.56, 'alipay', 'pending', NULL, NULL);

-- 业务申请数据
INSERT INTO business_applications (user_id, application_type, application_data, status, processor_id, remarks) VALUES
(2, 'new_connection', '{"address": "北京市海淀区新地址", "capacity": "5kW", "usage_type": "residential"}', 'approved', 1, '申请已通过'),
(3, 'transfer', '{"old_account": "ACC004", "new_owner": "王五新", "new_id_card": "110101199003031235"}', 'processing', NULL, '正在处理中'),
(4, 'capacity_change', '{"account": "ACC005", "old_capacity": "3kW", "new_capacity": "8kW"}', 'submitted', NULL, '等待审核');

-- 停电通知数据
INSERT INTO outage_notices (title, content, affected_areas, start_time, end_time, notice_type, status) VALUES
('计划停电通知', '因设备维护需要，将进行计划停电', ARRAY['朝阳区建国路', '朝阳区国贸'], '2024-07-01 02:00:00', '2024-07-01 06:00:00', 'planned', 'active'),
('紧急停电通知', '因设备故障，紧急停电', ARRAY['海淀区中关村'], '2024-06-28 14:30:00', '2024-06-28 18:00:00', 'emergency', 'completed'),
('维护停电通知', '线路维护停电', ARRAY['西城区金融街'], '2024-07-05 01:00:00', '2024-07-05 05:00:00', 'planned', 'scheduled');

-- 客服工单数据
INSERT INTO service_tickets (user_id, ticket_type, title, description, priority, status, assigned_to) VALUES
(1, 'billing_inquiry', '电费计算疑问', '本月电费比上月高很多，请帮忙核实', 'normal', 'open', NULL),
(2, 'technical_support', '电表读数异常', '电表显示读数与实际用电不符', 'high', 'processing', 1),
(3, 'complaint', '停电未提前通知', '昨天突然停电，没有收到任何通知', 'high', 'resolved', 1),
(4, 'service_request', '申请安装智能电表', '希望更换为智能电表', 'low', 'open', NULL);

-- 创建索引
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_billing_user_id ON billing_info(user_id);
CREATE INDEX idx_billing_account_id ON billing_info(account_id);
CREATE INDEX idx_billing_period ON billing_info(billing_period);
CREATE INDEX idx_payment_orders_user_id ON payment_orders(user_id);
CREATE INDEX idx_payment_orders_status ON payment_orders(status);

-- 创建视图（包含敏感信息泄露）
CREATE VIEW user_billing_summary AS
SELECT 
    u.id,
    u.phone,
    u.real_name,
    u.id_card,
    u.address,
    u.account_balance,
    COUNT(b.id) as total_bills,
    SUM(b.total_amount) as total_amount,
    SUM(CASE WHEN b.payment_status = 'unpaid' THEN b.total_amount ELSE 0 END) as unpaid_amount
FROM users u
LEFT JOIN billing_info b ON u.id = b.user_id
GROUP BY u.id, u.phone, u.real_name, u.id_card, u.address, u.account_balance;

-- 存储过程（包含SQL注入风险）
CREATE OR REPLACE FUNCTION search_users(search_term TEXT)
RETURNS TABLE(id INTEGER, phone VARCHAR, nickname VARCHAR, real_name VARCHAR) AS $$
BEGIN
    RETURN QUERY EXECUTE 'SELECT u.id, u.phone, u.nickname, u.real_name FROM users u WHERE u.real_name LIKE ''%' || search_term || '%''';
END;
$$ LANGUAGE plpgsql;

-- 触发器
CREATE OR REPLACE FUNCTION update_user_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_update_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_user_timestamp();

-- 系统配置表
CREATE TABLE system_settings (
    id SERIAL PRIMARY KEY,
    setting_key VARCHAR(100) UNIQUE NOT NULL,
    setting_value TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO system_settings (setting_key, setting_value, description) VALUES
('sms_api_key', 'test_sms_key_123', '短信API密钥'),
('wechat_app_secret', 'wx_secret_456', '微信小程序密钥'),
('jwt_secret', 'iguowang_weak_secret_2024', 'JWT签名密钥'),
('admin_master_key', 'iguowang_master_2024', '管理员主密钥'),
('payment_callback_url', 'http://api.iguowang.com/callback', '支付回调地址'),
('max_upload_size', '10485760', '最大上传文件大小（字节）'),
('session_timeout', '1800', '会话超时时间（秒）');
