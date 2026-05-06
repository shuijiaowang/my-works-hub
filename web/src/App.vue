<script setup>
import { RouterView } from 'vue-router'
</script>

<template>
  <div id="app">
    <header class="top-nav">
      <div class="left-links">
        <router-link to="/">主页</router-link>
        <router-link to="/admin">管理</router-link>
      </div>
    </header>
    <RouterView />
  </div>
</template>

<style scoped>
#app {
  min-height: 100vh;
}

.top-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 64px;
  padding: 0 24px;
  border-bottom: 1px solid var(--color-border);
}

.left-links {
  display: flex;
  gap: 16px;
  align-items: center;
}

a.router-link-active {
  color: hsla(160, 100%, 30%, 1);
  font-weight: 600;
}
</style>