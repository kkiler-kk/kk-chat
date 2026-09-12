// WS 端到端验证：presence 广播 / ping-pong / 私聊消息推送 / recent_updated / 群聊广播 / 登出踢人
import { readFileSync } from 'node:fs'

const API = 'http://localhost:9345'
const WS = 'ws://localhost:9345/api/v1/ws'
const t1 = readFileSync('/tmp/t1.txt', 'utf8').trim()
const t2 = readFileSync('/tmp/t2.txt', 'utf8').trim()

const results = []
function check(name, ok, extra = '') {
  results.push({ name, ok })
  console.log(`${ok ? '✅' : '❌'} ${name} ${extra}`)
}

function connect(token) {
  return new Promise((resolve, reject) => {
    const ws = new WebSocket(`${WS}?token=${encodeURIComponent(token)}`)
    const inbox = []
    const waiters = []
    ws.onmessage = (ev) => {
      const env = JSON.parse(ev.data)
      for (let i = waiters.length - 1; i >= 0; i--) {
        if (waiters[i].event === env.event) {
          const w = waiters.splice(i, 1)[0]
          w.resolve(env)
          return
        }
      }
      inbox.push(env)
    }
    ws.onopen = () =>
      resolve({
        ws,
        inbox,
        wait(event, timeout = 4000) {
          const hit = inbox.findIndex((e) => e.event === event)
          if (hit >= 0) return Promise.resolve(inbox.splice(hit, 1)[0])
          return new Promise((res, rej) => {
            const timer = setTimeout(() => rej(new Error(`等待 ${event} 超时`)), timeout)
            waiters.push({
              event,
              resolve: (env) => {
                clearTimeout(timer)
                res(env)
              },
            })
          })
        },
      })
    ws.onerror = (e) => reject(new Error('ws error ' + JSON.stringify(e)))
    setTimeout(() => reject(new Error('ws 连接超时')), 5000)
  })
}

async function api(method, path, token, body) {
  const res = await fetch(API + path, {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: body ? JSON.stringify(body) : undefined,
  })
  return { status: res.status, body: await res.json() }
}

try {
  // 1. kk001 先连接
  const c1 = await connect(t1)
  check('kk001 WS 连接', true)

  // 2. kk002 连接 → kk001 应收到 presence.changed(online)
  const c2 = await connect(t2)
  const presence = await c1.wait('presence.changed')
  check(
    'presence.changed 上线广播',
    presence.data.user_id === 2 && presence.data.online === true,
    JSON.stringify(presence.data),
  )

  // 3. ping → pong
  c2.ws.send(JSON.stringify({ event: 'ping' }))
  const pong = await c2.wait('pong')
  check('ping/pong 心跳', !!pong.ts, `ts=${pong.ts}`)

  // 4. kk001 HTTP 发私聊消息 → kk002 收到 chat.message + chat.recent_updated；kk001 自己也收到 recent_updated
  const content = `E2E测试消息-${Date.now()}`
  const sent = await api('POST', '/api/v1/messages', t1, {
    conversation_id: 'u:1_2',
    content,
    content_type: 'text',
  })
  check('HTTP 发消息', sent.body.code === 0 && sent.body.data.content === content)
  const msg = await c2.wait('chat.message')
  check(
    'kk002 收到 chat.message',
    msg.data.content === content && msg.data.sender_id === 1 && msg.data.conversation_id === 'u:1_2',
  )
  const recent2 = await c2.wait('chat.recent_updated')
  check(
    'kk002 收到 recent_updated(peer=发送者)',
    recent2.data.peer_id === 1 && recent2.data.last_message_content === content,
  )
  const recent1 = await c1.wait('chat.recent_updated')
  check('kk001 收到自己的 recent_updated(peer=对方)', recent1.data.peer_id === 2)

  // 5. 群聊：kk002 在 g:1 发言 → kk001 收到 chat.message
  const gcontent = `群E2E-${Date.now()}`
  const gsent = await api('POST', '/api/v1/messages', t2, {
    conversation_id: 'g:1',
    content: gcontent,
    content_type: 'text',
  })
  check('群聊 HTTP 发消息', gsent.body.code === 0)
  const gmsg = await c1.wait('chat.message')
  check(
    'kk001 收到群聊推送',
    gmsg.data.conversation_id === 'g:1' && gmsg.data.content === gcontent && gmsg.data.sender_id === 2,
  )

  // 6. kk001 登出 → 两端都应收到 system.kick（登出即全局踢下线）
  const lo = await api('POST', '/api/v1/auth/logout', t1)
  check('登出接口', lo.body.code === 0)
  const kick1 = await c1.wait('system.kick')
  check('kk001 收到 system.kick', kick1.data.reason === 'logout')

  // 7. 登出后 token 失效
  const me = await api('GET', '/api/v1/users/me', t1)
  check('登出后 token 失效(401)', me.status === 401 && me.body.code === 20003)

  c2.ws.close()
} catch (e) {
  console.error('❌ 异常:', e.message)
  results.push({ name: '异常:' + e.message, ok: false })
}

const failed = results.filter((r) => !r.ok)
console.log(`\n=== E2E 结果: ${results.length - failed.length}/${results.length} 通过 ===`)
process.exit(failed.length ? 1 : 0)
