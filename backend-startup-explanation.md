# PentAGI 后端启动方式详解

## 🔍 实际启动方式

**不是** `go run main.go`，而是：

### 1️⃣ 编译阶段 (Build)
```bash
cd backend
go build -o pentagi ./cmd/pentagi
```

这个命令做了什么：
- `go build`: Go编译命令
- `-o pentagi`: 指定输出的可执行文件名为 `pentagi`
- `./cmd/pentagi`: 指定要编译的包路径（包含main.go的目录）

### 2️⃣ 运行阶段 (Run)
```bash
DATABASE_URL="postgres://postgres:postgres@localhost:5432/pentagidb?sslmode=disable" ./pentagi
```

这个命令做了什么：
- `DATABASE_URL=...`: 设置环境变量
- `./pentagi`: 运行编译好的二进制可执行文件

## 📁 文件结构

```
backend/
├── cmd/
│   └── pentagi/
│       ├── main.go      ← 主入口文件
│       └── tools.go
├── pkg/                 ← 业务逻辑包
│   ├── config/
│   ├── server/
│   ├── database/
│   └── ...
├── go.mod              ← Go模块定义
├── go.sum              ← 依赖版本锁定
└── pentagi             ← 编译后的可执行文件 (81MB)
```

## 🔄 完整流程对比

### ❌ 不是这样 (解释型运行)
```bash
go run main.go          # 每次都要编译
go run ./cmd/pentagi    # 每次都要编译
```

### ✅ 实际是这样 (编译型运行)
```bash
# 第一步：编译 (只需要一次，除非代码改变)
go build -o pentagi ./cmd/pentagi

# 第二步：运行编译好的二进制文件
DATABASE_URL="..." ./pentagi
```

## 💡 为什么这样做？

### 优势：
1. **性能更好**: 编译后的二进制文件运行更快
2. **部署简单**: 只需要一个可执行文件，不需要Go环境
3. **启动更快**: 不需要每次启动时编译
4. **生产就绪**: 这是Go应用的标准部署方式

### 对比：
| 方式 | 启动时间 | 运行性能 | 部署要求 |
|------|----------|----------|----------|
| `go run` | 慢 (需编译) | 相同 | 需要Go环境 |
| 编译后运行 | 快 | 相同 | 只需二进制文件 |

## 🛠️ 开发 vs 生产

### 开发环境 (当前)
```bash
# 修改代码后重新编译
go build -o pentagi ./cmd/pentagi
# 运行
DATABASE_URL="..." ./pentagi
```

### 生产环境
```bash
# 通常会交叉编译或在目标环境编译
go build -o pentagi ./cmd/pentagi
# 直接运行，通常配合systemd等服务管理
./pentagi
```

## 📊 当前运行状态

根据进程信息：
```
PID: 9504
命令: ./pentagi
工作目录: /mnt/persist/workspace/backend
文件大小: 81MB
运行时间: 约30分钟
```

## 🔧 如果要重启后端

```bash
# 1. 停止当前进程 (Ctrl+C 或 kill)
kill 9504

# 2. 如果修改了代码，重新编译
cd backend
go build -o pentagi ./cmd/pentagi

# 3. 重新启动
DATABASE_URL="postgres://postgres:postgres@localhost:5432/pentagidb?sslmode=disable" ./pentagi
```

## 📝 总结

- **编译型语言**: Go是编译型语言，不是解释型
- **二进制文件**: `pentagi` 是编译后的可执行文件
- **环境变量**: 通过环境变量传递配置
- **生产就绪**: 这种方式适合生产环境部署
