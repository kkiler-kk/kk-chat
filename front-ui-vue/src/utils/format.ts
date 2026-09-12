function toDate(input: Date | string): Date {
  return input instanceof Date ? input : new Date(input)
}

function pad(n: number): string {
  return n < 10 ? `0${n}` : String(n)
}

/** 相对时间：刚刚 / x分钟前 / x小时前 / x天前 / MM-DD HH:mm / YYYY-MM-DD */
export function formatPast(input: Date | string): string {
  const d = toDate(input)
  const now = Date.now()
  const diff = now - d.getTime()
  if (Number.isNaN(diff)) return ''
  if (diff < 60_000) return '刚刚'
  if (diff < 3_600_000) return `${Math.floor(diff / 60_000)}分钟前`
  if (diff < 86_400_000) return `${Math.floor(diff / 3_600_000)}小时前`
  if (diff < 7 * 86_400_000) return `${Math.floor(diff / 86_400_000)}天前`
  if (d.getFullYear() === new Date().getFullYear()) {
    return `${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
  }
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}

/** 格式化时间，支持 YYYY MM DD HH mm ss 占位符 */
export function formatTime(input: Date | string, fmt = 'YYYY-MM-DD HH:mm:ss'): string {
  const d = toDate(input)
  if (Number.isNaN(d.getTime())) return ''
  return fmt
    .replace('YYYY', String(d.getFullYear()))
    .replace('MM', pad(d.getMonth() + 1))
    .replace('DD', pad(d.getDate()))
    .replace('HH', pad(d.getHours()))
    .replace('mm', pad(d.getMinutes()))
    .replace('ss', pad(d.getSeconds()))
}
