# kk-chat 全面重构设计文档

- 日期：2026-09-12
- 状态：已获用户确认
- 范围：`server-go`（Go 后端）+ `front-ui-vue`（Vue 前端），一次性大重构

## 1. 背景与目标

kk-chat 是一个私聊/群聊即时通讯项目（Go 后端 + Vue3 前端 + MySQL/Redis/MongoDB）。现有代码存在四类核心问题：

1. **目录/模块结构混乱**：后端 `common/utility` 大杂烩、`global.Init()` 直接跑起 HTTP 服务、两个 router 包、无 `cmd/` 入口规范；前端组件、工具、全局事件混杂。
2. **代码重复、职责不清**：handler 每个方法重复「绑定→校验→ErrorResp→WriteJsonExit」模板；service 层透传 `*gin.Context` 且 import websocket 包；前端 `onMessage` 与 `socket.onmessage` 整段重复、消息归属判断散落在多个组件。
3. **错误处理不规范**：裸 `errors.New` 中文消息直传前端、无错误码体系、HTTP 永远 200、大量 `_ =` 忽略错误、`fmt.Println` 调试残留、gorm `.Debug()` 生产泄漏 SQL。
4. **可测试性差**：包级单例（`var UserBasicService = &userBasicService{}`）、无任何接口抽象、依赖全局变量、实际业务测试为零。

**重构目标**：采用 Clean Architecture（六边形架构）重建后端，按职责重组前端并升级工程化，同时重设计 HTTP/WS 协议与存储键空间。

**用户确认的约束与决策**：

| 决策项 | 结论 |
|---|---|
| 接口协议 / 存储结构 | 可自由重新设计，**旧数据不迁移**，MySQL/Mongo/Redis 全部重建 |
| 后端技术栈 | Echo v4 + sqlc + redis/v9 + mongo-driver + slog，Go 1.22+ |
| 前端技术栈 | Vue 3.5+ / TS / Vite 6 / Pinia / ant-design-vue 4.x，移除 vuex 与冗余 emoji 库 |
| 测试 | 暂不写自动化测试；验收 = 编译/vet/type-check/build 通过 + 手工联调清单 |
| 节奏 | 一次性大重构，旧代码在新实现通过验证后删除 |
| 架构风格 | Clean Architecture / 六边形（用户明确选择，优先于推荐的务实分层） |

## 2. 后端架构

### 2.1 依赖规则

```
main(组装) → adapter(适配器) → usecase(用例/端口) → domain(领域核心)
```

依赖只能向内：`domain` 不 import 任何第三方或内部包；`usecase` 只依赖 `domain` 和自己定义的端口接口；Echo、sqlc、redis、mongo、websocket 均为可替换适配器。

### 2.2 目录结构

```
server-go/
├── cmd/server/main.go          # 组装根：配置→依赖图→启动→优雅关闭
├── internal/
│   ├── domain/                 # 领域核心（零依赖）
│   │   ├── user.go             # User 实体 + ErrUserNotFound/ErrEmailTaken 等哨兵错误
│   │   ├── message.go          # Message、ConversationID 值对象、ContentType
│   │   ├── friend.go           # 好友关系
│   │   ├── group.go            # Group、GroupMember、角色（owner/member）
│   │   └── presence.go         # 在线状态
│   ├── usecase/                # 应用层
│   │   ├── port/               # 端口接口（输出：repo/store/notifier/clock；输入：各 UseCase 接口）
│   │   ├── auth.go             # 注册/登录/登出/图形码/邮箱验证码
│   │   ├── user.go             # 资料查询/修改/搜索
│   │   ├── friend.go           # 加好友/好友列表/关系判断
│   │   ├── group.go            # 建群/加群/群列表/按名搜索
│   │   ├── chat.go             # 发消息/历史消息/最近会话
│   │   └── presence.go         # 心跳/上下线
│   ├── adapter/
│   │   ├── http/               # 入站适配器（Echo）
│   │   │   ├── server.go       # 实例/中间件链/静态文件/优雅关闭
│   │   │   ├── middleware/     # jwt 鉴权、recover、request-id、cors、统一错误处理
│   │   │   ├── handler/        # auth/user/friend/group/chat/file handler（薄）
│   │   │   └── dto/            # 请求/响应 DTO + 统一响应体
│   │   ├── ws/                 # 入站适配器（WebSocket 纯推送通道）
│   │   │   ├── hub.go          # uid→[]*Client 注册/注销/按用户·按会话广播
│   │   │   ├── client.go       # 读写循环、心跳超时
│   │   │   └── protocol.go     # 类型化事件信封
│   │   ├── persistence/        # 出站适配器
│   │   │   ├── mysql/          # sqlc 生成代码 + UserRepo/FriendRepo/GroupRepo 实现
│   │   │   ├── redis/          # TokenStore/CaptchaStore/CodeStore/PresenceStore/RecentChatStore
│   │   │   └── mongo/          # MessageRepo 实现
│   │   └── notify/             # Notifier 端口实现：usecase → ws hub 桥接
│   ├── config/                 # TOML + 环境变量覆盖
│   └── platform/               # slog、JWT、bcrypt、校验器、ID 生成
├── migrations/schema.sql
├── sqlc.yaml
└── go.mod                      # Go 1.22+
```

