# Swagger API 文档说明

## 概述

已为电池进销存管理系统添加了完整的 Swagger API 文档支持。

## 功能特性

- ✅ 完整的 API 接口文档
- ✅ 中文接口描述和注释
- ✅ 交互式 API 测试界面
- ✅ JWT 认证支持
- ✅ 请求/响应示例
- ✅ 数据模型定义

## 如何访问 Swagger 文档

1. **启动后端服务**：
   ```bash
   cd backend
   go run main.go
   ```

2. **访问 Swagger UI**：
   在浏览器中打开：http://localhost:8036/swagger/index.html

## 已添加的 API 文档

### 认证模块
- POST `/auth/login` - 用户登录

### 用户管理
- GET `/users` - 获取所有用户
- POST `/users` - 创建用户
- GET `/users/{id}` - 根据ID获取用户
- PUT `/users/{id}` - 更新用户
- DELETE `/users/{id}` - 删除用户

### 电池分类管理
- GET `/categories` - 获取所有电池分类
- POST `/categories` - 创建电池分类
- GET `/categories/{id}` - 根据ID获取电池分类
- PUT `/categories/{id}` - 更新电池分类
- DELETE `/categories/{id}` - 删除电池分类

### 入库管理
- GET `/inbound/orders` - 获取所有入库订单
- POST `/inbound/orders` - 创建入库订单
- GET `/inbound/orders/{id}` - 根据ID获取入库订单
- PUT `/inbound/orders/{id}` - 更新入库订单
- DELETE `/inbound/orders/{id}` - 删除入库订单

### 库存管理
- GET `/inventory` - 获取所有库存
- GET `/inventory/{categoryId}` - 根据分类ID获取库存

## 如何使用

### 1. 认证
大部分 API 需要 JWT token 认证：
1. 首先调用 `/auth/login` 接口获取 token
2. 在后续请求的 Header 中添加：`Authorization: Bearer <your-token>`

### 2. 测试 API
在 Swagger UI 界面中：
1. 点击 "Authorize" 按钮
2. 输入 `Bearer <your-token>`
3. 选择要测试的 API 接口
4. 填写参数并执行

## 开发说明

### 添加新的 API 文档

在控制器方法上添加 Swagger 注解：

```go
// MethodName godoc
// @Summary      接口简要描述
// @Description  接口详细描述
// @Tags         标签分组
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        paramName paramType dataType required "参数描述"
// @Success      200 {object} ResponseType "成功响应"
// @Failure      400 {object} models.Response "失败响应"
// @Router       /path [method]
func (ctrl *Controller) MethodName(c *gin.Context) {
    // 方法实现
}
```

### 重新生成文档

在 backend 目录下运行：
```bash
swag init
```

## 项目结构

```
backend/
├── docs/                 # 生成的 Swagger 文档
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/api/v1/      # API 控制器（包含 Swagger 注解）
└── main.go              # 主程序（包含 Swagger 配置）
```

## 注意事项

- 确保所有的数据模型都有正确的 JSON 标签
- API 路径要与路由配置保持一致
- 认证相关的接口需要添加 `@Security BearerAuth` 注解
- 修改 API 后需要重新运行 `swag init` 生成文档 