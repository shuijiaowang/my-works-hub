import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import AdminView from '@/views/AdminView.vue'

const DEFAULT_TITLE = 'My Works Hub'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'home', component: HomeView, meta: { title: '主页' } },
    { path: '/admin', name: 'admin', component: AdminView, meta: { title: '管理' } },
  ],
})

router.afterEach((to) => {
  const piece = to.meta.title
  document.title = piece ? `${piece} | ${DEFAULT_TITLE}` : DEFAULT_TITLE
})

export default router
