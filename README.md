# KK-Chat

基于 **Go（Echo + sqlc + Clean Architecture）** 与 **Vue 3（TypeScript + Pinia + ant-design-vue）** 的即时通讯应用：私聊、群聊、好友关系、在线状态、图片消息。

> 2026-09 完成全面重构。设计文档：[`docs/superpowers/specs/2026-09-12-kk-chat-refactor-design.md`](docs/superpowers/specs/2026-09-12-kk-chat-refactor-design.md)；实现计划：[`docs/superpowers/plans/`](docs/superpowers/plans/)

## 功能

- 注册 / 登录（邮箱验证码 + 图形验证码）、登出即全局踢下线
- 搜索用户；**陌生人可直接发起会话（限 3 条）**，添加好友后不限
- 私聊 / 群聊实时消息（WebSocket 推送），图片消息上传
- 群管理：建群、搜索加群、**群成员邀请**
- 最近会话列表（按最后消息时间排序）、历史消息游标分页
- 在线状态：心跳保活、上下线广播、多端登录
- 用户资料：昵称 / 手机号 / 邮箱（改邮箱需验证码）/ 签名 / 生日 / 头像

## 技术栈

| 层 | 技术 |
|---|---|
| 后端框架 | Echo v4（HTTP）+ gorilla/websocket（纯推送通道） |
| 架构 | Clean Architecture：`domain → usecase → adapter`，端口-适配器、手动依赖注入 |
| MySQL | sqlc（编译期类型安全 SQL） |
| Redis | go-redis/v9：token 会话 / 验证码 / 在线状态 / 最近会话 ZSET / 非好友限流 |
| MongoDB | 官方 driver：消息存储（会话键 + TTL 90 天） |
| 鉴权 / 日志 | JWT(HS256) + Redis 白名单；slog 结构化日志；bcrypt 密码哈希 |
| 前端 | Vue 3.5 / TypeScript(strict) / Vite 6 / Pinia / vue-router 4 / ant-design-vue 4 |
| 前端工程化 | ESLint 9 (flat config) + Prettier + husky + lint-staged |

## 目录结构

```
server-go/
├── cmd/server/            # 组装根：手动 DI、优雅关闭
├── migrations/schema.sql  # MySQL 建表脚本
├── sqlc.yaml              # sqlc 代码生成配置
├── scripts/ws_e2e.mjs     # WS 端到端验证脚本（Node 22+）
└── internal/
    ├── domain/            # 领域核心：实体、ConversationID 值对象、哨兵错误（零依赖）
    ├── usecase/           # 应用层：auth/user/friend/group/chat/presence
    │   ├── port/          #   端口接口与 DTO（repo/store/notifier/usecase 契约）
    │   └── apperror/      #   统一业务错误类型与错误码
    ├── adapter/
    │   ├── http/          # 入站：Echo 路由/handler/中间件（错误渲染、JWT、日志）
    │   ├── ws/            # 入站：WebSocket Hub/Client/升级入口
    │   ├── persistence/   # 出站：mysql(sqlc) / redis / mongo 适配器
    │   └── notify/        # Notifier 端口 → ws Hub 桥接
    ├── platform/          # JWT/bcrypt/图形码/邮件/时钟/slog
    └── config/            # TOML + 环境变量覆盖

front-ui-vue/src/
├── types/       # 与后端协议对齐的类型（DTO / WS 事件 / 错误码）
├── api/         # 全类型化 axios 客户端 + 六模块
├── ws/          # WSClient（类型化事件订阅、指数退避重连、心跳）+ store 接线
├── stores/      # Pinia：auth / chat / contacts / presence（单一事实来源）
├── components/  # 通用组件（消息气泡、用户卡片、搜索/建群/邀请弹窗…）
├── views/       # login / home（三栏布局）/ 错误页
└── scripts/ui_smoke.py  # Playwright UI 冒烟脚本
```

