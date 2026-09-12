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