### 2.3 领域模型与存储设计

**MySQL（sqlc 管理，schema 重建）**

- `users`：`id, identity(唯一账号), name, password_hash(bcrypt，去掉独立 salt 字段), email, phone, avatar, signature, birth_date, is_admin, status(1正常/2封禁), created_at, updated_at`
  - 移除 `login_time/heartbeat_time/login_out_time/is_logout`：在线状态属易变数据，改由 Redis TTL 管理（现状每次心跳写 MySQL 属明显错用）。
- `friends`：`(user_id, friend_id)` 复合主键，**双向各存一行**，查询无需 OR 两个方向。
- `groups`：`id, name, avatar, owner_id, created_at`；`group_members`：`(group_id, user_id)` 复合主键 + `role, joined_at`。

**MongoDB `messages` 集合**

- 引入 `conversation_id` 值对象：私聊 = `u:{小ID}_{大ID}`（双方共享同一会话键），群聊 = `g:{群ID}`。
- 取代现状 `type:1--1->2` / `type:1--2->1` 双向键（查历史需查两次再内存合并排序）；新设计一次查询完成。
- 文档：`{_id, conversation_id, type(private|group), sender_id, content, content_type(text|image), created_at}`。
- 索引：`(conversation_id, created_at DESC)`；TTL 索引控制消息保留期，默认 90 天。

**Redis 键空间**

| 键 | 类型 | 用途 | TTL |
|---|---|---|---|
| `token:{jti}` | string(JSON) | 登录会话（登出/踢人即时失效） | 同 token 有效期 |
| `captcha:{id}` | string | 图形验证码 | 5min |
| `verify:email:{email}` | string | 邮箱验证码 | 10min |
| `presence:{uid}` | string | 在线心跳 | 2min，过期即离线 |
| `recent:{uid}` | ZSET(member=会话ID, score=最后消息时间) | 最近会话列表 | 无 |
| `conv:{conversation_id}` | HASH | 会话摘要：`type`、最后一条消息 JSON（内容/发送者/时间） | 无 |
| `msg:limit:{conversation_id}:{uid}` | string(INCR) | 非好友消息限流计数 | 24h |

最近会话从「JSON 数组反序列化→遍历→写回」（含 float64 比较缺陷）改为 ZSET 单命令读写，天然按时间排序。

### 2.4 依赖注入

`main.go` 按依赖顺序显式构造：config → logger → mysql/sqlc → redis → mongo → repos → hub/notifier → usecases → handlers → echo server。无全局变量、无包级单例；所有 usecase 结构体持有端口接口，构造函数注入。`Clock` 端口封装时间获取。

### 2.5 优雅关闭

`signal.NotifyContext(SIGINT/SIGTERM)` → `echo.Shutdown(ctx)` → hub 关闭全部连接 → 释放 db/redis/mongo 连接池。

## 3. 协议设计

