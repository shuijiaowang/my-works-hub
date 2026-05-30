import { createRouter, createWebHistory } from 'vue-router'
import AdminView from '@/views/AdminView.vue'
import ProjectsAllView from '@/views/ProjectsAllView.vue'
import ProjectDetailView from '@/views/ProjectDetailView.vue'

const DEFAULT_TITLE = 'My Works Hub'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', redirect: '/projects' },
    { path: '/projects', name: 'projects', component: ProjectsAllView, meta: { title: '我的作品' } },
    { path: '/projects/:folderName', name: 'project-detail', component: ProjectDetailView, meta: { title: '项目详情' } },
    { path: '/admin', name: 'admin', component: AdminView, meta: { title: '管理' } },
  ],
})

router.afterEach((to) => {
  const piece = to.meta.title
  document.title = piece ? `${piece} | ${DEFAULT_TITLE}` : DEFAULT_TITLE
})

export default router
