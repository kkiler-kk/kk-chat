# kk-chat 前端 Vue3 重构实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 front-ui-vue 重组为 types/api/ws/stores/composables/components/views 分层结构，WebSocket 改为类型化 WSClient 类 + Pinia store 单一事实来源，并落地 ESLint/Prettier/husky 工程化。

**Architecture:** HTTP 走全类型化 api 层；WS 为纯推送通道（WSClient 类，指数退避重连+心跳），收到事件直接调 store action；消息归属、会话排序等业务逻辑全部收编进 Pinia store，组件只做渲染与交互。

**Tech Stack:** Vue 3.5+、TypeScript 5.x（strict）、Vite 6、Pinia、vue-router 4、ant-design-vue 4.x、axios、vue3-emoji-picker、ESLint 9（flat config）+ Prettier + husky + lint-staged。

**Spec:** `docs/superpowers/specs/2026-09-12-kk-chat-refactor-design.md`（§3 协议、§4 前端架构）
**前置依赖:** 后端计划 `docs/superpowers/plans/2026-09-12-backend-clean-architecture.md` 已全部完成（本计划的 API/WS 契约以新后端为准）。

## Global Constraints

- **不写自动化测试**（用户决策）；每个任务验证 = `yarn type-check` 通过，涉及构建的任务加 `yarn build`。
- 工作分支：继续用 `refactor/clean-architecture`（与后端同分支，前后端需同步部署）。
- 旧代码（现 `src/utils/websocket.ts`、旧组件等）在 Task F10 之前保持可编译；新代码放入新目录结构，**新文件禁止 import 旧 utils/websocket.ts**。
- API 路径、字段名（snake_case）、错误码、WS 事件名严格对齐 spec §3 与后端计划 B13/B14。
- 禁止裸 `any`（ESLint 规则 `@typescript-eslint/no-explicit-any: error`；第三方无类型处用 `unknown` + 收窄）。
- JSON 时间字段为 RFC3339 字符串，前端统一 `new Date(s)` 解析。
- 每个任务完成即 git commit。

---

### Task F1: 依赖升级与工程化基座

**Files:**
- Modify: `front-ui-vue/package.json`
- Modify: `front-ui-vue/tsconfig.json`、`tsconfig.app.json`、`tsconfig.node.json`
- Create: `front-ui-vue/eslint.config.js`
- Create: `front-ui-vue/.prettierrc.json`
- Create: `front-ui-vue/.husky/pre-commit`
- Modify: `front-ui-vue/vite.config.ts`（确认 `@` 别名保留）
- Modify: `front-ui-vue/.env.development`、`.env.production`、`.env.testenvironment`（`VITE_API_URL`、`VITE_WEBSOCKET_URL` 指向新后端 `/api/v1` 与 `/api/v1/ws`）

**Interfaces:**
- Produces: 可用的 `yarn lint` / `yarn format` / `yarn type-check` / `yarn build` 命令；husky pre-commit 钩子（lint-staged）。

- [ ] **Step 1: 升级/增删依赖**

```bash
cd front-ui-vue
yarn add vue@^3.5 ant-design-vue@^4 @ant-design/icons-vue@^7 pinia@^2 vue-router@^4 axios@^1 vue3-emoji-picker@^1
yarn remove vuex emoji-mart-vue v-emoji-picker js-pinyin unplugin-vue-components
yarn add -D vite@^6 @vitejs/plugin-vue@^5 typescript@~5.6 vue-tsc@^2 \
  eslint@^9 eslint-plugin-vue@^9 vue-eslint-parser @typescript-eslint/eslint-plugin @typescript-eslint/parser \
  prettier eslint-config-prettier eslint-plugin-prettier husky lint-staged @vue/eslint-config-typescript @vue/eslint-config-prettier
```

package.json scripts 更新为：

```json
{
  "dev": "vite --mode development",
  "build": "run-p type-check \"build-only {@}\" --",
  "build-only": "vite build",
  "type-check": "vue-tsc --build --force",
  "lint": "eslint . --fix",
  "format": "prettier --write src/",
  "prepare": "husky"
}
```

并添加：

```json
"lint-staged": {
  "src/**/*.{ts,vue}": ["eslint --fix", "prettier --write"],
  "src/**/*.{css,scss,json,md}": ["prettier --write"]
}
```

- [ ] **Step 2: 写 eslint.config.js（flat config）**

```js
import pluginVue from 'eslint-plugin-vue'
import vueTsEslintConfig from '@vue/eslint-config-typescript'
import skipFormatting from '@vue/eslint-config-prettier/skip-formatting'

export default [
  { name: 'app/files-to-lint', files: ['**/*.{ts,mts,tsx,vue}'] },
  { name: 'app/files-to-ignore', ignores: ['**/dist/**', '**/node_modules/**'] },
  ...pluginVue.configs['flat/recommended'],
  ...vueTsEslintConfig(),
  skipFormatting,
  {
    rules: {
      '@typescript-eslint/no-explicit-any': 'error',
      'vue/multi-word-component-names': 'off',
    },
  },
]
```

- [ ] **Step 3: 写 .prettierrc.json**

```json
{
  "semi": false,
  "singleQuote": true,
  "printWidth": 100,
  "trailingComma": "all"
}
```

- [ ] **Step 4: tsconfig 严格化**

在 `tsconfig.app.json` 的 `compilerOptions` 中确保：`"strict": true`、`"noUnusedLocals": true`、`"noUnusedParameters": true`、`"noFallthroughCasesInSwitch": true`。

- [ ] **Step 5: husky 钩子**

```bash
yarn prepare   # 初始化 .husky/
echo "npx lint-staged" > .husky/pre-commit && chmod +x .husky/pre-commit
```

注意：husky 钩子路径——若 git 仓库根在上一级（kk-chat/），需 `git config core.hooksPath front-ui-vue/.husky` 或把钩子放仓库根 `.husky/` 并将 lint-staged 命令改为 `cd front-ui-vue && npx lint-staged`。**本仓库 git 根在 kk-chat/，采用后者**：钩子放仓库根 `.husky/pre-commit`，内容 `cd front-ui-vue && npx lint-staged`。

- [ ] **Step 6: .env 文件更新**

`.env.development`：

```
VITE_API_URL=http://localhost:8080
VITE_WEBSOCKET_URL=ws://localhost:8080/api/v1/ws
```

（`.env.production`/`.env.testenvironment` 按部署地址同构修改。注意：api 模块内路径以 `/api/v1/...` 开头，故 `VITE_API_URL` 只到 host。）

- [ ] **Step 7: 验证 + Commit**

