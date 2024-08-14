import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/home/HomeView.vue'
import NotFound from '@/views/ErrorStatus/NotFound/index.vue'
import { Local } from "@/utils/storage";
import { checkAuth } from '@/utils';

const routes = [
  {
    path: '/',
    redirect: '/home',
  },
  {
    path: '/home',
    name: "Home",
    component: HomeView,
    meta: {
      title: "首页-KK-Chat",
      isLogin: true,
    }
  },
  {
    path: '/about',
    name: 'about',
    // route level code-splitting
    // this generates a separate chunk (About.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('@/views/AboutView.vue')
  },
  {
    path: '/login',
    name: 'login',
    // route level code-splitting
    // this generates a separate chunk (About.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('@/views/login/index.vue'),
    meta: {
      title: "登录",
      svg: "user",
      hidden: true,
      isLogin: false,
    },
  },
  {
    path: '/404',
    name: 'NotFound',
    component: NotFound,
    meta: {
      title: "404-页面不存在",
      isLogin: false,
    }
  },
  {
    path: "/forbidden",
    name: "Forbidden",
    component: () => import('@/views/ErrorStatus/Forbidden/index.vue')
  },
  {
    path: "/server-error",
    name: "ServerError",
    component: () => import('@/views/ErrorStatus/ErrorServer/index.vue')
  },
  // 所有未定义的路由全部重定到404页
  {
    path: "/:pathMatch(.*)",
    redirect: '/404',
  },
]
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})
// 添加全局前置守卫 用于检测用户是否已登录
router.beforeEach((to, from, next) => {
  console.log("to", to)
  if (to.meta.isLogin) {
    if (Local.get('token') == '') { // 用户尚未登录
      next('/login')
    } else { // 用户已经登录
      if (!checkAuth(to)) {  // 检查用户访问权限 返回到没有返回权限页面
        next({name: "Forbidden"})
      }else {
        next()
      }
    }
  } else {
    next()
  }
})
export default router
