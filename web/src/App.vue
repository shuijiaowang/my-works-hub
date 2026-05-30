<script setup>
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { adminLoggedIn } from '@/adminAuth'

const route = useRoute()
</script>

<template>
  <div id="app">
    <nav v-if="adminLoggedIn" class="nav">
      <RouterLink
        to="/projects"
        class="nav-link"
        :class="{ active: route.path === '/projects' || route.path.startsWith('/projects/') }"
      >
        项目全查
      </RouterLink>
      <RouterLink to="/admin" class="nav-link" :class="{ active: route.path === '/admin' }">管理</RouterLink>
    </nav>
    <RouterView />
  </div>
</template>

<style scoped>
#app {
  min-height: 100vh;
}
.nav {
  display: flex;
  gap: 16px;
  padding: 12px 24px;
  border-bottom: 1px solid var(--color-border, #e5e5e5);
  background: #fff;
}
.nav-link {
  color: rgba(0, 0, 0, 0.65);
  text-decoration: none;
}
.nav-link.active {
  color: var(--el-color-primary, #409eff);
  font-weight: 600;
}
</style>