Run: `cd front-ui-vue && yarn install && yarn lint`（旧代码可能报 no-explicit-any 等错误——**允许**：本步只要求 eslint 可运行并报出旧代码问题；旧代码在 F10 删除，新代码从 F2 起必须零告警）。`yarn type-check` 此时可能因旧代码+strict 报错，同样在 F10 前容忍，但需记录基线错误数。

```bash
git add front-ui-vue/package.json front-ui-vue/yarn.lock front-ui-vue/eslint.config.js front-ui-vue/.prettierrc.json front-ui-vue/tsconfig*.json front-ui-vue/.env* .husky
git commit -m "build(web): 依赖升级（Vue3.5/antd4/Vite6）与 ESLint+Prettier+husky 工程化基座"
```

---

### Task F2: 协议类型与基础工具

**Files:**
- Create: `front-ui-vue/src/types/api.ts`
- Create: `front-ui-vue/src/types/ws.ts`
- Create: `front-ui-vue/src/types/errorcode.ts`
- Create: `front-ui-vue/src/utils/storage.ts`（重写：类型化）
- Create: `front-ui-vue/src/utils/format.ts`（重写：从 formatTime.ts 精简）

**Interfaces:**
- Produces（后续所有任务的类型契约，与后端 port DTO 逐字段对齐）:
  - `types/api.ts`: `ApiResponse<T>`、`UserInfo`、`LoginResult`、`UserDetail`、`UserSearchItem`、`FriendItem`、`GroupItem`、`OutgoingMessage`、`RecentConversation`、`CaptchaResult`、`UploadResult`
  - `types/ws.ts`: `WsEnvelope<T>`、`WsServerEvents`（事件名→payload 映射）
  - `types/errorcode.ts`: `ErrCode` 常量对象（与后端 apperror 码一致）
  - `utils/storage.ts`: `Local.get<T>(key)/set<T>(key,v)/remove(key)/clear()`、`Session.*` 同构（localStorage/sessionStorage JSON 封装，键前缀 `kk:`）
  - `utils/format.ts`: `formatPast(d: Date | string): string`（刚刚/x分钟前/x小时前/x天前/日期）、`formatTime(d, fmt?)`

- [ ] **Step 1: 写 types/api.ts**

```ts
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
  request_id: string
}

export interface UserInfo {
  id: number
  identity: string
  name: string
  avatar: string
  email: string
  phone: string
  signature: string
  birth_date?: string
  created_at: string
}

export interface LoginResult {
  token: string
  user: UserInfo
}

export interface UserDetail extends UserInfo {
  is_friend: boolean
  is_self: boolean
}

export interface UserSearchItem {
  id: number
  identity: string
  name: string
  avatar: string
  is_friend: boolean
}

export interface FriendItem {
  id: number
  name: string
  avatar: string
  online: boolean
}

export interface GroupItem {
  id: number
  name: string
  avatar: string
  owner_id: number
  member_count: number
}

export interface OutgoingMessage {
  id: string
  conversation_id: string
  sender_id: number
  sender_name: string
  sender_avatar: string
  content: string
  content_type: 'text' | 'image'
  created_at: string
}

export interface RecentConversation {
  conversation_id: string
  type: 'private' | 'group'
  peer_id: number
  peer_name: string
  peer_avatar: string
  online: boolean
  last_message_content: string
  last_sender_name: string
  last_time: string
}

export interface CaptchaResult {
  captcha_id: string
  image: string
}

export interface UploadResult {
  url: string
}
```

- [ ] **Step 2: 写 types/ws.ts 与 types/errorcode.ts**

```ts
import type { OutgoingMessage, RecentConversation } from './api'

export interface WsEnvelope<T = unknown> {
  event: string
  data: T
  ts: number
}

/** 服务端 → 客户端事件与 payload 的映射（spec §3.4） */
export interface WsServerEvents {
  pong: undefined
  'chat.message': OutgoingMessage
  'chat.recent_updated': RecentConversation
  'presence.changed': { user_id: number; online: boolean }
  'group.updated': { group_id: number }
  'system.kick': { reason: string }
}

export type WsServerEvent = keyof WsServerEvents
```

```ts
/** 与后端 apperror 业务错误码一致（spec §3.2） */
export const ErrCode = {
  OK: 0,
  InvalidParam: 10001,
  Internal: 10002,
  CaptchaInvalid: 20001,
  BadCredentials: 20002,
  TokenInvalid: 20003,
  UserBanned: 20004,
  EmailTaken: 30001,
  IdentityTaken: 30002,
  UserNotFound: 30003,
  EmailCodeInvalid: 30004,
  AlreadyFriend: 40001,
  NotFriendLimit: 40002,
  GroupNotFound: 50001,
  NotGroupMember: 50002,
  ConversationInvalid: 60001,
} as const

export type ErrCodeValue = (typeof ErrCode)[keyof typeof ErrCode]
```

- [ ] **Step 3: 重写 utils/storage.ts（类型化，修复旧版 Session 未导出/any 问题）**

```ts
const PREFIX = 'kk:'

function createStorage(engine: Storage) {
  return {
    get<T>(key: string): T | null {
      const raw = engine.getItem(PREFIX + key)
      if (raw === null) return null
      try {
        return JSON.parse(raw) as T
      } catch {
        return null
      }
    },
    set<T>(key: string, value: T): void {
      engine.setItem(PREFIX + key, JSON.stringify(value))
    },
    remove(key: string): void {
      engine.removeItem(PREFIX + key)
    },
    clear(): void {
      // 只清理本应用前缀键，避免误伤同域其他数据
      const keys: string[] = []
      for (let i = 0; i < engine.length; i++) {
        const k = engine.key(i)
        if (k && k.startsWith(PREFIX)) keys.push(k)
      }
      keys.forEach((k) => engine.removeItem(k))
    },
  }
}

export const Local = createStorage(localStorage)
export const Session = createStorage(sessionStorage)

export const StorageKeys = {
  token: 'token',
  userInfo: 'userInfo',
} as const
```

- [ ] **Step 4: 写 utils/format.ts**

从旧 `formatTime.ts`（137 行）中仅保留并迁移 `formatPast` 与通用 `formatTime`，签名：

```ts
export function formatPast(input: Date | string): string
// <60s "刚刚"；<60min "x分钟前"；<24h "x小时前"；<7d "x天前"；今年内 "MM-DD HH:mm"；跨年 "YYYY-MM-DD"
export function formatTime(input: Date | string, fmt = 'YYYY-MM-DD HH:mm:ss'): string
```

（实现用原生 Date，不引入 dayjs；旧文件中其余未使用函数一律不迁移。）

- [ ] **Step 5: 验证 + Commit**

Run: `yarn type-check`——新文件必须零错误（旧代码错误忽略，记录数量）。`npx eslint src/types src/utils/storage.ts src/utils/format.ts` 零告警。

