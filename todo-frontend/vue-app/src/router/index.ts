import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/',         redirect: '/todos' },
    { path: '/login',    component: () => import('@/views/LoginView.vue'),    meta: { guest: true } },
    { path: '/register', component: () => import('@/views/RegisterView.vue'), meta: { guest: true } },
    { path: '/mfa',      component: () => import('@/views/MfaView.vue'),      meta: { guest: true } },
    {
      path: '/todos',
      component: () => import('@/views/TodosView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/todos/:id',
      component: () => import('@/views/TodoDetailView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/settings',
      component: () => import('@/views/SettingsView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach(to => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isAuthenticated) return '/login'
  if (to.meta.guest && auth.isAuthenticated)          return '/todos'
})

export default router
