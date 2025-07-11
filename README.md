# 净莲阁微信小程序后端

净莲阁微信小程序的Go后端服务，提供菜单管理、捐赠功德榜、用户认证等功能。

## 功能特性

### 🍽️ 菜单模块
- 菜单列表查询（支持模糊搜索、分页）
- 单个菜品详情查询
- 菜品点赞功能（支持用户点赞状态）
- 菜品评论功能
- 菜品上传和软删除（管理员功能）

### 💰 捐赠功德榜模块
- 捐款记录创建（金额、昵称、留言）
- 功德榜查询（支持年份、时间段、昵称筛选）
- 捐款统计（总金额、总人次）
- 多维度排序（按时间或金额）
- 分页查询支持

### 🔐 用户认证模块
- 微信小程序用户认证
- JWT token生成和验证
- 可选认证中间件（支持登录和匿名访问）
- 用户信息管理

### 📋 其他功能
- 完整的Swagger API文档
- 统一的错误处理和响应格式
- 详细的操作日志记录
- 参数验证和安全检查

## 技术栈

- **Web框架**: Gin
- **数据库**: MySQL + GORM
- **认证**: JWT
- **配置管理**: INI配置文件
- **API文档**: Swagger
- **日志**: 自定义日志组件
- **测试**: Go测试框架 + testify

## 快速开始

### 1. 环境要求
- Go 1.24.4+
- MySQL 5.7+

### 2. 安装依赖
```bash
go mod download
```

### 3. 配置数据库
复制配置文件并修改数据库连接信息：
```bash
cp conf/app.ini.example conf/app.ini
# 编辑 conf/app.ini 中的数据库配置
```

### 4. 运行服务
```bash
go run cmd/main.go
```

服务默认运行在 `:8000` 端口。

### 5. 查看API文档
访问 `http://localhost:8000/swagger/index.html` 查看完整的API文档。

## 测试

### 运行所有测试
```bash
./run_tests.sh
```

### 运行特定测试
```bash
# JWT功能测试
go test ./pkg/util/jwt_simple_test.go ./pkg/util/jwt.go -v

# API层测试
go test ./internal/router/api/v1/api_simple_test.go -v
```

更多测试信息请参考 [TESTING.md](TESTING.md)。

### 功能测试
项目提供了功能测试脚本：
```bash
# 测试菜单点赞功能
./test_menu_liked.sh

# 测试日志记录功能
./test_logging.sh
```

## API接口

### 菜单接口
- `POST /api/v1/menu/getMenu` - 获取菜单列表
- `POST /api/v1/menu/getMenuByID` - 获取单个菜品详情
- `POST /api/v1/menu/like` - 菜品点赞/取消点赞
- `POST /api/v1/menu/comment` - 添加菜品评论
- `POST /api/v1/menu/getMenuComments` - 获取菜品评论列表

### 捐赠接口
- `POST /api/v1/donation/getDonationList` - 获取功德榜列表
- `POST /api/v1/donation/createDonation` - 创建捐款记录
- `POST /api/v1/donation/getDonationStats` - 获取捐款统计

### 认证接口
- `POST /api/v1/auth/login` - 用户认证登录

### 其他接口
- `GET /api/v1/about/getDescription` - 获取关于信息

## 项目结构

```
├── cmd/                    # 应用入口
│   └── main.go
├── internal/               # 内部代码
│   ├── model/             # 数据模型
│   ├── repo/              # 数据访问层
│   └── router/            # 路由和控制器
├── pkg/                   # 工具包
│   ├── app/               # 应用层工具
│   ├── e/                 # 错误码定义
│   ├── logging/           # 日志组件
│   ├── setting/           # 配置管理
│   └── util/              # 通用工具
├── docs/                  # Swagger文档
├── conf/                  # 配置文件
├── runtime/               # 运行时文件
│   ├── logs/              # 日志文件
│   └── uploads/           # 上传文件
└── *.sh                   # 测试脚本
```

## 开发指南

### 添加新的API接口
1. 在 `internal/model/` 中定义数据模型
2. 在 `internal/repo/` 中实现数据访问逻辑
3. 在 `internal/router/api/v1/` 中实现API处理器
4. 添加Swagger文档注释
5. 更新路由配置
6. 编写单元测试

### 代码规范
- 使用有意义的变量和函数名
- 添加必要的注释和文档
- 遵循Go编码规范
- 所有公开函数都应该有测试

### 提交规范
- 提交前运行测试确保通过
- 提交信息使用中文，简洁明了
- 大功能分多次小提交

## 部署

### 编译
```bash
make build
# 或者
go build -o bin/jingliange_server cmd/main.go
```

### 配置
确保生产环境的配置文件 `conf/app.ini` 正确设置：
- 数据库连接信息
- JWT密钥
- 日志级别
- 服务端口

### 运行
```bash
./bin/jingliange_server
```

## 许可证

[MIT License](LICENSE)

## 贡献

欢迎提交Issue和Pull Request来帮助改进项目。

---

如有问题，请联系项目维护者。