```bash
git add front-ui-vue/src/types front-ui-vue/src/utils/storage.ts front-ui-vue/src/utils/format.ts
git commit -m "feat(web): 协议类型（api/ws/错误码）与类型化 storage/format 工具"
```

---

### Task F3: API 层（axios 客户端 + 六个模块）

**Files:**
- Create: `front-ui-vue/src/api/client.ts`
- Create: `front-ui-vue/src/api/auth.ts`
- Create: `front-ui-vue/src/api/user.ts`
- Create: `front-ui-vue/src/api/friend.ts`
- Create: `front-ui-vue/src/api/group.ts`
- Create: `front-ui-vue/src/api/chat.ts`
- Create: `front-ui-vue/src/api/file.ts`

**Interfaces:**
- Produces:
  - `client.ApiError` 类：`{ code: number; message: string; requestId: string; status: number }`
  - `client.setUnauthorizedHandler(fn: () => void): void`
  - 各模块函数（全部返回 `Promise<T>`，T 为 types/api.ts 类型）：
    - `auth`: `register(p: RegisterParams)`、`login(p: LoginParams): Promise<LoginResult>`、`logout()`、`getCaptcha(): Promise<CaptchaResult>`、`sendEmailCode(email: string)`
    - `user`: `getMe(): Promise<UserInfo>`、`updateMe(p: UpdateMeParams)`、`getDetail(id: number): Promise<UserDetail>`、`search(keyword: string): Promise<UserSearchItem[]>`
    - `friend`: `add(userId: number)`、`list(): Promise<FriendItem[]>`
    - `group`: `create(p: CreateGroupParams): Promise<GroupItem>`、`join(id: number)`、`list(): Promise<GroupItem[]>`、`search(keyword: string): Promise<GroupItem[]>`
    - `chat`: `sendMessage(p: SendMessageParams): Promise<OutgoingMessage>`、`conversations(): Promise<RecentConversation[]>`、`history(convId: string, cursor?: string, limit?: number): Promise<OutgoingMessage[]>`
    - `file`: `uploadImage(f: File): Promise<UploadResult>`
  - 参数类型（同文件或 types/api.ts 内）：`RegisterParams{identity,name,password,email,email_code}`、`LoginParams{account,password,captcha_id,captcha_answer}`、`UpdateMeParams{name?,phone?,email?,email_code?,avatar?,signature?,birth_date?}`、`CreateGroupParams{name,member_ids?}`、`SendMessageParams{conversation_id,content,content_type}`

- [ ] **Step 1: 写 api/client.ts**

```ts
import axios, { type AxiosResponse } from 'axios'
import { message } from 'ant-design-vue'
import { Local, StorageKeys } from '@/utils/storage'
import type { ApiResponse } from '@/types/api'
import { ErrCode } from '@/types/errorcode'

export class ApiError extends Error {
  code: number
  requestId: string
  status: number
  constructor(code: number, msg: string, requestId = '', status = 0) {
    super(msg)
    this.code = code
    this.requestId = requestId
    this.status = status
  }
}

let unauthorizedHandler: () => void = () => {}
export function setUnauthorizedHandler(fn: () => void): void {
  unauthorizedHandler = fn
}

const http = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' },
})

http.interceptors.request.use((config) => {
  const token = Local.get<string>(StorageKeys.token)
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

http.interceptors.response.use(
  (resp: AxiosResponse<ApiResponse<unknown>>) => {
    const body = resp.data
    if (body.code === ErrCode.OK) return resp
    // 业务错误：统一 toast + 抛出 ApiError（401 类不 toast，交给 unauthorized 流程）
    const err = new ApiError(body.code, body.message, body.request_id, resp.status)
    if (err.code !== ErrCode.TokenInvalid) message.error(body.message)
    return Promise.reject(err)
  },
  (error) => {
    if (axios.isAxiosError(error) && error.response) {
      const body = error.response.data as ApiResponse<unknown> | undefined
      const err = new ApiError(
        body?.code ?? error.response.status * 100,
        body?.message ?? error.message,
        body?.request_id ?? '',
        error.response.status,
      )
      if (err.status === 401 || err.code === ErrCode.TokenInvalid) {
        unauthorizedHandler()
      } else {
        message.error(err.message)
      }
      return Promise.reject(err)
    }
    // 网络层错误（超时/断网）
    const msg = axios.isAxiosError(error) && error.code === 'ECONNABORTED' ? '网络超时' : '网络连接错误'
    message.error(msg)
    return Promise.reject(new ApiError(ErrCode.Internal, msg, '', 0))
  },
)

/** request 泛型封装：直接 resolve data 字段 */
export async function request<T>(config: Parameters<typeof http.request>[0]): Promise<T> {
  const resp = await http.request<ApiResponse<T>>(config)
  return resp.data.data
}
```

- [ ] **Step 2: 写六个 API 模块（示例 auth.ts，其余同构）**

```ts
import { request } from './client'
import type { CaptchaResult, LoginResult } from '@/types/api'

export interface RegisterParams {
  identity: string
  name: string
  password: string
  email: string
  email_code: string
}
export interface LoginParams {
  account: string
  password: string
  captcha_id: string
  captcha_answer: string
}

export const register = (p: RegisterParams) =>
  request<null>({ url: '/api/v1/auth/register', method: 'POST', data: p })
export const login = (p: LoginParams) =>
  request<LoginResult>({ url: '/api/v1/auth/login', method: 'POST', data: p })
export const logout = () => request<null>({ url: '/api/v1/auth/logout', method: 'POST' })
export const getCaptcha = () =>
  request<CaptchaResult>({ url: '/api/v1/auth/captcha', method: 'POST' })
export const sendEmailCode = (email: string) =>
  request<null>({ url: '/api/v1/auth/email-code', method: 'POST', data: { email } })
```

其余模块 URL（与后端计划 B14 路由一一对应）：
- `user.ts`: GET `/api/v1/users/me`；PATCH `/api/v1/users/me`；GET `/api/v1/users/${id}`；GET `/api/v1/users` params `{search}`
- `friend.ts`: POST `/api/v1/friends` data `{user_id}`；GET `/api/v1/friends`
- `group.ts`: POST `/api/v1/groups`；POST `/api/v1/groups/${id}/members`；GET `/api/v1/groups`（列表）；GET `/api/v1/groups` params `{search}`（搜索——后端 List 内部按 search 参数分流）
- `chat.ts`: POST `/api/v1/messages`；GET `/api/v1/conversations`；GET `/api/v1/conversations/${encodeURIComponent(convId)}/messages` params `{cursor?, limit?}`
- `file.ts`: POST `/api/v1/files/images`，`FormData`（字段名 `file`），headers `{'Content-Type':'multipart/form-data'}`

