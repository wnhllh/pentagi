const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const app = express();

// 添加日志中间件
app.use((req, res, next) => {
  console.log(`📥 收到请求: ${req.method} ${req.url}`);
  next();
});

// 创建代理中间件实例 - 只创建一次
const apiProxy = createProxyMiddleware({
  target: 'http://localhost:8080',
  changeOrigin: true,
  ws: true,
  onProxyReq: (proxyReq, req, res) => {
    console.log(`🔄 API代理: ${req.method} ${req.url} -> http://localhost:8080${req.url}`);
  },
  onError: (err, req, res) => {
    console.error(`❌ API代理错误:`, err.message);
  }
});

const graphqlProxy = createProxyMiddleware({
  target: 'http://localhost:8080',
  changeOrigin: true,
  ws: true,
  onProxyReq: (proxyReq, req, res) => {
    console.log(`🔄 GraphQL代理: ${req.method} ${req.url} -> http://localhost:8080${req.url}`);
  },
  onError: (err, req, res) => {
    console.error(`❌ GraphQL代理错误:`, err.message);
  }
});

const frontendProxy = createProxyMiddleware({
  target: 'http://localhost:8001',
  changeOrigin: true,
  ws: true,
  onProxyReq: (proxyReq, req, res) => {
    console.log(`🎨 前端代理: ${req.method} ${req.url} -> http://localhost:8001${req.url}`);
  },
  onError: (err, req, res) => {
    console.error(`❌ 前端代理错误:`, err.message);
  }
});

// 使用代理中间件 - 顺序很重要！
// API路由必须在前端路由之前
app.use('/api', apiProxy);
app.use('/graphql', graphqlProxy);

// 前端路由作为fallback，但要排除API路径
app.use((req, res, next) => {
  // 如果是API或GraphQL请求，不应该到达这里
  if (req.url.startsWith('/api/') || req.url.startsWith('/graphql')) {
    return res.status(404).json({ error: 'API endpoint not found' });
  }
  // 其他所有请求都代理到前端
  return frontendProxy(req, res, next);
});

const PORT = 3000;
app.listen(PORT, () => {
  console.log(`🚀 代理服务器运行在 http://localhost:${PORT}`);
  console.log(`📡 后端API代理: http://localhost:${PORT}/api -> http://localhost:8080/api`);
  console.log(`🎨 前端代理: http://localhost:${PORT} -> http://localhost:8001`);
});
