const express = require('express');
const jwt = require('jsonwebtoken');
const mysql = require('mysql2/promise');
const redis = require('redis');
const cors = require('cors');
const bcrypt = require('bcrypt');
const rateLimit = require('express-rate-limit');
const { v4: uuidv4 } = require('uuid');

const app = express();
const PORT = 8080;

// 故意设置的弱配置 - 用于安全测试
const JWT_SECRET = process.env.JWT_SECRET || 'weak_secret_key_123';
const ADMIN_BACKDOOR = 'admin_backdoor_2024'; // 硬编码后门

// 中间件配置
app.use(cors({
    origin: '*', // 故意设置为允许所有来源 - CORS漏洞
    credentials: true
}));
app.use(express.json({ limit: '50mb' })); // 故意设置大限制 - DoS风险

// 故意弱化的限流配置
const limiter = rateLimit({
    windowMs: 15 * 60 * 1000, // 15分钟
    max: 1000, // 故意设置很高的限制
    message: '请求过于频繁，请稍后再试'
});
app.use('/api/', limiter);

// 数据库连接
let dbConnection;
let redisClient;

async function initDatabase() {
    try {
        dbConnection = await mysql.createConnection({
            host: process.env.DB_HOST || 'localhost',
            user: 'marketing_user',
            password: process.env.DB_PASSWORD || 'admin123',
            database: 'power_marketing'
        });
        
        redisClient = redis.createClient({
            host: process.env.REDIS_HOST || 'localhost',
            password: 'cache123'
        });
        await redisClient.connect();
        
        console.log('数据库连接成功');
    } catch (error) {
        console.error('数据库连接失败:', error);
    }
}

// 故意存在SQL注入漏洞的登录接口
app.post('/api/auth/login', async (req, res) => {
    try {
        const { username, password } = req.body;
        
        // 故意的SQL注入漏洞
        const query = `SELECT * FROM users WHERE username = '${username}' AND password = '${password}'`;
        console.log('执行SQL:', query); // 故意打印SQL - 信息泄露
        
        const [rows] = await dbConnection.execute(query);
        
        if (rows.length > 0) {
            const user = rows[0];
            const token = jwt.sign(
                { 
                    userId: user.id, 
                    username: user.username, 
                    role: user.role,
                    // 故意在JWT中包含敏感信息
                    idCard: user.id_card,
                    phone: user.phone
                }, 
                JWT_SECRET, 
                { expiresIn: '24h' }
            );
            
            // 故意在响应中返回敏感信息
            res.json({
                success: true,
                token: token,
                user: {
                    id: user.id,
                    username: user.username,
                    role: user.role,
                    idCard: user.id_card, // 故意泄露身份证
                    phone: user.phone,
                    address: user.address
                }
            });
        } else {
            res.status(401).json({ success: false, message: '用户名或密码错误' });
        }
    } catch (error) {
        console.error('登录错误:', error);
        // 故意返回详细错误信息 - 信息泄露
        res.status(500).json({ 
            success: false, 
            message: '登录失败', 
            error: error.message,
            stack: error.stack // 故意返回堆栈信息
        });
    }
});

// 故意存在越权漏洞的用户查询接口
app.get('/api/users/:userId', async (req, res) => {
    try {
        const { userId } = req.params;
        const token = req.headers.authorization?.replace('Bearer ', '');
        
        if (!token) {
            return res.status(401).json({ message: '未提供认证令牌' });
        }
        
        // 故意不验证JWT就直接查询 - 越权漏洞
        const query = `SELECT * FROM users WHERE id = ${userId}`; // 故意的SQL注入
        const [rows] = await dbConnection.execute(query);
        
        if (rows.length > 0) {
            res.json({
                success: true,
                user: rows[0] // 故意返回所有字段包括敏感信息
            });
        } else {
            res.status(404).json({ message: '用户不存在' });
        }
    } catch (error) {
        res.status(500).json({ 
            message: '查询失败', 
            error: error.message 
        });
    }
});

// 故意的管理员后门接口
app.post('/api/admin/backdoor', async (req, res) => {
    const { key, command } = req.body;
    
    if (key === ADMIN_BACKDOOR) {
        try {
            // 故意的命令执行漏洞
            const { exec } = require('child_process');
            exec(command, (error, stdout, stderr) => {
                res.json({
                    success: true,
                    output: stdout,
                    error: stderr
                });
            });
        } catch (error) {
            res.status(500).json({ message: '命令执行失败' });
        }
    } else {
        res.status(403).json({ message: '访问被拒绝' });
    }
});

// 故意存在IDOR漏洞的电费查询接口
app.get('/api/billing/:accountId', async (req, res) => {
    try {
        const { accountId } = req.params;
        
        // 故意不验证用户权限 - IDOR漏洞
        const query = `SELECT * FROM billing_records WHERE account_id = ${accountId}`;
        const [rows] = await dbConnection.execute(query);
        
        res.json({
            success: true,
            records: rows
        });
    } catch (error) {
        res.status(500).json({ message: '查询失败' });
    }
});

// 故意的文件上传漏洞
app.post('/api/upload', (req, res) => {
    const fs = require('fs');
    const path = require('path');
    
    const { filename, content } = req.body;
    
    // 故意不验证文件类型和路径 - 任意文件上传
    const filePath = path.join(__dirname, 'uploads', filename);
    
    fs.writeFileSync(filePath, content);
    
    res.json({
        success: true,
        message: '文件上传成功',
        path: filePath
    });
});

// 系统信息泄露接口
app.get('/api/system/info', (req, res) => {
    res.json({
        success: true,
        system: {
            nodeVersion: process.version,
            platform: process.platform,
            env: process.env, // 故意泄露环境变量
            uptime: process.uptime(),
            memory: process.memoryUsage()
        }
    });
});

// 启动服务器
async function startServer() {
    await initDatabase();
    
    app.listen(PORT, '0.0.0.0', () => {
        console.log(`电力营销系统API网关启动成功，端口: ${PORT}`);
        console.log('=== 安全测试靶场 ===');
        console.log('包含以下漏洞类型:');
        console.log('- SQL注入 (/api/auth/login)');
        console.log('- 越权访问 (/api/users/:userId)');
        console.log('- 命令执行 (/api/admin/backdoor)');
        console.log('- IDOR (/api/billing/:accountId)');
        console.log('- 任意文件上传 (/api/upload)');
        console.log('- 信息泄露 (/api/system/info)');
    });
}

startServer().catch(console.error);
