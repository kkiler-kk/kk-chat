/** 私聊会话键：u:{小ID}_{大ID}（与后端 domain.NewPrivateConversation 一致） */
export function privateConvId(a: number, b: number): string {
  const lo = Math.min(a, b)
  const hi = Math.max(a, b)
  return `u:${lo}_${hi}`
}

/** 群聊会话键：g:{群ID} */
export function groupConvId(gid: number): string {
  return `g:${gid}`
}