- [ ] **Step 3: 验证 + Commit**

Run: `yarn type-check` 新文件零错误；`npx eslint src/api` 零告警。

```bash
git add front-ui-vue/src/api
git commit -m "feat(web): 全类型化 API 层（axios 客户端 + 六模块）"
```

---

### Task F4: WSClient（类型化事件、指数退避重连、心跳）

**Files:**
- Create: `front-ui-vue/src/ws/WSClient.ts`

**Interfaces:**
- Consumes: `types/ws.ts` 的 `WsServerEvents/WsEnvelope`；`utils/storage`
- Produces:
  - `class WSClient`：`connect(token: string): void`、`close(): void`、`on<E extends WsServerEvent>(event: E, handler: (data: WsServerEvents[E]) => void): () => void`（返回解绑函数）、`get connected(): boolean`、`setKickHandler(fn: (reason: string) => void): void`
  - `export const wsClient: WSClient`（应用级单例）

- [ ] **Step 1: 写 WSClient.ts**

```ts
import type { WsEnvelope, WsServerEvent, WsServerEvents } from '@/types/ws'

type Handler<E extends WsServerEvent> = (data: WsServerEvents[E]) => void

const HEARTBEAT_INTERVAL = 60_000 // 与后端 presence.heartbeat_seconds 对齐
const RECONNECT_BASE = 1_000
const RECONNECT_MAX = 30_000

export class WSClient {
  private socket: WebSocket | null = null
  private handlers = new Map<string, Set<Handler<WsServerEvent>>>()
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null
  private reconnectAttempts = 0
  private manualClose = false
  private token = ''
  private kickHandler: (reason: string) => void = () => {}

  get connected(): boolean {
    return this.socket?.readyState === WebSocket.OPEN
  }

  setKickHandler(fn: (reason: string) => void): void {
    this.kickHandler = fn
  }

  connect(token: string): void {
    this.token = token
    this.manualClose = false
    this.open()
  }

  close(): void {
    this.manualClose = true
    this.clearTimers()
    this.socket?.close()
    this.socket = null
  }

  on<E extends WsServerEvent>(event: E, handler: Handler<E>): () => void {
    let set = this.handlers.get(event)
    if (!set) {
      set = new Set()
      this.handlers.set(event, set)
    }
    set.add(handler as Handler<WsServerEvent>)
    return () => set!.delete(handler as Handler<WsServerEvent>)
  }

  private open(): void {
    const url = `${import.meta.env.VITE_WEBSOCKET_URL}?token=${encodeURIComponent(this.token)}`
    const socket = new WebSocket(url)
    this.socket = socket

    socket.onopen = () => {
      this.reconnectAttempts = 0
      this.startHeartbeat()
    }
    socket.onmessage = (ev: MessageEvent<string>) => this.dispatch(ev.data)
    socket.onerror = () => {
      /* 交给 onclose 统一重连 */
    }
    socket.onclose = () => {
      this.stopHeartbeat()
      this.socket = null
      if (!this.manualClose) this.scheduleReconnect()
    }
  }

  private dispatch(raw: string): void {
    let env: WsEnvelope<unknown>
    try {
      env = JSON.parse(raw) as WsEnvelope<unknown>
    } catch {
      return // 非 JSON 消息直接忽略
    }
    if (env.event === 'system.kick') {
      const data = env.data as { reason: string }
      this.manualClose = true
      this.kickHandler(data?.reason ?? '')
      this.close()
      return
    }
    this.handlers.get(env.event)?.forEach((h) => {
      try {
        h(env.data as never)
      } catch (e) {
        console.error(`[ws] handler for "${env.event}" failed`, e)
      }
    })
  }

  private startHeartbeat(): void {
    this.stopHeartbeat()
    this.heartbeatTimer = setInterval(() => {
      if (this.connected) {
        this.socket?.send(JSON.stringify({ event: 'ping' }))
      }
    }, HEARTBEAT_INTERVAL)
  }

  private stopHeartbeat(): void {
    if (this.heartbeatTimer) clearInterval(this.heartbeatTimer)
    this.heartbeatTimer = null
  }

  private scheduleReconnect(): void {
    if (this.reconnectTimer) return
    const delay = Math.min(RECONNECT_BASE * 2 ** this.reconnectAttempts, RECONNECT_MAX)
    this.reconnectAttempts++
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null
      if (!this.manualClose && this.token) this.open()
    }, delay)
  }

  private clearTimers(): void {
    this.stopHeartbeat()
    if (this.reconnectTimer) clearTimeout(this.reconnectTimer)
    this.reconnectTimer = null
  }
}

export const wsClient = new WSClient()
```

- [ ] **Step 2: 验证 + Commit**

Run: `yarn type-check` 零新增错误；`npx eslint src/ws` 零告警。

```bash
git add front-ui-vue/src/ws
git commit -m "feat(web): WSClient——类型化事件订阅/指数退避重连/心跳"
```

---

### Task F5: Pinia stores（auth / presence / chat / contacts）

**Files:**
- Create: `front-ui-vue/src/stores/index.ts`（沿用现有 setupStore，追加导出）
- Create: `front-ui-vue/src/stores/auth.ts`
- Create: `front-ui-vue/src/stores/presence.ts`
- Create: `front-ui-vue/src/stores/chat.ts`
- Create: `front-ui-vue/src/stores/contacts.ts`
- Create: `front-ui-vue/src/ws/bindings.ts`（WS 事件 → store action 的接线）

**Interfaces:**
- Consumes: F3 api 模块、F4 wsClient、F2 types/storage
- Produces:
  - `useAuthStore()`：state `{token: string|null, user: UserInfo|null}`；getters `isLoggedIn`；actions `login(p)/register(p)/logout()/restore()/updateProfile(p)/setUnauthorized()`
  - `usePresenceStore()`：state `{onlineIds: Set<number>}`；actions `set(uid, online)`、`isOnline(uid): boolean`、`applyFriends(items: FriendItem[])`
  - `useChatStore()`：state `{conversations: RecentConversation[], messagesByConv: Record<string, OutgoingMessage[]>, activeConvId: string|null, historyCursor: Record<string, string>, loading: boolean}`；actions `loadConversations()/openConversation(convId)/loadMoreHistory(convId)/send(content, contentType)/handleIncoming(msg)/handleRecentUpdated(conv)/setActive(convId)`
  - `useContactsStore()`：state `{friends: FriendItem[], groups: GroupItem[], userSearch: UserSearchItem[], groupSearch: GroupItem[]}`；actions `loadFriends()/loadGroups()/addFriend(uid)/createGroup(name, memberIds)/joinGroup(gid)/searchUsers(kw)/searchGroups(kw)/handlePresenceChanged(uid, online)`
  - `bindings.ts`: `setupWsBindings(): () => void`（注册全部 WS 订阅，返回统一解绑函数）

