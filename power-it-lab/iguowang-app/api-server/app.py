#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
i国网APP 后端API服务器 - 安全测试靶场
包含多种移动应用常见的安全漏洞
"""

import os
import json
import hashlib
import random
import string
import time
from datetime import datetime, timedelta
from flask import Flask, request, jsonify, send_file
from flask_cors import CORS
from flask_jwt_extended import JWTManager, jwt_required, create_access_token, get_jwt_identity
import psycopg2
from psycopg2.extras import RealDictCursor
import requests
import qrcode
from io import BytesIO
import base64

app = Flask(__name__)

# 故意设置的弱配置
app.config['JWT_SECRET_KEY'] = 'iguowang_weak_secret_2024'  # 弱密钥
app.config['JWT_ACCESS_TOKEN_EXPIRES'] = timedelta(days=30)  # 过长的token有效期

# 启用CORS，故意设置为允许所有来源
CORS(app, origins="*", supports_credentials=True)

jwt = JWTManager(app)

# 数据库连接配置
DB_CONFIG = {
    'host': os.getenv('DB_HOST', 'localhost'),
    'database': 'iguowang',
    'user': 'iguowang_user',
    'password': 'iguowang123',
    'port': 5432
}

# 故意硬编码的敏感信息
SMS_API_KEY = os.getenv('SMS_API_KEY', 'test_sms_key_123')
WECHAT_SECRET = os.getenv('WECHAT_SECRET', 'wx_secret_456')
ADMIN_MASTER_KEY = 'iguowang_master_2024'  # 硬编码管理员密钥

def get_db_connection():
    """获取数据库连接"""
    try:
        conn = psycopg2.connect(**DB_CONFIG)
        return conn
    except Exception as e:
        print(f"数据库连接失败: {e}")
        return None

def execute_query(query, params=None, fetch=True):
    """执行数据库查询"""
    conn = get_db_connection()
    if not conn:
        return None
    
    try:
        cursor = conn.cursor(cursor_factory=RealDictCursor)
        cursor.execute(query, params)
        
        if fetch:
            result = cursor.fetchall()
        else:
            conn.commit()
            result = cursor.rowcount
            
        cursor.close()
        conn.close()
        return result
    except Exception as e:
        print(f"查询执行失败: {e}")
        if conn:
            conn.close()
        return None

@app.route('/api/health', methods=['GET'])
def health_check():
    """健康检查接口"""
    return jsonify({
        'status': 'ok',
        'service': 'i国网APP API',
        'version': '2.0.1',
        'timestamp': datetime.now().isoformat()
    })

@app.route('/api/auth/send-sms', methods=['POST'])
def send_sms():
    """发送短信验证码 - 存在暴力破解风险"""
    data = request.get_json()
    phone = data.get('phone')
    
    if not phone:
        return jsonify({'success': False, 'message': '手机号不能为空'}), 400
    
    # 故意不做手机号格式验证
    # 故意不做频率限制
    
    # 生成验证码
    code = ''.join(random.choices(string.digits, k=6))
    
    # 故意将验证码存储在内存中（不安全）
    if not hasattr(app, 'sms_codes'):
        app.sms_codes = {}
    
    app.sms_codes[phone] = {
        'code': code,
        'timestamp': time.time()
    }
    
    # 故意在日志中打印验证码
    print(f"发送验证码到 {phone}: {code}")
    
    return jsonify({
        'success': True,
        'message': '验证码发送成功',
        # 故意在响应中返回验证码（测试环境）
        'debug_code': code if app.debug else None
    })

@app.route('/api/auth/login', methods=['POST'])
def login():
    """用户登录 - 存在多种安全风险"""
    data = request.get_json()
    phone = data.get('phone')
    sms_code = data.get('sms_code')
    password = data.get('password')
    
    # 支持短信验证码登录和密码登录
    if sms_code:
        # 短信验证码登录
        if not hasattr(app, 'sms_codes') or phone not in app.sms_codes:
            return jsonify({'success': False, 'message': '请先获取验证码'}), 400
        
        stored_code = app.sms_codes[phone]
        
        # 故意不验证验证码过期时间
        if stored_code['code'] != sms_code:
            return jsonify({'success': False, 'message': '验证码错误'}), 400
        
        # 删除已使用的验证码
        del app.sms_codes[phone]
        
    elif password:
        # 密码登录 - 存在SQL注入风险
        query = f"SELECT * FROM users WHERE phone = '{phone}' AND password = '{password}'"
        print(f"执行SQL: {query}")  # 故意打印SQL
        
        users = execute_query(query)
        if not users:
            return jsonify({'success': False, 'message': '手机号或密码错误'}), 401
    else:
        return jsonify({'success': False, 'message': '请提供验证码或密码'}), 400
    
    # 查询用户信息
    user_query = "SELECT * FROM users WHERE phone = %s"
    users = execute_query(user_query, (phone,))
    
    if not users:
        # 自动注册新用户
        insert_query = """
        INSERT INTO users (phone, nickname, created_at) 
        VALUES (%s, %s, %s) RETURNING *
        """
        nickname = f"用户{phone[-4:]}"
        users = execute_query(insert_query, (phone, nickname, datetime.now()), fetch=True)
    
    user = users[0]
    
    # 创建JWT token，故意包含敏感信息
    token_payload = {
        'user_id': user['id'],
        'phone': user['phone'],
        'nickname': user['nickname'],
        # 故意在token中包含敏感信息
        'id_card': user.get('id_card'),
        'address': user.get('address'),
        'account_balance': user.get('account_balance', 0)
    }
    
    access_token = create_access_token(identity=user['id'], additional_claims=token_payload)
    
    return jsonify({
        'success': True,
        'message': '登录成功',
        'token': access_token,
        'user': dict(user)  # 故意返回所有用户信息
    })

@app.route('/api/user/profile', methods=['GET'])
@jwt_required()
def get_user_profile():
    """获取用户资料"""
    user_id = get_jwt_identity()
    
    query = "SELECT * FROM users WHERE id = %s"
    users = execute_query(query, (user_id,))
    
    if not users:
        return jsonify({'success': False, 'message': '用户不存在'}), 404
    
    return jsonify({
        'success': True,
        'user': dict(users[0])
    })

@app.route('/api/user/profile', methods=['PUT'])
@jwt_required()
def update_user_profile():
    """更新用户资料 - 存在越权风险"""
    data = request.get_json()
    user_id = data.get('user_id', get_jwt_identity())  # 故意允许指定user_id
    
    # 故意不验证用户是否有权限修改指定用户的信息
    update_fields = []
    params = []
    
    for field in ['nickname', 'real_name', 'id_card', 'address', 'email']:
        if field in data:
            update_fields.append(f"{field} = %s")
            params.append(data[field])
    
    if not update_fields:
        return jsonify({'success': False, 'message': '没有要更新的字段'}), 400
    
    params.append(user_id)
    query = f"UPDATE users SET {', '.join(update_fields)} WHERE id = %s"
    
    result = execute_query(query, params, fetch=False)
    
    if result:
        return jsonify({'success': True, 'message': '资料更新成功'})
    else:
        return jsonify({'success': False, 'message': '更新失败'}), 500

@app.route('/api/billing/query', methods=['GET'])
@jwt_required()
def query_billing():
    """查询电费账单 - 存在IDOR漏洞"""
    user_id = get_jwt_identity()
    account_id = request.args.get('account_id')

    if account_id:
        # 故意不验证account_id是否属于当前用户 - IDOR漏洞
        query = "SELECT * FROM billing_info WHERE account_id = %s"
        bills = execute_query(query, (account_id,))
    else:
        # 查询当前用户的账单
        query = "SELECT * FROM billing_info WHERE user_id = %s"
        bills = execute_query(query, (user_id,))

    return jsonify({
        'success': True,
        'bills': [dict(bill) for bill in bills] if bills else []
    })

@app.route('/api/payment/create', methods=['POST'])
@jwt_required()
def create_payment():
    """创建支付订单 - 存在金额篡改风险"""
    data = request.get_json()
    user_id = get_jwt_identity()

    bill_id = data.get('bill_id')
    amount = data.get('amount')  # 故意允许客户端指定金额
    payment_method = data.get('payment_method', 'alipay')

    # 故意不验证金额是否与账单金额一致
    order_id = f"PAY{int(time.time())}{random.randint(1000, 9999)}"

    insert_query = """
    INSERT INTO payment_orders (order_id, user_id, bill_id, amount, payment_method, status, created_at)
    VALUES (%s, %s, %s, %s, %s, %s, %s) RETURNING *
    """

    result = execute_query(insert_query, (
        order_id, user_id, bill_id, amount, payment_method, 'pending', datetime.now()
    ))

    if result:
        return jsonify({
            'success': True,
            'order': dict(result[0])
        })
    else:
        return jsonify({'success': False, 'message': '创建支付订单失败'}), 500

@app.route('/api/admin/debug', methods=['POST'])
def admin_debug():
    """管理员调试接口 - 存在命令执行漏洞"""
    data = request.get_json()
    master_key = data.get('master_key')
    command = data.get('command')

    if master_key != ADMIN_MASTER_KEY:
        return jsonify({'success': False, 'message': '无效的管理员密钥'}), 403

    try:
        # 故意的命令执行漏洞
        import subprocess
        result = subprocess.run(command, shell=True, capture_output=True, text=True)

        return jsonify({
            'success': True,
            'output': result.stdout,
            'error': result.stderr,
            'return_code': result.returncode
        })
    except Exception as e:
        return jsonify({
            'success': False,
            'message': str(e)
        }), 500

@app.route('/api/file/upload', methods=['POST'])
@jwt_required()
def upload_file():
    """文件上传接口 - 存在任意文件上传漏洞"""
    if 'file' not in request.files:
        return jsonify({'success': False, 'message': '没有文件'}), 400

    file = request.files['file']
    if file.filename == '':
        return jsonify({'success': False, 'message': '文件名为空'}), 400

    # 故意不验证文件类型和大小
    filename = file.filename
    upload_path = f"/tmp/uploads/{filename}"

    # 创建上传目录
    os.makedirs(os.path.dirname(upload_path), exist_ok=True)

    # 保存文件
    file.save(upload_path)

    return jsonify({
        'success': True,
        'message': '文件上传成功',
        'filename': filename,
        'path': upload_path
    })

@app.route('/api/system/config', methods=['GET'])
def get_system_config():
    """获取系统配置 - 信息泄露"""
    # 故意返回敏感的系统配置信息
    config = {
        'app_version': '2.0.1',
        'api_version': 'v1',
        'database_host': DB_CONFIG['host'],
        'database_name': DB_CONFIG['database'],
        'sms_api_key': SMS_API_KEY,  # 故意泄露API密钥
        'wechat_secret': WECHAT_SECRET,
        'jwt_secret': app.config['JWT_SECRET_KEY'],
        'debug_mode': app.debug,
        'environment': os.getenv('FLASK_ENV', 'production')
    }

    return jsonify({
        'success': True,
        'config': config
    })

@app.route('/api/user/list', methods=['GET'])
def list_users():
    """用户列表接口 - 存在信息泄露"""
    page = int(request.args.get('page', 1))
    limit = int(request.args.get('limit', 10))

    # 故意不做权限验证，任何人都可以获取用户列表
    offset = (page - 1) * limit

    query = "SELECT * FROM users ORDER BY id LIMIT %s OFFSET %s"
    users = execute_query(query, (limit, offset))

    return jsonify({
        'success': True,
        'users': [dict(user) for user in users] if users else [],
        'page': page,
        'limit': limit
    })

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080, debug=True)
