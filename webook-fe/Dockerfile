# 基础 image
FROM node:alpine AS app

# 设置项目目录
WORKDIR /app

# 安装项目依赖
COPY package.json ./
RUN npm install --registry=https://registry.npm.taobao.org

COPY . .

EXPOSE 3000
# 运行
CMD ["npm", "run", "dev"]