## 快速开始

### 依赖服务

需要 MySQL 8、Redis ≥ 6.2、MongoDB。Docker 示例：

```bash
docker run -d --name kk-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=<密码> mysql:8
docker run -d --name kk-redis -p 6379:6379 redis:7-alpine
docker run -d --name kk-mongo -p 27017:27017 mongo:7
```

### 后端（Go 1.22+）

```bash
# 1. 建库建表
mysql -uroot -p < server-go/migrations/schema.sql

# 2. 准备配置（填写本地 MySQL/Redis/Mongo/SMTP 真实值；可用 KK_* 环境变量覆盖）
cp server-go/config.example.toml server-go/config.toml

# 3. 启动（默认 :8080，config 可改）
cd server-go && go run ./cmd/server --config config.toml
```

健康检查：`GET /api/v1/healthz`。修改 SQL 查询后执行 `sqlc generate`。

### 前端（Node 20+）

```bash
cd front-ui-vue
yarn install          # 或 npm install
yarn dev              # 开发服务器（默认 :9234）
yarn build            # 构建（含 vue-tsc 类型检查）
yarn lint             # ESLint
```

环境地址在 `.env.development` / `.env.production`：`VITE_API_URL`（后端 host）、`VITE_WEBSOCKET_URL`（`ws://<host>/api/v1/ws`）。前端直连后端，CORS 已在服务端放开。

## 协议概览

**HTTP**：RESTful `/api/v1`，统一信封 `{"code":0,"message":"ok","data":...,"request_id":"..."}`；HTTP 状态码语义化（400/401/403/404/409/429/500），业务错误码分段（1xxxx 通用 / 2xxxx 认证 / 3xxxx 用户 / 4xxxx 好友 / 5xxxx 群 / 6xxxx 会话）。

| 组 | 端点 |
|---|---|
| auth | `POST /auth/{register,login,logout,captcha,email-code}` |
| users | `GET /users/me`、`PATCH /users/me`、`GET /users/:id`、`GET /users?search=` |
| friends | `POST /friends`、`GET /friends` |
| groups | `POST /groups`、`GET /groups?search=`、`POST /groups/:id/members`（加群）、`POST /groups/:id/invite`（邀请） |
| chat | `POST /messages`、`GET /conversations`、`GET /conversations/:id/messages?cursor=&limit=` |
| files | `POST /files/images` |

**WebSocket**：`GET /api/v1/ws?token=`，**仅做服务端推送**（发消息走 HTTP，请求-响应更可靠）。事件：`ping/pong`（心跳）、`chat.message`（新消息）、`chat.recent_updated`（会话摘要，按接收者视角）、`presence.changed`（上下线）、`group.updated`（群变更）、`system.kick`（踢下线）。

**会话键**：私聊 `u:{小ID}_{大ID}`（双方共享，历史一次查询）、群聊 `g:{群ID}`；Mongo 消息 TTL 默认 90 天。

## 验证脚本

```bash
# WS 端到端（需后端已启动 + 两个已注册账号的 token 写入 /tmp/t1.txt /tmp/t2.txt）
node server-go/scripts/ws_e2e.mjs        # 12 项断言：presence/心跳/推送/踢人/token失效

# UI 冒烟（需 pip install playwright && playwright install chromium；前后端已启动）
python3 front-ui-vue/scripts/ui_smoke.py # 16 项断言：登录→会话→聊天→联系人→群邀请
```

## 注意事项

- ⚠️ **BREAKING**：2026-09 重构后 HTTP/WS 协议与存储结构全部重新设计，**旧数据不兼容、前后端需同步部署**；部署前先执行 `migrations/schema.sql`。
- 非好友私聊限流 3 条/24h（`config.toml [message] non_friend_limit`）。
- 图片上传限 5MB，支持 jpg/png/gif/webp，存于 `server-go/static/uploads/{yyyyMM}/`。
