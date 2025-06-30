-- 电力ERP系统数据库初始化脚本 (SAP风格)
-- 包含故意设置的安全漏洞用于测试

-- SAP用户表
CREATE TABLE sap_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    department VARCHAR(100),
    email VARCHAR(100),
    client VARCHAR(10) DEFAULT '100',
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    failed_attempts INTEGER DEFAULT 0,
    locked BOOLEAN DEFAULT FALSE
);

-- HR员工表
CREATE TABLE hr_employees (
    id SERIAL PRIMARY KEY,
    employee_id VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    department VARCHAR(100),
    position VARCHAR(100),
    salary DECIMAL(12,2),
    id_card VARCHAR(18),
    phone VARCHAR(20),
    email VARCHAR(100),
    address TEXT,
    hire_date DATE,
    status VARCHAR(20) DEFAULT 'active',
    manager_id VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 财务数据表
CREATE TABLE financial_data (
    id SERIAL PRIMARY KEY,
    company_code VARCHAR(10) NOT NULL,
    account_number VARCHAR(20) NOT NULL,
    account_name VARCHAR(100),
    account_type VARCHAR(50),
    amount DECIMAL(15,2),
    currency VARCHAR(10) DEFAULT 'CNY',
    fiscal_year VARCHAR(4),
    period VARCHAR(3),
    posting_date DATE,
    document_number VARCHAR(20),
    reference VARCHAR(50),
    created_by VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 采购订单表
CREATE TABLE purchase_orders (
    id SERIAL PRIMARY KEY,
    po_number VARCHAR(20) UNIQUE NOT NULL,
    vendor_code VARCHAR(20),
    vendor_name VARCHAR(100),
    total_amount DECIMAL(12,2),
    currency VARCHAR(10) DEFAULT 'CNY',
    status VARCHAR(20) DEFAULT 'created',
    created_by VARCHAR(50),
    approved_by VARCHAR(50),
    created_date DATE,
    delivery_date DATE
);

-- 供应商表
CREATE TABLE vendors (
    id SERIAL PRIMARY KEY,
    vendor_code VARCHAR(20) UNIQUE NOT NULL,
    vendor_name VARCHAR(100) NOT NULL,
    contact_person VARCHAR(50),
    phone VARCHAR(20),
    email VARCHAR(100),
    address TEXT,
    bank_account VARCHAR(50),
    tax_number VARCHAR(30),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 物料主数据表
CREATE TABLE materials (
    id SERIAL PRIMARY KEY,
    material_code VARCHAR(20) UNIQUE NOT NULL,
    material_name VARCHAR(100) NOT NULL,
    material_type VARCHAR(50),
    unit VARCHAR(10),
    standard_price DECIMAL(10,2),
    currency VARCHAR(10) DEFAULT 'CNY',
    plant VARCHAR(10),
    storage_location VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 库存表
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    material_code VARCHAR(20),
    plant VARCHAR(10),
    storage_location VARCHAR(10),
    quantity DECIMAL(10,3),
    unit VARCHAR(10),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (material_code) REFERENCES materials(material_code)
);

-- 插入测试数据

-- SAP用户数据 (包含弱密码和默认账户)
INSERT INTO sap_users (username, password, role, department, email, client) VALUES
('ADMIN', 'ADMIN123', 'SUPER_ADMIN', 'IT', 'admin@power-erp.com', '100'),
('SAP*', '06071992', 'SYSTEM_ADMIN', 'BASIS', 'sap@power-erp.com', '000'), -- SAP经典默认账户
('DDIC', 'DDIC', 'DEVELOPER', 'BASIS', 'ddic@power-erp.com', '100'), -- 数据字典用户
('DEVELOPER', 'developer123', 'DEVELOPER', 'IT', 'dev@power-erp.com', '100'),
('FINANCE_MGR', 'finance2024', 'FINANCE_MANAGER', 'FINANCE', 'finance@power-erp.com', '100'),
('HR_MGR', 'hr123456', 'HR_MANAGER', 'HR', 'hr@power-erp.com', '100'),
('PURCHASE_MGR', 'purchase123', 'PURCHASE_MANAGER', 'PROCUREMENT', 'purchase@power-erp.com', '100'),
('USER001', 'password', 'USER', 'OPERATIONS', 'user001@power-erp.com', '100'),
('USER002', '123456', 'USER', 'OPERATIONS', 'user002@power-erp.com', '100'),
('GUEST', '', 'GUEST', 'GUEST', 'guest@power-erp.com', '100'); -- 空密码账户

-- HR员工数据
INSERT INTO hr_employees (employee_id, name, department, position, salary, id_card, phone, email, hire_date, manager_id) VALUES
('EMP001', '张三', '财务部', '财务经理', 15000.00, '110101198001011234', '13800138001', 'zhangsan@power-erp.com', '2020-01-15', NULL),
('EMP002', '李四', '人力资源部', 'HR经理', 14000.00, '110101198002021234', '13800138002', 'lisi@power-erp.com', '2020-03-01', NULL),
('EMP003', '王五', '采购部', '采购经理', 13000.00, '110101198003031234', '13800138003', 'wangwu@power-erp.com', '2020-05-10', NULL),
('EMP004', '赵六', '运营部', '运营专员', 8000.00, '110101198004041234', '13800138004', 'zhaoliu@power-erp.com', '2021-01-20', 'EMP001'),
('EMP005', '钱七', 'IT部', '系统管理员', 12000.00, '110101198005051234', '13800138005', 'qianqi@power-erp.com', '2021-03-15', NULL),
('EMP006', '孙八', '财务部', '会计', 7000.00, '110101198006061234', '13800138006', 'sunba@power-erp.com', '2021-06-01', 'EMP001'),
('EMP007', '周九', '采购部', '采购专员', 6500.00, '110101198007071234', '13800138007', 'zhoujiu@power-erp.com', '2022-01-10', 'EMP003'),
('EMP008', '吴十', '人力资源部', 'HR专员', 6000.00, '110101198008081234', '13800138008', 'wushi@power-erp.com', '2022-03-20', 'EMP002');

-- 财务数据
INSERT INTO financial_data (company_code, account_number, account_name, account_type, amount, fiscal_year, period, posting_date, document_number, created_by) VALUES
('1000', '1000000', '现金', 'ASSET', 5000000.00, '2024', '006', '2024-06-30', 'DOC001', 'FINANCE_MGR'),
('1000', '1100000', '银行存款', 'ASSET', 50000000.00, '2024', '006', '2024-06-30', 'DOC002', 'FINANCE_MGR'),
('1000', '1300000', '应收账款', 'ASSET', 8000000.00, '2024', '006', '2024-06-30', 'DOC003', 'FINANCE_MGR'),
('1000', '2000000', '应付账款', 'LIABILITY', -3000000.00, '2024', '006', '2024-06-30', 'DOC004', 'FINANCE_MGR'),
('1000', '3000000', '实收资本', 'EQUITY', -40000000.00, '2024', '006', '2024-06-30', 'DOC005', 'FINANCE_MGR'),
('1000', '4000000', '主营业务收入', 'REVENUE', -15000000.00, '2024', '006', '2024-06-30', 'DOC006', 'FINANCE_MGR'),
('1000', '5000000', '主营业务成本', 'EXPENSE', 8000000.00, '2024', '006', '2024-06-30', 'DOC007', 'FINANCE_MGR'),
('1000', '6000000', '管理费用', 'EXPENSE', 2000000.00, '2024', '006', '2024-06-30', 'DOC008', 'FINANCE_MGR');

-- 供应商数据
INSERT INTO vendors (vendor_code, vendor_name, contact_person, phone, email, address, bank_account, tax_number) VALUES
('V001', '北京电力设备有限公司', '张经理', '010-12345678', 'zhang@bjdl.com', '北京市朝阳区电力大厦', '1234567890123456789', '91110000123456789X'),
('V002', '上海智能电网科技公司', '李总监', '021-87654321', 'li@shzn.com', '上海市浦东新区科技园', '9876543210987654321', '91310000987654321Y'),
('V003', '广州电缆制造厂', '王主任', '020-11111111', 'wang@gzdl.com', '广州市天河区工业园', '1111111111111111111', '91440000111111111Z'),
('V004', '深圳新能源设备公司', '赵总', '0755-22222222', 'zhao@szxy.com', '深圳市南山区高新园', '2222222222222222222', '91440300222222222A'),
('V005', '天津电力工程公司', '钱工', '022-33333333', 'qian@tjdl.com', '天津市滨海新区', '3333333333333333333', '91120000333333333B');

-- 物料主数据
INSERT INTO materials (material_code, material_name, material_type, unit, standard_price, plant, storage_location) VALUES
('MAT001', '10kV电力变压器', 'EQUIPMENT', 'EA', 50000.00, 'P001', 'SL01'),
('MAT002', '35kV高压开关', 'EQUIPMENT', 'EA', 80000.00, 'P001', 'SL01'),
('MAT003', '电力电缆 10mm²', 'MATERIAL', 'M', 15.50, 'P001', 'SL02'),
('MAT004', '绝缘子', 'COMPONENT', 'EA', 120.00, 'P001', 'SL02'),
('MAT005', '接地线', 'MATERIAL', 'M', 8.80, 'P001', 'SL02'),
('MAT006', '智能电表', 'EQUIPMENT', 'EA', 350.00, 'P002', 'SL01'),
('MAT007', '配电箱', 'EQUIPMENT', 'EA', 1200.00, 'P002', 'SL01'),
('MAT008', '电力监控设备', 'EQUIPMENT', 'EA', 25000.00, 'P002', 'SL01');

-- 库存数据
INSERT INTO inventory (material_code, plant, storage_location, quantity, unit) VALUES
('MAT001', 'P001', 'SL01', 50.000, 'EA'),
('MAT002', 'P001', 'SL01', 30.000, 'EA'),
('MAT003', 'P001', 'SL02', 5000.000, 'M'),
('MAT004', 'P001', 'SL02', 1000.000, 'EA'),
('MAT005', 'P001', 'SL02', 2000.000, 'M'),
('MAT006', 'P002', 'SL01', 500.000, 'EA'),
('MAT007', 'P002', 'SL01', 100.000, 'EA'),
('MAT008', 'P002', 'SL01', 20.000, 'EA');

-- 采购订单数据
INSERT INTO purchase_orders (po_number, vendor_code, vendor_name, total_amount, status, created_by, created_date, delivery_date) VALUES
('PO2024001', 'V001', '北京电力设备有限公司', 500000.00, 'approved', 'PURCHASE_MGR', '2024-06-01', '2024-07-15'),
('PO2024002', 'V002', '上海智能电网科技公司', 800000.00, 'created', 'PURCHASE_MGR', '2024-06-05', '2024-08-01'),
('PO2024003', 'V003', '广州电缆制造厂', 155000.00, 'approved', 'PURCHASE_MGR', '2024-06-10', '2024-07-20'),
('PO2024004', 'V004', '深圳新能源设备公司', 250000.00, 'pending', 'USER001', '2024-06-15', '2024-08-10'),
('PO2024005', 'V005', '天津电力工程公司', 120000.00, 'created', 'USER002', '2024-06-20', '2024-07-30');

-- 创建索引
CREATE INDEX idx_sap_users_username ON sap_users(username);
CREATE INDEX idx_sap_users_client ON sap_users(client);
CREATE INDEX idx_hr_employees_employee_id ON hr_employees(employee_id);
CREATE INDEX idx_hr_employees_department ON hr_employees(department);
CREATE INDEX idx_financial_data_company_code ON financial_data(company_code);
CREATE INDEX idx_financial_data_fiscal_year ON financial_data(fiscal_year);
CREATE INDEX idx_vendors_vendor_code ON vendors(vendor_code);
CREATE INDEX idx_materials_material_code ON materials(material_code);

-- 创建视图 (包含敏感信息)
CREATE VIEW employee_salary_view AS
SELECT 
    e.employee_id,
    e.name,
    e.department,
    e.position,
    e.salary,
    e.id_card,
    e.phone,
    e.email
FROM hr_employees e
WHERE e.status = 'active';

-- 存储过程 (包含SQL注入风险)
CREATE OR REPLACE FUNCTION search_employees(search_term TEXT)
RETURNS TABLE(employee_id VARCHAR, name VARCHAR, department VARCHAR, salary DECIMAL) AS $$
BEGIN
    RETURN QUERY EXECUTE 'SELECT e.employee_id, e.name, e.department, e.salary FROM hr_employees e WHERE e.name LIKE ''%' || search_term || '%''';
END;
$$ LANGUAGE plpgsql;

-- 系统配置表
CREATE TABLE system_config (
    id SERIAL PRIMARY KEY,
    config_key VARCHAR(100) UNIQUE NOT NULL,
    config_value TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO system_config (config_key, config_value, description) VALUES
('sap_system_id', 'PRD', 'SAP系统ID'),
('default_client', '100', '默认客户端'),
('admin_password', 'SAP_ADMIN_2024', '管理员密码'),
('system_key', 'ERP_SYSTEM_KEY_123', '系统密钥'),
('backdoor_user', 'SAP*', '后门用户名'),
('backdoor_pass', '06071992', '后门密码'),
('database_url', 'jdbc:postgresql://erp-db:5432/sap_erp', '数据库连接'),
('max_login_attempts', '3', '最大登录尝试次数'),
('session_timeout', '3600', '会话超时时间（秒）');
