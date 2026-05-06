<script setup>
import { ref } from 'vue'
import { text as exampleText } from '@/api/example'

const result = ref('')
const loading = ref(false)

const callExample = async () => {
  loading.value = true
  try {
    const res = await exampleText({})
    result.value = typeof res?.data === 'string' ? res.data : JSON.stringify(res, null, 2)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="page">
    <h1 class="title">个人作品集</h1>
    <p class="desc">用于分享介绍与提供下载。主页无需登录。</p>

    <section class="card">
      <div class="row">
        <el-button type="primary" :loading="loading" @click="callExample">调用示例 API</el-button>
      </div>
      <pre v-if="result" class="pre">{{ result }}</pre>
    </section>
  </main>
</template>

<style scoped>
.page {
  max-width: 960px;
  margin: 0 auto;
  padding: 24px;
}
.title {
  font-size: 28px;
  margin: 0 0 8px;
}
.desc {
  margin: 0 0 16px;
  color: rgba(0, 0, 0, 0.65);
}
.card {
  padding: 16px;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  background: #fff;
}
.row {
  display: flex;
  gap: 12px;
  align-items: center;
}
.pre {
  margin-top: 12px;
  padding: 12px;
  background: #0b1020;
  color: #e6e6e6;
  border-radius: 10px;
  overflow: auto;
}
</style>