### 3.1 HTTP API（RESTful，前缀 `/api/v1`，🔒=需 JWT）

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/auth/register` | 注册（邮箱验证码） |
| POST | `/auth/login` | 登录（图形验证码）→ token + 用户信息 |
| POST | `/auth/logout` 🔒 | 登出 |
| POST | `/auth/captcha` | 图形验证码（id + base64 图） |
| POST | `/auth/email-code` | 发送邮箱验证码 |
| GET | `/users/me` 🔒 | 当前用户资料 |
| PATCH | `/users/me` 🔒 | 修改资料（改邮箱需验证码） |
| GET | `/users/{id}` | 查用户（隐私字段按好友关系过滤） |
| GET | `/users?search=xx` 🔒 | 搜索用户 |
| POST | `/friends` 🔒 | 添加好友 `{user_id}` |
| GET | `/friends` 🔒 | 好友列表 |
| POST | `/groups` 🔒 | 建群 |
| POST | `/groups/{id}/members` 🔒 | 加群 |
| GET | `/groups?search=xx` 🔒 | 搜群 |
| GET | `/groups` 🔒 | 我的群列表 |
| POST | `/messages` 🔒 | 发消息 `{conversation_id, content, content_type}`，同步返回消息体，服务端经 WS 推给对端 |
| GET | `/conversations` 🔒 | 最近会话列表 |
| GET | `/conversations/{id}/messages?cursor=&limit=` 🔒 | 历史消息（游标分页） |
| POST | `/files/images` 🔒 | 上传图片 |

**发送消息走 HTTP 而非 WS**：请求-响应可靠（失败可感知、可重试），WS 退化为纯推送通道。最近会话列表同理迁至 HTTP GET。

### 3.2 统一响应与错误码

响应体（成功/失败同构）：

```json
{ "code": 0, "message": "ok", "data": {}, "request_id": "..." }
```

HTTP 状态码使用真实语义：400 参数错、401 未认证、403 拒绝、404 不存在、429 限流（非好友消息超限）、500 内部错误。`code` 承载业务错误码：

- `0` 成功
- `1xxxx` 通用：10001 参数校验失败、10002 内部错误
- `2xxxx` 认证：20001 验证码错误、20002 账号或密码错误、20003 token 无效/过期、20004 账号被封禁
- `3xxxx` 用户、`4xxxx` 好友（40002 非好友消息超限）、`5xxxx` 群组、`6xxxx` 会话/消息

**错误处理管道**：

```
domain:   哨兵错误（ErrUserNotFound / ErrEmailTaken / ErrNotFriend …）
usecase:  apperror.Wrap(code, 用户可读消息, err)，保留内部错误链
adapter:  Echo HTTPErrorHandler 中间件统一捕获：AppError→状态码+响应体；
          未知错误→500 + slog 记录堆栈与 request_id
handler:  只 return err，不处理错误格式
```

### 3.3 鉴权与在线状态

- JWT HS256，payload：`sub(uid)`、`jti`、过期时间；secret/有效期进配置。
- Redis 白名单：登录写 `token:{jti}`；鉴权中间件校验签名 + Redis 存在性，支持登出即时失效与踢人（WS 推 `system.kick` 后断开）。
- 中间件将 `uid` 以类型化私有 key 写入 `context.Value`；usecase 只收 `ctx` + 显式参数，不感知 HTTP（替代 `c.Get("id")` + 裸断言）。
- Presence：WS 连接建立写 `presence:{uid}` 并随心跳续期；断开/TTL 过期删除，并经 hub 向好友广播 `presence.changed`。

### 3.4 WebSocket 协议（纯推送通道）

连接：`GET /api/v1/ws?token=xxx`（升级前鉴权，未认证拒绝）。信封：`{ "event": "...", "data": {...}, "ts": 1234567890 }`。

| 方向 | 事件 | 数据 |
|---|---|---|
| C→S | `ping` | — |
| S→C | `pong` | — |
| S→C | `chat.message` | 新消息体（conversation_id/sender/content/created_at） |
| S→C | `chat.recent_updated` | 会话摘要 |
| S→C | `presence.changed` | `{user_id, online}` |
| S→C | `group.updated` | 群变更通知 |
| S→C | `system.kick` | 踢下线原因 |

- Hub 维护 `uid → []*Client`（多端）；usecase 经 `Notifier` 端口下达推送（`ToUser` / `ToUsers`），**usecase 不 import ws 包**。
- 心跳：客户端每 60s `ping`，服务端刷新 presence TTL；连续 2 个周期无心跳判定离线。

## 4. 前端架构

### 4.1 目录结构

```
front-ui-vue/src/
├── types/          # 与后端协议对齐：DTO、WS 事件信封、错误码常量
├── api/            # client.ts（axios 封装）+ auth/user/friend/group/chat/file 模块，全类型化
├── ws/             # WSClient 类
├── stores/         # Pinia：auth / contacts / chat / presence（唯一事实来源）
├── composables/    # useChat、usePresence 等
├── components/     # MessageBubble、ConversationItem、UserCard、AvatarUpload…
├── views/          # login/（含 register）、home/（薄壳布局+子组件）、error/
├── router/         # 路由 + 全局守卫（无有效 token → 登录页）
├── utils/          # storage（类型化）、formatTime（清理 toolsValidate，只留用到的）
└── assets/
```

### 4.2 核心重构点

**① WSClient 类替代全局函数集**

- 单例类：`connect()/close()/on(event, handler)/off()`，`WsEvents` 映射类型约束每个事件 payload。
- 指数退避重连（1s→2s→4s…上限 30s，替代固定 5s + `location.reload()` 兜底）、心跳定时器内置。
- 消灭 `window.dispatchEvent("onmessageWS")` 模式：WSClient 收事件后直接调用注册的 store action，组件不再各自 `addEventListener` 解析原始消息；删除 `onMessage` 与 `socket.onmessage` 的重复代码。

**② chat store 收编消息逻辑**

```
state:  conversations / messagesByConv（按会话分组）/ activeConvId
action: handleIncomingMessage(msg) → 归入 messagesByConv[msg.conversation_id]，
        更新 conversations 排序与摘要（归属判断只写一次）
