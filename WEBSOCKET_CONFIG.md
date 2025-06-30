# PentAGI WebSocket 智能配置技术文档

## 🌐 概述

PentAGI系统实现了智能WebSocket配置，能够根据部署环境自动选择最佳的实时通信方式，确保在不同环境下都能提供良好的用户体验。

## 🔧 技术实现

### 环境检测机制

```typescript
// 检查是否在cloudflare隧道环境下
const isCloudflareEnvironment = window.location.hostname.includes('trycloudflare.com');
```

### 连接策略

#### 本地开发环境
- **检测条件**: `hostname` 不包含 `trycloudflare.com`
- **WebSocket URL**: `ws://localhost:8001/api/v1/graphql` 或 `wss://localhost:8001/api/v1/graphql`
- **连接方式**: 真正的WebSocket连接
- **实时性**: < 100ms 延迟
- **重试机制**: 指数退避，最大重试5次

#### Cloudflare隧道环境
- **检测条件**: `hostname` 包含 `trycloudflare.com`
- **连接方式**: HTTP轮询
- **轮询间隔**: 3000ms (3秒)
- **实时性**: ≤ 3秒延迟
- **原因**: Cloudflare免费隧道不支持WebSocket

### 代码实现

```typescript
const wsLink = isCloudflareEnvironment 
    ? null // 在cloudflare环境下禁用WebSocket
    : new GraphQLWsLink(
        createClient({
            url: `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}${baseUrl}/graphql`,
            retryAttempts: 5,
            connectionParams: () => {
                return {}; // Cookies are handled automatically
            },
            // ... 其他配置
        }),
    );

const link = wsLink 
    ? split(
        ({ query }) => {
            const definition = getMainDefinition(query);
            return definition.kind === 'OperationDefinition' && definition.operation === 'subscription';
        },
        wsLink,
        httpLink,
    )
    : httpLink; // 在cloudflare环境下只使用HTTP连接
```

### Apollo客户端配置

```typescript
const defaultOptions: DefaultOptions = {
    watchQuery: {
        fetchPolicy: 'cache-and-network',
        nextFetchPolicy: 'cache-first',
        notifyOnNetworkStatusChange: true,
        pollInterval: isCloudflareEnvironment ? 3000 : undefined, // 在cloudflare环境下使用轮询
    },
};
```

## 📊 性能对比

| 环境 | 连接方式 | 延迟 | 资源消耗 | 可靠性 |
|------|----------|------|----------|--------|
| 本地开发 | WebSocket | < 100ms | 低 | 高 |
| Cloudflare | HTTP轮询 | ≤ 3s | 中等 | 高 |

## 🔍 调试和监控

### 日志输出

系统会在控制台输出当前使用的连接模式：

```
🔌 运行在本地环境，使用WebSocket实时模式
🌐 运行在Cloudflare隧道环境，使用HTTP轮询模式
```

### WebSocket事件监控

在本地环境下，系统会记录以下WebSocket事件：
- `connected`: WebSocket连接成功
- `error`: WebSocket连接错误
- `closed`: WebSocket连接关闭
- `connecting`: WebSocket正在连接
- `ping/pong`: 心跳检测

### 故障排除

#### 常见问题

1. **消息不实时更新**
   - 检查当前环境模式
   - 在Cloudflare环境下等待最多3秒
   - 检查网络连接

2. **WebSocket连接失败**
   - 确认是否在支持WebSocket的环境
   - 检查防火墙设置
   - 验证后端GraphQL服务状态

3. **轮询频率过高**
   - 调整 `pollInterval` 参数
   - 考虑网络带宽和服务器负载

## ⚙️ 配置选项

### 轮询间隔调整

```typescript
// 在 frontend/src/lib/apollo.ts 中修改
pollInterval: isCloudflareEnvironment ? 5000 : undefined, // 改为5秒轮询
```

### WebSocket重试配置

```typescript
// 在 frontend/src/lib/apollo.ts 中修改
retryAttempts: 3, // 减少重试次数
retryWait: (retries) => new Promise((resolve) => {
    const timeout = Math.min(1000 * 2 ** retries, 5000); // 减少最大等待时间
    setTimeout(() => resolve(), timeout);
}),
```

### 环境检测自定义

```typescript
// 自定义环境检测逻辑
const isCloudflareEnvironment = 
    window.location.hostname.includes('trycloudflare.com') ||
    window.location.hostname.includes('your-custom-domain.com');
```

## 🚀 部署建议

### 开发环境
- 使用本地WebSocket连接
- 启用详细日志
- 配置较短的重试间隔

### 生产环境
- 如果支持WebSocket，优先使用WebSocket
- 如果不支持，使用HTTP轮询
- 配置合适的轮询间隔（建议3-5秒）
- 启用错误监控

### Cloudflare隧道
- 自动使用HTTP轮询
- 建议轮询间隔不少于3秒
- 监控服务器负载

## 🔮 未来改进

1. **智能轮询**: 根据用户活跃度动态调整轮询频率
2. **Server-Sent Events**: 作为WebSocket的替代方案
3. **WebSocket代理**: 通过专门的WebSocket代理服务
4. **缓存优化**: 减少不必要的轮询请求

## 📝 维护注意事项

1. **定期检查**: 验证不同环境下的连接状态
2. **性能监控**: 监控轮询对服务器的影响
3. **用户反馈**: 收集用户对实时性的反馈
4. **技术更新**: 关注Cloudflare对WebSocket支持的更新

## 🔗 相关文件

- `frontend/src/lib/apollo.ts`: 主要配置文件
- `frontend/vite.config.ts`: Vite代理配置
- `RESTART_GUIDE.md`: 重启和故障排除指南
- `backend/`: GraphQL服务器实现
