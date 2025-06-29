废旧电池回收进销存（ERP）系统
技术设计与实现规约
文档版本

V4.0 (Monorepo Final)

完成日期

2025年6月28日

目标读者
Cluade Code


Export to Sheets
目录
引言

1.1. 项目目标

1.2. 核心功能

1.3. 技术栈

系统架构

2.1. 总体架构

2.2. 用户角色与权限

数据库设计

3.1. 设计规约

3.2. 实体关系图 (ERD) 概念

3.3. 表结构详情

核心业务逻辑

4.1. 原子性操作（事务）

4.2. 库存管理核心逻辑

4.3. 单号生成规则

4.4. 关键数据计算

项目结构 (Monorepo)

5.1. 统一代码库结构

5.2. 后端 (Go) 内部结构

5.3. 前端 (React) 内部结构

API 规约

6.1. 通用规约

6.2. 通用响应结构

6.3. API 接口详解

6.3.1. 认证接口 (Auth)

6.3.2. 用户接口 (Users)

6.3.3. 品类接口 (Categories)

6.3.4. 采购入库接口 (Inbound)

6.3.5. 销售出库接口 (Outbound)

6.3.6. 库存接口 (Inventory)

6.3.7. 报表接口 (Reports)

部署与运维

7.1. Docker 容器化配置

7.1.1. 后端 Dockerfile

7.1.2. 前端 Dockerfile

7.1.3. Docker Compose (Monorepo)

7.2. 数据备份与灾备

总结与展望


1. 引言
1.1. 项目目标
本文档旨在为“废旧电池回收进销存（ERP）系统”提供全面、详细的技术设计和实现规约。项目的核心目标是开发一套高效、稳定、易用的管理系统，实现对废旧电池回收业务中采购、销售、库存等环节的数字化管理，从而提升运营效率、数据准确性和决策支持能力。

1.2. 核心功能
采购管理: 实现电池回收的入库流程，包括称重、计价、生成入库单据和打印。

销售管理: 实现电池销售的出库流程，记录销售信息、生成出库单据。

库存管理: 实时跟踪各类电池的库存重量，提供库存查询和预警功能。

报表分析: 按日、月、年等维度生成进销存报表，提供业务概览和数据洞察。

用户与权限管理: 支持多用户协作，通过角色划分保证操作安全。

1.3. 技术栈
后端: Go 1.23 (使用 Gin 框架)

前端: React (推荐使用 Vite)

数据库: MySQL 8.0+

容器化: Docker & Docker Compose

API 架构: RESTful API

认证机制: JWT (JSON Web Tokens)


2. 系统架构
2.1. 总体架构
系统采用主流的前后端分离架构，各部分职责明确，便于独立开发、测试和部署。

前端 (React): 负责用户界面（UI）和用户交互（UX），通过调用后端 API 完成数据操作。

前端设计需要简约美观。

后端 (Go): 作为无状态应用，负责处理业务逻辑、数据持久化，并通过 RESTful API 对外提供服务。

数据库 (MySQL): 作为系统的统一数据存储中心。

Docker Compose: 负责编排和管理前端、后端、数据库三个服务容器，实现一键化部署和环境隔离。

2.2. 用户角色与权限
系统内置两级权限体系，以满足不同岗位员工的操作需求。

角色

权限描述

超级管理员 (Super Admin)

拥有系统所有功能权限，包括用户管理、品类管理、查看所有单据与报表、系统配置等。

普通用户 (Normal User)

负责日常业务操作，包括创建入库/出库单、查看个人经手的单据、查询库存及公共报表。


Export to Sheets

3. 数据库设计
3.1. 设计规约
所有表必须使用 InnoDB 存储引擎，以支持事务和外键。

所有表必须包含 created_at 和 updated_at 两个字段，类型为 DATETIME，用于记录数据的创建和最后修改时间，遵循阿里云数据库设计规范。

表名和字段名采用小写字母和下划线 (snake_case) 的命名方式。

3.2. 实体关系图 (ERD) 概念
users 表记录系统用户。

battery_categories 表定义了所有可交易的电池品类。

inbound_orders (主表) 和 inbound_order_items (明细表) 共同记录一笔完整的采购入库业务。

outbound_orders 和 outbound_order_items 记录销售出库业务。

inventory 表是核心库存表，实时反映每种品类的当前库存重量。

3.3. 表结构详情
(为简洁起见，此处省略具体的表结构定义，其设计已在前期版本中明确，核心要点是所有表都包含 id, created_at, updated_at 等必要字段。)


4. 核心业务逻辑
4.1. 原子性操作（事务）
所有涉及多表数据变更的操作（尤其是入库和出库）必须置于数据库事务中，确保数据的一致性。若任一步骤失败，整个操作必须回滚。

4.2. 库存管理核心逻辑
防超卖: 执行出库操作时，必须锁定并校验库存记录，确保 current_weight_kg 大于或等于本次出库重量。若库存不足，则事务失败。

自动建档: 当一个新品类的电池首次入库时，系统应自动在 inventory 表中创建对应的库存记录。

4.3. 单号生成规则
单号应保证全局唯一且具备可读性，建议采用 “前缀-日期单位到纳秒秒-随机数（0-1000）” 格式。

入库单号: IN-YYYYMMDD-XXXX (e.g., IN-20250628-0001)

出库单号: OUT-YYYYMMDD-XXXX (e.g., OUT-20250628-0001)

4.4. 关键数据计算
净重计算: Net
Weight=Gross
Weight−Tare
Weight