- [ ] **Step 1: 写 stores/auth.ts**

```ts
import { defineStore } from 'pinia'
import * as authApi from '@/api/auth'
import * as userApi from '@/api/user'
import { ApiError, setUnauthorizedHandler } from '@/api/client'
import { Local, StorageKeys } from '@/utils/storage'
import { wsClient } from '@/ws/WSClient'
import type { LoginResult, UserInfo } from '@/types/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: Local.get<string>(StorageKeys.token),
    user: Local.get<UserInfo>(StorageKeys.userInfo),
  }),
  getters: {
    isLoggedIn: (s): boolean => !!s.token && !!s.user,
  },
  actions: {
    async login(account: string, password: string, captchaId: string, captchaAnswer: string) {
      const result: LoginResult = await authApi.login({
        account, password, captcha_id: captchaId, captcha_answer: captchaAnswer,
      })
      this.applySession(result.token, result.user)
    },
    async register(p: Parameters<typeof authApi.register>[0]) {
      await authApi.register(p)
    },
    applySession(token: string, user: UserInfo) {
      this.token = token
      this.user = user
      Local.set(StorageKeys.token, token)
      Local.set(StorageKeys.userInfo, user)
      wsClient.setKickHandler(() => this.setUnauthorized())
      wsClient.connect(token)
    },
    async logout() {
      try { await authApi.logout() } catch { /* 忽略登出接口错误，本地照清 */ }
      this.clearSession()
    },
    /** 401/token失效/被踢：清理本地并跳转登录（由 router 守卫兜底） */
    setUnauthorized() {
      this.clearSession()
      window.location.href = '/login'
    },
    clearSession() {
      wsClient.close()
      this.token = null
      this.user = null
      Local.remove(StorageKeys.token)
      Local.remove(StorageKeys.userInfo)
    },
    /** 刷新页面后恢复：校验本地 token 是否仍有效 */
    async restore(): Promise<boolean> {
      if (!this.isLoggedIn) return false
      try {
        this.user = await userApi.getMe()
        Local.set(StorageKeys.userInfo, this.user)
        setUnauthorizedHandler(() => this.setUnauthorized())
        wsClient.setKickHandler(() => this.setUnauthorized())
        wsClient.connect(this.token!)
        return true
      } catch (e) {
        if (e instanceof ApiError) this.clearSession()
        return false
      }
    },
    async updateProfile(p: Parameters<typeof userApi.updateMe>[0]) {
      await userApi.updateMe(p)
      this.user = await userApi.getMe()
      Local.set(StorageKeys.userInfo, this.user)
    },
  },
})
```

- [ ] **Step 2: 写 stores/presence.ts**

```ts
import { defineStore } from 'pinia'
import type { FriendItem } from '@/types/api'

export const usePresenceStore = defineStore('presence', {
  state: () => ({
    onlineIds: new Set<number>(),
  }),
  actions: {
    set(uid: number, online: boolean) {
      if (online) this.onlineIds.add(uid)
      else this.onlineIds.delete(uid)
      // 触发响应式：Set 变更需重新赋值（Pinia 对 Set 的 add/delete 可响应，但为兼容组件 computed，保留此写法亦可）
    },
    isOnline(uid: number): boolean {
      return this.onlineIds.has(uid)
    },
    applyFriends(items: FriendItem[]) {
      this.onlineIds = new Set(items.filter((f) => f.online).map((f) => f.id))
    },
  },
})
```

- [ ] **Step 3: 写 stores/chat.ts（消息归属逻辑的唯一所在地）**

```ts
import { defineStore } from 'pinia'
import * as chatApi from '@/api/chat'
import type { OutgoingMessage, RecentConversation } from '@/types/api'

export const useChatStore = defineStore('chat', {
  state: () => ({
    conversations: [] as RecentConversation[],
    messagesByConv: {} as Record<string, OutgoingMessage[]>,
    activeConvId: null as string | null,
    historyCursor: {} as Record<string, string>, // convId → 最旧一条的 created_at
    historyDone: {} as Record<string, boolean>,
    loading: false,
  }),
  getters: {
    activeMessages: (s): OutgoingMessage[] =>
      s.activeConvId ? (s.messagesByConv[s.activeConvId] ?? []) : [],
    activeConversation: (s): RecentConversation | undefined =>
      s.conversations.find((c) => c.conversation_id === s.activeConvId),
  },
  actions: {
    async loadConversations() {
      this.conversations = await chatApi.conversations()
    },
    async openConversation(convId: string) {
      this.activeConvId = convId
      if (this.messagesByConv[convId]) return // 已加载过
      this.loading = true
      try {
        const msgs = await chatApi.history(convId, undefined, 50)
        this.messagesByConv[convId] = msgs
        if (msgs.length > 0) this.historyCursor[convId] = msgs[0].created_at
        this.historyDone[convId] = msgs.length < 50
      } finally {
        this.loading = false
      }
    },
    async loadMoreHistory(convId: string) {
      const cursor = this.historyCursor[convId]
      if (!cursor || this.historyDone[convId]) return
      const msgs = await chatApi.history(convId, cursor, 50)
      this.messagesByConv[convId] = [...msgs, ...(this.messagesByConv[convId] ?? [])]
      if (msgs.length > 0) this.historyCursor[convId] = msgs[0].created_at
      this.historyDone[convId] = msgs.length < 50
    },
    async send(content: string, contentType: 'text' | 'image' = 'text') {
      if (!this.activeConvId) return
      const msg = await chatApi.sendMessage({
        conversation_id: this.activeConvId, content, content_type: contentType,
      })
      this.handleIncoming(msg) // 自己发的也走统一入口（后端不回推给自己）
    },
    /** WS chat.message → 唯一入口：归属判断只在这里做一次 */
    handleIncoming(msg: OutgoingMessage) {
      const list = this.messagesByConv[msg.conversation_id]
      if (list && !list.some((m) => m.id === msg.id)) list.push(msg)
      else if (!list) this.messagesByConv[msg.conversation_id] = [msg]
    },
    /** WS chat.recent_updated → 更新/置顶会话 */
    handleRecentUpdated(conv: RecentConversation) {
      const idx = this.conversations.findIndex(
        (c) => c.conversation_id === conv.conversation_id,
      )
      if (idx >= 0) this.conversations.splice(idx, 1, conv)
      else this.conversations.unshift(conv)
      // 后端已按 last_time 推送最新值，重排保证置顶
      this.conversations.sort(
        (a, b) => new Date(b.last_time).getTime() - new Date(a.last_time).getTime(),
      )
    },
    /** 退出登录时清空（auth.clearSession 调用） */
    $resetAll() {
      this.conversations = []
      this.messagesByConv = {}
      this.activeConvId = null
      this.historyCursor = {}
      this.historyDone = {}
    },
  },
})
```

