<script setup>
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { adminLoggedIn } from '@/adminAuth'

const route = useRoute()
</script>

<template>
  <div id="app">
    <nav v-if="adminLoggedIn" class="nav">
      <div class="nav-inner">
        <RouterLink to="/projects" class="nav-brand">My Works</RouterLink>
        <div class="nav-links">
          <RouterLink
            to="/projects"
            class="nav-link"
            :class="{ active: route.path === '/projects' || route.path.startsWith('/projects/') }"
          >
            作品集
          </RouterLink>
          <RouterLink to="/admin" class="nav-link" :class="{ active: route.path === '/admin' }">管理</RouterLink>
        </div>
      </div>
    </nav>
    <RouterView />
  </div>
</template>

<style scoped>
#app {
  min-height: 100vh;
}
.nav {
  border-bottom: 1px solid var(--portfolio-border);
  background: var(--portfolio-surface);
  box-shadow: var(--portfolio-shadow-sm);
}
.nav-inner {
  max-width: 1120px;
  margin: 0 auto;
  padding: 0 24px;
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.nav-brand {
  font-size: 14px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--portfolio-text);
  font-family: var(--portfolio-font-mono);
}
.nav-brand:hover {
  color: var(--portfolio-accent);
}
.nav-links {
  display: flex;
  gap: 4px;
}
.nav-link {
  padding: 6px 14px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  color: var(--portfolio-text-muted);
  text-decoration: none;
  transition: color 0.15s, background 0.15s;
}
.nav-link:hover {
  color: var(--portfolio-text);
  background: var(--portfolio-border-light);
}
.nav-link.active {
  color: var(--portfolio-accent);
  background: var(--portfolio-accent-soft);
  font-weight: 600;
}
</style>
