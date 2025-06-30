FROM node:18-alpine

WORKDIR /app

# 安装依赖
COPY package*.json ./
RUN npm install

# 复制源码
COPY . .

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["node", "server.js"]