注意：`auth.ts clearSession()` 中补充调用 `useChatStore().$resetAll()` 与 `useContactsStore().$reset()`（Pinia setup 之外调用 store 需确保 pinia 已安装——在 action 内调用是安全的）。

- [ ] **Step 4: 写 stores/contacts.ts**

state/actions 按 Interfaces 所列实现：
- `loadFriends()`：`friend.list()` → `this.friends = items`；`usePresenceStore().applyFriends(items)`
- `loadGroups()`：`group.list()`
- `addFriend(uid)`：`friend.add(uid)`；成功后 `loadFriends()`
- `createGroup(name, memberIds)`：`group.create({name, member_ids})`；成功后 `loadGroups()`
- `joinGroup(gid)`：`group.join(gid)`；成功后 `loadGroups()`
- `searchUsers(kw)` / `searchGroups(kw)`：写 `userSearch` / `groupSearch`
- `handlePresenceChanged(uid, online)`：`usePresenceStore().set(uid, online)`；同步更新 `this.friends` 中对应项的 `online` 字段（`find` 后直接赋值，保证列表 UI 响应）
- `$reset()`：清空四个数组

- [ ] **Step 5: 写 ws/bindings.ts（接线：WS 事件 → store action，全部订阅集中于此）**

```ts
import { wsClient } from './WSClient'
import { useChatStore } from '@/stores/chat'
import { useContactsStore } from '@/stores/contacts'
import { useAuthStore } from '@/stores/auth'

/** 在登录成功/页面恢复后调用一次；返回解绑函数（登出时调用）。 */
export function setupWsBindings(): () => void {
  const chat = useChatStore()
  const contacts = useContactsStore()
  const offs = [
    wsClient.on('chat.message', (msg) => chat.handleIncoming(msg)),
    wsClient.on('chat.recent_updated', (conv) => chat.handleRecentUpdated(conv)),
    wsClient.on('presence.changed', (d) => contacts.handlePresenceChanged(d.user_id, d.online)),
    wsClient.on('group.updated', () => contacts.loadGroups()),
  ]
  return () => offs.forEach((off) => off())
}
```

并在 `auth.ts` 的 `applySession` 与 `restore` 成功路径中调用 `setupWsBindings()`（保存返回的解绑函数到模块级变量，`clearSession` 时调用）。

- [ ] **Step 6: 验证 + Commit**

Run: `yarn type-check` 新文件零错误；`npx eslint src/stores src/ws` 零告警。

```bash
git add front-ui-vue/src/stores front-ui-vue/src/ws/bindings.ts
git commit -m "feat(web): Pinia stores（auth/presence/chat/contacts）与 WS 事件接线"
```

---

### Task F6: 路由、入口与 App 壳

**Files:**
- Rewrite: `front-ui-vue/src/router/index.ts`
- Rewrite: `front-ui-vue/src/main.ts`
- Rewrite: `front-ui-vue/src/App.vue`
- Modify: `front-ui-vue/src/stores/index.ts`

**Interfaces:**
- Produces: 路由表 `['/login', '/home'（守卫）, '/:pathMatch(.*)*' → NotFound]`；`main.ts` 装配顺序：pinia → router → antd → emoji-picker 样式 → mount。

- [ ] **Step 1: 写 router/index.ts（含异步守卫）**

```ts
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', redirect: '/home' },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/login/index.vue'),
      meta: { public: true },
    },
    {
      path: '/home',
      name: 'home',
      component: () => import('@/views/home/HomeView.vue'),
    },
    {
      path: '/403', name: 'forbidden', meta: { public: true },
      component: () => import('@/views/ErrorStatus/Forbidden/index.vue'),
    },
    {
      path: '/:pathMatch(.*)*', name: 'not-found', meta: { public: true },
      component: () => import('@/views/ErrorStatus/NotFound/index.vue'),
    },
  ],
})

let restorePromise: Promise<boolean> | null = null

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (to.meta.public) {
    // 已登录访问登录页 → 回 home
    if (to.name === 'login' && auth.isLoggedIn) return { name: 'home' }
    return true
  }
  if (auth.isLoggedIn) return true
  // 无内存会话时尝试用本地 token 恢复（只并发执行一次）
  restorePromise ??= auth.restore()
  const ok = await restorePromise
  if (ok) return true
  restorePromise = null
  return { name: 'login', query: { redirect: to.fullPath } }
})

export default router
```

- [ ] **Step 2: 写 main.ts 与 App.vue**

```ts
import { createApp } from 'vue'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'
import 'vue3-emoji-picker/css'
import App from './App.vue'
import router from './router'
import { setupStore } from './stores'

const app = createApp(App)
setupStore(app)
app.use(router)
app.use(Antd)
app.mount('#app')
```

App.vue 仅保留 `<router-view />` 与全局样式入口（删除旧版多余逻辑）。

- [ ] **Step 3: 验证 + Commit**

Run: `yarn type-check`（router/main 零错误；views 尚未重写会报错，属预期——记录基线）。

```bash
git add front-ui-vue/src/router front-ui-vue/src/main.ts front-ui-vue/src/App.vue front-ui-vue/src/stores/index.ts
git commit -m "feat(web): 路由（异步鉴权守卫）与应用入口装配"
```

---

### Task F7: 登录与注册页

**Files:**
- Rewrite: `front-ui-vue/src/views/login/index.vue`
- Rewrite: `front-ui-vue/src/views/login/component/login.vue`
- Rewrite: `front-ui-vue/src/views/login/component/register.vue`

**Interfaces:**
- Consumes: `useAuthStore`、`authApi.getCaptcha/sendEmailCode`、`ErrCode`
- Produces: 可用的 /login 页面（tab 切换 登录/注册）

- [ ] **Step 1: 写 login/index.vue**

布局壳：居中卡片 + `a-tabs`（登录/注册）切换两个子组件；保留旧版视觉风格（背景/品牌名 kk-chat）。

- [ ] **Step 2: 写 component/login.vue**

行为要点：
1. `onMounted` → `authApi.getCaptcha()` → 显示 `image`（base64，点击刷新），保存 `captcha_id`
2. 表单（a-form + rules）：account（identity 或邮箱）、password、captcha_answer
3. 提交 → `auth.login(...)` → 成功后 `router.push(route.query.redirect ?? '/home')`
4. 登录失败且 `e instanceof ApiError && (e.code === ErrCode.CaptchaInvalid || e.code === ErrCode.BadCredentials)` → 自动刷新图形验证码
5. 组件内不做 message.error（client 拦截器已统一 toast）

