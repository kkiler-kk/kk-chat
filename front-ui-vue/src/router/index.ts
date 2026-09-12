import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { wsClient } from '@/ws/WSClient'

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
      path: '/403',
      name: 'forbidden',
      meta: { public: true },
      component: () => import('@/views/ErrorStatus/Forbidden/index.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      meta: { public: true },
      component: () => import('@/views/ErrorStatus/NotFound/index.vue'),
    },
  ],
})

// 会话恢复只并发执行一次（校验 token + 重连 WS）
let restorePromise: Promise<boolean> | null = null

function ensureSession(): Promise<boolean> {
  restorePromise ??= useAuthStore()
    .restore()
    .finally(() => {
      // 失败后允许下次导航重试
      restorePromise?.then((ok) => {
        if (!ok) restorePromise = null
      })
    })
  return restorePromise
}

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (to.meta.public) {
    // 已登录访问登录页 → 回 home
    if (to.name === 'login' && auth.isLoggedIn && wsClient.connected) return { name: 'home' }
    return true
  }
  // 已登录且 WS 已连接：放行
  if (auth.isLoggedIn && wsClient.connected) return true
  // 有本地 token：恢复会话（校验 + 重连）
  if (auth.token) {
    const ok = await ensureSession()
    if (ok) return true
  }
  return { name: 'login', query: { redirect: to.fullPath } }
})

export default router
