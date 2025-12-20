# 营养追踪系统测试指南

## 1. 系统准备

### 1.1 数据库准备
确保PostgreSQL数据库已安装并运行，执行以下SQL创建数据库和表结构：

```sql
-- 创建数据库
CREATE DATABASE nutrition_tracker;

-- 切换到新数据库
\c nutrition_tracker;

-- 创建用户表
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    height DECIMAL(5,2),
    weight DECIMAL(5,2),
    age INTEGER,
    gender VARCHAR(10),
    activity_level VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建营养目标表
CREATE TABLE nutrition_goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    calories INTEGER NOT NULL,
    protein DECIMAL(5,2) NOT NULL,
    carbohydrates DECIMAL(5,2) NOT NULL,
    fat DECIMAL(5,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建食物表
CREATE TABLE foods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    calories DECIMAL(5,2) NOT NULL,
    protein DECIMAL(5,2) NOT NULL,
    carbohydrates DECIMAL(5,2) NOT NULL,
    fat DECIMAL(5,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建餐次记录表
CREATE TABLE meal_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    meal_type VARCHAR(20) NOT NULL,
    record_date DATE NOT NULL,
    total_calories DECIMAL(6,2) DEFAULT 0,
    total_protein DECIMAL(6,2) DEFAULT 0,
    total_carbohydrates DECIMAL(6,2) DEFAULT 0,
    total_fat DECIMAL(6,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建食物记录表
CREATE TABLE food_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    meal_record_id UUID REFERENCES meal_records(id) ON DELETE CASCADE,
    food_id UUID REFERENCES foods(id) ON DELETE CASCADE,
    portion DECIMAL(5,2) NOT NULL,
    calories DECIMAL(6,2) NOT NULL,
    protein DECIMAL(6,2) NOT NULL,
    carbohydrates DECIMAL(6,2) NOT NULL,
    fat DECIMAL(6,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_nutrition_goals_user_id ON nutrition_goals(user_id);
CREATE INDEX idx_meal_records_user_id ON meal_records(user_id);
CREATE INDEX idx_meal_records_user_date ON meal_records(user_id, record_date);
CREATE INDEX idx_food_records_user_id ON food_records(user_id);
CREATE INDEX idx_food_records_meal_id ON food_records(meal_record_id);
CREATE INDEX idx_food_records_food_id ON food_records(food_id);
```

### 1.2 服务器启动

1. 确保数据库连接配置正确（在 `main.go` 文件中）
2. 运行服务器：
   ```bash
   cd d:\Program Files (x86)\Golang\nutrition_tracker
   go run cmd/server/main.go
   ```

3. 检查服务器是否成功启动：
   - 访问健康检查接口：`http://localhost:8080/api/v1/health`
   - 访问数据库测试接口：`http://localhost:8080/api/v1/test/db`

## 2. API接口测试

### 2.1 认证相关接口

#### 注册用户
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

#### 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

**注意：** 登录成功后会返回 `token`，后续所有受保护接口都需要在请求头中添加：
```
Authorization: Bearer <your_token>
```

### 2.2 用户相关接口

#### 获取用户信息
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer <your_token>"
```

#### 更新用户信息
```bash
curl -X PUT http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"height":175.0,"weight":65.0,"age":25,"gender":"male","activity_level":"moderate"}'
```

### 2.3 营养目标接口

#### 设置营养目标
```bash
curl -X POST http://localhost:8080/api/v1/goals \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"calories":2000,"protein":150,"carbohydrates":250,"fat":67}'
```

#### 获取营养目标
```bash
curl -X GET http://localhost:8080/api/v1/goals \
  -H "Authorization: Bearer <your_token>"
```

#### 计算营养目标
```bash
curl -X POST http://localhost:8080/api/v1/goals/calculate \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"height":175.0,"weight":65.0,"age":25,"gender":"male","activity_level":"moderate","goal_type":"maintain"}'
```

### 2.4 餐次记录接口

#### 创建餐次记录
```bash
curl -X POST http://localhost:8080/api/v1/meals \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"meal_type":"breakfast","record_date":"2024-05-20"}'
```

#### 获取当日餐次记录
```bash
curl -X GET "http://localhost:8080/api/v1/meals?date=2024-05-20" \
  -H "Authorization: Bearer <your_token>"
```

#### 获取特定餐次记录
```bash
curl -X GET http://localhost:8080/api/v1/meals/<meal_id> \
  -H "Authorization: Bearer <your_token>"
```

#### 删除餐次记录
```bash
curl -X DELETE http://localhost:8080/api/v1/meals/<meal_id> \
  -H "Authorization: Bearer <your_token>"
```

### 2.5 食物记录接口

#### 创建食物记录
```bash
curl -X POST http://localhost:8080/api/v1/food-records \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"meal_record_id":"<meal_id>","food_id":"<food_id>","portion":100}'
```

#### 获取当日食物记录
```bash
curl -X GET "http://localhost:8080/api/v1/food-records?date=2024-05-20" \
  -H "Authorization: Bearer <your_token>"
```

#### 获取特定餐次的食物记录
```bash
curl -X GET "http://localhost:8080/api/v1/food-records/meal?meal_id=<meal_id>" \
  -H "Authorization: Bearer <your_token>"
```

#### 获取特定食物记录
```bash
curl -X GET http://localhost:8080/api/v1/food-records/<food_record_id> \
  -H "Authorization: Bearer <your_token>"
```

#### 更新食物记录
```bash
curl -X PUT http://localhost:8080/api/v1/food-records/<food_record_id> \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"portion":150}'
```

#### 删除食物记录
```bash
curl -X DELETE http://localhost:8080/api/v1/food-records/<food_record_id> \
  -H "Authorization: Bearer <your_token>"
```

## 3. 测试顺序建议

1. 先测试数据库连接和服务器启动
2. 注册新用户
3. 用户登录获取token
4. 更新用户个人信息
5. 计算或设置营养目标
6. 创建餐次记录
7. 为餐次记录添加食物记录
8. 测试查询、更新和删除功能

## 4. 常见问题排查

### 4.1 数据库连接失败
- 检查PostgreSQL服务是否启动
- 检查数据库连接配置是否正确
- 检查数据库密码是否正确

### 4.2 认证失败
- 检查token是否正确
- 检查token是否过期
- 检查请求头格式是否正确

### 4.3 接口返回500错误
- 查看服务器日志获取详细错误信息
- 检查数据库表结构是否正确
- 检查请求参数格式是否正确

### 4.4 接口返回400错误
- 检查请求参数是否符合要求
- 检查必填字段是否都已提供

## 5. 测试工具推荐

- **Postman**：图形化API测试工具
- **curl**：命令行工具（已内置示例）
- **Swagger UI**：API文档和测试工具（如需可自行集成）

## 6. 功能验证清单

- [ ] 用户认证功能正常工作
- [ ] 用户信息管理功能正常
- [ ] 营养目标计算和设置功能正常
- [ ] 餐次记录CRUD功能正常
- [ ] 食物记录CRUD功能正常
- [ ] 营养数据计算准确
- [ ] 权限控制有效
- [ ] 错误处理合理