- [ ] **Step 3: 写 component/register.vue**

行为要点：
1. 表单：identity、name、password、confirmPassword（前端校验一致）、email、email_code
2. "发送验证码"按钮 → `sendEmailCode(email)` → 60s 倒计时禁用（`setInterval`，onUnmounted 清理）
3. 提交 → `auth.register({...})` → `message.success('注册成功，请登录')` → 切回登录 tab（emit 事件给 index.vue）

- [ ] **Step 4: 验证 + Commit**

Run: `yarn type-check`；`npx eslint src/views/login`。可 `yarn dev` 手工打开 /login 检查渲染（后端未起时接口报错属预期）。

```bash
git add front-ui-vue/src/views/login
git commit -m "feat(web): 登录/注册页重写（图形码、邮箱验证码倒计时、错误联动刷新）"
```

---

### Task F8: Home 骨架、会话列表与聊天窗口

**Files:**
- Rewrite: `front-ui-vue/src/views/home/HomeView.vue`
- Create: `front-ui-vue/src/views/home/components/ConversationList.vue`（替代旧 viewsSoder/recentMessage.vue）
- Create: `front-ui-vue/src/views/home/components/ChatWindow.vue`（替代旧 content.vue）
- Create: `front-ui-vue/src/components/MessageBubble.vue`
- Create: `front-ui-vue/src/views/home/components/SideNav.vue`（替代旧 sider.vue：头像+导航图标列）
- Create: `front-ui-vue/src/views/home/components/TheHeader.vue`（替代旧 header.vue）

**Interfaces:**
- Consumes: `useChatStore/useAuthStore/usePresenceStore`、`formatPast`
- Produces: /home 三栏布局：SideNav（左窄栏）｜中栏（ConversationList / 联系人面板 由 nav 切换）｜右栏 ChatWindow

- [ ] **Step 1: 写 HomeView.vue**

要点：
1. `onMounted`：`chat.loadConversations()`、`contacts.loadFriends()`、`contacts.loadGroups()`（并发 `Promise.all`）
2. 布局：flex 三栏；中栏由 `activePanel: 'chats' | 'contacts' | 'groups' | 'settings'` 局部状态切换（SideNav emit）
3. `onUnmounted` 无需解绑 WS（bindings 生命周期跟随登录会话，在 auth store 管理）

- [ ] **Step 2: 写 ConversationList.vue**

要点：
1. 数据源：`chatStore.conversations`（computed）；点击项 → `chatStore.openConversation(c.conversation_id)`；高亮 `activeConvId`
2. 每项渲染：头像（a-avatar，group 用群头像）、peer_name、`last_message_content`（text 直接显示，image 显示 "[图片]"）、`formatPast(last_time)`、在线小绿点（`type==='private' && presence.isOnline(peer_id)` 或直接用 `c.online`）
3. 搜索框（可选输入）本地过滤 peer_name
4. 空态：a-empty

- [ ] **Step 3: 写 MessageBubble.vue**

props：`message: OutgoingMessage`、`self: boolean`。渲染：头像 + 昵称 + `formatPast(created_at)` + 内容气泡（`content_type==='image'` 渲染 `<img>` 最大宽 240px 可点击新窗口打开；text 渲染文本）；self 右对齐（旧版 `.self` 样式迁移进来，气泡样式全部内聚于本组件）。

- [ ] **Step 4: 写 ChatWindow.vue**

要点：
1. 顶部：`chatStore.activeConversation` 的 peer_name + 在线状态；无 activeConvId 时显示空态占位
2. 消息区：`v-for chatStore.activeMessages` → MessageBubble（`self = msg.sender_id === auth.user?.id`）；容器 ref 滚动到底（`watch(activeMessages.length)` + `nextTick` 后 `scrollTop = scrollHeight`；**函数命名避开 Vue 的 nextTick**，旧版同名遮蔽问题不再复现）
3. 上拉加载更早消息：消息区 `@scroll` 到顶部（scrollTop < 40）→ `chatStore.loadMoreHistory(activeConvId)`，加载后恢复滚动位置（记录加载前 scrollHeight 差值）
4. 输入区：a-textarea（Enter 发送、Shift+Enter 换行）+ vue3-emoji-picker（插入 emoji 到光标处）+ 图片上传按钮（`file.uploadImage` → 成功后 `chatStore.send(url, 'image')`）+ 发送按钮（`chatStore.send(text)`，空内容禁用）
5. 非好友限流提示：`send` 抛 `ApiError.code === ErrCode.NotFriendLimit` 时展示 a-alert 于输入区上方（client 已 toast，此处补充常驻提示）

- [ ] **Step 5: 写 SideNav.vue 与 TheHeader.vue**

SideNav：顶部当前用户头像（点击 → settings 面板）；图标列：消息（MessageOutlined）/联系人（TeamOutlined）/群组（UsergroupAddOutlined）/设置（SettingOutlined）；`v-model:active` 或 emit `update:activePanel`。TheHeader：中栏顶部标题条（当前面板名 + 右侧操作按钮插槽，如"发起群聊"）。

- [ ] **Step 6: 验证 + Commit**

Run: `yarn type-check`（新文件零错误）；`yarn dev` 手工冒烟：登录后会话列表渲染、点击会话拉历史、发消息上屏、WS 推送进对应会话。

```bash
git add front-ui-vue/src/views/home front-ui-vue/src/components/MessageBubble.vue
git commit -m "feat(web): Home 三栏骨架、会话列表、聊天窗口与消息气泡"
```

---

### Task F9: 联系人、群组、搜索与设置

**Files:**
- Create: `front-ui-vue/src/views/home/components/ContactsPanel.vue`（替代旧 contact.vue）
- Create: `front-ui-vue/src/views/home/components/GroupsPanel.vue`
- Create: `front-ui-vue/src/views/home/components/SettingsPanel.vue`（替代旧 updateUserInfo.vue）
- Create: `front-ui-vue/src/components/UserCard.vue`（重写旧 components/card/userInfo.vue）
- Create: `front-ui-vue/src/components/SearchUserModal.vue`（重写旧 plusFirendGroup.vue）
- Create: `front-ui-vue/src/components/CreateGroupModal.vue`（重写旧 createGroup.vue）
- Create: `front-ui-vue/src/components/AvatarUpload.vue`

**Interfaces:**
- Consumes: `useContactsStore/useAuthStore/useChatStore`、`user.getDetail`、`file.uploadImage`
- Produces: 中栏三个面板组件 + 四个通用组件

- [ ] **Step 1: ContactsPanel.vue**

