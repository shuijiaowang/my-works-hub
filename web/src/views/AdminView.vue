<script setup>
import { ref } from 'vue'
import { adminLoggedIn } from '@/adminAuth'
import AdminLoginCard from '@/components/admin/AdminLoginCard.vue'
import CreateProjectDialog from '@/components/admin/CreateProjectDialog.vue'

const createOpen = ref(false)
const lastCreated = ref(null)
</script>

<template>
  <main class="page">
    <AdminLoginCard />

    <section v-if="adminLoggedIn" class="toolbar">
      <el-button type="primary" @click="createOpen = true">新建项目</el-button>
      <span v-if="lastCreated" class="hint">最近创建：{{ lastCreated?.Name || lastCreated?.name }}</span>
    </section>

    <CreateProjectDialog v-model="createOpen" @created="(p) => (lastCreated = p)" />
  </main>
</template>

<style scoped>
.page {
  max-width: 960px;
  margin: 0 auto;
  padding: 24px;
}
.toolbar {
  margin-top: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
}
.hint {
  color: rgba(0, 0, 0, 0.6);
}
</style>