小计金额: Subtotal
Price=Net
Weight
timesUnit
Price

总金额: Total
Price=
sum(Subtotal
Price)


5. 项目结构 (Monorepo)
5.1. 统一代码库结构
前后端项目位于同一个 Git 仓库下的不同目录中，便于统一管理。

Plaintext

/battery-erp-system/
|
|-- backend/
|   |-- cmd/
|   |-- internal/
|   |-- go.mod
|   |-- Dockerfile  <-- 后端 Dockerfile
|   |-- ...
|
|-- frontend/
|   |-- src/
|   |-- public/
|   |-- package.json
|   |-- Dockerfile  <-- 前端 Dockerfile
|   |-- ...
|
|-- docker-compose.yml  <-- 根目录的 Docker Compose 文件
|
|-- .gitignore
|
|-- README.md
5.2. 后端 (Go) 内部结构
采用标准的 MVC 分层结构，实现高内聚、低耦合。

Plaintext

/backend/
|-- /cmd/server/main.go         # 程序入口
|-- /internal
|   |-- /api/v1                 # Controller层: 接口控制器
|   |-- /models                 # Model层: 数据库实体
|   |-- /services               # Service层: 核心业务逻辑
|   |-- /repository             # Repository层: 数据访问 (GORM)
|   |-- ...
5.3. 前端 (React) 内部结构
按功能模块组织文件，便于维护。

Plaintext

/frontend/
|-- /src
|   |-- /api                    # API请求服务
|   |-- /components             # 可复用UI组件
|   |-- /pages                  # 页面级组件 (按业务模块划分)
|   |-- ...

6. API 规约
6.1. 通用规约
API 前缀: 所有接口均以 /jxc/v1 为前缀。

认证: 除登录接口外，所有接口均需在 HTTP Header 中携带 Authorization: Bearer <token> 进行 JWT 认证。

数据格式: 请求体和响应体均使用 application/json 格式。

6.2. 通用响应结构
所有 API 响应都包裹在统一结构中，便于前端进行统一处理。

JSON

{
  "code": 0,          // 状态码: 0 表示成功, 非零表示错误
  "msg": "success",   // 提示信息
  "data": { ... }     // 响应数据
}
6.3. API 接口详解
(此处仅列出关键接口，详细请求/响应体已在前期版本中明确。)

6.3.1. 认证接口 (Auth)
POST /jxc/v1/auth/login: 用户登录，成功后返回 token 和用户信息。

6.3.2. 用户接口 (Users)
GET /jxc/v1/users: [Admin] 获取用户列表。

POST /jxc/v1/users: [Admin] 创建新用户。

6.3.3. 品类接口 (Categories)
GET /jxc/v1/categories: 获取所有电池品类。

POST /jxc/v1/categories: [Admin] 创建新品类。

6.3.4. 采购入库接口 (Inbound)
POST /jxc/v1/inbound/orders: 创建入库单，后端将执行原子性的数据写入和库存更新。

GET /jxc/v1/inbound/orders: 获取入库单列表（支持分页和筛选）。

GET /jxc/v1/inbound/orders/{id}: 获取单个入库单详情，用于查看和打印。

6.3.5. 销售出库接口 (Outbound)
POST /jxc/v1/outbound/orders: 创建出库单，后端将执行防超卖检查和库存扣减。

GET /jxc/v1/outbound/orders: 获取出库单列表。

GET /jxc/v1/outbound/orders/{id}: 获取单个出库单详情。

6.3.6. 库存接口 (Inventory)
GET /jxc/v1/inventory: 获取所有品类的当前库存列表。

6.3.7. 报表接口 (Reports)
GET /jxc/v1/reports/summary: 获取指定时间范围（或日/月/年）的业务总览数据。


7. 部署与运维
7.1. Docker 容器化配置
7.1.1. 后端 Dockerfile
文件路径: /backend/Dockerfile

Dockerfile

# ---- Build Stage ----
FROM golang:1.23-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/main ./cmd/server

# ---- Release Stage ----
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main /app/main
EXPOSE 8036
CMD ["/app/main"]
7.1.2. 前端 Dockerfile
文件路径: /frontend/Dockerfile

Dockerfile

# ---- Build Stage ----
FROM node:20-alpine AS builder
WORKDIR /app
RUN npm config set registry https://registry.npmmirror.com
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# ---- Production Stage ----
FROM nginx:stable-alpine
COPY --from=builder /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
7.1.3. Docker Compose (Monorepo)
文件路径: /docker-compose.yml (项目根目录)

YAML

version: '3.8'

services:
  backend:
    build:
      context: ./backend  # 指定后端目录
      dockerfile: Dockerfile
    container_name: battery_backend
    restart: unless-stopped
    ports:
      - "8036:8036"
    environment:
      - DB_HOST=mysql_db
      - DB_USER=root
      - DB_PASSWORD=your_strong_password
      - DB_NAME=battery_erp
      - JWT_SECRET=your_jwt_secret_key
    depends_on:
      - mysql_db
    networks:
      - app-network

  frontend:
    build:
      context: ./frontend # 指定前端目录
      dockerfile: Dockerfile
    container_name: battery_frontend
    restart: unless-stopped
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - app-network

  mysql_db:
    image: mysql:8.0
    container_name: battery_mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: your_strong_password
      MYSQL_DATABASE: battery_erp
    volumes:
      - mysql-data:/var/lib/mysql
    networks:
      - app-network

networks:
  app-network:
    driver: bridge

volumes:
  mysql-data:

