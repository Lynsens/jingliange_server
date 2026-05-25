# 净莲阁微信小程序后端

净莲阁微信小程序的 Go 后端服务，提供菜单、功德榜、关于信息、用户认证和管理员内容管理接口。

## 功能特性

### 🍽️ 菜单模块
- 菜单列表查询（支持模糊搜索、分页）
- 单个菜品详情查询
- 菜品点赞功能（支持用户点赞状态）
- 菜品评论功能
- 菜品上传、编辑和软删除（管理员功能）
- 今日推荐菜品标注：同一时间最多只有一个菜品可作为今日推荐
- 菜品下架：使用 archive 状态下架非当季菜品，区别于删除
- 评论展示微信昵称和头像；同一用户对同一菜品的点赞和评论共用一条反馈记录，避免重复点赞/重复评论

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
- 菜品新增、编辑、软删除和今日推荐标注
- 菜品 archive 下架功能：管理员可下架非当季菜品，也可重新上架
- 活动管理功能：管理员可新增、编辑、删除和置顶近期活动
- 评论管理功能：管理员可查看、搜索和删除用户评论；删除评论不会影响点赞状态
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
复制配置模板后编辑本地配置：

```bash
cp conf/app.example.ini conf/app.ini
cp conf/app_test.example.ini conf/app_test.ini
```

编辑 `conf/app.ini`：
- `[database]`：配置 MySQL 连接信息。
- `[app]`：配置 `JwtSecret`、`JwtExpire`、端口和运行时目录。
- `[admin]`：配置管理员用户名和 bcrypt 密码哈希。
- `[ops]`：配置只读维护后台可读取的 Nginx 和后端应用日志路径。

管理员账号：
- 用户名：`admin`
- 密码：不写入代码或文档；通过 `[admin].PasswordHash` 配置 bcrypt 哈希。

生产环境上线前必须使用单独的管理员密码哈希，禁止提交明文密码。

`conf/app.ini` 和 `conf/app_test.ini` 包含数据库密码、JWT secret 和管理员密码哈希，只保留在本地或服务器上，不进入版本控制。仓库中只提交 `conf/app.example.ini` 和 `conf/app_test.example.ini`。

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

注意：`internal/router/api/v1` 中部分测试会连接配置里的 MySQL 测试库。如果测试库不可访问，完整测试会因为 `ERROR_DB` 失败。

### 测试数据库
测试数据库已经迁移到服务器上的独立 database：

- MySQL 服务器：`49.234.22.169:3306`
- 生产库：`jlg`
- 测试库：`jlg_test`
- 测试账号：`jlg_test`
- 测试配置文件：`conf/app_test.ini`
- 测试库账号只授予 `jlg_test.*` 权限，避免测试误写生产数据

首次在新机器运行测试前，先从 `conf/app_test.example.ini` 复制出 `conf/app_test.ini`，再填入测试库密码、JWT secret 和管理员 bcrypt 哈希。不要把真实测试配置提交到 Git。

当前测试库初始化方式：

```bash
# 在服务器上创建测试库和测试账号
sudo mysql -e "CREATE DATABASE IF NOT EXISTS jlg_test CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;"
sudo mysql -e "CREATE USER IF NOT EXISTS 'jlg_test'@'%' IDENTIFIED BY '<test-password>';"
sudo mysql -e "GRANT ALL PRIVILEGES ON jlg_test.* TO 'jlg_test'@'%'; FLUSH PRIVILEGES;"

# 首次初始化时从生产库复制一份结构和基础数据
sudo mysqldump --no-tablespaces jlg | sudo mysql jlg_test
```

测试库需要保留这些基础数据，避免测试依赖生产库当前内容：

- `menu.id = 1`：用于菜单详情、点赞、评论测试。
- `user.id IN ('test_user', 'test_user_2', 'test_user_123', 'wulongcha_test')`：用于功德榜和认证相关测试。

可以用下面命令验证本地开发机能连接测试库：

```bash
mysql -h 49.234.22.169 -P 3306 -u jlg_test -p jlg_test \
  -e "SELECT DATABASE(), COUNT(*) FROM menu;"
```

后续任何数据库表结构变更，都必须同时处理生产库和测试库：

1. 先更新 `docs/sql/*.sql` 中的建表或变更语句。
2. 先在测试库 `jlg_test` 执行对应 `ALTER TABLE` 或迁移语句。
3. 跑 `go test ./...`，确认测试库结构和代码兼容。
4. 部署前在生产库 `jlg` 执行同一份结构变更。
5. 生产库变更后做一次公开接口和管理员接口健康检查。
6. 在 README 或部署记录中注明已同步 `prod` 和 `test`。

禁止只改生产库或只改测试库。两个库结构不一致时，测试结果会失真，线上部署也容易在 GORM 查询字段时失败。

当前菜单和活动管理相关表结构变更包括：

- `menu.category`：菜品分类，固定为 `前菜/小菜`、`主食`、`热食`、`甜品/饮品` 之一；为空或非法值时后端默认保存为 `热食`。
- `menu.is_recommended`：今日推荐标记，最多一个正常上架菜品为 `1`。
- `menu.is_archived`：菜品下架标记，`1` 表示非当季/已下架。
- `menu.archive_time`：菜品下架时间，可为空。
- `menu_feedback.user_nickname`：评论展示昵称。
- `menu_feedback.user_avatar_url`：评论展示头像 URL。
- `menu_feedback.uk_menu_feedback_menu_user`：`menu_id + user_id` 唯一索引，用于保证同一用户对同一菜品只有一条反馈记录。
- `activity.is_top`：活动置顶标记。
- `activity.event_time`：活动时间展示文案，例如 `周六 10:30`。
- `activity.place`：活动地点。

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
- `POST /api/admin/menu/list` - 管理员菜单列表，可返回已下架菜品
- `POST /api/admin/comment/list` - 管理员评论列表，可按评论、用户、菜品和 ID 搜索
- `DELETE /api/admin/comment/delete` - 删除评论内容，不影响用户点赞状态
- `PUT /api/admin/updateMenuItem` - 更新菜品，需要管理员 JWT
- `PUT /api/admin/recommendMenuItem` - 设置今日推荐，需要管理员 JWT；设置时自动取消其他菜品推荐
- `PUT /api/admin/archiveMenuItem` - 下架或重新上架菜品；下架今日推荐菜品时自动取消推荐
- `DELETE /api/admin/deleteMenuItem` - 删除菜品，需要管理员 JWT
- `POST /api/admin/activity/list` - 管理员活动列表
- `POST /api/admin/activity/create` - 新增活动
- `PUT /api/admin/activity/update` - 更新活动
- `DELETE /api/admin/activity/delete` - 删除活动
- `PUT /api/admin/activity/top` - 设置或取消活动置顶
- `GET /api/admin/ops/summary` - 运维流量概览，只读读取配置指定的 access log
- `GET /api/admin/ops/access-logs` - Nginx 访问日志列表
- `GET /api/admin/ops/error-logs` - Nginx 错误日志列表
- `GET /api/admin/ops/app-logs` - 后端应用日志列表

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
├── website/               # jingliange.com 静态首页源文件
├── conf/                  # 配置模板；真实 app.ini/app_test.ini 本地保存
├── runtime/               # 运行时文件；日志和上传文件不进入版本控制
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

### 菜单状态规则
- `status = 1` 表示未删除，`status = 0` 表示软删除。
- `category` 表示菜品分类，固定为 `前菜/小菜`、`主食`、`热食`、`甜品/饮品` 之一，供菜单管理和首页今日推荐编排使用。
- `is_archived = 0` 表示上架中，`is_archived = 1` 表示已下架/非当季。
- 普通用户菜单接口只返回 `status = 1 AND is_archived = 0`。
- 管理员菜单列表应返回 `status = 1` 的菜单，并支持按全部、上架中、已下架筛选。
- 今日推荐只允许设置在 `status = 1 AND is_archived = 0` 的菜品上。
- 下架今日推荐菜品时，后端必须自动将该菜品 `is_recommended` 设为 `0`。

### 活动置顶规则
- 首页近期活动优先展示 `is_top = 1` 的活动。
- 活动排序建议为 `is_top DESC, create_time DESC`。
- 活动置顶可支持多个活动，不强制单一置顶。
- 活动删除使用软删除：`status = 0`。

### 官网静态页
- `website/index.html` 是 `https://jingliange.com` 的静态首页源文件。
- `website/admin/` 是 `https://jingliange.com/admin/` 的只读维护后台源文件。
- 服务器部署路径是 `/var/www/jingliange.com/index.html`。
- 维护后台部署路径是 `/var/www/jingliange.com/admin/`。
- 更新线上文件前先备份当前 HTML，例如：

```bash
TS=$(date +%Y%m%d%H%M%S)
sudo cp /var/www/jingliange.com/index.html /var/www/jingliange.com/index.html.bak.$TS
sudo install -m 0644 -o root -g root website/index.html /var/www/jingliange.com/index.html
sudo nginx -t
```

- 当前只展示工信部 ICP 备案号；没有公安备案号时不要添加公安备案占位或虚假编号。

### 维护后台
- 维护后台复用管理员登录接口 `POST /api/admin/login` 和 AdminJWT。
- 第一版只读展示日志和流量，不提供重启服务、删除日志、清理文件或修改配置功能。
- 后端只读取 `[ops]` 配置中指定的日志路径，不接受任意文件路径参数。
- 生产环境需要确保运行后端服务的用户对 `/var/log/nginx/jingliange.com_access.log` 和 `/var/log/nginx/jingliange.com_error.log` 有只读权限。

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
  -d '{"username":"admin","password":"<admin-password>"}'
```

## 许可证

[MIT License](LICENSE)

## 贡献

欢迎提交Issue和Pull Request来帮助改进项目。

---

如有问题，请联系项目维护者。
