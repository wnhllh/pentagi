FROM openjdk:17-jdk-slim

WORKDIR /app

# 安装必要的工具
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    && rm -rf /var/lib/apt/lists/*

# 复制应用文件
COPY . .

# 编译Java应用
RUN javac -cp ".:lib/*" *.java

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["java", "-cp", ".:lib/*", "ERPServer"]
