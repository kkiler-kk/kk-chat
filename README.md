# KK-Chat

基于 **Go (Echo + sqlc + Clean Architecture)** 与 **Vue 3 (TypeScript + Pinia + ant-design-vue)** 的即时通讯应用，支持私聊、群聊、好友关系、在线状态等基础功能。

> 2026-09 全面重构：后端采用 Clean Architecture（domain/usecase/adapter 分层、端口-适配器、手动依赖注入），协议与存储键空间全部重新设计。设计文档见 [`docs/superpowers/specs/2026-09-12-kk-chat-refactor-design.md`](docs/superpowers/specs/2026-09-12-kk-chat-refactor-design.md)。

## 技术栈

### 后端（server-go/）

| 组件 | 技术 |
|---|---|
| Web 框架 | Echo v4（HTTP）+ gorilla/websocket（纯推送通道） |
| MySQL | sqlc（编译期类型安全 SQL） |
| Redis | go-redis/v9（token 会话/验证码/在线状态/最近会话 ZSET/非好友限流） |
| MongoDB | 官方 driver（消息存储，会话键 `u:{小ID}_{大ID}` / `g:{群ID}`，TTL 默认 90 天） |
| 鉴权 | JWT(HS256) + Redis 白名单（登出/踢人即时失效） |
| 日志/密码 | slog 结构化日志、bcrypt |

### 前端（front-ui-vue/）

Vue 3.5+ / TypeScript(strict) / Vite / Pinia / vue-router / ant-design-vue 4.x；WSClient 类型化事件订阅 + 指数退避重连；ESLint 9 + Prettier + husky。

## 运行

### 依赖服务

本地需要 MySQL 8、Redis ≥6.2、MongoDB（例如 `docker run -d --name kk-mongo -p 27017:27017 mongo:7`）。

### 后端

```bash
# 1. 初始化数据库
mysql -uroot -p < server-go/migrations/schema.sql

# 2. 准备配置（填本地真实的 MySQL/Redis/Mongo/邮箱配置）
cp server-go/config.example.toml server-go/config.toml

# 3. 启动（Go 1.22+）
cd server-go && go run ./cmd/server --config config.toml
```

服务默认监听 `:8080`（config.toml 可改），健康检查 `GET /api/v1/healthz`。
若修改了 sqlc 查询：`sqlc generate`（配置见 `server-go/sqlc.yaml`）。

### 前端

```bash
cd front-ui-vue
yarn install
yarn dev        # 开发
yarn build      # 构建（含 vue-tsc 类型检查）
```

环境变量见 `.env.development`（`VITE_API_URL`、`VITE_WEBSOCKET_URL`）。

## 目录结构

```
server-go/
├── cmd/server/          # 组装根：手动 DI、优雅关闭
├── migrations/          # MySQL schema
├── internal/
│   ├── domain/          # 领域核心：实体、ConversationID、哨兵错误（零依赖）
│   ├── usecase/         # 应用层：业务用例 + port/ 端口接口 + apperror/
│   ├── adapter/         # 适配器：http(Echo)/ws/persistence(mysql,redis,mongo)/notify
│   ├── platform/        # JWT/bcrypt/验证码/邮件/时钟/日志
│   └── config/          # TOML + 环境变量
front-ui-vue/src/
├── types/ api/ ws/ stores/ composables/ components/ views/ router/ utils/
```

## API 概览

RESTful `/api/v1`（统一信封 `{code, message, data, request_id}`）：auth（注册/登录/登出/验证码）、users、friends、groups、messages、conversations、files/images；WS `GET /api/v1/ws?token=` 仅做服务端推送（`chat.message`、`chat.recent_updated`、`presence.changed`、`group.updated`、`system.kick`）。详见设计文档 §3。

## 功能列表

- 登录 / 注册（图形验证码 + 邮箱验证码）
- 搜索用户 / 添加好友 / 好友列表（含在线状态）
- 私聊（非好友限流 3 条）/ 群聊 / 最近会话（Redis ZSET 排序）
- 历史消息（游标分页）/ 图片消息上传
- 多端登录 / 心跳与断线重连 / 登出踢下线
- 朋友圈、设置（未完成）