好友列表（`contacts.friends`）：头像+在线点+昵称；点击 → 构造私聊会话键 `u:{min}_{max}`（工具函数 `privateConvId(a: number, b: number): string`，放 `src/utils/conversation.ts` 并导出 `groupConvId(gid)`）→ `chatStore.openConversation(convId)` 且 HomeView 切回 chats 面板（emit('go-chat')）。顶部按钮：添加好友（打开 SearchUserModal）。

- [ ] **Step 2: GroupsPanel.vue**

我的群列表（`contacts.groups`）：群头像+名称+成员数；点击 → `chatStore.openConversation(groupConvId(g.id))` + emit('go-chat')。顶部按钮：发起群聊（CreateGroupModal）、搜索群（a-input-search → `contacts.searchGroups(kw)` 结果列表 + 加入按钮 → `contacts.joinGroup`）。

- [ ] **Step 3: UserCard.vue + SearchUserModal.vue**

UserCard props：`userId: number`。`onMounted/watch(userId)` → `user.getDetail(userId)` → 弹层展示：头像、昵称、identity、个性签名、注册时间；`is_friend || is_self` 时显示邮箱/手机号；底部按钮：非好友 →"添加好友"（`contacts.addFriend` 成功后刷新 detail）；好友 →"发消息"（跳私聊）。
SearchUserModal：输入关键词（防抖 300ms）→ `contacts.searchUsers(kw)` → 结果行（头像/昵称/identity/是否已是好友徽标）→ 点击行打开 UserCard。

- [ ] **Step 4: CreateGroupModal.vue**

群名输入 + 好友多选（a-checkbox 列表来自 `contacts.friends`）→ `contacts.createGroup(name, memberIds)` → 成功后关闭并 emit('created', groupItem) → 父组件直接 `chatStore.openConversation(groupConvId(g.id))`。

- [ ] **Step 5: SettingsPanel.vue + AvatarUpload.vue**

SettingsPanel：展示/编辑当前用户资料（a-form：昵称、手机号、个性签名、生日 a-date-picker）；修改邮箱需先点"发送验证码"（`authApi.sendEmailCode(newEmail)`，60s 倒计时）并填入验证码；保存 → `auth.updateProfile({...})`。底部：退出登录按钮（`auth.logout()` → router.push('/login')）。
AvatarUpload：a-upload（`customRequest` 调 `file.uploadImage`，接受 image/*，≤5MB 前端预检）→ 成功后 emit('uploaded', url) → SettingsPanel 写入表单 avatar 字段（保存时随 updateProfile 提交）。

- [ ] **Step 6: 验证 + Commit**

Run: `yarn type-check`；`yarn dev` 手工冒烟：加好友→私聊、建群→群聊、改资料、换头像、退出登录。

```bash
git add front-ui-vue/src/views/home/components front-ui-vue/src/components front-ui-vue/src/utils/conversation.ts
git commit -m "feat(web): 联系人/群组/搜索/设置面板与通用组件（UserCard/建群/头像上传）"
```

---

### Task F10: 旧代码清理、全量验证与文档

**Files:**
- Delete: `front-ui-vue/src/utils/websocket.ts`、`src/utils/request.ts`、`src/utils/toolsValidate.ts`、`src/utils/is/`、`src/utils/formatTime.ts`（formatPast 已迁入 format.ts）
- Delete: 旧视图/组件：`src/views/home/components/content.vue`、`sider.vue`、`header.vue`、`footer.vue`、`viewsSoder/`（整目录）、`src/components/card/`、`src/components/emoji/`、`src/components/HelloWorld.vue`、`src/components/__tests__/`、`src/components/base.css`/`main.css`（如仅被旧代码引用；逐个 grep 确认后删）、`src/components/icons/`（未使用的 Icon*.vue）
- Delete: 旧 api：`src/api/userBasic/`、`src/api/group/`、`src/api/userFriend/`、`src/api/file/`、`src/api/auth/`（被 F3 新文件取代；注意新旧同名的 auth/file 以 F3 版本为准）
- Modify: `README.md`（根）：前端启动/构建说明、环境变量说明、指向设计文档

- [ ] **Step 1: 删除旧代码**

删除前对每个目标执行 `grep -rn "文件名" src/` 确认无新代码引用；删除后：

```bash
yarn type-check && yarn lint && yarn build
```

Expected: 三者全部零错误零告警，dist 构建成功。

- [ ] **Step 2: 全链路手工联调（对照 spec §5 清单，需本地 MySQL/Redis/Mongo + 后端已启动）**

> 注册（邮箱验证码；无邮件服务时 redis-cli 手工注入验证码）→ 登录（图形码）→ 搜索用户 → 加好友（对方视角收到列表刷新）→ 私聊互发（双方实时上屏、会话列表置顶刷新）→ 非好友发消息 3 条后第 4 条被 429 拦截并提示 → 建群 → 加群 → 群聊 → 断网重连（关 wifi/杀后端 10s 内恢复，WS 自动重连且在线状态正确）→ 双端登录同一账号（两端都能实时收消息）→ 登出 → 踢人场景（任一端登出后，同账号其余在线端收到 system.kick 被断开并回到登录页）→ 图片上传并在聊天中显示

发现的任何前后端契约不一致：以后端计划/spec 为准修正前端，若属后端缺陷则回改后端并同步更新两份计划文档中的对应契约。

- [ ] **Step 3: 更新 README 并提交**

根 README.md 重写"运行"章节：后端（Go 1.22+、schema.sql、config.toml、`go run ./cmd/server`）、前端（Node 20+、yarn、env、`yarn dev`/`yarn build`）、依赖服务（MySQL/Redis/Mongo 版本与默认端口）、指向 `docs/superpowers/specs/2026-09-12-kk-chat-refactor-design.md`。

```bash
git add -A front-ui-vue README.md
git commit -m "refactor(web)!: 移除旧组件/工具/api，前端切换至新分层架构"
```

- [ ] **Step 4: 合并分支**

```bash
git checkout master && git merge refactor/clean-architecture
```

（若用户希望先自行体验再合并，停在上一任务，等用户指令。）

---

## 前端完成标准

- `yarn type-check && yarn lint && yarn build` 全部干净
- 旧文件（utils/websocket.ts、request.ts、toolsValidate.ts、viewsSoder/ 等）不复存在
- F10 Step 2 联调清单全通过
- package.json 无 vuex/emoji-mart-vue/v-emoji-picker/js-pinyin 依赖

## 两份计划的执行顺序

后端 B1→B16 全部完成后，才开始前端 F1→F10（前端联调依赖新后端跑通）。同一分支 `refactor/clean-architecture`，合并回 master 为最后一步。
