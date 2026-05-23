# 净莲阁微信小程序后端

净莲阁微信小程序的 Go 后端服务，提供菜单、功德榜、关于信息、用户认证和管理员内容管理接口。

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
- JWT token 生成和验证
- 可选认证中间件（支持登录和匿名访问）
- 管理员登录与管理员专用 JWT 鉴权

### 🛠️ 管理员模块
- 配置文件管理员账号
- bcrypt 密码哈希校验
- 菜品新增和软删除
- 管理员接口与普通用户 JWT 隔离

### 📋 其他功能
- 完整的 Swagger API 文档
- 统一的错误处理和响应格式
- 详细的操作日志记录
- 参数验证和安全检查

## 技术栈

- **Web 框架**: Gin
- **数据库**: MySQL + GORM
- **认证**: JWT
- **密码校验**: bcrypt
- **配置管理**: INI 配置文件
- **API 文档**: Swagger
- **日志**: 自定义日志组件
- **测试**: Go test + testify

## 快速开始

### 1. 环境要求
- Go 1.24.4+
- MySQL 5.7+

### 2. 安装依赖
```bash
go mod download
```

### 3. 配置服务
编辑 `conf/app.ini`：
- `[database]`：配置 MySQL 连接信息。
- `[app]`：配置 `JwtSecret`、`JwtExpire`、端口和运行时目录。
- `[admin]`：配置管理员用户名和 bcrypt 密码哈希。

本地默认管理员账号：
- 用户名：`admin`
- 密码：`jingliange-admin`

生产环境上线前必须替换默认管理员密码哈希。

### 4. 运行服务
```bash
go run cmd/main.go
```

服务默认运行在 `:8000` 端口。

### 5. 查看API文档
访问 `http://localhost:8000/swagger/index.html` 查看完整的API文档。

## 测试

### 快速测试
```bash
go test ./pkg/util ./internal/router/api/admin ./internal/router
```

### 完整测试
```bash
go test ./...
```

注意：`internal/router/api/v1` 中部分测试会连接配置里的 MySQL 测试库。如果本地没有对应数据库，完整测试会因为 `ERROR_DB` 失败。

### 功能测试
项目提供了功能测试脚本：
```bash
# 测试菜单点赞功能
./test_menu_liked.sh
```

## API接口

### 菜单接口
- `POST /api/v1/menu/getMenu` - 获取菜单列表
- `POST /api/v1/menu/getMenuByID` - 获取单个菜品详情
- `POST /api/v1/menu/like` - 菜品点赞/取消点赞
- `POST /api/v1/menu/comment` - 添加菜品评论
- `POST /api/v1/menu/getComments` - 获取菜品评论列表

### 捐赠接口
- `POST /api/v1/donation/getDonationList` - 获取功德榜列表
- `POST /api/v1/donation/createDonation` - 创建捐款记录
- `POST /api/v1/donation/getDonationStats` - 获取捐款统计

### 认证接口
- `POST /api/v1/auth/login` - 用户认证登录

### 管理员接口
- `POST /api/admin/login` - 管理员登录，返回管理员 JWT
- `POST /api/admin/uploadMenuItem` - 新增菜品，需要管理员 JWT
- `DELETE /api/admin/deleteMenuItem` - 删除菜品，需要管理员 JWT

### 其他接口
- `GET /api/v1/about/getDescription` - 获取关于信息
- `GET /api/v1/about/getTopImage` - 获取首页头图
- `POST /api/v1/about/getActivityList` - 获取活动列表
- `POST /api/v1/about/getImageList` - 获取图片列表

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

### 本地编译
```bash
make build
# 或者
go build -o bin/jingliange_server cmd/main.go
```

### Linux 生产编译
```bash
GOOS=linux GOARCH=amd64 go build -o bin/jingliange_server cmd/main.go
```

### 配置
确保生产环境的配置文件 `conf/app.ini` 正确设置：
- 数据库连接信息
- `JwtSecret`
- 管理员账号和 `PasswordHash`
- 服务端口
- `RuntimeRootPath`

### 运行
```bash
./bin/jingliange_server
```

### 当前生产部署
当前线上服务部署在 `49.234.22.169`：
- systemd 服务：`jlg.service`
- 工作目录：`/opt/jlg`
- 可执行文件：`/opt/jlg/jingliange_server`
- 配置文件：`/opt/jlg/conf/app.ini`
- nginx 将 `https://jingliange.com/api/` 转发到 `127.0.0.1:8000/api/`

部署时建议：
1. 先备份 `/opt/jlg/jingliange_server` 和 `/opt/jlg/conf/app.ini`。
2. 上传新的 Linux amd64 二进制文件。
3. 替换 `/opt/jlg/jingliange_server`。
4. 执行 `sudo systemctl restart jlg.service`。
5. 验证公开接口和管理员登录接口。

部署后健康检查：
```bash
curl https://jingliange.com/api/v1/about/getDescription
curl -X POST https://jingliange.com/api/admin/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"jingliange-admin"}'
```

## 许可证

[MIT License](LICENSE)

## 贡献

欢迎提交Issue和Pull Request来帮助改进项目。

---

如有问题，请联系项目维护者。
