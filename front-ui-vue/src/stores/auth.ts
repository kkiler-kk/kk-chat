import { defineStore } from 'pinia'
import * as authApi from '@/api/auth'
import type { RegisterParams } from '@/api/auth'
import type { UpdateMeParams } from '@/api/user'
import * as userApi from '@/api/user'
import { ApiError, setUnauthorizedHandler } from '@/api/client'
import { Local, StorageKeys } from '@/utils/storage'
import { wsClient } from '@/ws/WSClient'
import { setupWsBindings } from '@/ws/bindings'
import { useChatStore } from './chat'
import { useContactsStore } from './contacts'
import type { LoginResult, UserInfo } from '@/types/api'

let unbindWs: (() => void) | null = null

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
        account,
        password,
        captcha_id: captchaId,
        captcha_answer: captchaAnswer,
      })
      this.applySession(result.token, result.user)
    },
    async register(p: RegisterParams) {
      await authApi.register(p)
    },
    applySession(token: string, user: UserInfo) {
      this.token = token
      this.user = user
      Local.set(StorageKeys.token, token)
      Local.set(StorageKeys.userInfo, user)
      setUnauthorizedHandler(() => this.setUnauthorized())
      wsClient.setKickHandler(() => this.setUnauthorized())
      wsClient.connect(token)
      unbindWs?.()
      unbindWs = setupWsBindings()
    },
    async logout() {
      try {
        await authApi.logout()
      } catch {
        /* 忽略登出接口错误，本地照清 */
      }
      this.clearSession()
    },
    /** 401/token失效/被踢：清理本地并跳转登录 */
    setUnauthorized() {
      this.clearSession()
      window.location.href = '/login'
    },
    clearSession() {
      unbindWs?.()
      unbindWs = null
      wsClient.close()
      this.token = null
      this.user = null
      Local.remove(StorageKeys.token)
      Local.remove(StorageKeys.userInfo)
      useChatStore().$resetAll()
      useContactsStore().$reset()
    },
    /** 刷新页面后恢复：校验本地 token 是否仍有效 */
    async restore(): Promise<boolean> {
      if (!this.token) return false
      try {
        this.user = await userApi.getMe()
        Local.set(StorageKeys.userInfo, this.user)
        setUnauthorizedHandler(() => this.setUnauthorized())
        wsClient.setKickHandler(() => this.setUnauthorized())
        wsClient.connect(this.token)
        unbindWs?.()
        unbindWs = setupWsBindings()
        return true
      } catch (e) {
        if (e instanceof ApiError) {
          unbindWs?.()
          unbindWs = null
          wsClient.close()
          this.token = null
          this.user = null
          Local.remove(StorageKeys.token)
          Local.remove(StorageKeys.userInfo)
        }
        return false
      }
    },
    async updateProfile(p: UpdateMeParams) {
      await userApi.updateMe(p)
      this.user = await userApi.getMe()
      Local.set(StorageKeys.userInfo, this.user)
    },
  },
})
