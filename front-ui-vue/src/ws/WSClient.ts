import type { WsEnvelope, WsServerEvent, WsServerEvents } from '@/types/ws'

type Handler<E extends WsServerEvent> = (data: WsServerEvents[E]) => void

const HEARTBEAT_INTERVAL = 60_000 // 与后端 presence.heartbeat_seconds 对齐
const RECONNECT_BASE = 1_000
const RECONNECT_MAX = 30_000

/**
 * 类型化 WebSocket 客户端：
 * - on(event, handler) 返回解绑函数，事件 payload 由 WsServerEvents 约束
 * - 指数退避重连（1s→2s→4s…上限 30s），手动 close 不重连
 * - 内置心跳（每 60s 发应用层 ping）
 * - system.kick 触发 kickHandler 并停止重连
 */
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
    if (this.socket && (this.token !== token || !this.manualClose)) {
      // 已有连接：同 token 且非手动关闭则忽略；否则先断开
      if (this.token === token && !this.manualClose) return
      this.close()
    }
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
    return () => {
      set?.delete(handler as Handler<WsServerEvent>)
    }
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
      const data = env.data as { reason: string } | undefined
      this.manualClose = true
      this.close()
      this.kickHandler(data?.reason ?? '')
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