组件:   只读 store 的 computed，纯渲染
```

**③ 巨型组件拆分**

- `contact.vue`(279 行) → 搜索框 + 结果列表 + UserCard
- `updateUserInfo.vue`(226 行) → 表单子组件 + 弹窗容器
- `content.vue` → 消息列表 + MessageBubble（自己/他人样式内聚到 bubble）

### 4.3 依赖与工程化

| 动作 | 内容 |
|---|---|
| 升级 | Vue 3.5+、ant-design-vue 4.x、Pinia、vue-router、TypeScript 5.x、Vite 6 |
| 移除 | `vuex`、`emoji-mart-vue`、`v-emoji-picker`（保留 `vue3-emoji-picker`）、`js-pinyin`（确认引用后清理） |
| 新增 | ESLint 9（flat config：eslint-plugin-vue + @typescript-eslint）、Prettier、husky + lint-staged |
| 严格化 | tsconfig 开 `strict`，`vue-tsc` 纳入 build 门槛；API/WS 层禁止裸 `any` |
| 清理 | 脚手架残留（HelloWorld.spec.ts、默认 icons 等） |

## 5. 实施顺序与验证

**后端**（每阶段 `go build ./... && go vet ./...` 通过再继续）：

1. 骨架：go.mod(Go 1.22+)、config、platform(slog/jwt/bcrypt)、`migrations/schema.sql`、sqlc 配置+查询+生成
2. domain 实体与错误
3. ports + persistence 适配器（mysql/redis/mongo）
4. usecase 六模块（auth/user/friend/group/chat/presence）
5. ws hub + notify 桥接
6. http handler/middleware/server + main 组装 + 优雅关闭
7. 删除全部旧后端代码（含 main_test.go 中与项目无关的算法练习代码）

**前端**（`yarn type-check && yarn build` 通过为准）：

8. 依赖升级 + ESLint/Prettier/husky 落地
9. types + api client + WSClient
10. stores（auth/chat/contacts/presence）
11. views/components 按序重写：登录→注册→home 骨架→会话列表→聊天窗→联系人→群组→设置
12. 删除旧代码与脚手架残留

**手工联调清单**（需本地 MySQL/Redis/Mongo）：

> 注册（邮箱验证码）→ 登录（图形码）→ 搜索用户 → 加好友 → 私聊互发 → 非好友限流(3条) → 建群 → 加群 → 群聊 → 最近会话排序 → 心跳/断线重连/在线状态 → 多端登录 → 登出 → 踢人 → 图片上传

**测试策略**：本次不写自动化测试（用户决策）；但所有 usecase 依赖端口接口、handler/service 无全局状态，后续补测试无需调整结构。

## 6. 顺带修复的已知缺陷

| 缺陷 | 位置 | 处理 |
|---|---|---|
| 模糊搜索 phone 条件误用 `args[1]`（应为 `args[2]`） | dao/user_basic.go | 新 repo 层重写查询，消除位置参数 |
| `UpdateUserBasic` 绑定失败后缺 `return`，继续执行 | control/user_basic.go | 新 handler 统一 `return err` 模式 |
| WS 升级后 `clientManager.Register <- client` 发送两次 | websocket/init.go | 新 hub 单次注册 |
| `DelUsers` 遍历中删除且提前 return，多端场景失效 | websocket/client_manager.go | 新 hub 按 uid→clients 切片管理 |
| 前端引用未定义的 `Session` 导致 kick 分支报错 | utils/websocket.ts | 新 WSClient 统一 storage 封装 |
| `getSocket()` 内 `location.reload()` 副作用 | utils/websocket.ts | 移除，改为显式重连 |
| 明文密码策略：自建 salt+加密函数 | service/user_basic.go | 改 bcrypt |
| gorm `.Debug()` 生产打印 SQL | dao | sqlc + slog 分级日志 |

## 7. 风险与对策

- **一次性重写风险**：严格按阶段验证（build/vet/type-check），旧代码删除放在新实现验证通过之后；git 分支隔离整个重构。
- **ant-design-vue 3→4 破坏性变更**：涉及的组件 API 变化在 views 重写时逐个适配（本项目用量集中在表单/弹窗/头像/消息提示）。
- **sqlc 学习成本**：查询均为单表简单 SQL，sqlc 属最简场景；`sqlc generate` 产物纳入版本控制。
- **协议变更导致前后端必须同步部署**：个人项目、无第三方调用方，可接受。
